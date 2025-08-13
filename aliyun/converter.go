package aliyun

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func (a *Aliyun) ConvChatCompletionsRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {

	chatCompletionRequest := model.ChatCompletionRequest{}
	if err := gjson.Unmarshal(data, &chatCompletionRequest); err != nil {
		logger.Error(ctx, err)
		return chatCompletionRequest, err
	}

	return chatCompletionRequest, nil
}

func (a *Aliyun) ConvChatCompletionsResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionResponse := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &chatCompletionResponse); err != nil {
		logger.Error(ctx, err)
		return chatCompletionResponse, err
	}

	return chatCompletionResponse, nil
}

func (a *Aliyun) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionResponse := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &chatCompletionResponse); err != nil {
		logger.Error(ctx, err)
		return chatCompletionResponse, err
	}

	return chatCompletionResponse, nil
}

func (a *Aliyun) ConvChatResponsesRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Aliyun) ConvChatResponsesResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Aliyun) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
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

func (a *Aliyun) ConvImageEditsRequest(ctx context.Context, data []byte) (model.ImageEditRequest, error) {
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

func (a *Aliyun) ConvAudioTranscriptionsRequest(ctx context.Context, data []byte) (model.AudioRequest, error) {
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

func (a *Aliyun) ConvTextModerationsRequest(ctx context.Context, data []byte) (model.ModerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Aliyun) ConvTextModerationsResponse(ctx context.Context, data []byte) (model.ModerationResponse, error) {
	//TODO implement me
	panic("implement me")
}
