package xai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/common"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (x *XAI) ConvChatCompletionsRequest(ctx context.Context, data any) (request model.ChatCompletionRequest, err error) {

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

	if request.Stream {
		// 默认让流式返回usage
		if request.StreamOptions == nil {
			request.StreamOptions = &model.StreamOptions{
				IncludeUsage: true,
			}
		}
	}

	if x.IsSupportSystemRole != nil {
		request.Messages = common.HandleMessages(request.Messages, *x.IsSupportSystemRole)
	}

	return request, nil
}

func convChatCompletionsOfficialJSON(request model.ChatCompletionRequest) model.XAIChatCompletionReq {
	return convChatToXAIReq(request)
}

func (x *XAI) ConvChatCompletionsResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvChatCompletionsResponse time: %d", gtime.TimestampMilli()-now)
	}()

	var official model.XAIChatCompletionRes
	if err = json.Unmarshal(data, &official); err != nil {
		logger.Error(ctx, err)
		return response, err
	}
	return convXAIChatToUnified(official, data), nil
}

func (x *XAI) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvChatCompletionsStreamResponse time: %d", gtime.TimestampMilli()-now)
	}()

	var official model.XAIChatCompletionRes
	if err = json.Unmarshal(data, &official); err != nil {
		logger.Error(ctx, err)
		return response, err
	}
	response = convXAIChatToUnified(official, data)

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

func (x *XAI) ConvChatResponsesRequest(ctx context.Context, data []byte) (request model.ChatCompletionRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *XAI) ConvChatResponsesResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *XAI) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *XAI) ConvImageGenerationsRequest(ctx context.Context, data []byte) (request model.ImageGenerationRequest, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvImageGenerationsRequest time: %d", gtime.TimestampMilli()-now)
	}()

	if err = json.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	if request.AspectRatio == "" && request.Size != "" {
		request.AspectRatio = sizeToAspectRatio(request.Size)
	}

	if request.Image != nil {
		if converted := convImageContent(request.Image); converted != nil {
			request.Image = converted
		}
	}

	return request, nil
}

func convImageGenerationsOfficialJSON(request model.ImageGenerationRequest) model.XAIImageReq {
	return convImageToXAIReq(request)
}

func sanitizeImageGenerationsJSON(data []byte) []byte {

	raw := map[string]any{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return data
	}

	size, _ := raw["size"].(string)
	aspectRatio, _ := raw["aspect_ratio"].(string)
	if aspectRatio == "" && size != "" {
		if converted := sizeToAspectRatio(size); converted != "" {
			raw["aspect_ratio"] = converted
		}
	}
	delete(raw, "size")

	bytes, err := json.Marshal(raw)
	if err != nil {
		return data
	}
	return bytes
}

func (x *XAI) ConvImageGenerationsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvImageGenerationsResponse time: %d", gtime.TimestampMilli()-now)
	}()

	var official model.XAIImageRes
	if err = json.Unmarshal(data, &official); err != nil {
		logger.Error(ctx, err)
		return response, err
	}
	response = convXAIImageToUnified(official, data)
	if len(response.Data) == 0 {
		raw := map[string]any{}
		if json.Unmarshal(data, &raw) == nil {
			if images, ok := raw["images"].([]any); ok {
				for _, item := range images {
					imageMap, _ := item.(map[string]any)
					if imageMap == nil {
						continue
					}
					entry := model.ImageResponseData{}
					if url, _ := imageMap["url"].(string); url != "" {
						entry.Url = url
					}
					if b64, _ := imageMap["b64_json"].(string); b64 != "" {
						entry.B64Json = b64
					} else if b64, _ = imageMap["base64"].(string); b64 != "" {
						entry.B64Json = b64
					}
					response.Data = append(response.Data, entry)
				}
			}
		}
	}
	return response, nil
}

func (x *XAI) ConvImageGenerationsStreamResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

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
	var streamOfficial model.XAIImageRes
	if json.Unmarshal(data, &streamOfficial) == nil {
		if u := ConvertUsage(streamOfficial.Usage); u != nil {
			response.Usage = *u
		}
	}

	if streamResponse.B64Json != "" {
		response.Data = []model.ImageResponseData{
			{
				B64Json: streamResponse.B64Json,
			},
		}
	}

	return response, nil
}

func (x *XAI) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (data *bytes.Buffer, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvImageEditsRequest time: %d", gtime.TimestampMilli()-now)
	}()

	images := make([]any, 0)

	for _, item := range request.Images {
		if converted := convImageEditImage(item); converted != nil {
			images = append(images, converted)
		}
	}

	switch v := request.Image.(type) {
	case string:
		if converted := convImageContent(v); converted != nil {
			images = append(images, converted)
		}
	case []string:
		for _, u := range v {
			if converted := convImageContent(u); converted != nil {
				images = append(images, converted)
			}
		}
	case []any:
		for _, item := range v {
			if converted := convImageContent(item); converted != nil {
				images = append(images, converted)
			}
		}
	case []*multipart.FileHeader:
		for _, fh := range v {
			dataURL, convErr := fileHeaderToDataURL(fh)
			if convErr != nil {
				logger.Errorf(ctx, "ConvImageEditsRequest xAI model: %s, fileHeaderToDataURL error: %v", x.Model, convErr)
				return nil, convErr
			}
			if converted := convImageContent(dataURL); converted != nil {
				images = append(images, converted)
			}
		}
	}

	jsonReq := map[string]any{
		"prompt": request.Prompt,
		"model":  request.Model,
	}

	if len(images) == 1 {
		jsonReq["image"] = images[0]
	} else if len(images) > 1 {
		jsonReq["images"] = images
	}

	if request.N != 0 {
		jsonReq["n"] = request.N
	}
	if request.Quality != "" {
		jsonReq["quality"] = request.Quality
	}
	if request.ResponseFormat != "" {
		jsonReq["response_format"] = request.ResponseFormat
	}
	if request.User != "" {
		jsonReq["user"] = request.User
	}

	aspectRatio := request.AspectRatio
	if aspectRatio == "" && request.Size != "" {
		aspectRatio = sizeToAspectRatio(request.Size)
	}
	if aspectRatio != "" {
		jsonReq["aspect_ratio"] = aspectRatio
	}
	if request.Resolution != "" {
		jsonReq["resolution"] = request.Resolution
	}
	if request.StorageOptions != nil {
		jsonReq["storage_options"] = request.StorageOptions
	}

	jsonBytes, err := json.Marshal(jsonReq)
	if err != nil {
		logger.Errorf(ctx, "ConvImageEditsRequest xAI model: %s, json.Marshal error: %v", x.Model, err)
		return nil, err
	}

	x.header["Content-Type"] = "application/json"
	return bytes.NewBuffer(jsonBytes), nil
}

func fileHeaderToDataURL(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader == nil {
		return "", nil
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" || mimeType == "application/octet-stream" {
		mimeType = http.DetectContentType(data)
	}
	if mimeType == "application/octet-stream" {
		switch strings.ToLower(path.Ext(fileHeader.Filename)) {
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".png":
			mimeType = "image/png"
		case ".webp":
			mimeType = "image/webp"
		case ".gif":
			mimeType = "image/gif"
		case ".mp4":
			mimeType = "video/mp4"
		}
	}

	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(data)), nil
}

func convImageContent(v any) any {
	switch item := v.(type) {
	case nil:
		return nil
	case string:
		if item == "" {
			return nil
		}
		return map[string]any{"url": item}
	case model.ImageEditImage:
		return convImageEditImage(item)
	case map[string]any:
		return convImageMap(item)
	default:
		return v
	}
}

func convImageEditImage(item model.ImageEditImage) map[string]any {
	content := map[string]any{}
	if item.FileId != "" {
		content["file_id"] = item.FileId
	}
	if item.ImageUrl != "" {
		content["url"] = item.ImageUrl
	}
	if len(content) == 0 {
		return nil
	}
	return content
}

func convImageMap(item map[string]any) map[string]any {
	content := map[string]any{}
	if fileId, _ := item["file_id"].(string); fileId != "" {
		content["file_id"] = fileId
	}
	if url, _ := item["url"].(string); url != "" {
		content["url"] = url
	} else if imageUrl, _ := item["image_url"].(string); imageUrl != "" {
		content["url"] = imageUrl
	}
	if len(content) == 0 {
		return nil
	}
	return content
}

func (x *XAI) ConvImageEditsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	return x.ConvImageGenerationsResponse(ctx, data)
}

func (x *XAI) ConvAudioSpeechRequest(ctx context.Context, data []byte) (request model.SpeechRequest, err error) {

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

func (x *XAI) ConvAudioSpeechResponse(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvAudioSpeechResponse time: %d", gtime.TimestampMilli()-now)
	}()

	return model.SpeechResponse{
		Data: data,
	}, nil
}

func (x *XAI) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (data *bytes.Buffer, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvAudioTranscriptionsRequest time: %d", gtime.TimestampMilli()-now)
	}()

	data = &bytes.Buffer{}
	builder := util.NewFormBuilder(data)

	defer func() {
		if err := builder.Close(); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest xAI model: %s, builder.Close() error: %v", x.Model, err)
		}
	}()

	if err = builder.WriteField("model", request.Model); err != nil {
		logger.Errorf(ctx, "ConvAudioTranscriptionsRequest xAI model: %s, error: %v", x.Model, err)
		return data, err
	}

	if request.File != nil {
		if err = builder.CreateFormFileHeader("file", request.File); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest xAI model: %s, error: %v", x.Model, err)
			return data, err
		}
	}

	if request.Prompt != "" {
		if err = builder.WriteField("prompt", request.Prompt); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest xAI model: %s, error: %v", x.Model, err)
			return data, err
		}
	}

	if request.ResponseFormat != "" {
		if err = builder.WriteField("response_format", request.ResponseFormat); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest xAI model: %s, error: %v", x.Model, err)
			return data, err
		}
	}

	if request.Temperature != 0 {
		if err = builder.WriteField("temperature", fmt.Sprintf("%.2f", request.Temperature)); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest xAI model: %s, error: %v", x.Model, err)
			return data, err
		}
	}

	if request.Language != "" {
		if err = builder.WriteField("language", request.Language); err != nil {
			logger.Errorf(ctx, "ConvAudioTranscriptionsRequest xAI model: %s, error: %v", x.Model, err)
			return data, err
		}
	}

	if len(request.TimestampGranularities) > 0 {
		for _, timestampGranularitie := range request.TimestampGranularities {
			if err = builder.WriteField("timestamp_granularities[]", timestampGranularitie); err != nil {
				logger.Errorf(ctx, "ConvAudioTranscriptionsRequest xAI model: %s, error: %v", x.Model, err)
				return data, err
			}
		}
	}

	x.header["Content-Type"] = builder.FormDataContentType()

	return data, nil
}

func (x *XAI) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (response model.AudioResponse, err error) {

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

func (x *XAI) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (request model.EmbeddingRequest, err error) {

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

func (x *XAI) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error) {

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

func (x *XAI) ConvVideoCreateRequest(ctx context.Context, request model.VideoCreateRequest) (data *bytes.Buffer, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvVideoCreateRequest time: %d", gtime.TimestampMilli()-now)
	}()

	jsonReq := map[string]any{
		"model":  request.Model,
		"prompt": request.Prompt,
	}

	duration := request.Duration
	if duration == 0 && request.Seconds != "" {
		if parsed, parseErr := strconv.Atoi(strings.TrimSpace(request.Seconds)); parseErr == nil {
			duration = parsed
		}
	}
	if duration > 0 {
		jsonReq["duration"] = duration
	}

	aspectRatio := request.AspectRatio
	resolution := request.Resolution
	if request.Size != "" {
		size := strings.ToLower(request.Size)
		if strings.HasSuffix(size, "p") && resolution == "" {
			resolution = size
		} else if strings.Contains(request.Size, ":") && aspectRatio == "" {
			aspectRatio = request.Size
		} else if aspectRatio == "" {
			aspectRatio = sizeToAspectRatio(request.Size)
		}
	}
	if aspectRatio != "" {
		jsonReq["aspect_ratio"] = aspectRatio
	}
	if resolution != "" {
		jsonReq["resolution"] = resolution
	}

	if image := convImageContent(request.Image); image != nil {
		jsonReq["image"] = image
	} else if request.ImageUrl != "" {
		jsonReq["image"] = convImageContent(request.ImageUrl)
	} else if request.InputReference != nil {
		dataURL, convErr := fileHeaderToDataURL(request.InputReference)
		if convErr != nil {
			logger.Errorf(ctx, "ConvVideoCreateRequest xAI model: %s, fileHeaderToDataURL error: %v", x.Model, convErr)
			return nil, convErr
		}
		if image = convImageContent(dataURL); image != nil {
			jsonReq["image"] = image
		}
	}

	if video := convVideoContent(request.Video); video != nil {
		jsonReq["video"] = video
	} else if request.VideoUrl != "" {
		jsonReq["video"] = convVideoContent(request.VideoUrl)
	}

	if lastFrame := convImageContent(request.LastFrame); lastFrame != nil {
		jsonReq["last_frame"] = lastFrame
	} else if request.LastFrameUrl != "" {
		jsonReq["last_frame"] = convImageContent(request.LastFrameUrl)
	}

	if request.Keyframes != nil {
		jsonReq["keyframes"] = request.Keyframes
	}
	if request.ReferenceImages != nil {
		jsonReq["reference_images"] = request.ReferenceImages
	}
	if request.ReferenceAudios != nil {
		jsonReq["reference_audios"] = request.ReferenceAudios
	}
	if request.GenerateAudio != nil {
		jsonReq["generate_audio"] = *request.GenerateAudio
	}
	if request.StorageOptions != nil {
		jsonReq["storage_options"] = request.StorageOptions
	}
	if request.User != "" {
		jsonReq["user"] = request.User
	}

	jsonBytes, err := json.Marshal(jsonReq)
	if err != nil {
		logger.Errorf(ctx, "ConvVideoCreateRequest xAI model: %s, json.Marshal error: %v", x.Model, err)
		return nil, err
	}

	x.header["Content-Type"] = "application/json"
	return bytes.NewBuffer(jsonBytes), nil
}

func (x *XAI) ConvVideoListResponse(ctx context.Context, data []byte) (response model.VideoListResponse, err error) {

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

func (x *XAI) ConvVideoContentResponse(ctx context.Context, data []byte) (response model.VideoContentResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvVideoContentResponse time: %d", gtime.TimestampMilli()-now)
	}()

	return model.VideoContentResponse{
		Data: data,
	}, nil
}

func (x *XAI) ConvVideoJobResponse(ctx context.Context, data []byte) (response model.VideoJobResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvVideoJobResponse time: %d", gtime.TimestampMilli()-now)
	}()

	var official model.XAIVideoJobRes
	if err = json.Unmarshal(data, &official); err != nil {
		logger.Error(ctx, err)
		return response, err
	}
	response = convXAIVideoToUnified(official, data)

	raw := map[string]any{}
	if json.Unmarshal(data, &raw) == nil {
		response.Status = convVideoStatus(response.Status, raw)
		if response.Id == "" {
			if requestId, _ := raw["request_id"].(string); requestId != "" {
				response.Id = requestId
			}
		}
		if response.VideoUrl == "" {
			if video, ok := raw["video"].(map[string]any); ok {
				if url, _ := video["url"].(string); url != "" {
					response.VideoUrl = url
				}
			}
		}
	}
	return response, nil
}

func convVideoStatus(status string, raw map[string]any) string {
	if status == "" && raw != nil {
		if value, _ := raw["status"].(string); value != "" {
			status = value
		}
	}

	switch strings.ToLower(status) {
	case "done", "completed", "complete", "succeeded", "success":
		return "completed"
	case "pending", "processing", "running", "in_progress":
		return "in_progress"
	case "queued":
		return "queued"
	case "expired":
		return "expired"
	case "failed", "error":
		return "failed"
	case "cancelled", "canceled":
		return "cancelled"
	default:
		return status
	}
}

func convVideoContent(v any) any {
	switch item := v.(type) {
	case nil:
		return nil
	case string:
		if item == "" {
			return nil
		}
		return map[string]any{"url": item}
	case map[string]any:
		content := map[string]any{}
		if fileId, _ := item["file_id"].(string); fileId != "" {
			content["file_id"] = fileId
		}
		if url, _ := item["url"].(string); url != "" {
			content["url"] = url
		} else if videoUrl, _ := item["video_url"].(string); videoUrl != "" {
			content["url"] = videoUrl
		} else if imageUrl, _ := item["image_url"].(string); imageUrl != "" {
			content["url"] = imageUrl
		}
		if len(content) == 0 {
			return nil
		}
		return content
	default:
		return v
	}
}

func sizeToAspectRatio(size string) string {

	size = strings.TrimSpace(strings.ToLower(size))
	if size == "" || size == "auto" {
		return ""
	}

	normalized := strings.ReplaceAll(size, "x", ":")
	switch normalized {
	case "1:1", "3:4", "4:3", "9:16", "16:9", "2:3", "3:2", "9:19.5", "19.5:9", "9:20", "20:9", "1:2", "2:1", "21:9", "5:2":
		return normalized
	}

	parts := strings.Split(strings.ReplaceAll(size, ":", "x"), "x")
	if len(parts) != 2 {
		return ""
	}

	width, errW := strconv.Atoi(parts[0])
	height, errH := strconv.Atoi(parts[1])
	if errW != nil || errH != nil || width <= 0 || height <= 0 {
		return ""
	}

	ratio := float64(width) / float64(height)
	official := []struct {
		name  string
		value float64
	}{
		{"1:1", 1},
		{"3:4", 0.75},
		{"4:3", 4.0 / 3},
		{"9:16", 9.0 / 16},
		{"16:9", 16.0 / 9},
		{"2:3", 2.0 / 3},
		{"3:2", 1.5},
		{"9:19.5", 9.0 / 19.5},
		{"19.5:9", 19.5 / 9},
		{"9:20", 0.45},
		{"20:9", 20.0 / 9},
		{"1:2", 0.5},
		{"2:1", 2},
	}

	best := ""
	bestDiff := 1.0
	for _, item := range official {
		diff := ratio - item.value
		if diff < 0 {
			diff = -diff
		}
		if diff < bestDiff {
			bestDiff = diff
			best = item.name
		}
	}
	if bestDiff > 0.08 {
		return ""
	}
	return best
}

func (x *XAI) ConvFileUploadRequest(ctx context.Context, request model.FileUploadRequest) (data *bytes.Buffer, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvFileUploadRequest time: %d", gtime.TimestampMilli()-now)
	}()

	data = &bytes.Buffer{}
	builder := util.NewFormBuilder(data)

	defer func() {
		if err := builder.Close(); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest xAI model: %s, builder.Close() error: %v", x.Model, err)
		}
	}()

	if request.File != nil {
		if err = builder.CreateFormFileHeader("file", request.File); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest xAI model: %s, error: %v", x.Model, err)
			return data, err
		}
	}

	if request.Purpose != "" {
		if err = builder.WriteField("purpose", request.Purpose); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest xAI model: %s, error: %v", x.Model, err)
			return data, err
		}
	}

	if request.ExpiresAfter.Anchor != "" {
		if err = builder.WriteField("expires_after[anchor]", request.ExpiresAfter.Anchor); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest xAI model: %s, error: %v", x.Model, err)
			return data, err
		}
	}

	if request.ExpiresAfter.Seconds != "" {
		if err = builder.WriteField("expires_after[seconds]", request.ExpiresAfter.Seconds); err != nil {
			logger.Errorf(ctx, "ConvFileUploadRequest xAI model: %s, error: %v", x.Model, err)
			return data, err
		}
	}

	x.header["Content-Type"] = builder.FormDataContentType()

	return data, nil
}

func (x *XAI) ConvFileListResponse(ctx context.Context, data []byte) (response model.FileListResponse, err error) {

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

func (x *XAI) ConvFileContentResponse(ctx context.Context, data []byte) (response model.FileContentResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvFileContentResponse time: %d", gtime.TimestampMilli()-now)
	}()

	return model.FileContentResponse{
		Data: data,
	}, nil
}

func (x *XAI) ConvFileResponse(ctx context.Context, data []byte) (response model.FileResponse, err error) {

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

func (x *XAI) ConvBatchCreateRequest(ctx context.Context, request model.BatchCreateRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *XAI) ConvBatchListResponse(ctx context.Context, data []byte) (response model.BatchListResponse, err error) {

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

func (x *XAI) ConvBatchResponse(ctx context.Context, data []byte) (response model.BatchResponse, err error) {

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
