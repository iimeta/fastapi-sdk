package google

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/iimeta/fastapi-sdk/v2/common"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
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

func (g *Google) ConvImageGenerationsRequest(ctx context.Context, data []byte) (request model.ImageGenerationRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvImageGenerationsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvImageEditsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvAudioSpeechRequest(ctx context.Context, data []byte) (request model.SpeechRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvAudioSpeechResponse(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (response model.AudioResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (request model.EmbeddingRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvVideoCreateRequest(ctx context.Context, request model.VideoCreateRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvVideoListResponse(ctx context.Context, data []byte) (response model.VideoListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvVideoContentResponse(ctx context.Context, data []byte) (response model.VideoContentResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvVideoJobResponse(ctx context.Context, data []byte) (response model.VideoJobResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvFileUploadRequest(ctx context.Context, request model.FileUploadRequest) (data *bytes.Buffer, err error) {

	data = &bytes.Buffer{}
	builder := util.NewFormBuilder(data)

	defer func() {
		if err := builder.Close(); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest Google model: %s, builder.Close() error: %v", g.Model, err)
		}
	}()

	if request.File != nil {
		if err = builder.CreateFormFileHeader("file", request.File); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest Google model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	g.header["Content-Type"] = builder.FormDataContentType()

	return data, nil
}

func (g *Google) ConvFileListResponse(ctx context.Context, data []byte) (response model.FileListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvFileContentResponse(ctx context.Context, data []byte) (response model.FileContentResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvFileResponse(ctx context.Context, data []byte) (response model.FileResponse, err error) {

	fileRes := model.GoogleFileResponse{}
	if err = json.Unmarshal(data, &fileRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	response = model.FileResponse{
		Id:            strings.TrimPrefix(fileRes.File.Name, "files/"),
		Object:        "file",
		Purpose:       "upload",
		Bytes:         gconv.Int(fileRes.File.SizeBytes),
		CreatedAt:     fileRes.File.CreateTime.Unix(),
		ExpiresAt:     fileRes.File.ExpirationTime.Unix(),
		FileUrl:       fileRes.File.Uri,
		ResponseBytes: data,
	}

	if fileRes.File.State == "PROCESSING" {
		response.Status = "processing"
	} else if fileRes.File.State == "ACTIVE" {
		response.Status = "processed"
	} else {
		response.Status = strings.ToLower(fileRes.File.State)
	}

	return response, nil
}

func (g *Google) ConvBatchCreateRequest(ctx context.Context, request model.BatchCreateRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvBatchListResponse(ctx context.Context, data []byte) (response model.BatchListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvBatchResponse(ctx context.Context, data []byte) (response model.BatchResponse, err error) {
	//TODO implement me
	panic("implement me")
}
