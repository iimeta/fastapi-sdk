package anthropic

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func (a *Anthropic) ConvChatCompletionsRequest(ctx context.Context, data any) (request model.ChatCompletionRequest, err error) {

	if v, ok := data.(model.ChatCompletionRequest); ok {
		request = v
	} else if v, ok := data.([]byte); ok {
		if err = json.Unmarshal(v, &request); err != nil {
			logger.Error(ctx, err)
			return request, err
		}
	} else {
		if err = json.Unmarshal(gjson.MustEncode(data), &request); err != nil {
			logger.Error(ctx, err)
			return request, err
		}
	}

	if a.IsSupportSystemRole != nil {
		request.Messages = common.HandleMessages(request.Messages, *a.IsSupportSystemRole)
	} else {
		request.Messages = common.HandleMessages(request.Messages, true)
	}

	return request, nil
}

func (a *Anthropic) ConvChatCompletionsResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	chatCompletionRes := model.AnthropicChatCompletionRes{}
	if err = json.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if chatCompletionRes.Error != nil && chatCompletionRes.Error.Type != "" {
		logger.Errorf(ctx, "ConvChatCompletionsResponse Anthropic model: %s, chatCompletionRes: %s", a.Model, gjson.MustEncodeString(chatCompletionRes))

		err = a.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsResponse Anthropic model: %s, error: %v", a.Model, err)

		return response, err
	}

	response = model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Id,
		Object:  consts.COMPLETION_OBJECT,
		Created: gtime.Timestamp(),
		Model:   a.Model,
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
				FinishReason: consts.FinishReasonStop,
			})
		}
	}

	return response, nil
}

func (a *Anthropic) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	chatCompletionRes := model.AnthropicChatCompletionRes{}
	if err = json.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if chatCompletionRes.Error != nil && chatCompletionRes.Error.Type != "" {
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponse Anthropic model: %s, chatCompletionRes: %s", a.Model, gjson.MustEncodeString(chatCompletionRes))

		err = a.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponse Anthropic model: %s, error: %v", a.Model, err)

		return response, err
	}

	response = model.ChatCompletionResponse{
		Id:            chatCompletionRes.Message.Id,
		Object:        consts.COMPLETION_STREAM_OBJECT,
		Created:       gtime.Timestamp(),
		Model:         a.Model,
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
			CompletionTokens:         chatCompletionRes.Message.Usage.OutputTokens,
			TotalTokens:              chatCompletionRes.Message.Usage.InputTokens + chatCompletionRes.Message.Usage.OutputTokens,
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

func (a *Anthropic) ConvChatResponsesRequest(ctx context.Context, data []byte) (request model.ChatCompletionRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ConvChatResponsesResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
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

func (a *Anthropic) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (*bytes.Buffer, error) {
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

func (a *Anthropic) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (*bytes.Buffer, error) {
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
