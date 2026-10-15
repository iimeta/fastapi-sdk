package google

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"mime/multipart"
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/iimeta/fastapi-sdk/v2/common"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
)

func convInteractionUsage(u *model.GoogleInteractionUsage) *model.Usage {

	if u == nil {
		return nil
	}

	if u.UsageMetadata != nil {
		usage := &model.Usage{
			PromptTokens:     u.UsageMetadata.PromptTokenCount,
			CompletionTokens: u.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      u.UsageMetadata.TotalTokenCount,
			InputTokens:      u.UsageMetadata.PromptTokenCount,
			OutputTokens:     u.UsageMetadata.CandidatesTokenCount,
			PromptTokensDetails: model.PromptTokensDetails{
				CachedTokens: u.UsageMetadata.CachedContentTokenCount,
			},
			OutputTokensDetails: model.OutputTokensDetails{
				ReasoningTokens: u.UsageMetadata.ThoughtsTokenCount,
			},
		}

		for _, detail := range u.UsageMetadata.PromptTokensDetails {
			switch strings.ToUpper(detail.Modality) {
			case "TEXT":
				usage.PromptTokensDetails.TextTokens = detail.TokenCount
				usage.InputTokensDetails.TextTokens = detail.TokenCount
			case "IMAGE":
				usage.PromptTokensDetails.ImageTokens = detail.TokenCount
				usage.InputTokensDetails.ImageTokens = detail.TokenCount
			case "AUDIO":
				usage.PromptTokensDetails.AudioTokens = detail.TokenCount
			}
		}

		for _, detail := range u.UsageMetadata.CandidatesTokensDetails {
			switch strings.ToUpper(detail.Modality) {
			case "TEXT":
				usage.CompletionTokensDetails.TextTokens = detail.TokenCount
			case "IMAGE":
				usage.CompletionTokensDetails.ImageTokens = detail.TokenCount
			}
		}

		return usage
	}

	inputTokens := u.InputTokens
	if inputTokens == 0 {
		inputTokens = u.TotalInputTokens
	}

	outputTokens := u.OutputTokens
	if outputTokens == 0 {
		outputTokens = u.TotalOutputTokens
	}

	totalTokens := u.TotalTokens
	if totalTokens == 0 {
		totalTokens = inputTokens + outputTokens
	}

	cachedTokens := u.CachedTokens
	if cachedTokens == 0 {
		cachedTokens = u.TotalCachedTokens
	}

	thoughtsTokens := u.ThoughtsTokens
	if thoughtsTokens == 0 {
		thoughtsTokens = u.TotalThoughtTokens
	}

	usage := &model.Usage{
		PromptTokens:     inputTokens,
		CompletionTokens: outputTokens,
		TotalTokens:      totalTokens,
		InputTokens:      inputTokens,
		OutputTokens:     outputTokens,
		PromptTokensDetails: model.PromptTokensDetails{
			CachedTokens: cachedTokens,
		},
		OutputTokensDetails: model.OutputTokensDetails{
			ReasoningTokens: thoughtsTokens,
		},
	}

	applyModalityTokens(u.InputTokensByModality, &usage.PromptTokensDetails, &usage.InputTokensDetails, nil)
	applyModalityTokens(u.OutputTokensByModality, nil, nil, &usage.CompletionTokensDetails)

	if u.InputTokensDetails != nil {
		if u.InputTokensDetails.TextTokens > 0 {
			usage.PromptTokensDetails.TextTokens = u.InputTokensDetails.TextTokens
			usage.InputTokensDetails.TextTokens = u.InputTokensDetails.TextTokens
		}
		if u.InputTokensDetails.ImageTokens > 0 {
			usage.PromptTokensDetails.ImageTokens = u.InputTokensDetails.ImageTokens
			usage.InputTokensDetails.ImageTokens = u.InputTokensDetails.ImageTokens
		}
		if u.InputTokensDetails.AudioTokens > 0 {
			usage.PromptTokensDetails.AudioTokens = u.InputTokensDetails.AudioTokens
		}
		if u.InputTokensDetails.CachedTokens > 0 {
			usage.PromptTokensDetails.CachedTokens = u.InputTokensDetails.CachedTokens
		}
	}

	if u.OutputTokensDetails != nil {
		if u.OutputTokensDetails.TextTokens > 0 {
			usage.CompletionTokensDetails.TextTokens = u.OutputTokensDetails.TextTokens
		}
		if u.OutputTokensDetails.ImageTokens > 0 {
			usage.CompletionTokensDetails.ImageTokens = u.OutputTokensDetails.ImageTokens
		}
		if u.OutputTokensDetails.CachedTokens > 0 {
			usage.PromptTokensDetails.CachedTokens = u.OutputTokensDetails.CachedTokens
		}
	}

	return usage
}

func applyModalityTokens(items []model.GoogleInteractionModalityToken, prompt *model.PromptTokensDetails, input *model.InputTokensDetails, completion *model.CompletionTokensDetails) {

	for _, item := range items {

		count := item.Tokens
		if count == 0 {
			count = item.TokenCount
		}
		if count == 0 {
			count = item.TokenCountSnake
		}

		switch strings.ToUpper(item.Modality) {
		case "TEXT":
			if prompt != nil {
				prompt.TextTokens = count
			}
			if input != nil {
				input.TextTokens = count
			}
			if completion != nil {
				completion.TextTokens = count
			}
		case "IMAGE":
			if prompt != nil {
				prompt.ImageTokens = count
			}
			if input != nil {
				input.ImageTokens = count
			}
			if completion != nil {
				completion.ImageTokens = count
			}
		case "AUDIO":
			if prompt != nil {
				prompt.AudioTokens = count
			}
			if completion != nil {
				completion.AudioTokens = count
			}
		}
	}
}

func convInteractionToChatCompletion(data []byte, modelName string, stream bool) (model.ChatCompletionResponse, error) {

	interactionRes := model.GoogleInteractionRes{}
	if err := json.Unmarshal(data, &interactionRes); err != nil {
		return model.ChatCompletionResponse{}, err
	}

	if interactionRes.Error.Code != 0 {
		return model.ChatCompletionResponse{}, errorsNewApi(interactionRes)
	}

	object := consts.COMPLETION_OBJECT
	if stream {
		object = consts.COMPLETION_STREAM_OBJECT
	}

	response := model.ChatCompletionResponse{
		Id:            interactionRes.Id,
		Object:        object,
		Created:       gtime.Timestamp(),
		Model:         modelName,
		Usage:         convInteractionUsage(interactionRes.Usage),
		ResponseBytes: data,
	}

	if response.Id == "" {
		response.Id = consts.COMPLETION_ID_PREFIX + grand.S(29)
	}

	if stream {
		text := interactionDeltaText(&interactionRes)
		if text != "" {
			response.Choices = append(response.Choices, model.ChatCompletionChoice{
				Index: 0,
				Delta: &model.ChatCompletionStreamChoiceDelta{
					Role:    consts.ROLE_ASSISTANT,
					Content: text,
				},
			})
		}
		return response, nil
	}

	contents := interactionTextContents(interactionRes.Steps)
	if len(contents) == 0 {
		contents = []string{""}
	}

	for i, content := range contents {
		response.Choices = append(response.Choices, model.ChatCompletionChoice{
			Index: i,
			Message: &model.ChatCompletionMessage{
				Role:    consts.ROLE_ASSISTANT,
				Content: content,
			},
			FinishReason: consts.FinishReasonStop,
		})
	}

	return response, nil
}

func errorsNewApi(res model.GoogleInteractionRes) error {
	g := &Google{}
	return g.interactionErrorHandler(&res)
}

func interactionDeltaText(res *model.GoogleInteractionRes) string {

	if res == nil {
		return ""
	}

	if res.Delta != nil && res.Delta.Text != "" {
		return res.Delta.Text
	}

	return strings.Join(interactionTextContents(res.Steps), "")
}

func interactionTextContents(steps []model.GoogleInteractionStep) []string {

	texts := make([]string, 0)
	for _, step := range steps {
		if step.Type != "" && step.Type != "model_output" && step.Type != "thought" {
			continue
		}
		for _, content := range step.Content {
			if content.Text != "" && (content.Type == "" || content.Type == "text" || content.Type == "thought") {
				texts = append(texts, content.Text)
			}
		}
	}

	return texts
}

func extractInteractionImages(steps []model.GoogleInteractionStep) []model.ImageResponseData {

	images := make([]model.ImageResponseData, 0)
	for _, step := range steps {
		for _, content := range step.Content {
			if !isImageContent(content) {
				continue
			}
			item := model.ImageResponseData{}
			if content.Uri != "" && (strings.HasPrefix(content.Uri, "http") || strings.HasPrefix(content.Uri, "/")) {
				item.Url = content.Uri
			}
			if content.Data != "" {
				if strings.HasPrefix(content.Data, "http") || strings.HasPrefix(content.Data, "/") {
					item.Url = content.Data
				} else {
					item.B64Json = content.Data
				}
			}
			if item.Url != "" || item.B64Json != "" {
				images = append(images, item)
			}
		}
	}

	return images
}

func extractInteractionVideo(steps []model.GoogleInteractionStep) (uri, data, mimeType string) {

	for _, step := range steps {
		for _, content := range step.Content {
			if !isVideoContent(content) {
				continue
			}
			if mimeType == "" {
				mimeType = content.MimeType
			}
			if content.Uri != "" {
				uri = content.Uri
			}
			if content.Data != "" {
				if strings.HasPrefix(content.Data, "http") || strings.Contains(content.Data, "/files/") {
					uri = content.Data
				} else {
					data = content.Data
				}
			}
			if uri != "" || data != "" {
				return uri, data, mimeType
			}
		}
	}

	return "", "", mimeType
}

func isImageContent(content model.GoogleInteractionContent) bool {
	if content.Type == "image" {
		return true
	}
	return strings.HasPrefix(content.MimeType, "image/")
}

func isVideoContent(content model.GoogleInteractionContent) bool {
	if content.Type == "video" {
		return true
	}
	return strings.HasPrefix(content.MimeType, "video/")
}

func convVideoJobStatus(status string) string {

	switch strings.ToLower(status) {
	case "completed", "complete", "succeeded", "success":
		return "completed"
	case "failed", "error":
		return "failed"
	case "cancelled", "canceled":
		return "cancelled"
	case "expired":
		return "expired"
	case "queued", "pending":
		return "queued"
	case "in_progress", "running", "live", "processing":
		return "in_progress"
	default:
		if status == "" {
			return "in_progress"
		}
		return status
	}
}

func interactionMediaPart(value any, defaultType string) map[string]any {

	switch item := value.(type) {
	case nil:
		return nil
	case string:
		if item == "" {
			return nil
		}
		if strings.HasPrefix(item, "http") || strings.HasPrefix(item, "files/") {
			part := map[string]any{"type": defaultType, "uri": item}
			if defaultType == "image" {
				part["mime_type"] = "image/jpeg"
			} else if defaultType == "video" {
				part["mime_type"] = "video/mp4"
			}
			return part
		}
		mimeType, data := common.GetMime(item)
		if data == "" {
			data = item
		}
		if defaultType == "image" && (mimeType == "" || !strings.HasPrefix(mimeType, "image/")) {
			mimeType = "image/jpeg"
		}
		if defaultType == "video" && (mimeType == "" || !strings.HasPrefix(mimeType, "video/")) {
			mimeType = "video/mp4"
		}
		return map[string]any{
			"type":      defaultType,
			"data":      data,
			"mime_type": mimeType,
		}
	case map[string]any:
		part := map[string]any{"type": defaultType}
		if mimeType := gconv.String(item["mime_type"]); mimeType != "" {
			part["mime_type"] = mimeType
		} else if mimeType = gconv.String(item["mimeType"]); mimeType != "" {
			part["mime_type"] = mimeType
		}
		if uri := firstNonEmpty(gconv.String(item["uri"]), gconv.String(item["url"]), gconv.String(item["image_url"]), gconv.String(item["video_url"])); uri != "" {
			part["uri"] = uri
		}
		if data := gconv.String(item["data"]); data != "" {
			part["data"] = data
		}
		if len(part) == 1 {
			return nil
		}
		if _, ok := part["mime_type"]; !ok {
			if defaultType == "image" {
				part["mime_type"] = "image/jpeg"
			} else if defaultType == "video" {
				part["mime_type"] = "video/mp4"
			}
		}
		return part
	default:
		return interactionMediaPart(gconv.String(item), defaultType)
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func fileHeaderToDataURL(header *multipart.FileHeader) (string, error) {

	if header == nil {
		return "", nil
	}

	src, err := header.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	content, err := io.ReadAll(src)
	if err != nil {
		return "", err
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	return "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(content), nil
}

func sizeToVideoAspectRatio(size string) string {

	size = strings.TrimSpace(strings.ToLower(size))
	if size == "" || size == "auto" {
		return ""
	}

	normalized := strings.ReplaceAll(size, "x", ":")
	switch normalized {
	case "9:16", "16:9":
		return normalized
	}

	parts := strings.FieldsFunc(size, func(r rune) bool {
		return r == 'x' || r == ':' || r == '*' || r == '×'
	})
	if len(parts) != 2 {
		return ""
	}

	width := gconv.Int(parts[0])
	height := gconv.Int(parts[1])
	if width <= 0 || height <= 0 {
		return ""
	}

	if width >= height {
		return "16:9"
	}

	return "9:16"
}

func forceInteractionBackground(data []byte) []byte {

	body := make(map[string]any)
	if err := json.Unmarshal(data, &body); err != nil {
		return data
	}

	body["background"] = true

	if format, ok := body["response_format"].(map[string]any); ok {
		if gconv.String(format["delivery"]) == "" {
			format["delivery"] = "uri"
			body["response_format"] = format
		}
	}

	return gjson.MustEncode(body)
}

func ParseInteractionChatCompletion(data []byte, modelName string, stream bool) (model.ChatCompletionResponse, error) {
	return convInteractionToChatCompletion(data, modelName, stream)
}

func ParseInteractionUsage(u *model.GoogleInteractionUsage) *model.Usage {
	return convInteractionUsage(u)
}

func (g *Google) ConvInteractionsResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvInteractionsResponse time: %d", gtime.TimestampMilli()-now)
	}()

	return convInteractionToChatCompletion(data, g.Model, false)
}

func (g *Google) ConvInteractionsStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvInteractionsStreamResponse time: %d", gtime.TimestampMilli()-now)
	}()

	return convInteractionToChatCompletion(data, g.Model, true)
}

func convInteractionImageRequest(data []byte, modelName string) (model.ImageGenerationRequest, bool) {

	req := model.GoogleInteractionReq{}
	if err := json.Unmarshal(data, &req); err != nil {
		return model.ImageGenerationRequest{}, false
	}

	if req.Input == nil && req.ResponseFormat == nil {
		return model.ImageGenerationRequest{}, false
	}

	request := model.ImageGenerationRequest{
		Model: modelName,
	}

	switch input := req.Input.(type) {
	case string:
		request.Prompt = input
	case []any:
		for _, item := range input {
			if content, ok := item.(map[string]any); ok {
				if gconv.String(content["type"]) == "text" || gconv.String(content["text"]) != "" {
					request.Prompt += gconv.String(content["text"])
				}
			}
		}
	}

	if format, ok := req.ResponseFormat.(map[string]any); ok {
		request.Quality = firstNonEmpty(gconv.String(format["image_size"]), gconv.String(format["quality"]))
		request.AspectRatio = gconv.String(format["aspect_ratio"])
	}

	return request, request.Prompt != "" || request.Quality != "" || request.AspectRatio != ""
}
