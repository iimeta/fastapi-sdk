package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (o *OpenAI) ConvChatCompletionsRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {

	request := model.ChatCompletionRequest{}
	if err := json.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	//for _, message := range request.Messages {
	//	if message.Role == consts.ROLE_SYSTEM && (gstr.HasPrefix(request.Model, "o1") || gstr.HasPrefix(request.Model, "o3")) {
	//		message.Role = consts.ROLE_DEVELOPER
	//	}
	//}

	if request.Stream {
		// 默认让流式返回usage
		if request.StreamOptions == nil { // request.Tools == nil &&
			request.StreamOptions = &model.StreamOptions{
				IncludeUsage: true,
			}
		}
	}

	if gstr.HasPrefix(request.Model, "o") || gstr.HasPrefix(request.Model, "gpt-5") {
		if request.MaxCompletionTokens == 0 && request.MaxTokens != 0 {
			request.MaxCompletionTokens = request.MaxTokens
		}
		request.MaxTokens = 0
	}

	if o.IsSupportSystemRole != nil {
		request.Messages = common.HandleMessages(request.Messages, *o.IsSupportSystemRole)
	}

	return request, nil
}

func (o *OpenAI) ConvChatCompletionsResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	response = model.ChatCompletionResponse{
		ResponseBytes: data,
	}

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	for _, choice := range response.Choices {
		if choice.Message.Annotations == nil {
			choice.Message.Annotations = []any{}
		}
	}

	return response, nil
}

func (o *OpenAI) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	response = model.ChatCompletionResponse{
		//ResponseBytes: data,
	}

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
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

func (o *OpenAI) ConvChatResponsesRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvChatResponsesResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvImageGenerationsRequest(ctx context.Context, data []byte) (model.ImageGenerationRequest, error) {

	request := model.ImageGenerationRequest{}
	if err := json.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	return request, nil
}

func (o *OpenAI) ConvImageGenerationsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {

	response := model.ImageResponse{}
	if err := json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (*bytes.Buffer, error) {

	data := &bytes.Buffer{}
	builder := util.NewFormBuilder(data)

	if err := builder.WriteField("model", request.Model); err != nil {
		logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
		return data, err
	}

	if len(request.Image) > 0 {
		if len(request.Image) == 1 {
			if err := builder.CreateFormFileHeader("image", request.Image[0]); err != nil {
				logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
				return data, err
			}
		} else {
			for _, image := range request.Image {
				if err := builder.CreateFormFileHeader("image[]", image); err != nil {
					logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
					return data, err
				}
			}
		}
	}

	if err := builder.WriteField("prompt", request.Prompt); err != nil {
		logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
		return data, err
	}

	if request.Background != "" {
		if err := builder.WriteField("background", request.Background); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.Mask != nil {
		if err := builder.CreateFormFileHeader("mask", request.Mask); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.N != 0 {
		if err := builder.WriteField("n", strconv.Itoa(request.N)); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.Quality != "" {
		if err := builder.WriteField("quality", request.Quality); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.ResponseFormat != "" {
		if err := builder.WriteField("response_format", request.ResponseFormat); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.Size != "" {
		if err := builder.WriteField("size", request.Size); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.User != "" {
		if err := builder.WriteField("user", request.User); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if err := builder.Close(); err != nil {
		logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
		return data, err
	}

	o.header["Content-Type"] = builder.FormDataContentType()

	return data, nil
}

func (o *OpenAI) ConvImageEditsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {

	response := model.ImageResponse{}
	if err := json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvAudioSpeechRequest(ctx context.Context, data []byte) (model.SpeechRequest, error) {

	request := model.SpeechRequest{}
	if err := json.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	return request, nil
}

func (o *OpenAI) ConvAudioSpeechResponse(ctx context.Context, data []byte) (model.SpeechResponse, error) {
	return model.SpeechResponse{
		Data: data,
	}, nil
}

func (o *OpenAI) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (*bytes.Buffer, error) {

	data := &bytes.Buffer{}
	builder := util.NewFormBuilder(data)

	if err := builder.WriteField("model", request.Model); err != nil {
		logger.Errorf(ctx, "ConvAudioTranscriptionsRequest OpenAI model: %s, error: %v", o.Model, err)
		return data, err
	}

	if request.File != nil {
		if err := builder.CreateFormFileHeader("file", request.File); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.Prompt != "" {
		if err := builder.WriteField("prompt", request.Prompt); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.Format != "" {
		if err := builder.WriteField("response_format", request.Format); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.Temperature != 0 {
		if err := builder.WriteField("temperature", fmt.Sprintf("%.2f", request.Temperature)); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.Language != "" {
		if err := builder.WriteField("language", request.Language); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if len(request.TimestampGranularities) > 0 {
		for _, timestampGranularitie := range request.TimestampGranularities {
			if err := builder.WriteField("timestamp_granularities[]", timestampGranularitie); err != nil {
				logger.Errorf(ctx, "ConvAudioTranscriptionsRequest OpenAI model: %s, error: %v", o.Model, err)
				return data, err
			}
		}
	}

	if err := builder.Close(); err != nil {
		logger.Errorf(ctx, "ConvAudioTranscriptionsRequest OpenAI model: %s, error: %v", o.Model, err)
		return data, err
	}

	o.header["Content-Type"] = builder.FormDataContentType()

	return data, nil
}

func (o *OpenAI) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (model.AudioResponse, error) {

	response := model.AudioResponse{}
	if err := json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (model.EmbeddingRequest, error) {

	request := model.EmbeddingRequest{}
	if err := json.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	return request, nil
}

func (o *OpenAI) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (model.EmbeddingResponse, error) {

	response := model.EmbeddingResponse{}
	if err := json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}
