package anthropic

import (
	"bytes"
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
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
