package anthropic

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func (a *Anthropic) ConvChatCompletionsRequestOfficial(ctx context.Context, data []byte) ([]byte, error) {

	request, err := a.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	chatCompletionReq := model.AnthropicChatCompletionReq{
		Model:         request.Model,
		Messages:      request.Messages,
		MaxTokens:     request.MaxTokens,
		StopSequences: request.Stop,
		Stream:        request.Stream,
		Temperature:   request.Temperature,
		ToolChoice:    request.ToolChoice,
		TopK:          request.TopK,
		TopP:          request.TopP,
		Tools:         request.Tools,
	}

	if chatCompletionReq.Messages[0].Role == consts.ROLE_SYSTEM {
		chatCompletionReq.System = chatCompletionReq.Messages[0].Content
		chatCompletionReq.Messages = chatCompletionReq.Messages[1:]
	}

	if request.User != "" {
		chatCompletionReq.Metadata = &model.Metadata{
			UserId: request.User,
		}
	}

	if chatCompletionReq.MaxTokens == 0 {
		chatCompletionReq.MaxTokens = 4096
	}

	for _, message := range chatCompletionReq.Messages {

		if contents, ok := message.Content.([]interface{}); ok {

			for _, value := range contents {

				if content, ok := value.(map[string]interface{}); ok {

					if content["type"] == "image_url" {

						if imageUrl, ok := content["image_url"].(map[string]interface{}); ok {

							mimeType, data := common.GetMime(gconv.String(imageUrl["url"]))

							content["source"] = model.Source{
								Type:      "base64",
								MediaType: mimeType,
								Data:      data,
							}

							content["type"] = "image"
							delete(content, "image_url")
						}
					}
				}
			}
		}
	}

	if a.isGcp {
		chatCompletionReq.Model = ""
		chatCompletionReq.AnthropicVersion = "vertex-2023-10-16"
	}

	return gjson.MustEncode(chatCompletionReq), nil
}

func (a *Anthropic) ConvChatCompletionsResponseOfficial(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	chatCompletionRes := model.AnthropicChatCompletionRes{}
	if err = json.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if chatCompletionRes.Error != nil && chatCompletionRes.Error.Type != "" {
		logger.Errorf(ctx, "ConvChatCompletionsResponseOfficial Anthropic model: %s, chatCompletionRes: %s", a.model, gjson.MustEncodeString(chatCompletionRes))

		err = a.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsResponseOfficial Anthropic model: %s, error: %v", a.model, err)

		return response, err
	}

	response = model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Id,
		Object:  consts.COMPLETION_OBJECT,
		Created: gtime.Timestamp(),
		Model:   a.model,
		Usage: &model.Usage{
			PromptTokens:             chatCompletionRes.Usage.InputTokens,
			CompletionTokens:         chatCompletionRes.Usage.OutputTokens,
			TotalTokens:              chatCompletionRes.Usage.InputTokens + chatCompletionRes.Usage.OutputTokens,
			CacheCreationInputTokens: chatCompletionRes.Usage.CacheCreationInputTokens,
			CacheReadInputTokens:     chatCompletionRes.Usage.CacheReadInputTokens,
		},
		ResponseBytes: data,
	}

	for _, content := range chatCompletionRes.Content {
		if content.Type == consts.DELTA_TYPE_INPUT_JSON {
			response.Choices = append(response.Choices, model.ChatCompletionChoice{
				Delta: &model.ChatCompletionStreamChoiceDelta{
					Role: consts.ROLE_ASSISTANT,
					ToolCalls: []model.ToolCall{{
						Function: model.FunctionCall{
							Arguments: content.PartialJson,
						},
					}},
				},
			})
		} else {
			response.Choices = append(response.Choices, model.ChatCompletionChoice{
				Message: &model.ChatCompletionMessage{
					Role:    chatCompletionRes.Role,
					Content: content.Text,
				},
				FinishReason: "stop",
			})
		}
	}

	return response, nil
}

func (a *Anthropic) ConvChatCompletionsStreamResponseOfficial(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	chatCompletionRes := model.AnthropicChatCompletionRes{}
	if err = json.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if chatCompletionRes.Error != nil && chatCompletionRes.Error.Type != "" {
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponseOfficial Anthropic model: %s, chatCompletionRes: %s", a.model, gjson.MustEncodeString(chatCompletionRes))

		err = a.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponseOfficial Anthropic model: %s, error: %v", a.model, err)

		return response, err
	}

	response = model.ChatCompletionResponse{
		Id:            chatCompletionRes.Message.Id,
		Object:        consts.COMPLETION_STREAM_OBJECT,
		Created:       gtime.Timestamp(),
		Model:         a.model,
		ResponseBytes: data,
	}

	if chatCompletionRes.Usage != nil {
		response.Usage = &model.Usage{
			PromptTokens:             chatCompletionRes.Usage.InputTokens,
			CompletionTokens:         chatCompletionRes.Usage.OutputTokens,
			TotalTokens:              chatCompletionRes.Usage.InputTokens + chatCompletionRes.Usage.OutputTokens,
			CacheCreationInputTokens: chatCompletionRes.Usage.CacheCreationInputTokens,
			CacheReadInputTokens:     chatCompletionRes.Usage.CacheReadInputTokens,
		}
	}

	if chatCompletionRes.Message.Usage != nil {
		response.Usage = &model.Usage{
			PromptTokens:             chatCompletionRes.Message.Usage.InputTokens,
			CacheCreationInputTokens: chatCompletionRes.Message.Usage.CacheCreationInputTokens,
			CacheReadInputTokens:     chatCompletionRes.Message.Usage.CacheReadInputTokens,
		}
	}

	if chatCompletionRes.Delta.StopReason != "" {
		response.Choices = append(response.Choices, model.ChatCompletionChoice{
			FinishReason: consts.FinishReasonStop,
		})
	} else {
		if chatCompletionRes.Delta.Type == consts.DELTA_TYPE_INPUT_JSON {
			response.Choices = append(response.Choices, model.ChatCompletionChoice{
				Delta: &model.ChatCompletionStreamChoiceDelta{
					Role: consts.ROLE_ASSISTANT,
					ToolCalls: []model.ToolCall{{
						Function: model.FunctionCall{
							Arguments: chatCompletionRes.Delta.PartialJson,
						},
					}},
				},
			})
		} else {
			response.Choices = append(response.Choices, model.ChatCompletionChoice{
				Delta: &model.ChatCompletionStreamChoiceDelta{
					Role:    consts.ROLE_ASSISTANT,
					Content: chatCompletionRes.Delta.Text,
				},
			})
		}
	}

	return response, nil
}
