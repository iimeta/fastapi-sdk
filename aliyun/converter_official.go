package aliyun

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (a *Aliyun) ConvChatCompletionsRequestOfficial(ctx context.Context, request model.ChatCompletionRequest) ([]byte, error) {

	chatCompletionReq := model.AliyunChatCompletionReq{
		Model: request.Model,
		Input: model.Input{
			Messages: request.Messages,
		},
		Parameters: model.Parameters{
			ResultFormat:      "message",
			MaxTokens:         request.MaxTokens,
			Temperature:       request.Temperature,
			TopP:              request.TopP,
			TopK:              request.N,
			Stop:              request.Stop,
			RepetitionPenalty: request.FrequencyPenalty,
			Seed:              request.Seed,
			Tools:             request.Tools,
			IncrementalOutput: request.Stream,
		},
	}

	if request.ResponseFormat != nil {
		chatCompletionReq.Parameters.ResultFormat = gconv.String(request.ResponseFormat.Type)
	}

	return gjson.MustEncode(chatCompletionReq), nil
}

func (a *Aliyun) ConvChatCompletionsResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Aliyun) ConvChatCompletionsStreamResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
