package zhipuai

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/iimeta/fastapi-sdk/model"
)

func (z *ZhipuAI) ConvChatCompletionsRequestOfficial(ctx context.Context, request model.ChatCompletionRequest) ([]byte, error) {

	chatCompletionReq := model.ZhipuAIChatCompletionReq{
		Model:       z.Model,
		Messages:    request.Messages,
		MaxTokens:   request.MaxTokens,
		Temperature: request.Temperature,
		TopP:        request.TopP,
		Stream:      request.Stream,
		Stop:        request.Stop,
		Tools:       request.Tools,
		ToolChoice:  request.ToolChoice,
		UserId:      request.User,
	}

	if chatCompletionReq.TopP == 1 {
		chatCompletionReq.TopP -= 0.01
	} else if chatCompletionReq.TopP == 0 {
		chatCompletionReq.TopP += 0.01
	}

	if chatCompletionReq.Temperature == 1 {
		chatCompletionReq.Temperature -= 0.01
	} else if chatCompletionReq.Temperature == 0 {
		chatCompletionReq.Temperature += 0.01
	}

	if chatCompletionReq.MaxTokens == 1 {
		chatCompletionReq.MaxTokens = 2
	}

	return gjson.MustEncode(chatCompletionReq), nil
}

func (z *ZhipuAI) ConvChatCompletionsResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvChatCompletionsStreamResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
