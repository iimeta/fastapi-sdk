package model

import "github.com/sashabaranov/go-openai"

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	Model            string                               `json:"model"`
	Messages         []openai.ChatCompletionMessage       `json:"messages"`
	MaxTokens        int                                  `json:"max_tokens,omitempty"`
	Temperature      float32                              `json:"temperature,omitempty"`
	TopP             float32                              `json:"top_p,omitempty"`
	N                int                                  `json:"n,omitempty"`
	Stream           bool                                 `json:"stream,omitempty"`
	Stop             []string                             `json:"stop,omitempty"`
	PresencePenalty  float32                              `json:"presence_penalty,omitempty"`
	ResponseFormat   *openai.ChatCompletionResponseFormat `json:"response_format,omitempty"`
	Seed             *int                                 `json:"seed,omitempty"`
	FrequencyPenalty float32                              `json:"frequency_penalty,omitempty"`
	// LogitBias is must be a token id string (specified by their token ID in the tokenizer), not a word string.
	// incorrect: `"logit_bias":{"You": 6}`, correct: `"logit_bias":{"1639": 6}`
	// refs: https://platform.openai.com/docs/api-reference/chat/create#chat/create-logit_bias
	LogitBias map[string]int `json:"logit_bias,omitempty"`
	// LogProbs indicates whether to return log probabilities of the output tokens or not.
	// If true, returns the log probabilities of each output token returned in the content of message.
	// This option is currently not available on the gpt-4-vision-preview model.
	LogProbs bool `json:"logprobs,omitempty"`
	// TopLogProbs is an integer between 0 and 5 specifying the number of most likely tokens to return at each
	// token position, each with an associated log probability.
	// logprobs must be set to true if this parameter is used.
	TopLogProbs  int                         `json:"top_logprobs,omitempty"`
	User         string                      `json:"user,omitempty"`
	Functions    []openai.FunctionDefinition `json:"functions,omitempty"`
	FunctionCall any                         `json:"function_call,omitempty"`
	Tools        []openai.Tool               `json:"tools,omitempty"`
	// This can be either a string or an ToolChoice object.
	ToolChoice any `json:"tool_choice,omitempty"`
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	ID                string                    `json:"id"`
	Object            string                    `json:"object"`
	Created           int64                     `json:"created"`
	Model             string                    `json:"model"`
	Choices           []ChatCompletionChoice    `json:"choices"`
	Usage             *openai.Usage             `json:"usage"`
	SystemFingerprint string                    `json:"system_fingerprint,omitempty"`
	PromptAnnotations []openai.PromptAnnotation `json:"prompt_annotations,omitempty"`
	ConnTime          int64                     `json:"-"`
	Duration          int64                     `json:"-"`
	TotalTime         int64                     `json:"-"`
	Error             error                     `json:"-"`
}

type ChatCompletionChoice struct {
	Index                int                                     `json:"index"`
	Message              *openai.ChatCompletionMessage           `json:"message,omitempty"`
	Delta                *openai.ChatCompletionStreamChoiceDelta `json:"delta,omitempty"`
	LogProbs             *openai.LogProbs                        `json:"logprobs,omitempty"`
	FinishReason         openai.FinishReason                     `json:"finish_reason"`
	ContentFilterResults *openai.ContentFilterResults            `json:"content_filter_results,omitempty"`
}
