package common

import (
	"fmt"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/model"
)

func HandleMessages(messages []model.ChatCompletionMessage, isSupportSystemRole bool) []model.ChatCompletionMessage {

	var (
		newMessages       = make([]model.ChatCompletionMessage, 0)
		systemRoleMessage *model.ChatCompletionMessage
	)

	for _, message := range messages {
		if message.Content != "" {
			newMessages = append(newMessages, message)
		}
	}

	if isSupportSystemRole && newMessages[0].Role == consts.ROLE_SYSTEM {
		systemRoleMessage = &newMessages[0]
		newMessages = newMessages[1:]
	}

	if len(newMessages) != 0 && len(newMessages)%2 == 0 {
		newMessages = newMessages[1:]
	}

	for i := len(newMessages) - 1; i >= 0; i-- {
		if i%2 == 0 {
			newMessages[i].Role = consts.ROLE_USER
		} else {
			newMessages[i].Role = consts.ROLE_ASSISTANT
		}
	}

	if systemRoleMessage != nil {
		newMessages = append([]model.ChatCompletionMessage{*systemRoleMessage}, newMessages...)
	}

	return newMessages
}

func GetMime(url string) (mimeType string, data string) {

	base64 := gstr.Split(url, "base64,")
	if len(base64) > 1 {
		data = base64[1]
		if gstr.HasPrefix(url, "data:image/") {
			// data:image/jpeg;base64,
			mimeType = fmt.Sprintf("image/%s", gstr.Split(base64[0][11:], ";")[0])
		} else if gstr.HasPrefix(url, "data:text/") {
			// data:text/plain;base64,
			mimeType = fmt.Sprintf("text/%s", gstr.Split(base64[0][10:], ";")[0])
		}
	} else {
		data = url
	}

	if mimeType == "" {
		switch data[:3] {
		case "/9j":
			mimeType = consts.MIME_TYPE_MAP["jpg"]
		case "iVB":
			mimeType = consts.MIME_TYPE_MAP["png"]
		case "Ukl":
			mimeType = consts.MIME_TYPE_MAP["webp"]
		case "R0l":
			mimeType = consts.MIME_TYPE_MAP["gif"]
		case "JVB":
			mimeType = consts.MIME_TYPE_MAP["pdf"]
		default:
			mimeType = consts.MIME_TYPE_MAP["txt"]
		}
	}

	return mimeType, data
}
