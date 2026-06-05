package model

import (
	"mime/multipart"
	"net/http"
)

type ImageGenerationRequest struct {
	Prompt            string `json:"prompt,omitempty"`
	Background        string `json:"background,omitempty"`
	Model             string `json:"model,omitempty"`
	Moderation        string `json:"moderation,omitempty"`
	InputFidelity     string `json:"input_fidelity,omitempty"`
	N                 int    `json:"n,omitempty"`
	OutputCompression int    `json:"output_compression,omitempty"`
	OutputFormat      string `json:"output_format,omitempty"`
	PartialImages     int    `json:"partial_images,omitempty"`
	Quality           string `json:"quality,omitempty"`
	ResponseFormat    string `json:"response_format,omitempty"`
	Size              string `json:"size,omitempty"`
	Style             string `json:"style,omitempty"`
	User              string `json:"user,omitempty"`
	AspectRatio       string `json:"aspect_ratio,omitempty"`
	Stream            bool   `json:"stream,omitempty"`
}

type ImageResponse struct {
	Created         int64               `json:"created,omitempty"`
	Data            []ImageResponseData `json:"data,omitempty"`
	Usage           Usage               `json:"usage,omitempty"`
	ResponseBytes   []byte              `json:"-"`
	ResponseHeaders http.Header         `json:"-"` // 响应头
	Event           string              `json:"-"` // 流式事件类型
	ConnTime        int64               `json:"-"`
	Duration        int64               `json:"-"`
	TotalTime       int64               `json:"-"`
	Error           error               `json:"-"`
}

type ImageResponseData struct {
	Url           string `json:"url,omitempty"`
	B64Json       string `json:"b64_json,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

type ImageStreamResponse struct {
	CreatedAt         int64  `json:"created_at,omitempty"`
	Type              string `json:"type,omitempty"`
	B64Json           string `json:"b64_json,omitempty"`
	Background        string `json:"background,omitempty"`
	OutputFormat      string `json:"output_format,omitempty"`
	PartialImageIndex int    `json:"partial_image_index,omitempty"`
	Quality           string `json:"quality,omitempty"`
	SequenceNumber    int    `json:"sequence_number,omitempty"`
	Size              string `json:"size,omitempty"`
	Usage             Usage  `json:"usage,omitempty"`
}

type ImageEditRequest struct {
	Image          []*multipart.FileHeader `json:"image,omitempty"`
	Prompt         string                  `json:"prompt,omitempty"`
	Background     string                  `json:"background,omitempty"`
	Mask           *multipart.FileHeader   `json:"mask,omitempty"`
	Model          string                  `json:"model,omitempty"`
	InputFidelity  string                  `json:"input_fidelity,omitempty"`
	N              int                     `json:"n,omitempty"`
	PartialImages  int                     `json:"partial_images,omitempty"`
	Quality        string                  `json:"quality,omitempty"`
	ResponseFormat string                  `json:"response_format,omitempty"`
	Size           string                  `json:"size,omitempty"`
	User           string                  `json:"user,omitempty"`
	AspectRatio    string                  `json:"aspect_ratio,omitempty"`
	Stream         bool                    `json:"stream,omitempty"`
}
