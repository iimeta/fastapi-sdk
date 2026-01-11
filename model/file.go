package model

import (
	"mime/multipart"
)

type FileUploadRequest struct {
	File         *multipart.FileHeader `json:"file"`
	Purpose      string                `json:"purpose"`
	ExpiresAfter ExpiresAfter          `json:"expires_after"`
	Model        string                `json:"-"`
}

type ExpiresAfter struct {
	Anchor  string `json:"anchor"`
	Seconds string `json:"seconds"`
}

type FileListRequest struct {
	After   string `json:"after"`
	Limit   int64  `json:"limit"`
	Order   string `json:"order"`
	Purpose string `json:"purpose"`
}

type FileListResponse struct {
	Object    string         `json:"object"`
	Data      []FileResponse `json:"data"`
	FirstId   *string        `json:"first_id"`
	LastId    *string        `json:"last_id"`
	HasMore   bool           `json:"has_more"`
	TotalTime int64          `json:"-"`
}

type FileRetrieveRequest struct {
	FileId string `json:"file_id"`
}

type FileDeleteRequest struct {
	FileId string `json:"file_id"`
}

type FileContentRequest struct {
	FileId string `json:"file_id"`
}

type FileContentResponse struct {
	Data      []byte `json:"-"`
	TotalTime int64  `json:"-"`
}

type FileResponse struct {
	Id            string  `json:"id"`
	Object        string  `json:"object"`
	Purpose       string  `json:"purpose"`
	Filename      string  `json:"filename"`
	Bytes         int     `json:"bytes"`
	CreatedAt     int64   `json:"created_at"`
	ExpiresAt     int64   `json:"expires_at"`
	Status        string  `json:"status"`
	StatusDetails *string `json:"status_details"`
	Deleted       bool    `json:"deleted,omitempty"`
	ResponseBytes []byte  `json:"-"`
	TotalTime     int64   `json:"-"`
}
