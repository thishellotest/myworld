package biz

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib/gomail"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	ForwardTitle       = "Please execute once you get a chance. Thanks!"
	ForwardTitleRepeat = "Please execute ASAP, thanks!"
)

// MonitoredEmailsJobUsecase is responsible for the cron-driven workflow that
// processes monitored mailboxes (listen and auto-forward).
type MonitoredEmailsJobUsecase struct {
	log                   *log.Helper
	conf                  *conf.Data
	MonitoredEmailsStore  *MonitoredEmailsUsecase
	AttorneyUsecase       *AttorneyUsecase
	TasksUsecase          *MonitoredEmailsTasksUsecase
	ClientEnvelopeUsecase *ClientEnvelopeUsecase
}

func NewMonitoredEmailsJobUsecase(
	logger log.Logger,
	conf *conf.Data,
	MonitoredEmailsStore *MonitoredEmailsUsecase,
	AttorneyUsecase *AttorneyUsecase,
	TasksUsecase *MonitoredEmailsTasksUsecase,
	ClientEnvelopeUsecase *ClientEnvelopeUsecase,
) *MonitoredEmailsJobUsecase {
	return &MonitoredEmailsJobUsecase{
		log:                   log.NewHelper(logger),
		conf:                  conf,
		MonitoredEmailsStore:  MonitoredEmailsStore,
		AttorneyUsecase:       AttorneyUsecase,
		TasksUsecase:          TasksUsecase,
		ClientEnvelopeUsecase: ClientEnvelopeUsecase,
	}
}

// Handle is invoked by a cron schedule.
// Process monitored inboxes: compensation tasks first, then incremental email fetching.
func (c *MonitoredEmailsJobUsecase) Handle() error {
	c.log.Infof("MonitoredEmailsJobUsecase.Handle tick")

	// Load active monitored inboxes from store
	activeInboxes, err := c.MonitoredEmailsStore.ListActive()
	if err != nil {
		return fmt.Errorf("ListActive monitored inboxes error: %w", err)
	}

	c.log.Infof("Found %d active monitored inboxes", len(activeInboxes))

	// Process each inbox
	for _, inbox := range activeInboxes {
		if inbox == nil || strings.TrimSpace(inbox.EmailAddress) == "" {
			continue
		}

		c.log.Infof("Processing inbox: %s", inbox.EmailAddress)

		// 1. process forwarded but NOT signed emails
		if err := c.processReForwardTasks(); err != nil {
			c.log.Errorf("process re-forward tasks failed: %v", err)
		}

		// 2. Compensation processing for unfinished tasks
		if err := c.processCompensationTasks(inbox); err != nil {
			c.log.Errorf("compensation failed for inbox=%s: %v", inbox.EmailAddress, err)
		}

		// 3. Incremental fetch new emails
		if err := c.processNewEmails(inbox); err != nil {
			c.log.Errorf("process new emails failed for inbox=%s: %v", inbox.EmailAddress, err)
		}

		// 4. Update last scanned time
		if err := c.MonitoredEmailsStore.UpdateLastScannedAt(inbox.ID); err != nil {
			c.log.Errorf("update last scanned time failed for inbox=%s: %v", inbox.EmailAddress, err)
		}
	}

	return nil
}

// processCompensationTasks handles pending tasks (received or error status) for an inbox
func (c *MonitoredEmailsJobUsecase) processCompensationTasks(inbox *MonitoredEmailEntity) error {
	const maxTasks = 10

	pendingTasks, err := c.TasksUsecase.GetPendingTasks(inbox.EmailAddress, maxTasks)
	if err != nil {
		return fmt.Errorf("get pending tasks: %w", err)
	}

	if len(pendingTasks) == 0 {
		return nil
	}

	c.log.Infof("Found %d pending tasks for inbox=%s", len(pendingTasks), inbox.EmailAddress)

	for _, task := range pendingTasks {
		// Mark task as in progress
		if err := c.TasksUsecase.SetInProgress(task.ID); err != nil {
			c.log.Errorf("failed to set task %d in progress: %v", task.ID, err)
			continue
		}

		// Fetch the email by UID
		email, err := FetchEmailByUID(inbox.EmailAddress, inbox.EmailPassword, uint32(task.UID))
		if err != nil {
			c.log.Errorf("failed to fetch email UID=%d: %v", task.UID, err)
			c.TasksUsecase.UpdateStatus(task.ID, "error", fmt.Sprintf("fetch failed: %v", err))
			continue
		}

		if email == nil {
			c.log.Warnf("email UID=%d not found, marking as skipped", task.UID)
			c.TasksUsecase.UpdateStatus(task.ID, "skipped_no_match", "email not found")
			continue
		}

		// Process the email
		if err := c.processEmail(*email, inbox, task.ID); err != nil {
			c.log.Errorf("failed to process email UID=%d: %v", task.UID, err)
			c.TasksUsecase.UpdateStatus(task.ID, "error", fmt.Sprintf("process failed: %v", err))
			continue
		}
	}

	return nil
}

// processNewEmails handles incremental email fetching for an inbox
func (c *MonitoredEmailsJobUsecase) processNewEmails(inbox *MonitoredEmailEntity) error {
	const emailBatchSize = 20

	// Initialize last seen UID if it's 0 (first time)
	if inbox.LastSeenUID == 0 {
		c.log.Infof("Initializing inbox=%s with latest email", inbox.EmailAddress)
		return c.initializeInbox(inbox)
	}

	// Fetch emails after last seen UID
	emails, hasMore, err := FetchEmailsByUIDRange(
		inbox.EmailAddress,
		inbox.EmailPassword,
		uint32(inbox.LastSeenUID),
		emailBatchSize,
	)

	if err != nil {
		return fmt.Errorf("fetch emails by UID range: %w", err)
	}

	if len(emails) == 0 {
		c.log.Infof("No new emails for inbox=%s", inbox.EmailAddress)
		return nil
	}

	c.log.Infof("Found %d new emails for inbox=%s, hasMore=%v, lastSeenUID=%d", len(emails), inbox.EmailAddress, hasMore, inbox.LastSeenUID)

	var lastProcessedUID uint64
	for _, email := range emails {
		// Register email for idempotent processing
		isNew, err := c.TasksUsecase.RegisterTask(email, inbox.EmailAddress)
		if err != nil {
			c.log.Errorf("failed to register email UID=%d: %v", email.UID, err)
			continue
		}

		// Update last seen UID regardless of processing result
		if uint64(email.UID) > lastProcessedUID {
			lastProcessedUID = uint64(email.UID)
		}

		if !isNew {
			c.log.Debugf("Email UID=%d already exists, skipping", email.UID)
			continue
		}

		// Find the task ID for status updates
		var taskID uint64
		if task, err := c.TasksUsecase.GetByUID(inbox.EmailAddress, email.UID); err == nil && task != nil {
			taskID = task.ID
		}

		// Process the new email immediately
		if err := c.processEmail(email, inbox, taskID); err != nil {
			c.log.Errorf("failed to process new email UID=%d: %v", email.UID, err)
			continue
		}
	}

	// Update last seen UID
	if lastProcessedUID > 0 {
		if err := c.MonitoredEmailsStore.UpdateLastSeenUID(inbox.ID, lastProcessedUID); err != nil {
			c.log.Errorf("failed to update last seen UID for inbox=%s: %v", inbox.EmailAddress, err)
		} else {
			c.log.Infof("Updated last seen UID to %d for inbox=%s", lastProcessedUID, inbox.EmailAddress)
		}
	}

	return nil
}

// initializeInbox initializes an inbox by fetching the latest email and setting LastSeenUID
func (c *MonitoredEmailsJobUsecase) initializeInbox(inbox *MonitoredEmailEntity) error {
	emails, err := FetchRecentEmails(inbox.EmailAddress, inbox.EmailPassword, 1)
	if err != nil {
		return fmt.Errorf("fetch recent emails for initialization: %w", err)
	}

	if len(emails) == 0 {
		c.log.Infof("No emails found for initialization of inbox=%s", inbox.EmailAddress)
		return nil
	}

	latestEmail := emails[0]
	c.log.Infof("Initializing inbox=%s with UID=%d", inbox.EmailAddress, latestEmail.UID)

	// Register and process the latest email
	isNew, err := c.TasksUsecase.RegisterTask(latestEmail, inbox.EmailAddress)
	if err != nil {
		return fmt.Errorf("register initial email: %w", err)
	}

	if isNew {
		// Find the task ID for status updates
		var taskID uint64
		if task, err := c.TasksUsecase.GetByUID(inbox.EmailAddress, latestEmail.UID); err == nil && task != nil {
			taskID = task.ID
		}

		if err := c.processEmail(latestEmail, inbox, taskID); err != nil {
			c.log.Errorf("failed to process initial email UID=%d: %v", latestEmail.UID, err)
		}
	}

	// Set LastSeenUID
	if err := c.MonitoredEmailsStore.UpdateLastSeenUID(inbox.ID, uint64(latestEmail.UID)); err != nil {
		return fmt.Errorf("update last seen UID: %w", err)
	}

	return nil
}

// processEmail processes a single email: extract name, find attorney, and forward
// taskID is optional (0 for new emails, >0 for compensation tasks)
func (c *MonitoredEmailsJobUsecase) processEmail(email IncomingEmail, inbox *MonitoredEmailEntity, taskID uint64) error {
	// Check if subject contains required text
	if !strings.Contains(strings.ToLower(email.Subject), strings.ToLower("Your VA Representation Agreement with August Miles")) {
		c.log.Infof("Subject does not contain required text: %q UID=%d", email.Subject, email.UID)
		if taskID > 0 {
			c.TasksUsecase.UpdateStatus(taskID, "skipped_no_match", "")
		}
		return nil
	}

	// Extract target full name from subject by taking text after the last " - "
	fullName := FindNameInSubject(email.Subject)
	if fullName == "" {
		c.log.Infof("No matched name found for subject=%q UID=%d", email.Subject, email.UID)
		if taskID > 0 {
			c.TasksUsecase.UpdateStatus(taskID, "skipped_no_match", "")
		}
		return nil
	}

	// Find attorney by name using AttorneyUsecase (GetByName handles splitting internally)
	target, err := c.AttorneyUsecase.GetByName(fullName)
	if err != nil {
		errMsg := fmt.Sprintf("GetByName error for name=%q: %v", fullName, err)
		c.log.Errorf(errMsg)
		if taskID > 0 {
			c.TasksUsecase.UpdateStatus(taskID, "error", errMsg)
		}
		return fmt.Errorf(errMsg)
	}

	if target == nil || strings.TrimSpace(target.ForwardEmails) == "" {
		c.log.Infof("No matching attorney or missing ForwardEmails for subject=%q UID=%d", email.Subject, email.UID)
		if taskID > 0 {
			c.TasksUsecase.UpdateStatus(taskID, "skipped_no_match", "")
		}
		return nil
	}

	// Forward the email using the common forwarding logic
	forwardErr := c.forwardEmailToAttorney(email, inbox, ForwardTitle)

	// Update task status based on forwarding result
	if forwardErr != nil {
		errMsg := fmt.Sprintf("Forward failed: %v", forwardErr)
		c.log.Errorf(errMsg)
		if taskID > 0 {
			c.TasksUsecase.UpdateStatus(taskID, "error", errMsg)
		}
		return fmt.Errorf(errMsg)
	}

	if taskID > 0 {
		c.TasksUsecase.UpdateStatus(taskID, "forwarded", "")

		// Extract envelopId from email and save to data field
		envelopId := extractEnvelopId(email.HTMLBody)
		if envelopId != "" {
			data := map[string]string{
				"envelopId": envelopId,
			}
			dataJSON, err := json.Marshal(data)
			if err != nil {
				c.log.Errorf("Failed to marshal envelopId data for task %d: %v", taskID, err)
			} else {
				if err := c.TasksUsecase.UpdateData(taskID, string(dataJSON)); err != nil {
					c.log.Errorf("Failed to update envelopId data for task %d: %v", taskID, err)
				} else {
					c.log.Infof("Saved envelopId %s for task %d", envelopId, taskID)
				}
			}
		} else {
			c.log.Infof("No envelopId found in email UID=%d", email.UID)
		}
	}

	return nil
}

// NormalizeRecipientList converts a comma-separated email list into the semicolon-separated
// format expected by the mail sender. It trims spaces and removes empty entries.
func NormalizeRecipientList(input string) string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return ""
	}
	// Already using semicolons; normalize consecutive separators and spaces
	if strings.Contains(trimmed, ";") && !strings.Contains(trimmed, ",") {
		parts := strings.Split(trimmed, ";")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				out = append(out, p)
			}
		}
		return strings.Join(out, ";")
	}
	// Replace commas with semicolons, and trim each component
	parts := strings.Split(strings.ReplaceAll(trimmed, ",", ";"), ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, ";")
}

// FindNameInSubject extracts the attorney name from a subject line by taking the
// substring after the last occurrence of " - ". Returns empty string if not present.
func FindNameInSubject(subject string) string {
	const separator = " - "
	idx := strings.LastIndex(subject, separator)
	if idx < 0 {
		return ""
	}
	return strings.TrimSpace(subject[idx+len(separator):])
}

// extractSenderDisplayName returns the display name part from a header like
// "John Doe <john@example.com>". If no name is present, it tries to return
// the local part of the email address; otherwise returns an empty string.
func extractSenderDisplayName(fromHeader string) string {
	fromHeader = strings.TrimSpace(fromHeader)
	// Case: "Name <email@host>"
	if idx := strings.Index(fromHeader, "<"); idx >= 0 {
		name := strings.TrimSpace(fromHeader[:idx])
		name = strings.Trim(name, `"`)
		if name != "" {
			return name
		}
	}
	// Fallback: try local part before '@'
	if at := strings.Index(fromHeader, "@"); at > 0 {
		local := fromHeader[:at]
		local = strings.Trim(local, `"`)
		if local != "" {
			return local
		}
	}
	return ""
}

// getForwardCcEmails returns static CC list for auto-forwarded emails.
func (c *MonitoredEmailsJobUsecase) GetForwardCcEmails() []string {
	return c.conf.EmailMonitor.Cc
}

// forwardAsHTML forwards the message with HTML content in the body, preserving formatting, images, and attachments.
func (c *MonitoredEmailsJobUsecase) forwardAsHTML(
	mailer *MailUsecase,
	cfg *MailServiceConfig,
	msg IncomingEmail,
	target *AttorneyEntity,
	forwardTo string,
	forwardSubject string,
	title string,
) error {
	// Build HTML body with original email content
	htmlBody := c.buildHTMLForwardBody(msg, target, title)

	// Prepare attachments (original attachments + inline resources as attachments if needed)
	var attachmentInputs MailAttachmentInputs

	// Add original attachments
	for _, att := range msg.Attachments {
		headers := map[string][]string{
			"Content-Type":        {att.ContentType},
			"Content-Disposition": {fmt.Sprintf("attachment; filename=\"%s\"", att.Filename)},
		}
		attachmentInputs = append(attachmentInputs, MailAttachmentInput{
			Name:     att.Filename,
			Reader:   bytes.NewReader(att.Content),
			Settings: []gomail.FileSetting{gomail.SetHeader(headers)},
		})
	}

	// Add inline resources as attachments if not embedded in HTML
	for _, res := range msg.InlineResources {
		filename := fmt.Sprintf("inline_%s%s", res.ContentID, getExtensionFromContentType(res.ContentType))
		headers := map[string][]string{
			"Content-Type":        {res.ContentType},
			"Content-Disposition": {fmt.Sprintf("attachment; filename=\"%s\"", filename)},
		}
		attachmentInputs = append(attachmentInputs, MailAttachmentInput{
			Name:     filename,
			Reader:   bytes.NewReader(res.Content),
			Settings: []gomail.FileSetting{gomail.SetHeader(headers)},
		})
	}

	outMsg := &MailMessage{
		To:       forwardTo,
		MailType: "html", // Set to HTML to preserve formatting
		Subject:  "Fwd: " + forwardSubject,
		Body:     htmlBody,
	}
	if cc := c.GetForwardCcEmails(); len(cc) > 0 {
		outMsg.Cc = cc
	}

	c.log.Infof("Forwarding to=%s subject=%q", forwardTo, forwardSubject)

	// Send with attachments - use MailAttach_No and pass attachments via mailAttachmentInputs
	if err := mailer.SendEmail(cfg, outMsg, MailAttach_No, attachmentInputs); err != nil {
		return err
	}
	c.log.Infof("Forwarded to=%s subject=%q", forwardTo, forwardSubject)
	return nil
}

// buildHTMLForwardBody constructs an HTML forward body with original context and formatting.
func (c *MonitoredEmailsJobUsecase) buildHTMLForwardBody(msg IncomingEmail, target *AttorneyEntity, title string) string {
	var htmlBody strings.Builder

	// White background container start
	htmlBody.WriteString("<div style='padding: 10px 20px 30px 20px; background-color: #ffffff; color: #000000;'>")

	// Add title if provided
	if title != "" {
		htmlBody.WriteString(fmt.Sprintf("<div style='margin: 10px 0; color: #000000; font-size:14px;'>%s</div>", title))
	}

	// Build forwarding header
	htmlBody.WriteString("<div style='margin: 0 0; color: #000000; font-size: 14px'>")
	htmlBody.WriteString("<div>---------- Forwarded message ----------</div>")
	htmlBody.WriteString(fmt.Sprintf("<div><strong>From:</strong> %s</div>", msg.From))
	htmlBody.WriteString(fmt.Sprintf("<div><strong>To:</strong> %s</div>", msg.To))
	htmlBody.WriteString(fmt.Sprintf("<div><strong>Date:</strong> %s</div>", msg.Date.In(configs.GetVBCDefaultLocation()).Format("2006-01-02 15:04:05 MST")))
	htmlBody.WriteString("</div>")

	// White background container end
	htmlBody.WriteString("</div>")

	// Add separator
	// htmlBody.WriteString("<hr style='margin: 20px 0; border: none; border-top: 1px solid #ddd;'>")

	// Add original email content
	htmlBody.WriteString("<div style='margin: 20px 0;'>")

	// Prefer HTML content if available, otherwise use plain text
	if strings.TrimSpace(msg.HTMLBody) != "" {
		// Process HTML content to embed inline images as base64
		processedHTML := c.processHTMLContent(msg.HTMLBody, msg.InlineResources)
		htmlBody.WriteString(processedHTML)
	} else if strings.TrimSpace(msg.Body) != "" {
		// Convert plain text to HTML with proper formatting
		htmlBody.WriteString("<pre style='white-space: pre-wrap; font-family: Arial, sans-serif;'>")
		htmlBody.WriteString(strings.ReplaceAll(strings.ReplaceAll(msg.Body, "&", "&amp;"), "<", "&lt;"))
		htmlBody.WriteString("</pre>")
	} else {
		htmlBody.WriteString("<p><em>No content available</em></p>")
	}

	htmlBody.WriteString("</div>")

	return htmlBody.String()
}

// processHTMLContent processes HTML content to embed inline images as base64 data URIs
func (c *MonitoredEmailsJobUsecase) processHTMLContent(htmlContent string, inlineResources []InlineResource) string {
	processedHTML := htmlContent

	// Create a map for quick lookup of inline resources by Content-ID
	resourceMap := make(map[string]InlineResource)
	for _, res := range inlineResources {
		resourceMap[res.ContentID] = res
	}

	// Pattern to match cid: references in HTML
	cidPattern := regexp.MustCompile(`cid:([^"'\s>]+)`)

	// Replace cid: references with base64 data URIs
	processedHTML = cidPattern.ReplaceAllStringFunc(processedHTML, func(match string) string {
		// Extract the Content-ID from cid:content-id
		contentID := strings.TrimPrefix(match, "cid:")

		if resource, exists := resourceMap[contentID]; exists {
			// Convert to base64 data URI
			base64Data := base64.StdEncoding.EncodeToString(resource.Content)
			return fmt.Sprintf("data:%s;base64,%s", resource.ContentType, base64Data)
		}

		// If resource not found, keep original reference
		return match
	})

	return processedHTML
}

// extractEnvelopId extracts the envelopId from email HTML content
// It looks for links matching /sign/document/{envelopId} pattern
// e.g. <a href="https://example.com/sign/document/abc123" xxx>xxx</a> -> returns "abc123"
// e.g. <a href="/sign/document/xyz789?param=value">xx</a> -> returns "xyz789"
func extractEnvelopId(htmlContent string) string {
	// Pattern to match /sign/document/{envelopId} in href attributes
	pattern := regexp.MustCompile(`href=["'][^"']*/sign/document/([^/"'\s?]+)[^"']*["']`)
	matches := pattern.FindStringSubmatch(htmlContent)

	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractEnvelopIdFromTaskData extracts the envelopId from task's data field
func extractEnvelopIdFromTaskData(data string) string {
	if data == "" {
		return ""
	}

	var dataMap map[string]string
	if err := json.Unmarshal([]byte(data), &dataMap); err != nil {
		return ""
	}

	return dataMap["envelopId"]
}

// processReForwardTasks handles re-forwarding logic for tasks older than configured seconds
func (c *MonitoredEmailsJobUsecase) processReForwardTasks() error {
	// Get all forwarded tasks older than configured seconds
	tasks, err := c.TasksUsecase.GetForwardedTasksOlderThanConfigSeconds()
	if err != nil {
		return fmt.Errorf("get forwarded tasks older than configured seconds: %w", err)
	}

	if len(tasks) == 0 {
		c.log.Infof("No tasks found for re-forwarding")
		return nil
	}

	c.log.Infof("Found %d tasks for re-forwarding", len(tasks))

	for _, task := range tasks {
		if err := c.processReForwardTask(task); err != nil {
			c.log.Errorf("failed to re-forward task %d: %v", task.ID, err)
			continue
		}
	}

	return nil
}

// processReForwardTask processes a single task for re-forwarding
func (c *MonitoredEmailsJobUsecase) processReForwardTask(task *MonitoredEmailsTasksEntity) error {
	// Extract envelopId from task data
	envelopId := extractEnvelopIdFromTaskData(task.Data)
	if envelopId == "" {
		c.log.Warnf("No envelopId found in task %d data, marking as signed", task.ID)
		c.TasksUsecase.UpdateStatus(task.ID, "signed", "No envelopId found")
		return nil
	}

	// Query client_envelops table to check if signed
	clientEnvelope, err := c.ClientEnvelopeUsecase.GetByEnvelopeId("box", envelopId)
	if err != nil {
		return fmt.Errorf("get client envelope by envelopId %s: %w", envelopId, err)
	}

	if clientEnvelope == nil {
		c.log.Warnf("No client envelope found for envelopId %s, task %d, treating as signed", envelopId, task.ID)
		return c.TasksUsecase.UpdateStatus(task.ID, "signed", "No client envelope found")
	}

	// Only forward when is_signed == 0 (not signed)
	if clientEnvelope.IsSigned != ClientEnvelope_IsSigned_No {
		var msg string
		if clientEnvelope.IsSigned != ClientEnvelope_IsSigned_Yes {
			msg = fmt.Sprintf("is_signed=%d", clientEnvelope.IsSigned)
		}
		c.log.Infof("Envelope %s is signed (is_signed=%d), marking task %d as signed", envelopId, clientEnvelope.IsSigned, task.ID)
		return c.TasksUsecase.UpdateStatus(task.ID, "signed", msg)
	}

	// Only forward when is_signed == 0
	c.log.Infof("Envelope %s is not signed (is_signed=0), re-forwarding task %d", envelopId, task.ID)

	// Get the monitored email inbox
	inbox, err := c.MonitoredEmailsStore.GetByEmailAddress(task.InboxEmail)
	if err != nil {
		return fmt.Errorf("get inbox by email %s: %w", task.InboxEmail, err)
	}

	if inbox == nil {
		return fmt.Errorf("inbox not found for email %s", task.InboxEmail)
	}

	// Fetch the email by UID
	email, err := FetchEmailByUID(inbox.EmailAddress, inbox.EmailPassword, uint32(task.UID))
	if err != nil {
		return fmt.Errorf("fetch email UID=%d: %w", task.UID, err)
	}

	if email == nil {
		c.log.Warnf("Email UID=%d not found for task %d, marking as signed", task.UID, task.ID)
		c.TasksUsecase.UpdateStatus(task.ID, "signed", "Email not found")
		return nil
	}

	// Re-forward the email (reusing existing logic but without parsing envelopId again)
	return c.reForwardEmail(*email, inbox, task.ID)
}

// forwardEmailToAttorney handles the core email forwarding logic used by both initial and re-forwarding
func (c *MonitoredEmailsJobUsecase) forwardEmailToAttorney(email IncomingEmail, inbox *MonitoredEmailEntity, title string) error {
	// Extract target full name from subject by taking text after the last " - "
	fullName := FindNameInSubject(email.Subject)
	if fullName == "" {
		return fmt.Errorf("no matched name found for subject=%q UID=%d", email.Subject, email.UID)
	}

	// Find attorney by name using AttorneyUsecase
	target, err := c.AttorneyUsecase.GetByName(fullName)
	if err != nil {
		return fmt.Errorf("GetByName error for name=%q: %v", fullName, err)
	}

	if target == nil || strings.TrimSpace(target.ForwardEmails) == "" {
		return fmt.Errorf("no matching attorney or missing ForwardEmails for subject=%q UID=%d", email.Subject, email.UID)
	}

	// Forward logic: use receiver inbox account to forward to attorney.ForwardEmails.
	forwardTo := NormalizeRecipientList(target.ForwardEmails)
	forwardSubject := email.Subject // keep original subject

	// Build mail service config from the current inbox and original sender name
	dynamicMailServiceConfig := &MailServiceConfig{
		Name:        extractSenderDisplayName(email.From), // original sender name
		Host:        "smtp.gmail.com",
		Port:        587,
		Username:    inbox.EmailAddress,
		Password:    inbox.EmailPassword,
		FromAddress: inbox.EmailAddress, // send from the monitored inbox
	}

	mailer := &MailUsecase{}

	// Attempt forwarding
	forwardErr := c.forwardAsHTML(mailer, dynamicMailServiceConfig, email, target, forwardTo, forwardSubject, title)

	// Update task status based on forwarding result
	if forwardErr != nil {
		return fmt.Errorf("forward failed to=%s subject=%q: %v", forwardTo, forwardSubject, forwardErr)
	}

	// Success
	c.log.Infof("Successfully forwarded email UID=%d to=%s subject=%q", email.UID, forwardTo, forwardSubject)
	return nil
}

// reForwardEmail forwards an email again without parsing envelopId
func (c *MonitoredEmailsJobUsecase) reForwardEmail(email IncomingEmail, inbox *MonitoredEmailEntity, taskID uint64) error {
	// Forward the email using the common forwarding logic
	if err := c.forwardEmailToAttorney(email, inbox, ForwardTitleRepeat); err != nil {
		return err
	}

	// Success - update status to forwarded (which will reset the timer)
	return c.TasksUsecase.UpdateStatus(taskID, "forwarded", "")
}
