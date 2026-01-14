package model

import "mime/multipart"

type ImageGenerationRequest struct {
	Prompt            string `json:"prompt,omitempty"`
	Background        string `json:"background,omitempty"`
	Model             string `json:"model,omitempty"`
	Moderation        string `json:"moderation,omitempty"`
	N                 int    `json:"n,omitempty"`
	OutputCompression int    `json:"output_compression,omitempty"`
	OutputFormat      string `json:"output_format,omitempty"`
	Quality           string `json:"quality,omitempty"`
	ResponseFormat    string `json:"response_format,omitempty"`
	Size              string `json:"size,omitempty"`
	Style             string `json:"style,omitempty"`
	User              string `json:"user,omitempty"`
	AspectRatio       string `json:"aspect_ratio,omitempty"`
}

type ImageResponse struct {
	Created   int64               `json:"created,omitempty"`
	Data      []ImageResponseData `json:"data,omitempty"`
	Usage     Usage               `json:"usage,omitempty"`
	TotalTime int64               `json:"-"`
}

type ImageResponseData struct {
	Url           string `json:"url,omitempty"`
	B64Json       string `json:"b64_json,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

type ImageEditRequest struct {
	Image          []*multipart.FileHeader `json:"image,omitempty"`
	Prompt         string                  `json:"prompt,omitempty"`
	Background     string                  `json:"background,omitempty"`
	Mask           *multipart.FileHeader   `json:"mask,omitempty"`
	Model          string                  `json:"model,omitempty"`
	N              int                     `json:"n,omitempty"`
	Quality        string                  `json:"quality,omitempty"`
	ResponseFormat string                  `json:"response_format,omitempty"`
	Size           string                  `json:"size,omitempty"`
	User           string                  `json:"user,omitempty"`
	AspectRatio    string                  `json:"aspect_ratio,omitempty"`
}
