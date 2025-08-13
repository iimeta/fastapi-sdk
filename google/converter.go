package google

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func (g *Google) ConvChatCompletionsRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {

	chatCompletionRequest := model.ChatCompletionRequest{}
	if err := gjson.Unmarshal(data, &chatCompletionRequest); err != nil {
		logger.Error(ctx, err)
		return chatCompletionRequest, err
	}

	return chatCompletionRequest, nil
}

func (g *Google) ConvChatCompletionsResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionResponse := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &chatCompletionResponse); err != nil {
		logger.Error(ctx, err)
		return chatCompletionResponse, err
	}

	return chatCompletionResponse, nil
}

func (g *Google) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionResponse := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &chatCompletionResponse); err != nil {
		logger.Error(ctx, err)
		return chatCompletionResponse, err
	}

	return chatCompletionResponse, nil
}

func (g *Google) ConvChatResponsesRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvChatResponsesResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvImageGenerationsRequest(ctx context.Context, data []byte) (model.ImageGenerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvImageGenerationsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvImageEditsRequest(ctx context.Context, data []byte) (model.ImageEditRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvImageEditsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvAudioSpeechRequest(ctx context.Context, data []byte) (model.SpeechRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvAudioSpeechResponse(ctx context.Context, data []byte) (model.SpeechResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvAudioTranscriptionsRequest(ctx context.Context, data []byte) (model.AudioRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (model.AudioResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (model.EmbeddingRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (model.EmbeddingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvTextModerationsRequest(ctx context.Context, data []byte) (model.ModerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvTextModerationsResponse(ctx context.Context, data []byte) (model.ModerationResponse, error) {
	//TODO implement me
	panic("implement me")
}
