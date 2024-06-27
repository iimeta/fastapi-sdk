package common

import (
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
