package zhipuai

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func (z *ZhipuAI) ConvChatCompletionsRequestOfficial(ctx context.Context, data []byte) ([]byte, error) {

	request, err := z.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	chatCompletionReq := model.ZhipuAIChatCompletionReq{
		Model:       z.model,
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

func (z *ZhipuAI) ConvChatCompletionsResponseOfficial(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionRes := model.ZhipuAIChatCompletionRes{}
	if err := gjson.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return model.ChatCompletionResponse{}, err
	}

	if chatCompletionRes.Error.Code != "" && chatCompletionRes.Error.Code != "200" {
		logger.Errorf(ctx, "ConvChatCompletionsResponseOfficial ZhipuAI model: %s, chatCompletionRes: %s", z.model, gjson.MustEncodeString(chatCompletionRes))

		err := z.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsResponseOfficial ZhipuAI model: %s, error: %v", z.model, err)

		return model.ChatCompletionResponse{}, err
	}

	response := model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Id,
		Object:  consts.COMPLETION_OBJECT,
		Created: chatCompletionRes.Created,
		Model:   z.model,
		Usage:   chatCompletionRes.Usage,
	}

	for _, choice := range chatCompletionRes.Choices {
		response.Choices = append(response.Choices, model.ChatCompletionChoice{
			Index: choice.Index,
			Message: &model.ChatCompletionMessage{
				Role:         choice.Message.Role,
				Content:      choice.Message.Content,
				Refusal:      choice.Message.Refusal,
				Name:         choice.Message.Name,
				FunctionCall: choice.Message.FunctionCall,
				ToolCalls:    choice.Message.ToolCalls,
				ToolCallID:   choice.Message.ToolCallID,
				Audio:        choice.Message.Audio,
			},
			FinishReason: choice.FinishReason,
		})
	}

	return response, nil
}

func (z *ZhipuAI) ConvChatCompletionsStreamResponseOfficial(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionRes := new(model.ZhipuAIChatCompletionRes)
	if err := gjson.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return model.ChatCompletionResponse{}, err
	}

	if chatCompletionRes.Error.Code != "" && chatCompletionRes.Error.Code != "200" {
		logger.Errorf(ctx, "ChatCompletionsStream ZhipuAI model: %s, chatCompletionRes: %s", z.model, gjson.MustEncodeString(chatCompletionRes))

		err := z.apiErrorHandler(chatCompletionRes)
		logger.Errorf(ctx, "ChatCompletionsStream ZhipuAI model: %s, error: %v", z.model, err)

		return model.ChatCompletionResponse{}, err
	}

	response := model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Id,
		Object:  consts.COMPLETION_STREAM_OBJECT,
		Created: chatCompletionRes.Created,
		Model:   z.model,
		Usage:   chatCompletionRes.Usage,
	}

	for _, choice := range chatCompletionRes.Choices {
		response.Choices = append(response.Choices, model.ChatCompletionChoice{
			Index: choice.Index,
			Delta: &model.ChatCompletionStreamChoiceDelta{
				Content:      choice.Delta.Content,
				Role:         choice.Delta.Role,
				FunctionCall: choice.Delta.FunctionCall,
				ToolCalls:    choice.Delta.ToolCalls,
				Refusal:      choice.Delta.Refusal,
				Audio:        choice.Delta.Audio,
			},
			FinishReason: choice.FinishReason,
		})
	}

	return response, nil
}
