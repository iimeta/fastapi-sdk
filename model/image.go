package model

import "github.com/sashabaranov/go-openai"

type ImageResponse struct {
	Created   int64                           `json:"created,omitempty"`
	Data      []openai.ImageResponseDataInner `json:"data,omitempty"`
	TotalTime int64                           `json:"-"`
}
