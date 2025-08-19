package aliyun

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func (a *Aliyun) ConvChatCompletionsRequestOfficial(ctx context.Context, data []byte) ([]byte, error) {

	request, err := a.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

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

func (a *Aliyun) ConvChatCompletionsResponseOfficial(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionRes := model.AliyunChatCompletionRes{}
	if err := gjson.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return model.ChatCompletionResponse{}, err
	}

	if chatCompletionRes.Code != "" {
		logger.Errorf(ctx, "ConvChatCompletionsResponseOfficial Aliyun model: %s, chatCompletionRes: %s", a.model, gjson.MustEncodeString(chatCompletionRes))

		err := a.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsResponseOfficial Aliyun model: %s, error: %v", a.model, err)

		return model.ChatCompletionResponse{}, err
	}

	response := model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.RequestId,
		Object:  consts.COMPLETION_OBJECT,
		Created: gtime.Timestamp(),
		Model:   a.model,
		Choices: []model.ChatCompletionChoice{{
			Message: &model.ChatCompletionMessage{
				Role:    consts.ROLE_ASSISTANT,
				Content: chatCompletionRes.Output.Text,
			},
		}},
		Usage: &model.Usage{
			PromptTokens:     chatCompletionRes.Usage.InputTokens,
			CompletionTokens: chatCompletionRes.Usage.OutputTokens,
			TotalTokens:      chatCompletionRes.Usage.InputTokens + chatCompletionRes.Usage.OutputTokens,
		},
	}

	return response, nil
}

func (a *Aliyun) ConvChatCompletionsStreamResponseOfficial(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	var (
		usage   *model.Usage
		created = gtime.Timestamp()
		id      = consts.COMPLETION_ID_PREFIX + grand.S(29)
	)

	chatCompletionRes := new(model.AliyunChatCompletionRes)
	if err := gjson.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return model.ChatCompletionResponse{}, err
	}

	if chatCompletionRes.Code != "" {
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponseOfficial Aliyun model: %s, chatCompletionRes: %s", a.model, gjson.MustEncodeString(chatCompletionRes))

		err := a.apiErrorHandler(chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponseOfficial Aliyun model: %s, error: %v", a.model, err)

		return model.ChatCompletionResponse{}, err
	}

	id = consts.COMPLETION_ID_PREFIX + chatCompletionRes.RequestId
	usage = &model.Usage{
		PromptTokens:     chatCompletionRes.Usage.InputTokens,
		CompletionTokens: chatCompletionRes.Usage.OutputTokens,
		TotalTokens:      chatCompletionRes.Usage.InputTokens + chatCompletionRes.Usage.OutputTokens,
	}

	response := model.ChatCompletionResponse{
		Id:      id,
		Object:  consts.COMPLETION_STREAM_OBJECT,
		Created: created,
		Model:   a.model,
		Choices: []model.ChatCompletionChoice{{
			Delta: &model.ChatCompletionStreamChoiceDelta{
				Role:    consts.ROLE_ASSISTANT,
				Content: chatCompletionRes.Output.Text,
			},
		}},
		Usage: usage,
	}

	if response.Usage != nil {
		if len(response.Choices) == 0 {
			response.Choices = append(response.Choices, model.ChatCompletionChoice{
				Delta:        new(model.ChatCompletionStreamChoiceDelta),
				FinishReason: consts.FinishReasonStop,
			})
		}
	}

	return response, nil
}
