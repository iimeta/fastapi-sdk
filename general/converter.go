package general

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/iimeta/fastapi-sdk/v2/common"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
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
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvImageGenerationsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvImageEditsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvAudioSpeechRequest(ctx context.Context, data []byte) (request model.SpeechRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvAudioSpeechResponse(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (response model.AudioResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (request model.EmbeddingRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvVideoCreateRequest(ctx context.Context, request model.VideoCreateRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvVideoListResponse(ctx context.Context, data []byte) (response model.VideoListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvVideoContentResponse(ctx context.Context, data []byte) (response model.VideoContentResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvVideoJobResponse(ctx context.Context, data []byte) (response model.VideoJobResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvFileUploadRequest(ctx context.Context, request model.FileUploadRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvFileListResponse(ctx context.Context, data []byte) (response model.FileListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvFileContentResponse(ctx context.Context, data []byte) (response model.FileContentResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvFileResponse(ctx context.Context, data []byte) (response model.FileResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvBatchCreateRequest(ctx context.Context, request model.BatchCreateRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvBatchListResponse(ctx context.Context, data []byte) (response model.BatchListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ConvBatchResponse(ctx context.Context, data []byte) (response model.BatchResponse, err error) {
	//TODO implement me
	panic("implement me")
}
