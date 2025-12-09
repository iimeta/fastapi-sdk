package model

import (
	"mime/multipart"
)

type VideoCreateRequest struct {
	Model          string                `json:"model"`
	Prompt         string                `json:"prompt"`
	InputReference *multipart.FileHeader `json:"input_reference"`
	Seconds        string                `json:"seconds"`
	Size           string                `json:"size"`
}

type VideoRemixRequest struct {
	VideoId string `json:"video_id"`
	Prompt  string `json:"prompt"`
}

type VideoListRequest struct {
	After string `json:"after"`
	Limit int    `json:"limit"`
	Order string `json:"order"`
}

type VideoListResponse struct {
	Data      []VideoJobResponse `json:"data"`
	Object    string             `json:"object"`
	TotalTime int64              `json:"-"`
}

type VideoRetrieveRequest struct {
	VideoId string `json:"video_id"`
}

type VideoDeleteRequest struct {
	VideoId string `json:"video_id"`
}

type VideoContentRequest struct {
	VideoId string `json:"video_id"`
	Variant string `json:"variant"`
}

type VideoContentResponse struct {
	Data      []byte
	TotalTime int64 `json:"-"`
}

type VideoJobResponse struct {
	Id                 string      `json:"id"`
	Object             string      `json:"object"`
	Model              string      `json:"model"`
	Status             string      `json:"status"`
	Progress           int         `json:"progress"`
	CreatedAt          int         `json:"created_at"`
	CompletedAt        int         `json:"completed_at"`
	ExpiresAt          int         `json:"expires_at"`
	Size               string      `json:"size"`
	Prompt             string      `json:"prompt"`
	Seconds            string      `json:"seconds"`
	Quality            string      `json:"quality"`
	RemixedFromVideoId string      `json:"remixed_from_video_id"`
	Error              *VideoError `json:"error"`
	TotalTime          int64       `json:"-"`
}

type VideoError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
