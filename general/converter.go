package general

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/iimeta/fastapi-sdk/v2/common"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (g *General) ConvChatCompletionsRequest(ctx context.Context, data any) (request model.ChatCompletionRequest, err error) {

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
		request.Messages = common.HandleMessages(request.Messages, true)
	}

	return request, nil
}

func (g *General) ConvChatCompletionsResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	response.ResponseBytes = data

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (g *General) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

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

func (g *General) ConvChatResponsesRequest(ctx context.Context, data []byte) (request model.ChatCompletionRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvChatResponsesResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvImageGenerationsRequest(ctx context.Context, data []byte) (request model.ImageGenerationRequest, err error) {

	if err = json.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	return request, nil
}

func (g *General) ConvImageGenerationsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (g *General) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (data *bytes.Buffer, err error) {

	data = &bytes.Buffer{}
	builder := util.NewFormBuilder(data)

	defer func() {
		if err := builder.Close(); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest General model: %s, builder.Close() error: %v", g.Model, err)
		}
	}()

	if err = builder.WriteField("model", request.Model); err != nil {
		logger.Errorf(ctx, "ConvImageEditsRequest General model: %s, error: %v", g.Model, err)
		return data, err
	}

	if len(request.Image) > 0 {
		if len(request.Image) == 1 {
			if err = builder.CreateFormFileHeader("image", request.Image[0]); err != nil {
				logger.Errorf(ctx, "ConvImageEditsRequest General model: %s, error: %v", g.Model, err)
				return data, err
			}
		} else {
			for _, image := range request.Image {
				if err = builder.CreateFormFileHeader("image[]", image); err != nil {
					logger.Errorf(ctx, "ConvImageEditsRequest General model: %s, error: %v", g.Model, err)
					return data, err
				}
			}
		}
	}

	if err = builder.WriteField("prompt", request.Prompt); err != nil {
		logger.Errorf(ctx, "ConvImageEditsRequest General model: %s, error: %v", g.Model, err)
		return data, err
	}

	if request.Background != "" {
		if err = builder.WriteField("background", request.Background); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if request.Mask != nil {
		if err = builder.CreateFormFileHeader("mask", request.Mask); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if request.N != 0 {
		if err = builder.WriteField("n", strconv.Itoa(request.N)); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if request.Quality != "" {
		if err = builder.WriteField("quality", request.Quality); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if request.ResponseFormat != "" {
		if err = builder.WriteField("response_format", request.ResponseFormat); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if request.Size != "" {
		if err = builder.WriteField("size", request.Size); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if request.User != "" {
		if err = builder.WriteField("user", request.User); err != nil {
			logger.Errorf(ctx, "ConvImageEditsRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	g.header["Content-Type"] = builder.FormDataContentType()

	return data, nil
}

func (g *General) ConvImageEditsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (g *General) ConvAudioSpeechRequest(ctx context.Context, data []byte) (request model.SpeechRequest, err error) {

	if err = json.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	return request, nil
}

func (g *General) ConvAudioSpeechResponse(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {
	return model.SpeechResponse{
		Data: data,
	}, nil
}

func (g *General) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (data *bytes.Buffer, err error) {

	data = &bytes.Buffer{}
	builder := util.NewFormBuilder(data)

	defer func() {
		if err := builder.Close(); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest General model: %s, builder.Close() error: %v", g.Model, err)
		}
	}()

	if err = builder.WriteField("model", request.Model); err != nil {
		logger.Errorf(ctx, "ConvAudioTranscriptionsRequest General model: %s, error: %v", g.Model, err)
		return data, err
	}

	if request.File != nil {
		if err = builder.CreateFormFileHeader("file", request.File); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if request.Prompt != "" {
		if err = builder.WriteField("prompt", request.Prompt); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if request.ResponseFormat != "" {
		if err = builder.WriteField("response_format", request.ResponseFormat); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if request.Temperature != 0 {
		if err = builder.WriteField("temperature", fmt.Sprintf("%.2f", request.Temperature)); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if request.Language != "" {
		if err = builder.WriteField("language", request.Language); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if len(request.TimestampGranularities) > 0 {
		for _, timestampGranularitie := range request.TimestampGranularities {
			if err = builder.WriteField("timestamp_granularities[]", timestampGranularitie); err != nil {
				logger.Errorf(ctx, "ConvAudioTranscriptionsRequest General model: %s, error: %v", g.Model, err)
				return data, err
			}
		}
	}

	g.header["Content-Type"] = builder.FormDataContentType()

	return data, nil
}

func (g *General) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (response model.AudioResponse, err error) {

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (g *General) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (request model.EmbeddingRequest, err error) {

	if err = json.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	return request, nil
}

func (g *General) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error) {

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (g *General) ConvVideoCreateRequest(ctx context.Context, request model.VideoCreateRequest) (data *bytes.Buffer, err error) {

	data = &bytes.Buffer{}
	builder := util.NewFormBuilder(data)

	defer func() {
		if err := builder.Close(); err != nil {
			logger.Errorf(ctx, "ConvVideoCreateRequest General model: %s, builder.Close() error: %v", g.Model, err)
		}
	}()

	if err = builder.WriteField("model", request.Model); err != nil {
		logger.Errorf(ctx, "ConvVideoCreateRequest General model: %s, error: %v", g.Model, err)
		return data, err
	}

	if request.Prompt != "" {
		if err = builder.WriteField("prompt", request.Prompt); err != nil {
			logger.Errorf(ctx, "ConvVideoCreateRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if request.InputReference != nil {
		if err = builder.CreateFormFileHeader("input_reference", request.InputReference); err != nil {
			logger.Errorf(ctx, "ConvVideoCreateRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if request.Seconds != "" {
		if err = builder.WriteField("seconds", request.Seconds); err != nil {
			logger.Errorf(ctx, "ConvVideoCreateRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if request.Size != "" {
		if err = builder.WriteField("size", request.Size); err != nil {
			logger.Errorf(ctx, "ConvVideoCreateRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	g.header["Content-Type"] = builder.FormDataContentType()

	return data, nil
}

func (g *General) ConvVideoListResponse(ctx context.Context, data []byte) (response model.VideoListResponse, err error) {

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (g *General) ConvVideoContentResponse(ctx context.Context, data []byte) (response model.VideoContentResponse, err error) {
	return model.VideoContentResponse{
		Data: data,
	}, nil
}

func (g *General) ConvVideoJobResponse(ctx context.Context, data []byte) (response model.VideoJobResponse, err error) {

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (g *General) ConvFileUploadRequest(ctx context.Context, request model.FileUploadRequest) (data *bytes.Buffer, err error) {

	data = &bytes.Buffer{}
	builder := util.NewFormBuilder(data)

	defer func() {
		if err := builder.Close(); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest General model: %s, builder.Close() error: %v", g.Model, err)
		}
	}()

	if request.File != nil {
		if err = builder.CreateFormFileHeader("file", request.File); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if request.Purpose != "" {
		if err = builder.WriteField("purpose", request.Purpose); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if request.ExpiresAfter.Anchor != "" {
		if err = builder.WriteField("expires_after[anchor]", request.ExpiresAfter.Anchor); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	if request.ExpiresAfter.Seconds != "" {
		if err = builder.WriteField("expires_after[seconds]", request.ExpiresAfter.Seconds); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest General model: %s, error: %v", g.Model, err)
			return data, err
		}
	}

	g.header["Content-Type"] = builder.FormDataContentType()

	return data, nil
}

func (g *General) ConvFileListResponse(ctx context.Context, data []byte) (response model.FileListResponse, err error) {

	response.ResponseBytes = data

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (g *General) ConvFileContentResponse(ctx context.Context, data []byte) (response model.FileContentResponse, err error) {
	return model.FileContentResponse{
		Data: data,
	}, nil
}

func (g *General) ConvFileResponse(ctx context.Context, data []byte) (response model.FileResponse, err error) {

	response.ResponseBytes = data

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (g *General) ConvBatchCreateRequest(ctx context.Context, request model.BatchCreateRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvBatchListResponse(ctx context.Context, data []byte) (response model.BatchListResponse, err error) {

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}

func (g *General) ConvBatchResponse(ctx context.Context, data []byte) (response model.BatchResponse, err error) {

	response.ResponseBytes = data

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	return response, nil
}
