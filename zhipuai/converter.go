package zhipuai

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func (z *ZhipuAI) ConvChatCompletionsRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {

	chatCompletionRequest := model.ChatCompletionRequest{}
	if err := gjson.Unmarshal(data, &chatCompletionRequest); err != nil {
		logger.Error(ctx, err)
		return chatCompletionRequest, err
	}

	if z.isSupportSystemRole != nil {
		chatCompletionRequest.Messages = common.HandleMessages(chatCompletionRequest.Messages, *z.isSupportSystemRole)
	} else {
		chatCompletionRequest.Messages = common.HandleMessages(chatCompletionRequest.Messages, true)
	}

	return chatCompletionRequest, nil
}

func (z *ZhipuAI) ConvChatCompletionsResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionResponse := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &chatCompletionResponse); err != nil {
		logger.Error(ctx, err)
		return chatCompletionResponse, err
	}

	return chatCompletionResponse, nil
}

func (z *ZhipuAI) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionResponse := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &chatCompletionResponse); err != nil {
		logger.Error(ctx, err)
		return chatCompletionResponse, err
	}

	return chatCompletionResponse, nil
}

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

func (z *ZhipuAI) ConvChatResponsesRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvChatResponsesResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvImageGenerationsRequest(ctx context.Context, data []byte) (model.ImageGenerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvImageGenerationsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvImageEditsRequest(ctx context.Context, data []byte) (model.ImageEditRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvImageEditsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvAudioSpeechRequest(ctx context.Context, data []byte) (model.SpeechRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvAudioSpeechResponse(ctx context.Context, data []byte) (model.SpeechResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvAudioTranscriptionsRequest(ctx context.Context, data []byte) (model.AudioRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (model.AudioResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (model.EmbeddingRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (model.EmbeddingResponse, error) {
	//TODO implement me
	panic("implement me")
}
