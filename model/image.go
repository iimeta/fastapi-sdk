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
	Async             bool   `json:"async,omitempty"`
	Image             any    `json:"image,omitempty"`
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

type ImageJobResponse struct {
	Id              string              `json:"id"`
	Object          string              `json:"object"`
	Model           string              `json:"model"`
	Status          string              `json:"status"`
	Progress        int                 `json:"progress"`
	CreatedAt       int64               `json:"created_at"`
	CompletedAt     *int64              `json:"completed_at"`
	ExpiresAt       *int64              `json:"expires_at"`
	Size            string              `json:"size"`
	Quality         string              `json:"quality"`
	N               int                 `json:"n"`
	Prompt          string              `json:"prompt"`
	OutputFormat    string              `json:"output_format"`
	Data            []ImageResponseData `json:"data,omitempty"`
	ImageUrl        string              `json:"image_url,omitempty"`
	Deleted         bool                `json:"deleted,omitempty"`
	Error           *ImageError         `json:"error"`
	Usage           *Usage              `json:"usage,omitempty"`
	ResponseBytes   []byte              `json:"-"`
	ResponseHeaders http.Header         `json:"-"`
	TotalTime       int64               `json:"-"`
}

type ImageError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ImageListRequest struct {
	After string `json:"after"`
	Limit int64  `json:"limit"`
	Order string `json:"order"`
}

type ImageListResponse struct {
	Object    string             `json:"object"`
	Data      []ImageJobResponse `json:"data"`
	FirstId   *string            `json:"first_id"`
	LastId    *string            `json:"last_id"`
	HasMore   bool               `json:"has_more"`
	TotalTime int64              `json:"-"`
}

type ImageRetrieveRequest struct {
	ImageId string `json:"image_id"`
}

type ImageDeleteRequest struct {
	ImageId string `json:"image_id"`
}

type ImageContentRequest struct {
	ImageId string `json:"image_id"`
	Index   int    `json:"index"`
}

type ImageContentResponse struct {
	Data      []byte `json:"-"`
	TotalTime int64  `json:"-"`
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
