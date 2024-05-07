package model

// ImageRequest represents the request structure for the image API.
type ImageRequest struct {
	Prompt         string `json:"prompt,omitempty"`
	Model          string `json:"model,omitempty"`
	N              int    `json:"n,omitempty"`
	Quality        string `json:"quality,omitempty"`
	Size           string `json:"size,omitempty"`
	Style          string `json:"style,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

type ImageResponse struct {
	Created   int64                    `json:"created,omitempty"`
	Data      []ImageResponseDataInner `json:"data,omitempty"`
	TotalTime int64                    `json:"-"`
}

// ImageResponseDataInner represents a response data structure for image API.
type ImageResponseDataInner struct {
	URL           string `json:"url,omitempty"`
	B64JSON       string `json:"b64_json,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}
