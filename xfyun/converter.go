package xfyun

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func (x *Xfyun) ConvChatCompletionsRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {

	chatCompletionRequest := model.ChatCompletionRequest{}
	if err := gjson.Unmarshal(data, &chatCompletionRequest); err != nil {
		logger.Error(ctx, err)
		return chatCompletionRequest, err
	}

	return chatCompletionRequest, nil
}

func (x *Xfyun) ConvChatCompletionsResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionResponse := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &chatCompletionResponse); err != nil {
		logger.Error(ctx, err)
		return chatCompletionResponse, err
	}

	return chatCompletionResponse, nil
}

func (x *Xfyun) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionResponse := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &chatCompletionResponse); err != nil {
		logger.Error(ctx, err)
		return chatCompletionResponse, err
	}

	return chatCompletionResponse, nil
}

func (x *Xfyun) ConvChatResponsesRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvChatResponsesResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvImageGenerationsRequest(ctx context.Context, data []byte) (model.ImageGenerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvImageGenerationsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvImageEditsRequest(ctx context.Context, data []byte) (model.ImageEditRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvImageEditsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvAudioSpeechRequest(ctx context.Context, data []byte) (model.SpeechRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvAudioSpeechResponse(ctx context.Context, data []byte) (model.SpeechResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvAudioTranscriptionsRequest(ctx context.Context, data []byte) (model.AudioRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (model.AudioResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (model.EmbeddingRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (model.EmbeddingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvTextModerationsRequest(ctx context.Context, data []byte) (model.ModerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvTextModerationsResponse(ctx context.Context, data []byte) (model.ModerationResponse, error) {
	//TODO implement me
	panic("implement me")
}
