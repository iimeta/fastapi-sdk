package model

import "github.com/sashabaranov/go-openai"

type ChatCompletionResponse struct {
	ID                string                        `json:"id"`
	Object            string                        `json:"object"`
	Created           int64                         `json:"created"`
	Model             string                        `json:"model"`
	Choices           []openai.ChatCompletionChoice `json:"choices"`
	Usage             openai.Usage                  `json:"usage"`
	SystemFingerprint string                        `json:"system_fingerprint"`
	TotalTime         int64                         `json:"-"`
}

type ChatCompletionStreamResponse struct {
	ID                string                       `json:"id"`
	Object            string                       `json:"object"`
	Created           int64                        `json:"created"`
	Model             string                       `json:"model"`
	SystemFingerprint interface{}                  `json:"system_fingerprint"`
	Choices           []ChatCompletionStreamChoice `json:"choices"`
	PromptAnnotations []openai.PromptAnnotation    `json:"prompt_annotations,omitempty"`
	ConnTime          int64                        `json:"-"`
	Duration          int64                        `json:"-"`
	TotalTime         int64                        `json:"-"`
}

type ChatCompletionStreamChoice struct {
	Index        int                                    `json:"index"`
	Delta        openai.ChatCompletionStreamChoiceDelta `json:"delta"`
	LogProbs     *openai.LogProbs                       `json:"logprobs"`
	FinishReason openai.FinishReason                    `json:"finish_reason"`
	//ContentFilterResults openai.ContentFilterResults            `json:"content_filter_results,omitempty"`
}
