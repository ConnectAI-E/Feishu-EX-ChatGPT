package domain

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/pkg/consts"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

var reTrimMentionText = regexp.MustCompile(`@[^ ]*`)

type LarkMessage larkim.EventMessage

func (m LarkMessage) ID() string {
	return *m.MessageId
}

func (m LarkMessage) IsMentionAt(name string) bool {

	for _, mention := range m.Mentions {

		if *mention.Name == name {
			return true
		}
	}

	return false
}

func (m LarkMessage) GetChatType() consts.ChatType {
	return consts.ChatType(*m.ChatType)
}

func (m LarkMessage) GetText() string {

	content := m.parseContent()

	return strings.TrimSpace(content)
}

func (m LarkMessage) parseContent() string {
	content := *m.Content

	//"{\"text\":\"@_user_1  hahaha\"}",
	//only get text content hahaha
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
	}
	if contentMap["text"] == nil {
		return ""
	}
	text := contentMap["text"].(string)
	return m.trimMentionText(text)
}

func (m LarkMessage) trimMentionText(msg string) string {
	//replace @到下一个非空的字段 为 ''
	return reTrimMentionText.ReplaceAllString(msg, "")
}
