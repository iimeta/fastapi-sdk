package model

import "github.com/sashabaranov/go-openai"

type ChatCompletionStreamResponse struct {
	ID                string                              `json:"id"`
	Object            string                              `json:"object"`
	Created           int64                               `json:"created"`
	Model             string                              `json:"model"`
	Choices           []openai.ChatCompletionStreamChoice `json:"choices"`
	PromptAnnotations []openai.PromptAnnotation           `json:"prompt_annotations,omitempty"`
	Usage             openai.Usage                        `json:"usage"`
}
