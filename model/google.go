package model

import "github.com/iimeta/go-openai"

type GoogleChatCompletionReq struct {
	Contents         []Content        `json:"contents"`
	GenerationConfig GenerationConfig `json:"generationConfig,omitempty"`
	Tools            any              `json:"tools,omitempty"`
}

type GoogleChatCompletionRes struct {
	Candidates    []Candidate    `json:"candidates"`
	UsageMetadata *UsageMetadata `json:"usageMetadata"`
	ModelVersion  string         `json:"modelVersion"`
	Error         struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
		Details []struct {
			Type     string `json:"@type"`
			Reason   string `json:"reason"`
			Domain   string `json:"domain"`
			Metadata struct {
				Service string `json:"service"`
			} `json:"metadata"`
		} `json:"details"`
	} `json:"error"`
	ResponseBytes []byte `json:"-"`
	ConnTime      int64  `json:"-"`
	Duration      int64  `json:"-"`
	TotalTime     int64  `json:"-"`
	Err           error  `json:"-"`
}

type Content struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

type Part struct {
	Text             string      `json:"text,omitempty"`
	InlineData       *InlineData `json:"inline_data,omitempty"`
	FileData         *FileData   `json:"file_data,omitempty"`
	FunctionCall     any         `json:"functionCall,omitempty"`
	FunctionResponse any         `json:"functionResponse,omitempty"`
}

type InlineData struct {
	MimeType string `json:"mime_type,omitempty"`
	Data     string `json:"data,omitempty"`
}

type FileData struct {
	FileUri  string `json:"file_uri,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
}

type Candidate struct {
	Content       Content             `json:"content"`
	FinishReason  openai.FinishReason `json:"finishReason"`
	Index         int                 `json:"index"`
	SafetyRatings []SafetyRating      `json:"safetyRatings"`
}

type SafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

type UsageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

type GenerationConfig struct {
	StopSequences   []string `json:"stopSequences,omitempty"`
	CandidateCount  int      `json:"candidateCount,omitempty"`
	MaxOutputTokens int      `json:"maxOutputTokens,omitempty"`
	Temperature     float32  `json:"temperature,omitempty"`
	TopP            float32  `json:"topP,omitempty"`
	TopK            int      `json:"topK,omitempty"`
}
