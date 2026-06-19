package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strconv"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/v2/common"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (o *OpenAI) ConvChatCompletionsRequest(ctx context.Context, data any) (request model.ChatCompletionRequest, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvChatCompletionsRequest time: %d", gtime.TimestampMilli()-now)
	}()

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

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvChatCompletionsResponse time: %d", gtime.TimestampMilli()-now)
	}()

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

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvChatCompletionsStreamResponse time: %d", gtime.TimestampMilli()-now)
	}()

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

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvImageGenerationsRequest time: %d", gtime.TimestampMilli()-now)
	}()

	if err = json.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	return request, nil
}

func (o *OpenAI) ConvImageGenerationsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvImageGenerationsResponse time: %d", gtime.TimestampMilli()-now)
	}()

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvImageGenerationsStreamResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvImageGenerationsStreamResponse time: %d", gtime.TimestampMilli()-now)
	}()

	response.ResponseBytes = data

	streamResponse := model.ImageStreamResponse{}
	if err = json.Unmarshal(data, &streamResponse); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	response.Created = streamResponse.CreatedAt
	response.Usage = streamResponse.Usage

	if streamResponse.B64Json != "" {
		response.Data = []model.ImageResponseData{
			{
				B64Json: streamResponse.B64Json,
			},
		}
	}

	return response, nil
}

func (o *OpenAI) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (data *bytes.Buffer, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvImageEditsRequest time: %d", gtime.TimestampMilli()-now)
	}()

	if len(request.Images) > 0 {
		return o.convImageEditsRequestJSON(ctx, request)
	}

	switch v := request.Image.(type) {
	case string:
		request.Images = []model.ImageEditImage{{ImageUrl: v}}
		return o.convImageEditsRequestJSON(ctx, request)
	case []string:
		for _, u := range v {
			request.Images = append(request.Images, model.ImageEditImage{ImageUrl: u})
		}
		return o.convImageEditsRequestJSON(ctx, request)
	case []any:
		for _, item := range v {
			if u, ok := item.(string); ok {
				request.Images = append(request.Images, model.ImageEditImage{ImageUrl: u})
			}
		}
		if len(request.Images) > 0 {
			return o.convImageEditsRequestJSON(ctx, request)
		}
	}

	return o.convImageEditsRequestForm(ctx, request)
}

func (o *OpenAI) convImageEditsRequestJSON(ctx context.Context, request model.ImageEditRequest) (data *bytes.Buffer, err error) {

	jsonReq := map[string]any{
		"prompt": request.Prompt,
		"model":  request.Model,
	}

	jsonReq["images"] = request.Images

	if request.Background != "" {
		jsonReq["background"] = request.Background
	}

	if request.InputFidelity != "" {
		jsonReq["input_fidelity"] = request.InputFidelity
	}

	if request.N != 0 {
		jsonReq["n"] = request.N
	}

	if request.OutputCompression != 0 {
		jsonReq["output_compression"] = request.OutputCompression
	}

	if request.OutputFormat != "" {
		jsonReq["output_format"] = request.OutputFormat
	}

	if request.PartialImages != 0 {
		jsonReq["partial_images"] = request.PartialImages
	}

	if request.Quality != "" {
		jsonReq["quality"] = request.Quality
	}

	if request.ResponseFormat != "" {
		jsonReq["response_format"] = request.ResponseFormat
	}

	if request.Size != "" {
		jsonReq["size"] = gstr.Replace(request.Size, ":", "x")
	}

	if request.User != "" {
		jsonReq["user"] = request.User
	}

	if request.AspectRatio != "" {
		jsonReq["aspect_ratio"] = request.AspectRatio
	}

	if request.Stream {
		jsonReq["stream"] = true
	}

	if request.Async {
		jsonReq["async"] = true
	}

	if request.Mask != nil {
		switch v := request.Mask.(type) {
		case string:
			if v != "" {
				jsonReq["mask"] = model.ImageEditImage{ImageUrl: v}
			}
		case model.ImageEditImage:
			jsonReq["mask"] = v
		default:
			jsonReq["mask"] = v
		}
	}

	jsonBytes, err := json.Marshal(jsonReq)
	if err != nil {
		logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, json.Marshal error: %v", o.Model, err)
		return nil, err
	}

	o.header["Content-Type"] = "application/json"

	return bytes.NewBuffer(jsonBytes), nil
}

func (o *OpenAI) convImageEditsRequestForm(ctx context.Context, request model.ImageEditRequest) (data *bytes.Buffer, err error) {

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

	if request.InputFidelity != "" {
		if err = builder.WriteField("input_fidelity", request.InputFidelity); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if fileHeaders, ok := request.Image.([]*multipart.FileHeader); ok && len(fileHeaders) > 0 {
		if len(fileHeaders) == 1 {
			if err = builder.CreateFormFileHeader("image", fileHeaders[0]); err != nil {
				logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
				return data, err
			}
		} else {
			for _, image := range fileHeaders {
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
		if maskFileHeader, ok := request.Mask.(*multipart.FileHeader); ok {
			if err = builder.CreateFormFileHeader("mask", maskFileHeader); err != nil {
				logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
				return data, err
			}
		}
	}

	if request.N != 0 {
		if err = builder.WriteField("n", strconv.Itoa(request.N)); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.OutputCompression != 0 {
		if err = builder.WriteField("output_compression", strconv.Itoa(request.OutputCompression)); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.OutputFormat != "" {
		if err = builder.WriteField("output_format", request.OutputFormat); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	if request.PartialImages != 0 {
		if err = builder.WriteField("partial_images", strconv.Itoa(request.PartialImages)); err != nil {
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
		if err = builder.WriteField("size", gstr.Replace(request.Size, ":", "x")); err != nil {
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

	if request.Stream {
		if err = builder.WriteField("stream", "true"); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	o.header["Content-Type"] = builder.FormDataContentType()

	return data, nil
}

func (o *OpenAI) ConvImageEditsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvImageEditsResponse time: %d", gtime.TimestampMilli()-now)
	}()

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvAudioSpeechRequest(ctx context.Context, data []byte) (request model.SpeechRequest, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvAudioSpeechRequest time: %d", gtime.TimestampMilli()-now)
	}()

	if err = json.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	return request, nil
}

func (o *OpenAI) ConvAudioSpeechResponse(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvAudioSpeechResponse time: %d", gtime.TimestampMilli()-now)
	}()

	return model.SpeechResponse{
		Data: data,
	}, nil
}

func (o *OpenAI) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (data *bytes.Buffer, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvAudioTranscriptionsRequest time: %d", gtime.TimestampMilli()-now)
	}()

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

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvAudioTranscriptionsResponse time: %d", gtime.TimestampMilli()-now)
	}()

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (request model.EmbeddingRequest, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvTextEmbeddingsRequest time: %d", gtime.TimestampMilli()-now)
	}()

	if err = json.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	return request, nil
}

func (o *OpenAI) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvTextEmbeddingsResponse time: %d", gtime.TimestampMilli()-now)
	}()

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvVideoCreateRequest(ctx context.Context, request model.VideoCreateRequest) (data *bytes.Buffer, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvVideoCreateRequest time: %d", gtime.TimestampMilli()-now)
	}()

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
		if err = builder.WriteField("size", gstr.Replace(request.Size, ":", "x")); err != nil {
			logger.Errorf(ctx, "ConvVideoCreateRequest OpenAI model: %s, error: %v", o.Model, err)
			return data, err
		}
	}

	o.header["Content-Type"] = builder.FormDataContentType()

	return data, nil
}

func (o *OpenAI) ConvVideoListResponse(ctx context.Context, data []byte) (response model.VideoListResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvVideoListResponse time: %d", gtime.TimestampMilli()-now)
	}()

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvVideoContentResponse(ctx context.Context, data []byte) (response model.VideoContentResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvVideoContentResponse time: %d", gtime.TimestampMilli()-now)
	}()

	return model.VideoContentResponse{
		Data: data,
	}, nil
}

func (o *OpenAI) ConvVideoJobResponse(ctx context.Context, data []byte) (response model.VideoJobResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvVideoJobResponse time: %d", gtime.TimestampMilli()-now)
	}()

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvFileUploadRequest(ctx context.Context, request model.FileUploadRequest) (data *bytes.Buffer, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvFileUploadRequest time: %d", gtime.TimestampMilli()-now)
	}()

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

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvFileListResponse time: %d", gtime.TimestampMilli()-now)
	}()

	response.ResponseBytes = data

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvFileContentResponse(ctx context.Context, data []byte) (response model.FileContentResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvFileContentResponse time: %d", gtime.TimestampMilli()-now)
	}()

	return model.FileContentResponse{
		Data: data,
	}, nil
}

func (o *OpenAI) ConvFileResponse(ctx context.Context, data []byte) (response model.FileResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvFileResponse time: %d", gtime.TimestampMilli()-now)
	}()

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

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvBatchListResponse time: %d", gtime.TimestampMilli()-now)
	}()

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ConvBatchResponse(ctx context.Context, data []byte) (response model.BatchResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvBatchResponse time: %d", gtime.TimestampMilli()-now)
	}()

	response.ResponseBytes = data

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}
