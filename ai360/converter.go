package ai360

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func (a *AI360) ConvChatCompletionsRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {

	chatCompletionRequest := model.ChatCompletionRequest{}
	if err := gjson.Unmarshal(data, &chatCompletionRequest); err != nil {
		logger.Error(ctx, err)
		return chatCompletionRequest, err
	}

	return chatCompletionRequest, nil
}

func (a *AI360) ConvChatCompletionsResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionResponse := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &chatCompletionResponse); err != nil {
		logger.Error(ctx, err)
		return chatCompletionResponse, err
	}

	return chatCompletionResponse, nil
}

func (a *AI360) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionResponse := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &chatCompletionResponse); err != nil {
		logger.Error(ctx, err)
		return chatCompletionResponse, err
	}

	return chatCompletionResponse, nil
}

func (a *AI360) ConvChatResponsesRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AI360) ConvChatResponsesResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AI360) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AI360) ConvImageGenerationsRequest(ctx context.Context, data []byte) (model.ImageGenerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AI360) ConvImageGenerationsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AI360) ConvImageEditsRequest(ctx context.Context, data []byte) (model.ImageEditRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AI360) ConvImageEditsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AI360) ConvAudioSpeechRequest(ctx context.Context, data []byte) (model.SpeechRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AI360) ConvAudioSpeechResponse(ctx context.Context, data []byte) (model.SpeechResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AI360) ConvAudioTranscriptionsRequest(ctx context.Context, data []byte) (model.AudioRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AI360) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (model.AudioResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AI360) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (model.EmbeddingRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AI360) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (model.EmbeddingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AI360) ConvTextModerationsRequest(ctx context.Context, data []byte) (model.ModerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AI360) ConvTextModerationsResponse(ctx context.Context, data []byte) (model.ModerationResponse, error) {
	//TODO implement me
	panic("implement me")
}
