package google

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func (g *Google) ConvChatCompletionsRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {

	chatCompletionRequest := model.ChatCompletionRequest{}
	if err := gjson.Unmarshal(data, &chatCompletionRequest); err != nil {
		logger.Error(ctx, err)
		return chatCompletionRequest, err
	}

	if g.isSupportSystemRole != nil {
		chatCompletionRequest.Messages = common.HandleMessages(chatCompletionRequest.Messages, *g.isSupportSystemRole)
	} else {
		chatCompletionRequest.Messages = common.HandleMessages(chatCompletionRequest.Messages, false)
	}

	return chatCompletionRequest, nil
}

func (g *Google) ConvChatCompletionsResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionResponse := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &chatCompletionResponse); err != nil {
		logger.Error(ctx, err)
		return chatCompletionResponse, err
	}

	return chatCompletionResponse, nil
}

func (g *Google) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionResponse := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &chatCompletionResponse); err != nil {
		logger.Error(ctx, err)
		return chatCompletionResponse, err
	}

	return chatCompletionResponse, nil
}

func (g *Google) ConvChatCompletionsRequestOfficial(ctx context.Context, data []byte) ([]byte, error) {

	request, err := g.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	contents := make([]model.Content, 0)
	for _, message := range request.Messages {

		role := message.Role

		if role == consts.ROLE_ASSISTANT {
			role = consts.ROLE_MODEL
		}

		parts := make([]model.Part, 0)

		if contents, ok := message.Content.([]interface{}); ok {

			for _, value := range contents {

				if content, ok := value.(map[string]interface{}); ok {

					if content["type"] == "image_url" {

						if imageUrl, ok := content["image_url"].(map[string]interface{}); ok {

							mimeType, data := common.GetMime(gconv.String(imageUrl["url"]))

							parts = append(parts, model.Part{
								InlineData: &model.InlineData{
									MimeType: mimeType,
									Data:     data,
								},
							})
						}

					} else if content["type"] == "video_url" {
						if videoUrl, ok := content["video_url"].(map[string]interface{}); ok {

							url := gconv.String(videoUrl["url"])
							format := gconv.String(videoUrl["format"])

							parts = append(parts, model.Part{
								FileData: &model.FileData{
									MimeType: "video/" + format,
									FileUri:  url,
								},
							})
						}
					} else {
						parts = append(parts, model.Part{
							Text: gconv.String(content["text"]),
						})
					}
				}
			}

		} else {
			parts = append(parts, model.Part{
				Text: gconv.String(message.Content),
			})
		}

		contents = append(contents, model.Content{
			Role:  role,
			Parts: parts,
		})
	}

	chatCompletionReq := model.GoogleChatCompletionReq{
		Contents: contents,
		GenerationConfig: model.GenerationConfig{
			MaxOutputTokens: request.MaxTokens,
			Temperature:     request.Temperature,
			TopP:            request.TopP,
		},
		Tools: request.Tools,
	}

	return gjson.MustEncode(chatCompletionReq), nil
}

func (g *Google) ConvChatCompletionsResponseOfficial(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	chatCompletionRes := model.GoogleChatCompletionRes{}
	if err := gjson.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return model.ChatCompletionResponse{}, err
	}

	if chatCompletionRes.Error.Code != 0 || (chatCompletionRes.Candidates[0].FinishReason != "STOP" && chatCompletionRes.Candidates[0].FinishReason != "MAX_TOKENS") {
		logger.Errorf(ctx, "ConvChatCompletionsResponseOfficial Google model: %s, chatCompletionRes: %s", g.model, gjson.MustEncodeString(chatCompletionRes))

		err := g.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsResponseOfficial Google model: %s, error: %v", g.model, err)

		return model.ChatCompletionResponse{}, err
	}

	response := model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + grand.S(29),
		Object:  consts.COMPLETION_OBJECT,
		Created: gtime.Timestamp(),
		Model:   g.model,
		Usage: &model.Usage{
			PromptTokens:     chatCompletionRes.UsageMetadata.PromptTokenCount,
			CompletionTokens: chatCompletionRes.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      chatCompletionRes.UsageMetadata.TotalTokenCount,
		},
	}

	for i, part := range chatCompletionRes.Candidates[0].Content.Parts {
		response.Choices = append(response.Choices, model.ChatCompletionChoice{
			Index: i,
			Message: &model.ChatCompletionMessage{
				Role:    consts.ROLE_ASSISTANT,
				Content: part.Text,
			},
			FinishReason: consts.FinishReasonStop,
		})
	}

	return response, nil
}

func (g *Google) ConvChatCompletionsStreamResponseOfficial(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	var (
		usage   *model.Usage
		created = gtime.Timestamp()
		id      = consts.COMPLETION_ID_PREFIX + grand.S(29)
	)

	chatCompletionRes := new(model.GoogleChatCompletionRes)
	if err := gjson.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return model.ChatCompletionResponse{}, err
	}

	if chatCompletionRes.Error.Code != 0 {
		logger.Errorf(ctx, "ChatCompletionsStream Google model: %s, chatCompletionRes: %s", g.model, gjson.MustEncodeString(chatCompletionRes))

		err := g.apiErrorHandler(chatCompletionRes)
		logger.Errorf(ctx, "ChatCompletionsStream Google model: %s, error: %v", g.model, err)

		return model.ChatCompletionResponse{}, err
	}

	if chatCompletionRes.UsageMetadata != nil {
		usage = &model.Usage{
			PromptTokens:     chatCompletionRes.UsageMetadata.PromptTokenCount,
			CompletionTokens: chatCompletionRes.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      chatCompletionRes.UsageMetadata.TotalTokenCount,
		}
	}

	response := model.ChatCompletionResponse{
		Id:      id,
		Object:  consts.COMPLETION_STREAM_OBJECT,
		Created: created,
		Model:   g.model,
		Usage:   usage,
	}

	for _, candidate := range chatCompletionRes.Candidates {
		response.Choices = append(response.Choices, model.ChatCompletionChoice{
			Index: candidate.Index,
			Delta: &model.ChatCompletionStreamChoiceDelta{
				Role:    consts.ROLE_ASSISTANT,
				Content: candidate.Content.Parts[0].Text,
			},
		})
	}

	return response, nil
}

func (g *Google) ConvChatResponsesRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvChatResponsesResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvImageGenerationsRequest(ctx context.Context, data []byte) (model.ImageGenerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvImageGenerationsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvImageEditsRequest(ctx context.Context, data []byte) (model.ImageEditRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvImageEditsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvAudioSpeechRequest(ctx context.Context, data []byte) (model.SpeechRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvAudioSpeechResponse(ctx context.Context, data []byte) (model.SpeechResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvAudioTranscriptionsRequest(ctx context.Context, data []byte) (model.AudioRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (model.AudioResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (model.EmbeddingRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (model.EmbeddingResponse, error) {
	//TODO implement me
	panic("implement me")
}
