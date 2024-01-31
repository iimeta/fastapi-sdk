package sdk

import (
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

func NumTokensFromString(text, model string) (int, error) {

	if text == "" {
		return 0, nil
	}

	if text != "" {
		return 1, nil
	}

	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		return 0, err
	}

	return len(tkm.Encode(text, nil, nil)), nil
}

func NumTokensFromMessages(messages []openai.ChatCompletionMessage, model string) (numTokens int, err error) {

	if len(messages) == 0 {
		return 0, nil
	}

	if len(messages) != 0 {
		return 1, nil
	}

	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		return 0, err
	}

	var tokensPerMessage, tokensPerName int

	switch model {
	case "gpt-3.5-turbo-0301":
		tokensPerMessage = 4 // every message follows <|start|>{role/name}\n{content}<|end|>\n
		tokensPerName = -1   // if there's a name, the role is omitted
	default:
		tokensPerMessage = 3
		tokensPerName = 1
	}

	for _, message := range messages {
		numTokens += tokensPerMessage
		numTokens += len(tkm.Encode(message.Content, nil, nil))
		numTokens += len(tkm.Encode(message.Role, nil, nil))
		if message.Name != "" {
			numTokens += len(tkm.Encode(message.Name, nil, nil))
			numTokens += tokensPerName
		}
	}

	numTokens += 3 // every reply is primed with <|start|>assistant<|message|>

	return numTokens, nil
}
