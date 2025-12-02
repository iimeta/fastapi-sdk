package model

import (
	"mime/multipart"
)

type VideoRequest struct {
	Model          string                `json:"model"`
	Prompt         string                `json:"prompt"`
	InputReference *multipart.FileHeader `json:"input_reference"`
	Seconds        string                `json:"seconds"`
	Size           string                `json:"size"`
}

type VideoResponse struct {
	Id                 string                `json:"id"`
	Object             string                `json:"object"`
	Model              string                `json:"model"`
	Status             string                `json:"status"`
	Progress           int                   `json:"progress"`
	CreatedAt          int                   `json:"created_at"`
	CompletedAt        int                   `json:"completed_at"`
	ExpiresAt          int                   `json:"expires_at"`
	Size               string                `json:"size"`
	Prompt             string                `json:"prompt"`
	Seconds            string                `json:"seconds"`
	Quality            string                `json:"quality"`
	RemixedFromVideoId string                `json:"remixed_from_video_id"`
	Error              *OpenAIResponsesError `json:"error"`
	TotalTime          int64                 `json:"-"`
}
