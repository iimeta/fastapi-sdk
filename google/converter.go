package google

import (
	"bytes"
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

func (g *Google) ConvChatCompletionsRequest(ctx context.Context, data any) (request model.ChatCompletionRequest, err error) {

	if v, ok := data.(model.ChatCompletionRequest); ok {
		request = v
	} else if v, ok := data.([]byte); ok {
		if err = json.Unmarshal(v, &request); err != nil {
			logger.Error(ctx, err)
			return request, err
		}
	} else {
		if err = json.Unmarshal(gjson.MustEncode(data), &request); err != nil {
			logger.Error(ctx, err)
			return request, err
		}
	}

	if g.IsSupportSystemRole != nil {
		request.Messages = common.HandleMessages(request.Messages, *g.IsSupportSystemRole)
	} else {
		request.Messages = common.HandleMessages(request.Messages, false)
	}

	return request, nil
}

func (g *Google) ConvChatCompletionsResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	chatCompletionRes := model.GoogleChatCompletionRes{}
	if err = json.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if chatCompletionRes.Error.Code != 0 || (chatCompletionRes.Candidates[0].FinishReason != "STOP" && chatCompletionRes.Candidates[0].FinishReason != "MAX_TOKENS") {
		logger.Errorf(ctx, "ConvChatCompletionsResponse Google model: %s, chatCompletionRes: %s", g.Model, gjson.MustEncodeString(chatCompletionRes))

		err = g.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsResponse Google model: %s, error: %v", g.Model, err)

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
			OutputTokensDetails: model.OutputTokensDetails{
				ReasoningTokens: chatCompletionRes.UsageMetadata.ThoughtsTokenCount,
			},
		},
		ResponseBytes: data,
	}

	for i, part := range chatCompletionRes.Candidates[0].Content.Parts {

		message := &model.ChatCompletionMessage{
			Role:    consts.ROLE_ASSISTANT,
			Content: part.Text,
		}

		if part.FunctionCall != nil {
			if functionCall, ok := part.FunctionCall.(map[string]any); ok {
				message.ToolCalls = []any{
					map[string]any{
						"id":   "call_" + grand.S(24),
						"type": "function",
						"function": map[string]any{
							"name":      functionCall["name"],
							"arguments": gconv.String(functionCall["args"]),
						},
						"extra_content": map[string]any{
							"google": map[string]any{
								"thought_signature": part.ThoughtSignature,
							},
						},
					},
				}
			}
		}

		response.Choices = append(response.Choices, model.ChatCompletionChoice{
			Index:        i,
			Message:      message,
			FinishReason: consts.FinishReasonStop,
		})
	}

	for _, promptTokensDetail := range chatCompletionRes.UsageMetadata.PromptTokensDetails {
		if promptTokensDetail.Modality == "TEXT" {
			response.Usage.PromptTokensDetails.TextTokens = promptTokensDetail.TokenCount
		}
	}

	for _, candidatesTokensDetail := range chatCompletionRes.UsageMetadata.CandidatesTokensDetails {
		if candidatesTokensDetail.Modality == "TEXT" {
			response.Usage.CompletionTokensDetails.TextTokens = candidatesTokensDetail.TokenCount
		}
	}

	return response, nil
}

func (g *Google) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	chatCompletionRes := model.GoogleChatCompletionRes{}
	if err = json.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if chatCompletionRes.Error.Code != 0 {
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponse Google model: %s, chatCompletionRes: %s", g.Model, gjson.MustEncodeString(chatCompletionRes))

		err = g.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponse Google model: %s, error: %v", g.Model, err)

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
			OutputTokensDetails: model.OutputTokensDetails{
				ReasoningTokens: chatCompletionRes.UsageMetadata.ThoughtsTokenCount,
			},
		}

		for _, promptTokensDetail := range chatCompletionRes.UsageMetadata.PromptTokensDetails {
			if promptTokensDetail.Modality == "TEXT" {
				response.Usage.PromptTokensDetails.TextTokens = promptTokensDetail.TokenCount
			}
		}

		for _, candidatesTokensDetail := range chatCompletionRes.UsageMetadata.CandidatesTokensDetails {
			if candidatesTokensDetail.Modality == "TEXT" {
				response.Usage.CompletionTokensDetails.TextTokens = candidatesTokensDetail.TokenCount
			}
		}
	}

	return response, nil
}

func (g *Google) ConvChatResponsesRequest(ctx context.Context, data []byte) (request model.ChatCompletionRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvChatResponsesResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
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

func (g *Google) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (*bytes.Buffer, error) {
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

func (g *Google) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (*bytes.Buffer, error) {
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

func (g *Google) ConvVideoCreateRequest(ctx context.Context, request model.VideoCreateRequest) (*bytes.Buffer, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvVideoListResponse(ctx context.Context, data []byte) (model.VideoListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvVideoContentResponse(ctx context.Context, data []byte) (model.VideoContentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvVideoJobResponse(ctx context.Context, data []byte) (model.VideoJobResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvFileUploadRequest(ctx context.Context, request model.FileUploadRequest) (*bytes.Buffer, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvFileListResponse(ctx context.Context, data []byte) (model.FileListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvFileContentResponse(ctx context.Context, data []byte) (model.FileContentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvFileResponse(ctx context.Context, data []byte) (model.FileResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvBatchCreateRequest(ctx context.Context, request model.BatchCreateRequest) (*bytes.Buffer, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvBatchListResponse(ctx context.Context, data []byte) (model.BatchListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvBatchResponse(ctx context.Context, data []byte) (model.BatchResponse, error) {
	//TODO implement me
	panic("implement me")
}
