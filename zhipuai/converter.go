package zhipuai

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func (z *ZhipuAI) ConvChatCompletionsRequest(ctx context.Context, data any) (request model.ChatCompletionRequest, err error) {

	request = model.ChatCompletionRequest{}

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

	if z.IsSupportSystemRole != nil {
		request.Messages = common.HandleMessages(request.Messages, *z.IsSupportSystemRole)
	} else {
		request.Messages = common.HandleMessages(request.Messages, true)
	}

	return request, nil
}

func (z *ZhipuAI) ConvChatCompletionsResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	chatCompletionRes := model.ZhipuAIChatCompletionRes{}
	if err = json.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if chatCompletionRes.Error.Code != "" && chatCompletionRes.Error.Code != "200" {
		logger.Errorf(ctx, "ConvChatCompletionsResponse ZhipuAI model: %s, chatCompletionRes: %s", z.Model, gjson.MustEncodeString(chatCompletionRes))

		err = z.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsResponse ZhipuAI model: %s, error: %v", z.Model, err)

		return response, err
	}

	response = model.ChatCompletionResponse{
		Id:            consts.COMPLETION_ID_PREFIX + chatCompletionRes.Id,
		Object:        consts.COMPLETION_OBJECT,
		Created:       chatCompletionRes.Created,
		Model:         z.Model,
		Choices:       chatCompletionRes.Choices,
		Usage:         chatCompletionRes.Usage,
		ResponseBytes: data,
	}

	return response, nil
}

func (z *ZhipuAI) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	chatCompletionRes := model.ZhipuAIChatCompletionRes{}
	if err = json.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if chatCompletionRes.Error.Code != "" && chatCompletionRes.Error.Code != "200" {
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponse ZhipuAI model: %s, chatCompletionRes: %s", z.Model, gjson.MustEncodeString(chatCompletionRes))

		err = z.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponse ZhipuAI model: %s, error: %v", z.Model, err)

		return response, err
	}

	response = model.ChatCompletionResponse{
		Id:            consts.COMPLETION_ID_PREFIX + chatCompletionRes.Id,
		Object:        consts.COMPLETION_STREAM_OBJECT,
		Created:       chatCompletionRes.Created,
		Model:         z.Model,
		Choices:       chatCompletionRes.Choices,
		Usage:         chatCompletionRes.Usage,
		ResponseBytes: data,
	}

	return response, nil
}

func (z *ZhipuAI) ConvChatResponsesRequest(ctx context.Context, data []byte) (request model.ChatCompletionRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvChatResponsesResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvImageGenerationsRequest(ctx context.Context, data []byte) (model.ImageGenerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvImageGenerationsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (*bytes.Buffer, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvImageEditsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvAudioSpeechRequest(ctx context.Context, data []byte) (model.SpeechRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvAudioSpeechResponse(ctx context.Context, data []byte) (model.SpeechResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (*bytes.Buffer, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (model.AudioResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (model.EmbeddingRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (model.EmbeddingResponse, error) {
	//TODO implement me
	panic("implement me")
}
