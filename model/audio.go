package model

import (
	"mime/multipart"
)

type SpeechRequest struct {
	Model          string  `json:"model"`
	Input          string  `json:"input"`
	Voice          string  `json:"voice"`
	ResponseFormat string  `json:"response_format,omitempty"`
	Speed          float64 `json:"speed,omitempty"`
}

type SpeechResponse struct {
	Data      []byte
	TotalTime int64 `json:"-"`
}

type AudioRequest struct {
	Model                  string                `json:"model"`
	File                   *multipart.FileHeader `json:"file"`
	Prompt                 string                `json:"prompt"`
	Temperature            float32               `json:"temperature"`
	Language               string                `json:"language"`
	Format                 string                `json:"response_format"`
	TimestampGranularities []string              `json:"timestamp_granularities"`
}

type AudioResponse struct {
	Task     string    `json:"task,omitempty"`
	Language string    `json:"language,omitempty"`
	Duration float64   `json:"duration,omitempty"`
	Segments []Segment `json:"segments,omitempty"`
	Words    []Word    `json:"words,omitempty"`
	Text     string    `json:"text,omitempty"`
	Usage    struct {
		Type    string `json:"type"`
		Seconds int    `json:"seconds"`
	} `json:"usage"`
	TotalTime int64 `json:"-"`
}

type Segment struct {
	Id               int     `json:"id"`
	Seek             int     `json:"seek"`
	Start            float64 `json:"start"`
	End              float64 `json:"end"`
	Text             string  `json:"text"`
	Tokens           []int   `json:"tokens"`
	Temperature      float64 `json:"temperature"`
	AvgLogprob       float64 `json:"avg_logprob"`
	CompressionRatio float64 `json:"compression_ratio"`
	NoSpeechProb     float64 `json:"no_speech_prob"`
	Transient        bool    `json:"transient"`
}

type Word struct {
	Word  string  `json:"word"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}
