package baidu

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func (b *Baidu) ConvChatCompletionsRequestOfficial(ctx context.Context, data []byte) ([]byte, error) {

	request, err := b.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

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

func (b *Baidu) ConvChatCompletionsResponseOfficial(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionRes := model.BaiduChatCompletionRes{}
	if err := gjson.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return model.ChatCompletionResponse{}, err
	}

	if chatCompletionRes.ErrorCode != 0 {
		logger.Errorf(ctx, "ConvChatCompletionsResponseOfficial Baidu model: %s, chatCompletionRes: %s", b.model, gjson.MustEncodeString(chatCompletionRes))

		err := b.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsResponseOfficial Baidu model: %s, error: %v", b.model, err)

		return model.ChatCompletionResponse{}, err
	}

	response := model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Id,
		Object:  consts.COMPLETION_OBJECT,
		Created: chatCompletionRes.Created,
		Model:   b.model,
		Choices: []model.ChatCompletionChoice{{
			Message: &model.ChatCompletionMessage{
				Role:    consts.ROLE_ASSISTANT,
				Content: chatCompletionRes.Result,
			},
		}},
		Usage: chatCompletionRes.Usage,
	}

	return response, nil
}

func (b *Baidu) ConvChatCompletionsStreamResponseOfficial(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionRes := new(model.BaiduChatCompletionRes)
	if err := gjson.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return model.ChatCompletionResponse{}, err
	}

	if chatCompletionRes.ErrorCode != 0 {
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponseOfficial Baidu model: %s, chatCompletionRes: %s", b.model, gjson.MustEncodeString(chatCompletionRes))

		err := b.apiErrorHandler(chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponseOfficial Baidu model: %s, error: %v", b.model, err)

		return model.ChatCompletionResponse{}, err
	}

	response := model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Id,
		Object:  consts.COMPLETION_STREAM_OBJECT,
		Created: chatCompletionRes.Created,
		Model:   b.model,
		Choices: []model.ChatCompletionChoice{{
			Index: chatCompletionRes.SentenceId,
			Delta: &model.ChatCompletionStreamChoiceDelta{
				Role:    consts.ROLE_ASSISTANT,
				Content: chatCompletionRes.Result,
			},
		}},
		Usage: chatCompletionRes.Usage,
	}

	return response, nil
}
