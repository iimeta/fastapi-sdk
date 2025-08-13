package deepseek

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func (d *DeepSeek) ConvChatCompletionsRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {

	chatCompletionRequest := model.ChatCompletionRequest{}
	if err := gjson.Unmarshal(data, &chatCompletionRequest); err != nil {
		logger.Error(ctx, err)
		return chatCompletionRequest, err
	}

	return chatCompletionRequest, nil
}

func (d *DeepSeek) ConvChatCompletionsResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionResponse := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &chatCompletionResponse); err != nil {
		logger.Error(ctx, err)
		return chatCompletionResponse, err
	}

	return chatCompletionResponse, nil
}

func (d *DeepSeek) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionResponse := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &chatCompletionResponse); err != nil {
		logger.Error(ctx, err)
		return chatCompletionResponse, err
	}

	return chatCompletionResponse, nil
}

func (d *DeepSeek) ConvChatResponsesRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvChatResponsesResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvImageGenerationsRequest(ctx context.Context, data []byte) (model.ImageGenerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvImageGenerationsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvImageEditsRequest(ctx context.Context, data []byte) (model.ImageEditRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvImageEditsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvAudioSpeechRequest(ctx context.Context, data []byte) (model.SpeechRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvAudioSpeechResponse(ctx context.Context, data []byte) (model.SpeechResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvAudioTranscriptionsRequest(ctx context.Context, data []byte) (model.AudioRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (model.AudioResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (model.EmbeddingRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (model.EmbeddingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvTextModerationsRequest(ctx context.Context, data []byte) (model.ModerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvTextModerationsResponse(ctx context.Context, data []byte) (model.ModerationResponse, error) {
	//TODO implement me
	panic("implement me")
}
