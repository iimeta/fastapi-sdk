package deepseek

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/v2/common"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (d *DeepSeek) ConvChatCompletionsRequest(ctx context.Context, data any) (request model.ChatCompletionRequest, err error) {

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

	if d.IsSupportSystemRole != nil {
		request.Messages = common.HandleMessages(request.Messages, *d.IsSupportSystemRole)
	}

	return request, nil
}

func (d *DeepSeek) ConvChatCompletionsResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	response.ResponseBytes = data

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if !gstr.HasPrefix(response.Id, consts.COMPLETION_ID_PREFIX) {
		response.Id = consts.COMPLETION_ID_PREFIX + response.Id
	}

	return response, nil
}

func (d *DeepSeek) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	response.ResponseBytes = data

	if err = json.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if !gstr.HasPrefix(response.Id, consts.COMPLETION_ID_PREFIX) {
		response.Id = consts.COMPLETION_ID_PREFIX + response.Id
	}

	return response, nil
}

func (d *DeepSeek) ConvChatResponsesRequest(ctx context.Context, data []byte) (request model.ChatCompletionRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvChatResponsesResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvImageGenerationsRequest(ctx context.Context, data []byte) (request model.ImageGenerationRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvImageGenerationsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvImageEditsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvAudioSpeechRequest(ctx context.Context, data []byte) (request model.SpeechRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvAudioSpeechResponse(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (response model.AudioResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (request model.EmbeddingRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvVideoCreateRequest(ctx context.Context, request model.VideoCreateRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvVideoListResponse(ctx context.Context, data []byte) (response model.VideoListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvVideoContentResponse(ctx context.Context, data []byte) (response model.VideoContentResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvVideoJobResponse(ctx context.Context, data []byte) (response model.VideoJobResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvFileUploadRequest(ctx context.Context, request model.FileUploadRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvFileListResponse(ctx context.Context, data []byte) (response model.FileListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvFileContentResponse(ctx context.Context, data []byte) (response model.FileContentResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvFileResponse(ctx context.Context, data []byte) (response model.FileResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvBatchCreateRequest(ctx context.Context, request model.BatchCreateRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvBatchListResponse(ctx context.Context, data []byte) (response model.BatchListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ConvBatchResponse(ctx context.Context, data []byte) (response model.BatchResponse, err error) {
	//TODO implement me
	panic("implement me")
}
