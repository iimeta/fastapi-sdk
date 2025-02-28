package model

import "github.com/iimeta/go-openai"

type EmbeddingRequest struct {
	Input          any                            `json:"input"`
	Model          openai.EmbeddingModel          `json:"model"`
	User           string                         `json:"user"`
	EncodingFormat openai.EmbeddingEncodingFormat `json:"encoding_format,omitempty"`
	// Dimensions The number of dimensions the resulting output embeddings should have.
	// Only supported in text-embedding-3 and later models.
	Dimensions int `json:"dimensions,omitempty"`
}

type EmbeddingResponse struct {
	Object    string                `json:"object"`
	Data      []any                 `json:"data"`
	Model     openai.EmbeddingModel `json:"model"`
	Usage     *Usage                `json:"usage"`
	TotalTime int64                 `json:"-"`
}
