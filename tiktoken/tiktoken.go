package tiktoken

import (
	"encoding/json"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/tiktoken-go"
)

func NumTokensFromString(model, text string) (int, error) {

	if text == "" {
		return 0, nil
	}

	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		return 0, err
	}

	return len(tkm.Encode(text, nil, nil)), nil
}

func NumTokensFromMessages(model string, messages []model.ChatCompletionMessage) (numTokens int, err error) {

	if len(messages) == 0 {
		return 0, nil
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
		numTokens += NumTokensFromContent(tkm, model, message.Content)
		numTokens += len(tkm.Encode(message.Role, nil, nil))

		if message.Name != "" {
			numTokens += len(tkm.Encode(message.Name, nil, nil))
			numTokens += tokensPerName
		}
	}

	numTokens += 3 // every reply is primed with <|start|>assistant<|message|>

	return numTokens, nil
}

func EncodingForModel(model string) (*tiktoken.Tiktoken, error) {
	return tiktoken.EncodingForModel(model)
}

func IsEncodingForModel(model string) bool {
	return tiktoken.IsEncodingForModel(model)
}

func NumTokensFromContent(tkm *tiktoken.Tiktoken, model string, content any) (numTokens int) {

	text := gconv.String(content)

	// 传入base64图片内容
	if gstr.Contains(text, "data:image/") {

		var data interface{}
		if err := json.Unmarshal([]byte(text), &data); err != nil {
			return len(tkm.Encode(text, nil, nil))
		}

		if result, ok := data.([]interface{}); ok {
			for _, value := range result {
				if content, ok := value.(map[string]interface{}); ok {
					if content["type"] == "text" {
						numTokens += len(tkm.Encode(gconv.String(content["text"]), nil, nil))
					} else if content["type"] == "image_url" || content["type"] == "image" {
						// 兼容目前计算错误情况
						if model == "gpt-4o-mini" {
							numTokens += 1023 * 36
						} else {
							// 其它多模态模型
							numTokens += 1023
						}
					} else {
						numTokens += len(tkm.Encode(gconv.String(content), nil, nil))
					}
				}
			}
		}

		return numTokens
	}

	return len(tkm.Encode(text, nil, nil))
}
