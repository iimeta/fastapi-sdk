package xfyun

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func (x *Xfyun) ConvChatCompletionsRequestOfficial(ctx context.Context, data []byte) ([]byte, error) {

	request, err := x.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

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

func (x *Xfyun) ConvChatCompletionsResponseOfficial(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionRes := model.XfyunChatCompletionRes{}
	if err := gjson.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return model.ChatCompletionResponse{}, err
	}

	if chatCompletionRes.Header.Code != 0 {
		logger.Errorf(ctx, "ChatCompletions Xfyun model: %s, chatCompletionRes: %s", x.model, gjson.MustEncodeString(chatCompletionRes))

		err := x.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ChatCompletions Xfyun model: %s, error: %v", x.model, err)

		return model.ChatCompletionResponse{}, err
	}

	response := model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Header.Sid,
		Object:  consts.COMPLETION_OBJECT,
		Created: gtime.Timestamp(),
		Model:   x.model,
		Choices: []model.ChatCompletionChoice{{
			Index: chatCompletionRes.Payload.Choices.Seq,
			Message: &model.ChatCompletionMessage{
				Role:         chatCompletionRes.Payload.Choices.Text[0].Role,
				FunctionCall: chatCompletionRes.Payload.Choices.Text[0].FunctionCall,
			},
		}},
		Usage: &model.Usage{
			PromptTokens:     chatCompletionRes.Payload.Usage.Text.PromptTokens,
			CompletionTokens: chatCompletionRes.Payload.Usage.Text.CompletionTokens,
			TotalTokens:      chatCompletionRes.Payload.Usage.Text.TotalTokens,
		},
	}

	return response, nil
}

func (x *Xfyun) ConvChatCompletionsStreamResponseOfficial(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionRes := new(model.XfyunChatCompletionRes)
	if err := gjson.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return model.ChatCompletionResponse{}, err
	}

	if chatCompletionRes.Header.Code != 0 {
		logger.Errorf(ctx, "ChatCompletionsStream Xfyun model: %s, chatCompletionRes: %s", x.model, gjson.MustEncodeString(chatCompletionRes))

		err := x.apiErrorHandler(chatCompletionRes)
		logger.Errorf(ctx, "ChatCompletionsStream Xfyun model: %s, error: %v", x.model, err)

		return model.ChatCompletionResponse{}, err
	}

	response := model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Header.Sid,
		Object:  consts.COMPLETION_STREAM_OBJECT,
		Created: gtime.Timestamp(),
		Model:   x.model,
		Choices: []model.ChatCompletionChoice{{
			Index: chatCompletionRes.Payload.Choices.Seq,
			Delta: &model.ChatCompletionStreamChoiceDelta{
				Role:         chatCompletionRes.Payload.Choices.Text[0].Role,
				Content:      chatCompletionRes.Payload.Choices.Text[0].Content,
				FunctionCall: chatCompletionRes.Payload.Choices.Text[0].FunctionCall,
			},
		}},
	}

	if chatCompletionRes.Payload.Usage != nil {
		response.Usage = &model.Usage{
			PromptTokens:     chatCompletionRes.Payload.Usage.Text.PromptTokens,
			CompletionTokens: chatCompletionRes.Payload.Usage.Text.CompletionTokens,
			TotalTokens:      chatCompletionRes.Payload.Usage.Text.TotalTokens,
		}
	}

	return response, nil
}
