package openai

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func ConvResponsesToChatCompletionsRequest(request *ghttp.Request, isChatCompletions bool) model.ChatCompletionRequest {

	if isChatCompletions {
		chatCompletionRequest := model.ChatCompletionRequest{}
		if err := gjson.Unmarshal(request.GetBody(), &chatCompletionRequest); err != nil {
			logger.Error(request.GetCtx(), err)
			return model.ChatCompletionRequest{}
		}
		return chatCompletionRequest
	}

	responsesReq := model.OpenAIResponsesReq{}
	if err := gjson.Unmarshal(request.GetBody(), &responsesReq); err != nil {
		logger.Error(request.GetCtx(), err)
		return model.ChatCompletionRequest{}
	}

	chatCompletionRequest := model.ChatCompletionRequest{
		Model:               responsesReq.Model,
		MaxCompletionTokens: responsesReq.MaxOutputTokens,
		Temperature:         responsesReq.Temperature,
		TopP:                responsesReq.TopP,
		Stream:              responsesReq.Stream,
		User:                responsesReq.User,
		Tools:               responsesReq.Tools,
		ToolChoice:          responsesReq.ToolChoice,
		ParallelToolCalls:   responsesReq.ParallelToolCalls,
		Store:               responsesReq.Store,
		Metadata:            responsesReq.Metadata,
	}

	if responsesReq.Input != nil {
		if value, ok := responsesReq.Input.([]interface{}); ok {

			inputs := make([]model.OpenAIResponsesInput, 0)
			if err := gjson.Unmarshal(gjson.MustEncode(value), &inputs); err != nil {
				logger.Error(request.GetCtx(), err)
				return chatCompletionRequest
			}

			for _, input := range inputs {
				chatCompletionRequest.Messages = append(chatCompletionRequest.Messages, model.ChatCompletionMessage{
					Role:    input.Role,
					Content: input.Content,
				})
			}

		} else {
			chatCompletionRequest.Messages = []model.ChatCompletionMessage{{
				Role:    "user",
				Content: responsesReq.Input,
			}}
		}
	}

	if responsesReq.Reasoning != nil {
		chatCompletionRequest.ReasoningEffort = responsesReq.Reasoning.Effort
	}

	return chatCompletionRequest
}

func ConvResponsesToChatCompletionsResponse(ctx context.Context, res model.OpenAIResponsesRes) model.ChatCompletionResponse {

	responsesRes := model.OpenAIResponsesRes{
		Model:         res.Model,
		Usage:         res.Usage,
		ResponseBytes: res.ResponseBytes,
		ConnTime:      res.ConnTime,
		Duration:      res.Duration,
		TotalTime:     res.TotalTime,
		Err:           res.Err,
	}

	if res.ResponseBytes != nil {
		if err := gjson.Unmarshal(res.ResponseBytes, &responsesRes); err != nil {
			logger.Error(ctx, err)
		}
	}

	chatCompletionResponse := model.ChatCompletionResponse{
		ID:            responsesRes.Id,
		Object:        responsesRes.Object,
		Created:       responsesRes.CreatedAt,
		Model:         responsesRes.Model,
		ResponseBytes: responsesRes.ResponseBytes,
		ConnTime:      responsesRes.ConnTime,
		Duration:      responsesRes.Duration,
		TotalTime:     responsesRes.TotalTime,
		Error:         responsesRes.Err,
	}

	for _, output := range responsesRes.Output {
		if len(output.Content) > 0 {
			chatCompletionResponse.Choices = append(chatCompletionResponse.Choices, model.ChatCompletionChoice{
				Message: &model.ChatCompletionMessage{
					Role:    output.Role,
					Content: output.Content[0].Text,
				},
				FinishReason: "stop",
			})
		}
	}

	if responsesRes.Tools != nil && len(gconv.Interfaces(responsesRes.Tools)) > 0 {
		chatCompletionResponse.Choices = append(chatCompletionResponse.Choices, model.ChatCompletionChoice{
			Message: &model.ChatCompletionMessage{
				Role:      "assistant",
				ToolCalls: responsesRes.Tools,
			},
			FinishReason: "tool_calls",
		})
	}

	if responsesRes.Usage != nil {
		chatCompletionResponse.Usage = &model.Usage{
			PromptTokens:     responsesRes.Usage.InputTokens,
			CompletionTokens: responsesRes.Usage.OutputTokens,
			TotalTokens:      responsesRes.Usage.TotalTokens,
			PromptTokensDetails: model.PromptTokensDetails{
				CachedTokens: responsesRes.Usage.InputTokensDetails.CachedTokens,
				TextTokens:   responsesRes.Usage.InputTokensDetails.TextTokens,
			},
			CompletionTokensDetails: model.CompletionTokensDetails{
				ReasoningTokens: responsesRes.Usage.OutputTokensDetails.ReasoningTokens,
			},
		}
	}

	return chatCompletionResponse
}

func ConvResponsesStreamToChatCompletionsResponse(ctx context.Context, res model.OpenAIResponsesStreamRes) model.ChatCompletionResponse {

	responsesStreamRes := model.OpenAIResponsesStreamRes{
		ResponseBytes: res.ResponseBytes,
		ConnTime:      res.ConnTime,
		Duration:      res.Duration,
		TotalTime:     res.TotalTime,
		Err:           res.Err,
	}

	if res.ResponseBytes != nil {
		if err := gjson.Unmarshal(res.ResponseBytes, &responsesStreamRes); err != nil {
			logger.Error(ctx, err)
		}
	}

	chatCompletionResponse := model.ChatCompletionResponse{
		ID:            responsesStreamRes.Response.Id,
		Object:        responsesStreamRes.Response.Object,
		Created:       responsesStreamRes.Response.CreatedAt,
		Model:         responsesStreamRes.Response.Model,
		ResponseBytes: responsesStreamRes.ResponseBytes,
		ConnTime:      responsesStreamRes.ConnTime,
		Duration:      responsesStreamRes.Duration,
		TotalTime:     responsesStreamRes.TotalTime,
		Error:         responsesStreamRes.Err,
	}

	if chatCompletionResponse.ID == "" {
		chatCompletionResponse.ID = responsesStreamRes.Item.Id
	}

	if chatCompletionResponse.ID == "" {
		chatCompletionResponse.ID = responsesStreamRes.ItemId
	}

	chatCompletionChoice := model.ChatCompletionChoice{
		Delta: &model.ChatCompletionStreamChoiceDelta{
			Content: responsesStreamRes.Delta,
		},
	}

	if "response.completed" == responsesStreamRes.Type {
		chatCompletionChoice.FinishReason = "stop"
	}

	chatCompletionResponse.Choices = append(chatCompletionResponse.Choices, chatCompletionChoice)

	if responsesStreamRes.Response.Usage != nil {
		chatCompletionResponse.Usage = &model.Usage{
			PromptTokens:     responsesStreamRes.Response.Usage.InputTokens,
			CompletionTokens: responsesStreamRes.Response.Usage.OutputTokens,
			TotalTokens:      responsesStreamRes.Response.Usage.TotalTokens,
			PromptTokensDetails: model.PromptTokensDetails{
				CachedTokens: responsesStreamRes.Response.Usage.InputTokensDetails.CachedTokens,
				TextTokens:   responsesStreamRes.Response.Usage.InputTokensDetails.TextTokens,
			},
			CompletionTokensDetails: model.CompletionTokensDetails{
				ReasoningTokens: responsesStreamRes.Response.Usage.OutputTokensDetails.ReasoningTokens,
			},
		}
	}

	return chatCompletionResponse
}

func ConvChatCompletionsToResponsesRequest(request *ghttp.Request) model.OpenAIResponsesReq {

	chatCompletionRequest := model.ChatCompletionRequest{}
	if err := gjson.Unmarshal(request.GetBody(), &chatCompletionRequest); err != nil {
		logger.Error(request.GetCtx(), err)
		return model.OpenAIResponsesReq{}
	}

	responsesReq := model.OpenAIResponsesReq{
		Model:             chatCompletionRequest.Model,
		Stream:            chatCompletionRequest.Stream,
		MaxOutputTokens:   chatCompletionRequest.MaxTokens,
		Metadata:          chatCompletionRequest.Metadata,
		ParallelToolCalls: chatCompletionRequest.ParallelToolCalls != nil,
		Store:             chatCompletionRequest.Store,
		Temperature:       chatCompletionRequest.Temperature,
		Tools:             chatCompletionRequest.Tools,
		ToolChoice:        chatCompletionRequest.ToolChoice,
		TopP:              chatCompletionRequest.TopP,
		User:              chatCompletionRequest.User,
	}

	input := make([]model.OpenAIResponsesInput, 0)

	for _, message := range chatCompletionRequest.Messages {

		responsesContent := make([]model.OpenAIResponsesContent, 0)

		if multiContent, ok := message.Content.([]interface{}); ok {
			for _, value := range multiContent {
				if content, ok := value.(map[string]interface{}); ok {

					if content["type"] == "text" {
						responsesContent = append(responsesContent, model.OpenAIResponsesContent{
							Type: "input_text",
							Text: gconv.String(content["text"]),
						})
					} else if content["type"] == "image_url" {

						imageContent := model.OpenAIResponsesContent{
							Type: "input_image",
						}

						if imageUrl, ok := content["image_url"].(map[string]interface{}); ok {
							imageContent.ImageUrl = gconv.String(imageUrl["url"])
						}

						responsesContent = append(responsesContent, imageContent)
					}
				}
			}
		} else {
			responsesContent = append(responsesContent, model.OpenAIResponsesContent{
				Type: "input_text",
				Text: gconv.String(message.Content),
			})
		}

		input = append(input, model.OpenAIResponsesInput{
			Role:    message.Role,
			Content: responsesContent,
		})
	}

	responsesReq.Input = input

	if chatCompletionRequest.ReasoningEffort != "" {
		responsesReq.Reasoning = &model.OpenAIResponsesReasoning{
			Effort: chatCompletionRequest.ReasoningEffort,
		}
	}

	return responsesReq
}

func (o *OpenAI) ConvChatCompletionsRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {

	chatCompletionRequest := model.ChatCompletionRequest{}
	if err := gjson.Unmarshal(data, &chatCompletionRequest); err != nil {
		logger.Error(ctx, err)
		return chatCompletionRequest, err
	}

	for _, message := range chatCompletionRequest.Messages {
		if message.Role == consts.ROLE_SYSTEM && (gstr.HasPrefix(chatCompletionRequest.Model, "o1") || gstr.HasPrefix(chatCompletionRequest.Model, "o3")) {
			message.Role = consts.ROLE_DEVELOPER
		}
	}

	if chatCompletionRequest.Stream {
		// 默认让流式返回usage
		if chatCompletionRequest.StreamOptions == nil { // request.Tools == nil &&
			chatCompletionRequest.StreamOptions = &model.StreamOptions{
				IncludeUsage: true,
			}
		}
	}

	if gstr.HasPrefix(chatCompletionRequest.Model, "o") || gstr.HasPrefix(chatCompletionRequest.Model, "gpt-5") {
		if chatCompletionRequest.MaxCompletionTokens == 0 && chatCompletionRequest.MaxTokens != 0 {
			chatCompletionRequest.MaxCompletionTokens = chatCompletionRequest.MaxTokens
		}
		chatCompletionRequest.MaxTokens = 0
	}

	return chatCompletionRequest, nil
}

func (o *OpenAI) ConvChatCompletionsResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionResponse := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &chatCompletionResponse); err != nil {
		logger.Error(ctx, err)
		return chatCompletionResponse, err
	}

	for _, choice := range chatCompletionResponse.Choices {
		if choice.Message.Annotations == nil {
			choice.Message.Annotations = []any{}
		}
	}

	return chatCompletionResponse, nil
}

func (o *OpenAI) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionResponse := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &chatCompletionResponse); err != nil {
		logger.Error(ctx, err)
		return chatCompletionResponse, err
	}

	return chatCompletionResponse, nil
}

func (o *OpenAI) ConvChatResponsesRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvChatResponsesResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvImageGenerationsRequest(ctx context.Context, data []byte) (model.ImageGenerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvImageGenerationsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvImageEditsRequest(ctx context.Context, data []byte) (model.ImageEditRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvImageEditsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvAudioSpeechRequest(ctx context.Context, data []byte) (model.SpeechRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvAudioSpeechResponse(ctx context.Context, data []byte) (model.SpeechResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvAudioTranscriptionsRequest(ctx context.Context, data []byte) (model.AudioRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (model.AudioResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (model.EmbeddingRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (model.EmbeddingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvTextModerationsRequest(ctx context.Context, data []byte) (model.ModerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvTextModerationsResponse(ctx context.Context, data []byte) (model.ModerationResponse, error) {
	//TODO implement me
	panic("implement me")
}
