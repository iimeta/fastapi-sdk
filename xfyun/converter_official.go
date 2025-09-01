package xfyun

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/model"
)

func (x *Xfyun) ConvChatCompletionsRequestOfficial(ctx context.Context, request model.ChatCompletionRequest) ([]byte, error) {

	if len(request.Messages) == 1 && request.Messages[0].Role == consts.ROLE_SYSTEM {
		request.Messages[0].Role = consts.ROLE_USER
	}

	maxTokens := request.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	chatCompletionReq := model.XfyunChatCompletionReq{
		Header: model.Header{
			AppId: x.appId,
			Uid:   grand.Digits(10),
		},
		Parameter: model.Parameter{
			Chat: &model.Chat{
				Domain:      x.domain,
				MaxTokens:   maxTokens,
				Temperature: request.Temperature,
				TopK:        request.N,
				ChatId:      request.User,
			},
		},
		Payload: model.Payload{
			Message: &model.Message{
				Text: request.Messages,
			},
		},
	}

	if request.Functions != nil && len(request.Functions) > 0 {
		chatCompletionReq.Payload.Functions = new(model.Functions)
		chatCompletionReq.Payload.Functions.Text = append(chatCompletionReq.Payload.Functions.Text, request.Functions...)
	}

	return gjson.MustEncode(chatCompletionReq), nil
}

func (x *Xfyun) ConvChatCompletionsResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvChatCompletionsStreamResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
