package volcengine

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

func (v *VolcEngine) ConvChatCompletionsRequest(ctx context.Context, data any) (request model.ChatCompletionRequest, err error) {

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

	if v.IsSupportSystemRole != nil {
		request.Messages = common.HandleMessages(request.Messages, *v.IsSupportSystemRole)
	} else {
		request.Messages = common.HandleMessages(request.Messages, true)
	}

	return request, nil
}

func (v *VolcEngine) ConvChatCompletionsResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

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

func (v *VolcEngine) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

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

func (v *VolcEngine) ConvChatResponsesRequest(ctx context.Context, data []byte) (request model.ChatCompletionRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvChatResponsesResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvImageGenerationsRequest(ctx context.Context, data []byte) (request model.ImageGenerationRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvImageGenerationsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvImageEditsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvAudioSpeechRequest(ctx context.Context, data []byte) (request model.SpeechRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvAudioSpeechResponse(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (response model.AudioResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (request model.EmbeddingRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvVideoCreateRequest(ctx context.Context, request model.VideoCreateRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvVideoListResponse(ctx context.Context, data []byte) (response model.VideoListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvVideoContentResponse(ctx context.Context, data []byte) (response model.VideoContentResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvVideoJobResponse(ctx context.Context, data []byte) (response model.VideoJobResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvFileUploadRequest(ctx context.Context, request model.FileUploadRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvFileListResponse(ctx context.Context, data []byte) (response model.FileListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvFileContentResponse(ctx context.Context, data []byte) (response model.FileContentResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvFileResponse(ctx context.Context, data []byte) (response model.FileResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvBatchCreateRequest(ctx context.Context, request model.BatchCreateRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvBatchListResponse(ctx context.Context, data []byte) (response model.BatchListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvBatchResponse(ctx context.Context, data []byte) (response model.BatchResponse, err error) {
	//TODO implement me
	panic("implement me")
}
