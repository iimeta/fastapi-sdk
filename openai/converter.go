package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/v2/common"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (o *OpenAI) ConvChatCompletionsRequest(ctx context.Context, data any) (request model.ChatCompletionRequest, err error) {

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

	for _, message := range request.Messages {
		if message.Role == consts.ROLE_SYSTEM && (gstr.HasPrefix(request.Model, "o1") || gstr.HasPrefix(request.Model, "o3")) {
			message.Role = consts.ROLE_DEVELOPER
		}
	}

	if request.Stream {
		// 默认让流式返回usage
		if request.StreamOptions == nil {
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

	response.ResponseBytes = data

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

	response.ResponseBytes = data

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

func (o *OpenAI) ConvChatResponsesRequest(ctx context.Context, data []byte) (request model.ChatCompletionRequest, err error) {
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

func (o *OpenAI) ConvImageGenerationsRequest(ctx context.Context, data []byte) (request model.ImageGenerationRequest, err error) {

	if err = json.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	return request, nil
}

func (o *OpenAI) ConvImageGenerationsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (data *bytes.Buffer, err error) {

	data = &bytes.Buffer{}
	builder := util.NewFormBuilder(data)

	defer func() {
		if err := builder.Close(); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, builder.Close() error: %v", o.Model, err)
		}
	}()

	if err = builder.WriteField("model", request.Model); err != nil {
		logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
		return data, err
	}

	if len(request.Image) > 0 {
		if len(request.Image) == 1 {
			if err = builder.CreateFormFileHeader("image", request.Image[0]); err != nil {
				logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
				return data, err
			}
		} else {
			for _, image := range request.Image {
				if err = builder.CreateFormFileHeader("image[]", image); err != nil {
					logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
					return data, err
				}
			}
		}
	}

	if err = builder.WriteField("prompt", request.Prompt); err != nil {
		logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
		return data, err
	}

	if request.Background != "" {
		if err = builder.WriteField("background", request.Background); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.Mask != nil {
		if err = builder.CreateFormFileHeader("mask", request.Mask); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.N != 0 {
		if err = builder.WriteField("n", strconv.Itoa(request.N)); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.Quality != "" {
		if err = builder.WriteField("quality", request.Quality); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.ResponseFormat != "" {
		if err = builder.WriteField("response_format", request.ResponseFormat); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.Size != "" {
		if err = builder.WriteField("size", request.Size); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.User != "" {
		if err = builder.WriteField("user", request.User); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	o.header["Content-Type"] = builder.FormDataContentType()

	return data, nil
}

func (o *OpenAI) ConvImageEditsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvAudioSpeechRequest(ctx context.Context, data []byte) (request model.SpeechRequest, err error) {

	if err = json.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	return request, nil
}

func (o *OpenAI) ConvAudioSpeechResponse(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {
	return model.SpeechResponse{
		Data: data,
	}, nil
}

func (o *OpenAI) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (data *bytes.Buffer, err error) {

	data = &bytes.Buffer{}
	builder := util.NewFormBuilder(data)

	defer func() {
		if err := builder.Close(); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest OpenAI model: %s, builder.Close() error: %v", o.Model, err)
		}
	}()

	if err = builder.WriteField("model", request.Model); err != nil {
		logger.Errorf(ctx, "ConvAudioTranscriptionsRequest OpenAI model: %s, error: %v", o.Model, err)
		return data, err
	}

	if request.File != nil {
		if err = builder.CreateFormFileHeader("file", request.File); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.Prompt != "" {
		if err = builder.WriteField("prompt", request.Prompt); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.ResponseFormat != "" {
		if err = builder.WriteField("response_format", request.ResponseFormat); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.Temperature != 0 {
		if err = builder.WriteField("temperature", fmt.Sprintf("%.2f", request.Temperature)); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.Language != "" {
		if err = builder.WriteField("language", request.Language); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if len(request.TimestampGranularities) > 0 {
		for _, timestampGranularitie := range request.TimestampGranularities {
			if err = builder.WriteField("timestamp_granularities[]", timestampGranularitie); err != nil {
				logger.Errorf(ctx, "ConvAudioTranscriptionsRequest OpenAI model: %s, error: %v", o.Model, err)
				return data, err
			}
		}
	}

	o.header["Content-Type"] = builder.FormDataContentType()

	return data, nil
}

func (o *OpenAI) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (response model.AudioResponse, err error) {

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (request model.EmbeddingRequest, err error) {

	if err = json.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	return request, nil
}

func (o *OpenAI) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error) {

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvVideoCreateRequest(ctx context.Context, request model.VideoCreateRequest) (data *bytes.Buffer, err error) {

	data = &bytes.Buffer{}
	builder := util.NewFormBuilder(data)

	defer func() {
		if err := builder.Close(); err != nil {
			logger.Errorf(ctx, "ConvVideoCreateRequest OpenAI model: %s, builder.Close() error: %v", o.Model, err)
		}
	}()

	if err = builder.WriteField("model", request.Model); err != nil {
		logger.Errorf(ctx, "ConvVideoCreateRequest OpenAI model: %s, error: %v", o.Model, err)
		return data, err
	}

	if request.Prompt != "" {
		if err = builder.WriteField("prompt", request.Prompt); err != nil {
			logger.Errorf(ctx, "ConvVideoCreateRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.InputReference != nil {
		if err = builder.CreateFormFileHeader("input_reference", request.InputReference); err != nil {
			logger.Errorf(ctx, "ConvVideoCreateRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.Seconds != "" {
		if err = builder.WriteField("seconds", request.Seconds); err != nil {
			logger.Errorf(ctx, "ConvVideoCreateRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.Size != "" {
		if err = builder.WriteField("size", request.Size); err != nil {
			logger.Errorf(ctx, "ConvVideoCreateRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	o.header["Content-Type"] = builder.FormDataContentType()

	return data, nil
}

func (o *OpenAI) ConvVideoListResponse(ctx context.Context, data []byte) (response model.VideoListResponse, err error) {

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvVideoContentResponse(ctx context.Context, data []byte) (response model.VideoContentResponse, err error) {
	return model.VideoContentResponse{
		Data: data,
	}, nil
}

func (o *OpenAI) ConvVideoJobResponse(ctx context.Context, data []byte) (response model.VideoJobResponse, err error) {

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvFileUploadRequest(ctx context.Context, request model.FileUploadRequest) (data *bytes.Buffer, err error) {

	data = &bytes.Buffer{}
	builder := util.NewFormBuilder(data)

	defer func() {
		if err := builder.Close(); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest OpenAI model: %s, builder.Close() error: %v", o.Model, err)
		}
	}()

	if request.File != nil {
		if err = builder.CreateFormFileHeader("file", request.File); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.Purpose != "" {
		if err = builder.WriteField("purpose", request.Purpose); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.ExpiresAfter.Anchor != "" {
		if err = builder.WriteField("expires_after[anchor]", request.ExpiresAfter.Anchor); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.ExpiresAfter.Seconds != "" {
		if err = builder.WriteField("expires_after[seconds]", request.ExpiresAfter.Seconds); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	o.header["Content-Type"] = builder.FormDataContentType()

	return data, nil
}

func (o *OpenAI) ConvFileListResponse(ctx context.Context, data []byte) (response model.FileListResponse, err error) {

	response.ResponseBytes = data

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvFileContentResponse(ctx context.Context, data []byte) (response model.FileContentResponse, err error) {
	return model.FileContentResponse{
		Data: data,
	}, nil
}

func (o *OpenAI) ConvFileResponse(ctx context.Context, data []byte) (response model.FileResponse, err error) {

	response.ResponseBytes = data

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvBatchCreateRequest(ctx context.Context, request model.BatchCreateRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAI) ConvBatchListResponse(ctx context.Context, data []byte) (response model.BatchListResponse, err error) {

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvBatchResponse(ctx context.Context, data []byte) (response model.BatchResponse, err error) {

	response.ResponseBytes = data

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}
