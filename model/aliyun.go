package model

import "github.com/sashabaranov/go-openai"

type QwenChatCompletionMessage struct {
	User string `json:"user"`
	Bot  string `json:"bot"`
}
type QwenChatCompletionReq struct {
	Model      string `json:"model"`
	Input      Input  `json:"input"`
	Parameters struct {
	} `json:"parameters"`
}
type Input struct {
	Messages []openai.ChatCompletionMessage `json:"messages"`
}
type QwenChatCompletionRes struct {
	Output struct {
		FinishReason string `json:"finish_reason"`
		Text         string `json:"text"`
	} `json:"output"`
	Usage struct {
		OutputTokens int `json:"output_tokens"`
		InputTokens  int `json:"input_tokens"`
	} `json:"usage"`
	RequestId string `json:"request_id"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}
