package model

type AnthropicChatCompletionReq struct {
	Model            string                  `json:"model,omitempty"`
	Messages         []ChatCompletionMessage `json:"messages"`
	MaxTokens        int                     `json:"max_tokens,omitempty"`
	Metadata         *Metadata               `json:"metadata,omitempty"`
	StopSequences    []string                `json:"stop_sequences,omitempty"`
	Stream           bool                    `json:"stream,omitempty"`
	System           any                     `json:"system,omitempty"`
	Temperature      float32                 `json:"temperature,omitempty"`
	ToolChoice       any                     `json:"tool_choice,omitempty"`
	Tools            []AnthropicTool         `json:"tools,omitempty"`
	TopK             int                     `json:"top_k,omitempty"`
	TopP             float32                 `json:"top_p,omitempty"`
	AnthropicVersion string                  `json:"anthropic_version,omitempty"`
}

type AnthropicChatCompletionRes struct {
	Id           string             `json:"id"`
	Type         string             `json:"type"`
	Role         string             `json:"role"`
	Content      []AnthropicContent `json:"content"`
	Model        string             `json:"model"`
	StopReason   string             `json:"stop_reason"`
	StopSequence string             `json:"stop_sequence"`
	Message      AnthropicMessage   `json:"message"`
	Index        int                `json:"index"`
	Delta        AnthropicContent   `json:"delta"`
	Usage        *AnthropicUsage    `json:"usage,omitempty"`
	Error        *AnthropicError    `json:"error,omitempty"`
}

type Metadata struct {
	UserId string `json:"user_id,omitempty"`
}

type Source struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

type AnthropicMessage struct {
	Id           string          `json:"id"`
	Type         string          `json:"type"`
	Role         string          `json:"role"`
	Model        string          `json:"model"`
	Content      []interface{}   `json:"content"`
	StopReason   interface{}     `json:"stop_reason"`
	StopSequence interface{}     `json:"stop_sequence"`
	Usage        *AnthropicUsage `json:"usage"`
}

type AnthropicContent struct {
	Type         string       `json:"type"`
	Text         string       `json:"text"`
	PartialJson  string       `json:"partial_json"`
	ContentBlock ContentBlock `json:"content_block,omitempty"`
	StopReason   string       `json:"stop_reason,omitempty"`
	StopSequence string       `json:"stop_sequence,omitempty"`
}

type ContentBlock struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Id    string `json:"id"`
	Name  string `json:"name"`
	Input any    `json:"input"`
}

type AnthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type AnthropicError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type AnthropicErrorResponse struct {
	Error *AnthropicError `json:"error,omitempty"`
}

type AnthropicTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"input_schema"`
}
