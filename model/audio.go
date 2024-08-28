package model

import (
	"github.com/iimeta/go-openai"
	"io"
)

type SpeechRequest struct {
	Model          openai.SpeechModel          `json:"model"`
	Input          string                      `json:"input"`
	Voice          openai.SpeechVoice          `json:"voice"`
	ResponseFormat openai.SpeechResponseFormat `json:"response_format,omitempty"` // Optional, default to mp3
	Speed          float64                     `json:"speed,omitempty"`           // Optional, default to 1.0
}

type SpeechResponse struct {
	io.ReadCloser
	TotalTime int64 `json:"-"`
}

// AudioRequest represents a request structure for audio API.
type AudioRequest struct {
	Model string `json:"model"`

	// FilePath is either an existing file in your filesystem or a filename representing the contents of Reader.
	FilePath string

	// Reader is an optional io.Reader when you do not want to use an existing file.
	Reader io.Reader

	Prompt                 string                                     `json:"prompt"`
	Temperature            float32                                    `json:"temperature"`
	Language               string                                     `json:"language"` // Only for transcription.
	Format                 openai.AudioResponseFormat                 `json:"response_format"`
	TimestampGranularities []openai.TranscriptionTimestampGranularity `json:"timestamp_granularities"` // Only for transcription.
}

// AudioResponse represents a response structure for audio API.
type AudioResponse struct {
	Task     string  `json:"task"`
	Language string  `json:"language"`
	Duration float64 `json:"duration"`
	Segments []struct {
		ID               int     `json:"id"`
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
	} `json:"segments"`
	Words []struct {
		Word  string  `json:"word"`
		Start float64 `json:"start"`
		End   float64 `json:"end"`
	} `json:"words"`
	Text      string `json:"text"`
	TotalTime int64  `json:"-"`
}
