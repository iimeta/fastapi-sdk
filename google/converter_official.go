package google

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

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

func (g *Google) ConvChatCompletionsResponseOfficial(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	chatCompletionRes := model.GoogleChatCompletionRes{}
	if err = json.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if chatCompletionRes.Error.Code != 0 || (chatCompletionRes.Candidates[0].FinishReason != "STOP" && chatCompletionRes.Candidates[0].FinishReason != "MAX_TOKENS") {
		logger.Errorf(ctx, "ConvChatCompletionsResponseOfficial Google model: %s, chatCompletionRes: %s", g.Model, gjson.MustEncodeString(chatCompletionRes))

		err = g.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsResponseOfficial Google model: %s, error: %v", g.Model, err)

		return response, err
	}

	response = model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + grand.S(29),
		Object:  consts.COMPLETION_OBJECT,
		Created: gtime.Timestamp(),
		Model:   g.Model,
		Usage: &model.Usage{
			PromptTokens:     chatCompletionRes.UsageMetadata.PromptTokenCount,
			CompletionTokens: chatCompletionRes.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      chatCompletionRes.UsageMetadata.TotalTokenCount,
		},
		ResponseBytes: data,
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

func (g *Google) ConvChatCompletionsStreamResponseOfficial(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	chatCompletionRes := model.GoogleChatCompletionRes{}
	if err = json.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if chatCompletionRes.Error.Code != 0 {
		logger.Errorf(ctx, "ChatCompletionsStream Google model: %s, chatCompletionRes: %s", g.Model, gjson.MustEncodeString(chatCompletionRes))

		err = g.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ChatCompletionsStream Google model: %s, error: %v", g.Model, err)

		return response, err
	}

	response = model.ChatCompletionResponse{
		Id:            consts.COMPLETION_ID_PREFIX + grand.S(29),
		Object:        consts.COMPLETION_STREAM_OBJECT,
		Created:       gtime.Timestamp(),
		Model:         g.Model,
		ResponseBytes: data,
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

	if chatCompletionRes.UsageMetadata != nil {
		response.Usage = &model.Usage{
			PromptTokens:     chatCompletionRes.UsageMetadata.PromptTokenCount,
			CompletionTokens: chatCompletionRes.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      chatCompletionRes.UsageMetadata.TotalTokenCount,
		}
	}

	return response, nil
}
