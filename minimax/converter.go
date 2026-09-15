package minimax

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (m *MiniMax) ConvChatCompletionsRequest(ctx context.Context, data any) (request model.ChatCompletionRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvChatCompletionsResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvChatResponsesRequest(ctx context.Context, data []byte) (request model.ChatCompletionRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvChatResponsesResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvImageGenerationsRequest(ctx context.Context, data []byte) (request model.ImageGenerationRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvImageGenerationsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvImageEditsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvAudioSpeechRequest(ctx context.Context, data []byte) (request model.SpeechRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvAudioSpeechResponse(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (response model.AudioResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (request model.EmbeddingRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error) {
	//TODO implement me
	panic("implement me")
}

// OpenAI 格式视频创建请求转为 MiniMax 官方格式请求体
func (m *MiniMax) ConvVideoCreateRequest(ctx context.Context, request model.VideoCreateRequest) (data *bytes.Buffer, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvVideoCreateRequest time: %d", gtime.TimestampMilli()-now)
	}()

	model_ := request.Model
	if model_ == "" {
		model_ = m.Model
	}

	req := model.MiniMaxVideoCreateReq{
		Model:      model_,
		Resolution: "768P",
		Duration:   5,
		Ratio:      "16:9",
		Content: []model.MiniMaxContentItem{
			{Type: "text", Text: request.Prompt},
		},
	}

	if request.Seconds != "" {
		duration := gconv.Int(request.Seconds)
		if duration > 0 {
			req.Duration = duration
		}
	}

	if request.Size != "" {
		resolution, ratio := convSizeToResolutionRatio(request.Size)
		if resolution != "" {
			req.Resolution = resolution
		}
		if ratio != "" {
			req.Ratio = ratio
		}
	}

	if request.InputReference != nil {
		dataUri, e := fileHeaderToDataUri(request.InputReference)
		if e != nil {
			logger.Errorf(ctx, "ConvVideoCreateRequest MiniMax input_reference error: %v", e)
			return nil, e
		}
		req.Content = append(req.Content, model.MiniMaxContentItem{
			Type:     "image_url",
			ImageUrl: &model.MiniMaxMediaUrl{Url: dataUri},
			Role:     "first_frame",
		})
		req.Ratio = "adaptive"
	}

	bytes_, err := json.Marshal(req)
	if err != nil {
		logger.Errorf(ctx, "ConvVideoCreateRequest MiniMax json.Marshal error: %v", err)
		return nil, err
	}

	return bytes.NewBuffer(bytes_), nil
}

func (m *MiniMax) ConvVideoListResponse(ctx context.Context, data []byte) (response model.VideoListResponse, err error) {
	return response, errors.New("MiniMax video does not support list")
}

func (m *MiniMax) ConvVideoContentResponse(ctx context.Context, data []byte) (response model.VideoContentResponse, err error) {
	return model.VideoContentResponse{Data: data}, nil
}

// 官方任务响应(创建/查询共用)转为系统标准响应
func (m *MiniMax) ConvVideoJobResponse(ctx context.Context, data []byte) (response model.VideoJobResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvVideoJobResponse time: %d", gtime.TimestampMilli()-now)
	}()

	var raw struct {
		TaskId string                    `json:"task_id"`
		Task   *model.MiniMaxVideoTask   `json:"task"`
		Type   string                    `json:"type"`
		Error  *model.MiniMaxErrorDetail `json:"error"`
	}
	if err = json.Unmarshal(data, &raw); err != nil {
		logger.Errorf(ctx, "ConvVideoJobResponse MiniMax json.Unmarshal error: %v", err)
		return response, err
	}

	if raw.Error != nil {
		err = m.apiErrorHandler(gconv.Int(raw.Error.HttpCode), raw.Error)
		logger.Errorf(ctx, "ConvVideoJobResponse MiniMax model: %s, error: %v", m.Model, err)
		return response, err
	}

	response = model.VideoJobResponse{
		Object:        "video",
		Model:         m.Model,
		CreatedAt:     gtime.Now().Unix(),
		ResponseBytes: data,
	}

	if raw.Task != nil {
		response.Id = raw.Task.Id
		response.Model = raw.Task.Model
		response.Status = convMiniMaxStatus(raw.Task.Status)
		response.CreatedAt = raw.Task.CreatedAt
		response.Seconds = gconv.String(raw.Task.Duration)

		if raw.Task.Resolution != "" {
			response.Size = raw.Task.Resolution
			if raw.Task.Ratio != "" && raw.Task.Ratio != "adaptive" {
				response.Size = raw.Task.Resolution + " " + raw.Task.Ratio
			}
		}

		expiresAt := raw.Task.CreatedAt + 7*24*60*60
		response.ExpiresAt = &expiresAt

		if raw.Task.Status == "succeeded" {
			response.Progress = 100
			if raw.Task.UpdatedAt > 0 {
				response.CompletedAt = &raw.Task.UpdatedAt
			}
		}

		if raw.Task.Content != nil {
			response.VideoUrl = raw.Task.Content.Url
		}

		if raw.Task.Usage != nil {
			if response.Seconds == "" || response.Seconds == "0" {
				if raw.Task.Usage.OutputSeconds > 0 {
					response.Seconds = gconv.String(raw.Task.Usage.OutputSeconds)
				} else if raw.Task.Usage.TotalSeconds > 0 {
					response.Seconds = gconv.String(raw.Task.Usage.TotalSeconds)
				}
			}
			response.Usage = &model.Usage{
				PromptTokens:     raw.Task.Usage.PromptTokens,
				CompletionTokens: raw.Task.Usage.CompletionTokens,
				TotalTokens:      raw.Task.Usage.TotalTokens,
			}
		}

		if raw.Task.Error != nil {
			response.Error = &model.VideoError{
				Code:    raw.Task.Error.Code,
				Message: raw.Task.Error.Message,
			}
		}

		return response, nil
	}

	if raw.TaskId != "" {
		response.Id = raw.TaskId
		response.Status = "queued"
		return response, nil
	}

	return response, nil
}

func (m *MiniMax) ConvFileUploadRequest(ctx context.Context, request model.FileUploadRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvFileListResponse(ctx context.Context, data []byte) (response model.FileListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvFileContentResponse(ctx context.Context, data []byte) (response model.FileContentResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvFileResponse(ctx context.Context, data []byte) (response model.FileResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvBatchCreateRequest(ctx context.Context, request model.BatchCreateRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvBatchListResponse(ctx context.Context, data []byte) (response model.BatchListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ConvBatchResponse(ctx context.Context, data []byte) (response model.BatchResponse, err error) {
	//TODO implement me
	panic("implement me")
}

// 将上传文件编码为 data:{MIME};base64,{data}
func fileHeaderToDataUri(fh *multipart.FileHeader) (string, error) {

	f, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	mimeType := fh.Header.Get("Content-Type")
	if mimeType == "" || mimeType == "application/octet-stream" {
		mimeType = http.DetectContentType(data)
	}

	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(data)), nil
}

// 将 "1280x720" 或 "768P"/"2K"/"480P" 转为 MiniMax 的 resolution 与 ratio
func convSizeToResolutionRatio(size string) (resolution, ratio string) {

	upper := gstr.ToUpper(size)
	switch upper {
	case "480P", "768P", "2K":
		return upper, "16:9"
	}

	var width, height int
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
	case short >= 1440:
		resolution = "2K"
	case short >= 720:
		resolution = "768P"
	default:
		resolution = "480P"
	}

	g := gcd(width, height)
	ratio = fmt.Sprintf("%d:%d", width/g, height/g)

	switch ratio {
	case "21:9", "16:9", "4:3", "1:1", "3:4", "9:16":
	default:
		ratio = "16:9"
	}

	return resolution, ratio
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// 将 MiniMax 任务状态映射到系统标准状态
func convMiniMaxStatus(status string) string {
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
	default:
		return status
	}
}
