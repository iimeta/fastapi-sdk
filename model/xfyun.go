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
	Domain      string  `json:"domain"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float32 `json:"temperature,omitempty"`
	TopK        float32 `json:"top_k,omitempty"`
	ChatId      string  `json:"chat_id,omitempty"`
}

type Payload struct {
	// req
	Message   *Message   `json:"message"`
	Functions *Functions `json:"functions,omitempty"`
	// res
	Choices *Choices `json:"choices,omitempty"`
	Usage   *Usage   `json:"usage,omitempty"`
}

type Message struct {
	// req
	Text []openai.ChatCompletionMessage `json:"text"`
}

type Functions struct {
	// req
	Text []openai.FunctionDefinition `json:"text"`
}

type Text struct {
	// req res
	Role    string `json:"role"`
	Content string `json:"content"`

	// Choices
	Index        int                  `json:"index,omitempty"`
	ContentType  string               `json:"content_type"`
	FunctionCall *openai.FunctionCall `json:"function_call"`

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
