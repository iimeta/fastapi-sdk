package volcengine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
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

func (v *VolcEngine) ConvVideoCreateRequest(ctx context.Context, request model.VideoCreateRequest) (*bytes.Buffer, error) {

	req := model.VolcVideoCreateReq{
		Model: v.Model,
	}

	if request.Prompt != "" {
		req.Content = append(req.Content, model.VolcVideoContent{Type: "text", Text: request.Prompt})
	}

	if request.Seconds != "" {
		dur := gconv.Int(request.Seconds)
		req.Duration = &dur
	}

	if request.Size != "" {
		req.Ratio, req.Resolution = convSizeToRatioResolution(request.Size)
	}

	b, err := json.Marshal(req)
	if err != nil {
		logger.Errorf(ctx, "ConvVideoCreateRequest VolcEngine json.Marshal error: %v", err)
		return nil, err
	}

	return bytes.NewBuffer(b), nil
}

func (v *VolcEngine) ConvVideoJobResponse(ctx context.Context, data []byte) (response model.VideoJobResponse, err error) {

	var raw model.VolcVideoTaskRes
	if err = json.Unmarshal(data, &raw); err != nil {
		logger.Errorf(ctx, "ConvVideoJobResponse VolcEngine json.Unmarshal error: %v", err)
		return response, err
	}

	response = model.VideoJobResponse{
		Id:            raw.Id,
		Object:        "video",
		Model:         raw.Model,
		Status:        convVolcStatus(raw.Status),
		CreatedAt:     raw.CreatedAt,
		ResponseBytes: data,
	}

	expiresAt := raw.CreatedAt + int64(raw.ExecutionExpiresAfter)
	response.ExpiresAt = &expiresAt

	if raw.Duration != nil && *raw.Duration > 0 {
		response.Seconds = gconv.String(*raw.Duration)
	} else if raw.Frames != nil && raw.FramesPerSecond != nil && *raw.Frames > 0 && *raw.FramesPerSecond > 0 {
		response.Seconds = gconv.String(*raw.Frames / *raw.FramesPerSecond)
	}

	if raw.Ratio != "" {
		response.Size = raw.Ratio
	}

	if raw.UpdatedAt > 0 && raw.Status == "succeeded" {
		response.CompletedAt = &raw.UpdatedAt
	}

	if raw.Content != nil {
		response.VideoUrl = raw.Content.VideoUrl
	}

	if raw.Usage != nil {
		response.Usage = &model.Usage{
			CompletionTokens: raw.Usage.CompletionTokens,
			TotalTokens:      raw.Usage.TotalTokens,
		}
	}

	if raw.Error != nil {
		response.Error = &model.VideoError{
			Code:    raw.Error.Code,
			Message: raw.Error.Message,
		}
	}

	return response, nil
}

func (v *VolcEngine) ConvVideoListResponse(ctx context.Context, data []byte) (response model.VideoListResponse, err error) {

	var raw model.VolcVideoListRes
	if err = json.Unmarshal(data, &raw); err != nil {
		logger.Errorf(ctx, "ConvVideoListResponse VolcEngine json.Unmarshal error: %v", err)
		return response, err
	}

	response.Object = "list"
	for i := range raw.Items {
		item, e := v.ConvVideoJobResponse(ctx, gjson.MustEncode(raw.Items[i]))
		if e != nil {
			logger.Errorf(ctx, "ConvVideoListResponse VolcEngine ConvVideoJobResponse error: %v", e)
			continue
		}
		response.Data = append(response.Data, item)
	}

	if len(response.Data) > 0 {
		first := response.Data[0].Id
		last := response.Data[len(response.Data)-1].Id
		response.FirstId = &first
		response.LastId = &last
	}

	return response, nil
}

func (v *VolcEngine) ConvVideoContentResponse(ctx context.Context, data []byte) (model.VideoContentResponse, error) {
	return model.VideoContentResponse{Data: data}, nil
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

// 将火山引擎状态映射到系统标准状态
func convVolcStatus(status string) string {
	switch status {
	case "queued":
		return "queued"
	case "running":
		return "in_progress"
	case "succeeded":
		return "completed"
	case "failed":
		return "failed"
	case "cancelled":
		return "deleted"
	case "expired":
		return "expired"
	default:
		return status
	}
}

// 将 "1280x720" 格式转换为 ratio("16:9") 和 resolution("720p")
func convSizeToRatioResolution(size string) (ratio, resolution string) {

	var width, height int
	// 支持 x / X / × / * 分隔
	for _, sep := range []string{"x", "X", "×", "*"} {
		parts := gstr.Split(size, sep)
		if len(parts) == 2 {
			width = gconv.Int(parts[0])
			height = gconv.Int(parts[1])
			break
		}
	}

	if width == 0 || height == 0 {
		return "", ""
	}

	short := height
	if width < height {
		short = width
	}
	switch {
	case short >= 2160:
		resolution = "4k"
	case short >= 1080:
		resolution = "1080p"
	case short >= 720:
		resolution = "720p"
	case short >= 480:
		resolution = "480p"
	default:
		resolution = fmt.Sprintf("%dp", short)
	}

	g := gcd(width, height)
	ratio = fmt.Sprintf("%d:%d", width/g, height/g)

	return ratio, resolution
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
