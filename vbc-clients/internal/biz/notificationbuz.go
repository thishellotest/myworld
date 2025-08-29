package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"regexp"
	"strings"
	"sync"
	"vbc/internal/conf"
	"vbc/lib"
)

type NotificationbuzUsecase struct {
	log                 *log.Helper
	conf                *conf.Data
	CommonUsecase       *CommonUsecase
	NotificationUsecase *NotificationUsecase
	UserUsecase         *UserUsecase
	NotesUsecase        *NotesUsecase
}

func NewNotificationbuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	NotificationUsecase *NotificationUsecase,
	UserUsecase *UserUsecase,
	NotesUsecase *NotesUsecase,
) *NotificationbuzUsecase {
	uc := &NotificationbuzUsecase{
		log:                 log.NewHelper(logger),
		CommonUsecase:       CommonUsecase,
		conf:                conf,
		NotificationUsecase: NotificationUsecase,
		UserUsecase:         UserUsecase,
		NotesUsecase:        NotesUsecase,
	}

	return uc
}

func (c *NotificationbuzUsecase) Tidy(records []NotificationEntity) (NotificationList, error) {

	var notificationList NotificationList

	var userGids []string
	var notesGids []string
	var once sync.Once
	for _, v := range records {
		if v.FromType == Notification_FromType_Notes {
			notesGids = append(notesGids, v.FromGid)
		} else if v.FromType == Notification_FromType_PW {

		}
		once.Do(func() {
			userGids = append(userGids, v.ReceiverGid)
		})

		userGids = append(userGids, v.SenderGid)
	}

	userGids = lib.RemoveDuplicates(userGids)
	userFacadesMap, err := c.UserUsecase.GetUserFacadesByGids(userGids)
	//lib.DPrintln(userFacadesMap)
	if err != nil {
		return nil, err
	}

	notesGids = lib.RemoveDuplicates(notesGids)
	noteFacadesMap, err := c.NotesUsecase.GetNoteFacadesByGids(notesGids)
	if err != nil {
		return nil, err
	}

	for _, v := range records {
		var item NotificationItem
		item.Id = v.ID
		item.Title = "Notes"
		item.NotificationTime = int32(v.CreatedAt)
		item.Gid = v.Gid
		item.Content = v.Content
		if user, ok := userFacadesMap[v.SenderGid]; ok {
			u := user.ToFabUser()
			item.TriggerUser = u
		}
		item.Unread = v.Unread

		if v.FromType == Notification_FromType_Notes {
			if note, ok := noteFacadesMap[v.FromGid]; ok {
				kind := note.CustomFields.TextValueByNameBasic(Notes_FieldName_kind)
				kindGid := note.CustomFields.TextValueByNameBasic(Notes_FieldName_kind_gid)
				url := fmt.Sprintf("/tab/%s/%s", KindConvertToModule(kind), kindGid)
				item.Url = url
			}
		} else if v.FromType == Notification_FromType_PW {
			url := fmt.Sprintf("/ps/%s", v.FromGid)
			item.Url = url
			item.OpenNewWindow = true
		}

		notificationList = append(notificationList, item)
	}

	return notificationList, nil
}

// NotificationTextExtractUTF8SubstringPrefix 提取 UTF-8 字符串的子串，确保不会截断字符
func NotificationTextExtractUTF8SubstringPrefix(s string, prefixLen int) (r string, hasMore bool) {
	runes := []rune(s)
	end := len(runes)
	start := max(end-prefixLen, 0)

	if start >= len(runes) {
		return "", false
	}
	if end > len(runes) {
		end = len(runes)
	}
	if start != 0 {
		hasMore = true
	}
	return string(runes[start:end]), hasMore
}

// NotificationTextExtractUTF8SubstringSuffix 提取 UTF-8 字符串的子串，确保不会截断字符
func NotificationTextExtractUTF8SubstringSuffix(s string, suffixLen int) (r string, hasMore bool) {
	runes := []rune(s)
	end := min(len(runes), suffixLen)
	start := 0
	if start >= len(runes) {
		return "", false
	}
	if end > len(runes) {
		end = len(runes)
	}
	if end != len(runes) {
		hasMore = true
	}
	return string(runes[start:end]), hasMore
}

func NotificationTextRemoveMentions(input string) string {
	// 定义匹配模式，匹配类似 @[名称](数字) 的格式
	re := regexp.MustCompile(`@\[[^\]]+\]\(.+?\)`)
	// 使用空字符串替换匹配到的内容
	result := re.ReplaceAllString(input, "")
	return result
}

type NotificationTextExtractContextList []NotificationTextExtractContextItem
type NotificationTextExtractContextItem struct {
	MatchResult   string
	UserGid       string
	PrefixText    string
	PrefixHasMore bool
	SuffixText    string
	SuffixHasMore bool
}

func (c *NotificationTextExtractContextItem) ToText() (r string) {
	if c.PrefixText != "" {
		if c.PrefixHasMore {
			r += "..."
		}
		r += c.PrefixText
	}
	r += c.MatchResult
	if c.SuffixText != "" {
		r += c.SuffixText
		if c.SuffixHasMore {
			r += "..."
		}
	}
	return
}

func NotificationTextExtractContext(input string, prefixLen, suffixLen int) (r NotificationTextExtractContextList) {
	// 定义正则表达式，匹配 @[xxx](yyyyyyyyyyyyyyyyy) 格式
	re := regexp.MustCompile(`@\[[^\]]+\]\((.+?)\)`)
	matches := re.FindAllStringSubmatch(input, -1)
	//lib.DPrintln(matches)
	existIds := make(map[string]bool)
	for _, match := range matches {
		if _, ok := existIds[match[1]]; ok {
			continue
		}
		existIds[match[1]] = true
		start := strings.Index(input, match[0])
		prefixInput := input[:start]
		prefixInput = NotificationTextRemoveMentions(prefixInput)
		suffixInput := input[start+len(match[0]):]
		suffixInput = NotificationTextRemoveMentions(suffixInput)
		// 提取前后字符
		prefix, prefixHasMore := NotificationTextExtractUTF8SubstringPrefix(prefixInput, prefixLen)
		suffix, suffixHasMore := NotificationTextExtractUTF8SubstringSuffix(suffixInput, suffixLen)

		item := NotificationTextExtractContextItem{
			MatchResult:   match[0],
			UserGid:       match[1],
			PrefixText:    prefix,
			PrefixHasMore: prefixHasMore,
			SuffixText:    suffix,
			SuffixHasMore: suffixHasMore,
		}
		r = append(r, item)

		//a := fmt.Sprintf("前: %q prefixHasMore:%v 后: %q suffixHasMore:%v", prefix, prefixHasMore, suffix, suffixHasMore)
		//lib.DPrintln(a)
		//result = append(result, fmt.Sprintf("前: %q 后: %q", prefix, suffix))
	}
	return r
}
