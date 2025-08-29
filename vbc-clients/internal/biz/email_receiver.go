package biz

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"mime"
	"net"
	"regexp"
	"strings"
	"time"

	"vbc/configs"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	gomail "github.com/emersion/go-message/mail"
	"golang.org/x/net/proxy"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// EmailAttachment represents an email attachment
type EmailAttachment struct {
	Filename    string
	ContentType string
	Content     []byte
}

// InlineResource represents an inline resource (like embedded images)
type InlineResource struct {
	ContentID   string // Content-ID header for referencing in HTML
	ContentType string
	Content     []byte
}

// IncomingEmail represents a simplified incoming email fetched from IMAP.
type IncomingEmail struct {
	// UID is the IMAP unique identifier for the message within the mailbox.
	UID uint32
	// MessageID is the RFC 5322 Message-Id header value, if present.
	MessageID string
	From      string
	To        string
	Subject   string
	Body      string
	// HTMLBody contains the HTML version of the email body
	HTMLBody string
	// Attachments contains file attachments
	Attachments []EmailAttachment
	// InlineResources contains inline resources like embedded images
	InlineResources []InlineResource
	// Raw holds the full RFC822 message bytes for high-fidelity forwarding.
	Raw []byte
	// Date is the email sending time from Envelope.Date
	Date time.Time
}

// FetchRecentEmails connects to Gmail IMAP with the given credentials and
// fetches up to maxCount most recent messages in INBOX. It returns simplified
// email structures for downstream processing.
//
// The function is designed to be unit-testable by isolating the IMAP interaction
// and returning deterministic data structures.
func FetchRecentEmails(username string, password string, maxCount uint32) ([]IncomingEmail, error) {
	if maxCount == 0 {
		maxCount = 1
	}

	c, mbox, err := connectAndSelectINBOX(username, password, false)
	if err != nil {
		return nil, err
	}
	defer c.Logout()

	if mbox.Messages == 0 {
		return nil, nil
	}

	// Determine the range for the latest messages
	var start uint32
	if mbox.Messages >= maxCount {
		start = mbox.Messages - (maxCount - 1)
	} else {
		start = 1
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(start, mbox.Messages)

	// Use the shared fetch function
	results, err := fetchMessagesBySeqSet(c, seqSet, username)
	if err != nil {
		return nil, fmt.Errorf("fetch messages: %w", err)
	}

	return results, nil
}

// FetchEmailsByUIDRange fetches emails by UID range for incremental processing.
// It uses UID search to find messages with UID > startUID and returns up to limit emails.
// Returns the emails list and whether there are more emails available.
func FetchEmailsByUIDRange(username string, password string, startUID uint32, limit int) ([]IncomingEmail, bool, error) {
	if limit <= 0 {
		limit = 10 // default limit
	}

	c, mbox, err := connectAndSelectINBOX(username, password, false)
	if err != nil {
		return nil, false, err
	}
	defer c.Logout()

	if mbox.Messages == 0 {
		return nil, false, nil
	}

	// Get the highest UID in the mailbox to avoid search issues
	// when startUID is greater than the actual max UID
	if startUID >= mbox.UidNext-1 {
		// startUID is already at or beyond the latest message
		return nil, false, nil
	}

	// Use UID SEARCH ALL first to get all UIDs, then filter
	criteria := &imap.SearchCriteria{}
	allUIDs, err := c.UidSearch(criteria)
	if err != nil {
		return nil, false, fmt.Errorf("UID search all: %w", err)
	}

	if len(allUIDs) == 0 {
		return nil, false, nil
	}

	// Filter UIDs that are greater than startUID
	var filteredUIDs []uint32
	for _, uid := range allUIDs {
		if uid > startUID {
			filteredUIDs = append(filteredUIDs, uid)
		}
	}

	if len(filteredUIDs) == 0 {
		return nil, false, nil
	}

	// Limit the number of UIDs to process
	hasMore := len(filteredUIDs) > limit
	if hasMore {
		filteredUIDs = filteredUIDs[:limit]
	}

	// Create sequence set for the UIDs
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(filteredUIDs...)

	// Fetch messages by UID
	emails, err := fetchMessagesByUIDSet(c, seqSet, username)
	if err != nil {
		return nil, false, fmt.Errorf("fetch messages by UID: %w", err)
	}

	return emails, hasMore, nil
}

// FetchEmailByUID fetches a single email by its UID
func FetchEmailByUID(username string, password string, uid uint32) (*IncomingEmail, error) {
	c, _, err := connectAndSelectINBOX(username, password, false)
	if err != nil {
		return nil, err
	}
	defer c.Logout()

	// Create sequence set for the single UID
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uid)

	emails, err := fetchMessagesByUIDSet(c, seqSet, username)
	if err != nil {
		return nil, fmt.Errorf("fetch message by UID %d: %w", uid, err)
	}

	if len(emails) == 0 {
		return nil, nil // Not found
	}

	return &emails[0], nil
}

// GetMailboxUIDNext returns the next UID that would be assigned to a new message
func GetMailboxUIDNext(username string, password string) (uint32, error) {
	c, mbox, err := connectAndSelectINBOX(username, password, false)
	if err != nil {
		return 0, err
	}
	defer c.Logout()

	return mbox.UidNext, nil
}

// ValidateEmailCredentials tests if the given credentials can successfully connect to IMAP
func ValidateEmailCredentials(username string, password string) error {
	c, _, err := connectAndSelectINBOX(username, password, true) // Read-only mode for validation
	if err != nil {
		return fmt.Errorf("cannot access: %w", err)
	}
	defer c.Logout()

	return nil
}

// connectAndSelectINBOX creates an IMAP connection and selects INBOX, returns client and mailbox status
func connectAndSelectINBOX(username string, password string, readonly bool) (*client.Client, *imap.MailboxStatus, error) {
	c, err := connectIMAP(username, password)
	if err != nil {
		return nil, nil, err
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", readonly)
	if err != nil {
		c.Logout()
		return nil, nil, fmt.Errorf("select INBOX: %w", err)
	}

	return c, mbox, nil
}

// connectIMAP creates an IMAP connection with proxy support (reused logic)
func connectIMAP(username string, password string) (*client.Client, error) {
	tlsConfig := &tls.Config{ServerName: "imap.gmail.com"}

	var (
		c   *client.Client
		err error
	)

	if configs.IsDev() {
		baseDialer := &net.Dialer{Timeout: 15 * time.Second}
		p, perr := proxy.SOCKS5("tcp", "127.0.0.1:7890", nil, baseDialer)
		if perr != nil {
			return nil, fmt.Errorf("init proxy: %w", perr)
		}
		conn, derr := p.Dial("tcp", "imap.gmail.com:993")
		if derr != nil {
			return nil, fmt.Errorf("dial via proxy: %w", derr)
		}
		tlsConn := tls.Client(conn, tlsConfig)
		if herr := tlsConn.Handshake(); herr != nil {
			return nil, fmt.Errorf("tls handshake: %w", herr)
		}
		c, err = client.New(tlsConn)
		if err != nil {
			return nil, fmt.Errorf("init imap client: %w", err)
		}
	} else {
		c, err = client.DialTLS("imap.gmail.com:993", tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("connect IMAP: %w", err)
		}
	}

	if err := c.Login(username, password); err != nil {
		return nil, fmt.Errorf("login IMAP: %w", err)
	}

	return c, nil
}

// fetchMessagesBySeqSet fetches messages using a sequence set (for sequence numbers, not UIDs)
func fetchMessagesBySeqSet(c *client.Client, seqSet *imap.SeqSet, username string) ([]IncomingEmail, error) {
	return fetchMessages(c, seqSet, username, false) // false = use regular fetch
}

// fetchMessagesByUIDSet fetches messages using a UID sequence set (reused logic)
func fetchMessagesByUIDSet(c *client.Client, seqSet *imap.SeqSet, username string) ([]IncomingEmail, error) {
	return fetchMessages(c, seqSet, username, true) // true = use UID fetch
}

// fetchMessages is the common function for fetching messages, either by sequence numbers or UIDs
func fetchMessages(c *client.Client, seqSet *imap.SeqSet, username string, useUID bool) ([]IncomingEmail, error) {
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid, section.FetchItem()}
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	// Choose fetch method based on useUID parameter
	if useUID {
		go func() { done <- c.UidFetch(seqSet, items, messages) }()
	} else {
		go func() { done <- c.Fetch(seqSet, items, messages) }()
	}

	var results []IncomingEmail

	for msg := range messages {
		if msg == nil || msg.Envelope == nil {
			continue
		}

		email := parseMessageToIncomingEmail(msg, username)
		results = append(results, email)
	}

	if err := <-done; err != nil {
		return results, fmt.Errorf("fetch messages: %w", err)
	}

	return results, nil
}

// parseMessageToIncomingEmail converts an IMAP message to IncomingEmail struct
func parseMessageToIncomingEmail(msg *imap.Message, username string) IncomingEmail {
	env := msg.Envelope
	subject := env.Subject

	// Build From string
	var fromDisplay string
	if len(env.From) > 0 {
		fromDisplay = formatEmailAddress(env.From[0])
	}

	// Build To string (first recipient). If empty, fallback to username.
	var toDisplay string
	if len(env.To) > 0 {
		toDisplay = formatEmailAddress(env.To[0])
	}
	if strings.TrimSpace(toDisplay) == "" {
		toDisplay = username
	}

	// Extract body and raw content
	bodyContent, htmlContent, attachments, inlineResources, rawBytes := extractEmailContent(msg)

	// Get email date from envelope
	emailDate := time.Now() // fallback to now if no date
	if !env.Date.IsZero() {
		emailDate = env.Date
	}

	return IncomingEmail{
		UID:             msg.Uid,
		MessageID:       env.MessageId,
		From:            fromDisplay,
		To:              toDisplay,
		Subject:         subject,
		Body:            bodyContent,
		HTMLBody:        htmlContent,
		Attachments:     attachments,
		InlineResources: inlineResources,
		Raw:             rawBytes,
		Date:            emailDate,
	}
}

// formatEmailAddress formats an IMAP address to display string
func formatEmailAddress(addr *imap.Address) string {
	mailbox := addr.MailboxName
	host := addr.HostName

	if addr.PersonalName != "" {
		return fmt.Sprintf("%s <%s@%s>", addr.PersonalName, mailbox, host)
	}
	return fmt.Sprintf("%s@%s", mailbox, host)
}

// normalizeEmailCharset preprocesses raw email content to handle charset issues
func normalizeEmailCharset(rawBytes []byte) []byte {
	rawStr := string(rawBytes)

	// Replace iso-8859-1 charset with utf-8 in headers
	// This is a workaround for go-message library's limited charset support
	charsetRegex := regexp.MustCompile(`(?i)charset=["']?iso-8859-1["']?`)
	normalizedStr := charsetRegex.ReplaceAllString(rawStr, `charset="utf-8"`)

	// Convert iso-8859-1 encoded content to UTF-8
	// Find parts with iso-8859-1 encoding and convert them
	if strings.Contains(strings.ToLower(rawStr), "charset=iso-8859-1") {
		// Convert the entire content from ISO-8859-1 to UTF-8
		decoder := charmap.ISO8859_1.NewDecoder()
		converted, _, err := transform.Bytes(decoder, rawBytes)
		if err == nil {
			// Replace charset declaration in the converted content
			convertedStr := string(converted)
			normalizedStr = charsetRegex.ReplaceAllString(convertedStr, `charset="utf-8"`)
			return []byte(normalizedStr)
		}
	}

	return []byte(normalizedStr)
}

// extractEmailContent extracts body content, HTML content, attachments, inline resources and raw bytes from IMAP message
func extractEmailContent(msg *imap.Message) (bodyContent string, htmlContent string, attachments []EmailAttachment, inlineResources []InlineResource, rawBytes []byte) {
	section := &imap.BodySectionName{}

	if r := msg.GetBody(section); r != nil {
		// Read full RFC822 into memory for later high-fidelity forwarding
		if b, err := io.ReadAll(r); err == nil {
			rawBytes = b

			// Normalize charset to handle go-message library limitations
			normalizedBytes := normalizeEmailCharset(b)

			if mr, err := gomail.CreateReader(bytes.NewReader(normalizedBytes)); err == nil {
				for {
					part, err := mr.NextPart()
					if err == io.EOF {
						break
					}
					if err != nil {
						break
					}

					switch h := part.Header.(type) {
					case *gomail.InlineHeader:
						contentType, params, _ := h.ContentType()
						disposition, dispParams, _ := h.ContentDisposition()

						// Read part content
						partContent, readErr := io.ReadAll(part.Body)
						if readErr != nil {
							continue
						}

						// Handle different content types and dispositions
						// Parse the base content type (ignore charset and other parameters)
						baseContentType := strings.ToLower(strings.Split(contentType, ";")[0])
						baseContentType = strings.TrimSpace(baseContentType)

						switch {
						case baseContentType == "text/plain" && bodyContent == "":
							bodyContent = string(partContent)
						case baseContentType == "text/html" && htmlContent == "":
							htmlContent = string(partContent)
						case disposition == "attachment" || (disposition == "" && shouldTreatAsAttachment(contentType)):
							// Handle attachment
							filename := getFilename(dispParams, params, contentType)
							attachments = append(attachments, EmailAttachment{
								Filename:    filename,
								ContentType: contentType,
								Content:     partContent,
							})
						case disposition == "inline" && strings.HasPrefix(contentType, "image/"):
							// Handle inline resource (e.g., embedded image)
							contentID := strings.Trim(h.Get("Content-Id"), "<>")
							if contentID != "" {
								inlineResources = append(inlineResources, InlineResource{
									ContentID:   contentID,
									ContentType: contentType,
									Content:     partContent,
								})
							}
						}
					case *gomail.AttachmentHeader:
						// Handle attachment headers
						contentType, params, _ := h.ContentType()
						filename := getFilename(nil, params, contentType)

						partContent, readErr := io.ReadAll(part.Body)
						if readErr != nil {
							continue
						}

						attachments = append(attachments, EmailAttachment{
							Filename:    filename,
							ContentType: contentType,
							Content:     partContent,
						})
					}
				}
			}
		}
	}

	return bodyContent, htmlContent, attachments, inlineResources, rawBytes
}

// shouldTreatAsAttachment determines if a content type should be treated as an attachment
func shouldTreatAsAttachment(contentType string) bool {
	// Common attachment content types
	attachmentTypes := []string{
		"application/", "image/", "audio/", "video/",
	}

	for _, prefix := range attachmentTypes {
		if strings.HasPrefix(contentType, prefix) {
			// Exclude inline images which should be handled separately
			if contentType == "image/gif" || contentType == "image/jpeg" ||
				contentType == "image/png" || contentType == "image/webp" {
				return false
			}
			return true
		}
	}
	return false
}

// getFilename extracts filename from Content-Disposition or Content-Type parameters
func getFilename(dispParams map[string]string, typeParams map[string]string, contentType string) string {
	decoder := mime.WordDecoder{}

	// Try Content-Disposition filename first
	if dispParams != nil {
		if filename := dispParams["filename"]; filename != "" {
			if decoded, err := decoder.DecodeHeader(filename); err == nil {
				return decoded
			}
			return filename
		}
	}

	// Try Content-Type name parameter
	if typeParams != nil {
		if name := typeParams["name"]; name != "" {
			if decoded, err := decoder.DecodeHeader(name); err == nil {
				return decoded
			}
			return name
		}
	}

	// Generate filename based on content type
	if ext := getExtensionFromContentType(contentType); ext != "" {
		return "attachment" + ext
	}

	return "attachment"
}

// getExtensionFromContentType returns file extension for common content types
func getExtensionFromContentType(contentType string) string {
	extensions := map[string]string{
		"text/plain":         ".txt",
		"text/html":          ".html",
		"application/pdf":    ".pdf",
		"application/msword": ".doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
		"application/vnd.ms-excel": ".xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ".xlsx",
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/gif":  ".gif",
		"image/webp": ".webp",
	}

	if ext, exists := extensions[contentType]; exists {
		return ext
	}

	// Try to guess from MIME type
	if parts := strings.SplitN(contentType, "/", 2); len(parts) == 2 {
		return "." + parts[1]
	}

	return ""
}
