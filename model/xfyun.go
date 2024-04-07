package model

import "github.com/sashabaranov/go-openai"

type Header struct {
	// req
	AppId string `json:"app_id"`
	Uid   string `json:"uid"`
	// res
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Sid     string `json:"sid,omitempty"`
	Status  int    `json:"status,omitempty"`
}

type Parameter struct {
	// req
	Chat *Chat `json:"chat"`
}

type Chat struct {
	// req
	Domain          string `json:"domain"`
	RandomThreshold int    `json:"random_threshold"`
	MaxTokens       int    `json:"max_tokens"`
}

type Payload struct {
	// req
	Message *Message `json:"message"`
	// res
	Choices *Choices `json:"choices,omitempty"`
	Usage   *Usage   `json:"usage,omitempty"`
}

type Message struct {
	// req
	Text []openai.ChatCompletionMessage `json:"text"`
}

type Text struct {
	// req res
	Role    string `json:"role"`
	Content string `json:"content"`

	// Choices
	Index int `json:"index,omitempty"`

	// Usage
	QuestionTokens   int `json:"question_tokens,omitempty"`
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

type Choices struct {
	// res
	Status int    `json:"status,omitempty"`
	Seq    int    `json:"seq,omitempty"`
	Text   []Text `json:"text,omitempty"`
}

type Usage struct {
	// res
	Text *Text `json:"text,omitempty"`
}

type SparkReq struct {
	Header    Header    `json:"header"`
	Parameter Parameter `json:"parameter"`
	Payload   Payload   `json:"payload"`
}
type SparkRes struct {
	Content string  `json:"content"`
	Header  Header  `json:"header"`
	Payload Payload `json:"payload"`
}
