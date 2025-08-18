package anthropic

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/go-openai"
)

func (a *Anthropic) ConvChatCompletionsRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {

	request := model.ChatCompletionRequest{}
	if err := gjson.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	if a.isSupportSystemRole != nil {
		request.Messages = common.HandleMessages(request.Messages, *a.isSupportSystemRole)
	} else {
		request.Messages = common.HandleMessages(request.Messages, true)
	}

	return request, nil
}

func (a *Anthropic) ConvChatCompletionsResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	response := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (a *Anthropic) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	response := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

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

func (a *Anthropic) ConvChatCompletionsResponseOfficial(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionRes := model.AnthropicChatCompletionRes{}
	if err := gjson.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return model.ChatCompletionResponse{}, err
	}

	if chatCompletionRes.Error != nil && chatCompletionRes.Error.Type != "" {
		logger.Errorf(ctx, "ConvChatCompletionsResponseOfficial Anthropic model: %s, chatCompletionRes: %s", a.model, gjson.MustEncodeString(chatCompletionRes))

		err := a.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsResponseOfficial Anthropic model: %s, error: %v", a.model, err)

		return model.ChatCompletionResponse{}, err
	}

	response := model.ChatCompletionResponse{
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
	}

	for _, content := range chatCompletionRes.Content {
		if content.Type == consts.DELTA_TYPE_INPUT_JSON {
			response.Choices = append(response.Choices, model.ChatCompletionChoice{
				Delta: &model.ChatCompletionStreamChoiceDelta{
					Role: consts.ROLE_ASSISTANT,
					ToolCalls: []openai.ToolCall{{
						Function: openai.FunctionCall{
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

func (a *Anthropic) ConvChatCompletionsStreamResponseOfficial(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionRes := new(model.AnthropicChatCompletionRes)
	if err := gjson.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return model.ChatCompletionResponse{}, err
	}

	if chatCompletionRes.Error != nil && chatCompletionRes.Error.Type != "" {
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponseOfficial Anthropic model: %s, chatCompletionRes: %s", a.model, gjson.MustEncodeString(chatCompletionRes))

		err := a.apiErrorHandler(chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponseOfficial Anthropic model: %s, error: %v", a.model, err)

		return model.ChatCompletionResponse{}, err
	}

	response := model.ChatCompletionResponse{
		Id:      chatCompletionRes.Message.Id,
		Object:  consts.COMPLETION_STREAM_OBJECT,
		Created: gtime.Timestamp(),
		Model:   a.model,
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

func (a *Anthropic) ConvChatResponsesRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ConvChatResponsesResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ConvImageGenerationsRequest(ctx context.Context, data []byte) (model.ImageGenerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ConvImageGenerationsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ConvImageEditsRequest(ctx context.Context, data []byte) (model.ImageEditRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ConvImageEditsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ConvAudioSpeechRequest(ctx context.Context, data []byte) (model.SpeechRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ConvAudioSpeechResponse(ctx context.Context, data []byte) (model.SpeechResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ConvAudioTranscriptionsRequest(ctx context.Context, data []byte) (model.AudioRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (model.AudioResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (model.EmbeddingRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (model.EmbeddingResponse, error) {
	//TODO implement me
	panic("implement me")
}
