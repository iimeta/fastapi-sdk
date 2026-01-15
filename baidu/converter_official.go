package baidu

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (b *Baidu) ConvChatCompletionsRequestOfficial(ctx context.Context, request model.ChatCompletionRequest) ([]byte, error) {

	chatCompletionReq := model.BaiduChatCompletionReq{
		Messages:        request.Messages,
		MaxOutputTokens: request.MaxTokens,
		Temperature:     request.Temperature,
		TopP:            request.TopP,
		Stream:          request.Stream,
		Stop:            request.Stop,
		PenaltyScore:    request.FrequencyPenalty,
		UserId:          request.User,
	}

	if chatCompletionReq.Messages[0].Role == consts.ROLE_SYSTEM {
		chatCompletionReq.System = gconv.String(chatCompletionReq.Messages[0].Content)
		chatCompletionReq.Messages = chatCompletionReq.Messages[1:]
	}

	if request.ResponseFormat != nil {
		chatCompletionReq.ResponseFormat = gconv.String(request.ResponseFormat.Type)
	}

	return gjson.MustEncode(chatCompletionReq), nil
}

func (b *Baidu) ConvChatCompletionsResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvChatCompletionsStreamResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvImageGenerationsRequestOfficial(ctx context.Context, request model.ImageGenerationRequest) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvImageGenerationsResponseOfficial(ctx context.Context, response model.ImageResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvImageEditsRequestOfficial(ctx context.Context, request model.ImageEditRequest) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvImageEditsResponseOfficial(ctx context.Context, response model.ImageResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
