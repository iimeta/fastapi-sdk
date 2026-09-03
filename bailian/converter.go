package bailian

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (b *Bailian) ConvChatCompletionsRequest(ctx context.Context, data any) (request model.ChatCompletionRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvChatCompletionsResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvChatResponsesRequest(ctx context.Context, data []byte) (request model.ChatCompletionRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvChatResponsesResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

// 请求体既可能是 OpenAI 兼容格式, 也可能是百炼官方格式(带 input.messages), 统一转换为 ImageGenerationRequest
func (b *Bailian) ConvImageGenerationsRequest(ctx context.Context, data []byte) (request model.ImageGenerationRequest, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvImageGenerationsRequest time: %d", gtime.TimestampMilli()-now)
	}()

	var official model.BailianImageGenerationReq
	if err = json.Unmarshal(data, &official); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	if len(official.Input.Messages) > 0 {
		return convOfficialToImageGenerationRequest(official), nil
	}

	if err = json.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	return request, nil
}

func (b *Bailian) ConvImageGenerationsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvImageGenerationsResponse time: %d", gtime.TimestampMilli()-now)
	}()

	res := model.BailianImageGenerationRes{}
	if err = json.Unmarshal(data, &res); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if res.Code != "" {
		err = b.apiErrorHandler(0, res.Code, res.Message)
		logger.Errorf(ctx, "ConvImageGenerationsResponse Bailian model: %s, error: %v", b.Model, err)
		return response, err
	}

	response = convOfficialToImageResponse(&res)
	response.ResponseBytes = data

	return response, nil
}

func (b *Bailian) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (data *bytes.Buffer, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvImageEditsRequest time: %d", gtime.TimestampMilli()-now)
	}()

	official, err := b.ConvImageEditsRequestOfficial(ctx, request)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(official), nil
}

func (b *Bailian) ConvImageEditsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	return b.ConvImageGenerationsResponse(ctx, data)
}

func (b *Bailian) ConvAudioSpeechRequest(ctx context.Context, data []byte) (request model.SpeechRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvAudioSpeechResponse(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (response model.AudioResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (request model.EmbeddingRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error) {
	//TODO implement me
	panic("implement me")
}

// OpenAI 格式视频创建请求转为百炼官方格式请求体
func (b *Bailian) ConvVideoCreateRequest(ctx context.Context, request model.VideoCreateRequest) (data *bytes.Buffer, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvVideoCreateRequest time: %d", gtime.TimestampMilli()-now)
	}()

	model_ := request.Model
	if model_ == "" {
		model_ = b.Model
	}

	req := model.BailianVideoCreateReq{
		Model: model_,
		Input: model.BailianVideoInput{
			Prompt: request.Prompt,
		},
	}

	// OpenAI 格式的参考图作为首帧传入, 以 Base64 data URI 形式提交
	if request.InputReference != nil {
		dataUri, e := fileHeaderToDataUri(request.InputReference)
		if e != nil {
			logger.Errorf(ctx, "ConvVideoCreateRequest Bailian input_reference error: %v", e)
			return nil, e
		}
		req.Input.Media = append(req.Input.Media, model.BailianVideoMedia{Type: "first_frame", Url: dataUri})
	}

	parameters := &model.BailianVideoParameters{}

	if request.Seconds != "" {
		duration := gconv.Int(request.Seconds)
		parameters.Duration = &duration
	}

	if request.Size != "" {
		parameters.Resolution, parameters.Ratio = convSizeToResolutionRatio(request.Size)
	}

	req.Parameters = parameters

	bytes_, err := json.Marshal(req)
	if err != nil {
		logger.Errorf(ctx, "ConvVideoCreateRequest Bailian json.Marshal error: %v", err)
		return nil, err
	}

	return bytes.NewBuffer(bytes_), nil
}

func (b *Bailian) ConvVideoListResponse(ctx context.Context, data []byte) (response model.VideoListResponse, err error) {
	return response, errors.New("Bailian video does not support list")
}

func (b *Bailian) ConvVideoContentResponse(ctx context.Context, data []byte) (response model.VideoContentResponse, err error) {
	return model.VideoContentResponse{Data: data}, nil
}

// 官方任务响应(创建/查询共用)转为系统标准响应
func (b *Bailian) ConvVideoJobResponse(ctx context.Context, data []byte) (response model.VideoJobResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvVideoJobResponse time: %d", gtime.TimestampMilli()-now)
	}()

	var raw model.BailianVideoTaskRes
	if err = json.Unmarshal(data, &raw); err != nil {
		logger.Errorf(ctx, "ConvVideoJobResponse Bailian json.Unmarshal error: %v", err)
		return response, err
	}

	// 请求级失败(如鉴权/参数错误), 不存在任务
	if raw.Code != "" {
		err = b.apiErrorHandler(0, raw.Code, raw.Message)
		logger.Errorf(ctx, "ConvVideoJobResponse Bailian model: %s, error: %v", b.Model, err)
		return response, err
	}

	response = model.VideoJobResponse{
		Object:        "video",
		Model:         b.Model,
		CreatedAt:     gtime.Now().Unix(),
		ResponseBytes: data,
	}

	if raw.Output == nil {
		return response, nil
	}

	response.Id = raw.Output.TaskId
	response.Status = convBailianStatus(raw.Output.TaskStatus)
	response.Prompt = raw.Output.OrigPrompt
	response.VideoUrl = raw.Output.VideoUrl

	if submitTime := parseBailianTime(raw.Output.SubmitTime); submitTime > 0 {
		response.CreatedAt = submitTime
	}

	// 任务与视频链接有效期均为24小时
	expiresAt := response.CreatedAt + 24*60*60
	response.ExpiresAt = &expiresAt

	if response.Status == "completed" {
		response.Progress = 100
		if endTime := parseBailianTime(raw.Output.EndTime); endTime > 0 {
			response.CompletedAt = &endTime
		}
	}

	if raw.Usage != nil {

		seconds := raw.Usage.OutputVideoDuration
		if seconds == 0 {
			seconds = raw.Usage.Duration
		}
		if seconds > 0 {
			response.Seconds = gconv.String(int(math.Ceil(seconds)))
		}

		if raw.Usage.SR > 0 {
			ratio := raw.Usage.Ratio
			if ratio == "" {
				ratio = "16:9"
			}
			response.Size = fmt.Sprintf("%dP%s", raw.Usage.SR, ratio)
		}
	}

	// 任务级失败(task_status=FAILED)时错误信息位于 output 内
	if raw.Output.Code != "" {
		response.Error = &model.VideoError{
			Code:    raw.Output.Code,
			Message: raw.Output.Message,
		}
	}

	return response, nil
}

func (b *Bailian) ConvFileUploadRequest(ctx context.Context, request model.FileUploadRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvFileListResponse(ctx context.Context, data []byte) (response model.FileListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvFileContentResponse(ctx context.Context, data []byte) (response model.FileContentResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvFileResponse(ctx context.Context, data []byte) (response model.FileResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvBatchCreateRequest(ctx context.Context, request model.BatchCreateRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvBatchListResponse(ctx context.Context, data []byte) (response model.BatchListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvBatchResponse(ctx context.Context, data []byte) (response model.BatchResponse, err error) {
	//TODO implement me
	panic("implement me")
}

// 判断请求体是否为百炼官方格式(含 input.messages)
func isOfficialRequest(data []byte) bool {
	var official model.BailianImageGenerationReq
	if err := json.Unmarshal(data, &official); err != nil {
		return false
	}
	return len(official.Input.Messages) > 0
}

// 官方请求转为系统标准请求, 供计费与日志使用
func convOfficialToImageGenerationRequest(official model.BailianImageGenerationReq) model.ImageGenerationRequest {

	request := model.ImageGenerationRequest{
		Model: official.Model,
	}

	for _, message := range official.Input.Messages {
		for _, content := range message.Content {
			if content.Text != "" {
				request.Prompt += content.Text
			}
			if content.Image != "" {
				request.Images = append(request.Images, model.ImageEditImage{ImageUrl: content.Image})
			}
		}
	}

	// 官方默认分辨率为 2K
	request.Size = "2K"

	if official.Parameters != nil {

		if official.Parameters.Size != "" {
			request.Size = official.Parameters.Size
		}

		request.N = official.Parameters.N
		request.Watermark = official.Parameters.Watermark

		if official.Parameters.EnableSequential != nil && *official.Parameters.EnableSequential {
			request.SequentialImageGeneration = "auto"
		}
	}

	return request
}

// 官方响应转为系统标准响应
func convOfficialToImageResponse(res *model.BailianImageGenerationRes) model.ImageResponse {

	response := model.ImageResponse{
		Created: gtime.Now().Unix(),
	}

	var size string
	if res.Usage != nil {
		size = res.Usage.Size
		response.Usage = model.Usage{
			InputTokens:  res.Usage.InputTokens,
			OutputTokens: res.Usage.OutputTokens,
			TotalTokens:  res.Usage.TotalTokens,
		}
	}

	if res.Output != nil {
		for _, choice := range res.Output.Choices {
			var revisedPrompt string
			for _, content := range choice.Message.Content {
				if content.Text != "" {
					revisedPrompt += content.Text
				}
			}
			for _, content := range choice.Message.Content {
				if content.Image == "" {
					continue
				}
				response.Data = append(response.Data, model.ImageResponseData{
					Url:           content.Image,
					Size:          size,
					RevisedPrompt: revisedPrompt,
					OutputFormat:  "png",
				})
			}
		}
	}

	return response
}

// 将 OpenAI 格式的 size(1024x1024) 转为官方格式(1024*1024), 1K/2K/4K 保持不变
func convSizeToOfficial(size, quality string) string {

	if size == "" {
		if gstr.HasSuffix(quality, "K") {
			return quality
		}
		return ""
	}

	if gstr.HasSuffix(size, "K") {
		return size
	}

	return gstr.ReplaceByMap(size, map[string]string{"×": "*", "x": "*", "X": "*"})
}

// 提取 OpenAI 格式请求中的输入图像(URL 或 Base64)
func collectImages(image any, images []model.ImageEditImage) []string {

	var urls []string

	for _, item := range images {
		if item.ImageUrl != "" {
			urls = append(urls, item.ImageUrl)
		}
	}

	switch v := image.(type) {
	case string:
		if v != "" {
			urls = append(urls, v)
		}
	case []string:
		for _, s := range v {
			if s != "" {
				urls = append(urls, s)
			}
		}
	case []any:
		for _, item := range v {
			if s, ok := item.(string); ok && s != "" {
				urls = append(urls, s)
			}
		}
	case model.ImageEditImage:
		if v.ImageUrl != "" {
			urls = append(urls, v.ImageUrl)
		}
	case map[string]any:
		if s, ok := v["image_url"].(string); ok && s != "" {
			urls = append(urls, s)
		}
	}

	return urls
}

func buildOfficialRequest(model_ string, prompt string, imageUrls []string, n int, size, quality string, watermark *bool, sequential string) model.BailianImageGenerationReq {

	content := make([]model.BailianImageContent, 0, len(imageUrls)+1)
	for _, url := range imageUrls {
		content = append(content, model.BailianImageContent{Image: url})
	}
	if prompt != "" {
		content = append(content, model.BailianImageContent{Text: prompt})
	}

	req := model.BailianImageGenerationReq{
		Model: model_,
		Input: model.BailianImageInput{
			Messages: []model.BailianImageMessage{{
				Role:    consts.ROLE_USER,
				Content: content,
			}},
		},
	}

	parameters := &model.BailianImageParameters{
		N:         n,
		Size:      convSizeToOfficial(size, quality),
		Watermark: watermark,
	}

	if sequential == "auto" {
		enable := true
		parameters.EnableSequential = &enable
	}

	req.Parameters = parameters

	return req
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

// 将 "1280x720" 格式转换为 resolution("720P") 和 ratio("16:9"); 已是 1080P/720P/480P 档位时直接使用, 比例自适应
func convSizeToResolutionRatio(size string) (resolution, ratio string) {

	upper := gstr.ToUpper(size)
	if upper == "1080P" || upper == "720P" || upper == "480P" {
		return upper, "adaptive"
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
	case short >= 1080:
		resolution = "1080P"
	case short >= 720:
		resolution = "720P"
	default:
		resolution = "480P"
	}

	g := width
	for b := height; b != 0; {
		g, b = b, g%b
	}
	ratio = fmt.Sprintf("%d:%d", width/g, height/g)

	// 官方仅支持固定几种比例, 其余交由模型自适应
	switch ratio {
	case "16:9", "4:3", "1:1", "3:4", "9:16":
	default:
		ratio = "adaptive"
	}

	return resolution, ratio
}

// 将百炼任务状态映射到系统标准状态
func convBailianStatus(status string) string {
	switch status {
	case "PENDING":
		return "queued"
	case "RUNNING":
		return "in_progress"
	case "SUCCEEDED":
		return "completed"
	case "FAILED":
		return "failed"
	case "CANCELED":
		return "deleted"
	case "UNKNOWN":
		return "expired"
	default:
		return status
	}
}

// 解析 "2026-08-06 10:01:35.452" 格式为秒级时间戳, 解析失败返回 0
func parseBailianTime(s string) int64 {

	if s == "" {
		return 0
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05.000", s, time.Local)
	if err != nil {
		if t, err = time.ParseInLocation(time.DateTime, s, time.Local); err != nil {
			return 0
		}
	}

	return t.Unix()
}
