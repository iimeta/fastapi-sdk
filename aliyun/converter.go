package aliyun

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

func (a *Aliyun) ConvChatCompletionsRequest(ctx context.Context, data any) (request model.ChatCompletionRequest, err error) {

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

func (a *Aliyun) ConvChatCompletionsResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	chatCompletionRes := model.AliyunChatCompletionRes{}
	if err = json.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if chatCompletionRes.Code != "" {
		logger.Errorf(ctx, "ConvChatCompletionsResponse Aliyun model: %s, chatCompletionRes: %s", a.Model, gjson.MustEncodeString(chatCompletionRes))

		err = a.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsResponse Aliyun model: %s, error: %v", a.Model, err)

		return response, err
	}

	response = model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.RequestId,
		Object:  consts.COMPLETION_OBJECT,
		Created: gtime.Timestamp(),
		Model:   a.Model,
		Choices: []model.ChatCompletionChoice{{
			Message: &model.ChatCompletionMessage{
				Role:    consts.ROLE_ASSISTANT,
				Content: chatCompletionRes.Output.Text,
			},
		}},
		Usage: &model.Usage{
			PromptTokens:     chatCompletionRes.Usage.InputTokens,
			CompletionTokens: chatCompletionRes.Usage.OutputTokens,
			TotalTokens:      chatCompletionRes.Usage.InputTokens + chatCompletionRes.Usage.OutputTokens,
		},
		ResponseBytes: data,
	}

	return response, nil
}

func (a *Aliyun) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	chatCompletionRes := model.AliyunChatCompletionRes{}
	if err = json.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if chatCompletionRes.Code != "" {
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponse Aliyun model: %s, chatCompletionRes: %s", a.Model, gjson.MustEncodeString(chatCompletionRes))

		err = a.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponse Aliyun model: %s, error: %v", a.Model, err)

		return response, err
	}

	response = model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.RequestId,
		Object:  consts.COMPLETION_STREAM_OBJECT,
		Created: gtime.Timestamp(),
		Model:   a.Model,
		Choices: []model.ChatCompletionChoice{{
			Delta: &model.ChatCompletionStreamChoiceDelta{
				Role:    consts.ROLE_ASSISTANT,
				Content: chatCompletionRes.Output.Text,
			},
		}},
		Usage: &model.Usage{
			PromptTokens:     chatCompletionRes.Usage.InputTokens,
			CompletionTokens: chatCompletionRes.Usage.OutputTokens,
			TotalTokens:      chatCompletionRes.Usage.InputTokens + chatCompletionRes.Usage.OutputTokens,
		},
		ResponseBytes: data,
	}

	// todo
	if response.Usage != nil {
		if len(response.Choices) == 0 {
			response.Choices = append(response.Choices, model.ChatCompletionChoice{
				Delta:        new(model.ChatCompletionStreamChoiceDelta),
				FinishReason: consts.FinishReasonStop,
			})
		}
	}

	return response, nil
}

func (a *Aliyun) ConvChatResponsesRequest(ctx context.Context, data []byte) (request model.ChatCompletionRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *Aliyun) ConvChatResponsesResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *Aliyun) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *Aliyun) ConvImageGenerationsRequest(ctx context.Context, data []byte) (model.ImageGenerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Aliyun) ConvImageGenerationsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Aliyun) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (*bytes.Buffer, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Aliyun) ConvImageEditsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Aliyun) ConvAudioSpeechRequest(ctx context.Context, data []byte) (model.SpeechRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Aliyun) ConvAudioSpeechResponse(ctx context.Context, data []byte) (model.SpeechResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Aliyun) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (*bytes.Buffer, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Aliyun) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (model.AudioResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Aliyun) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (model.EmbeddingRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Aliyun) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (model.EmbeddingResponse, error) {
	//TODO implement me
	panic("implement me")
}
