package model

type AnthropicChatCompletionReq struct {
	Model     string                  `json:"model"`
	Messages  []ChatCompletionMessage `json:"messages"`
	MaxTokens int                     `json:"max_tokens,omitempty"`
	Metadata  struct {
		UserId string `json:"user_id,omitempty"`
	} `json:"metadata,omitempty"`
	StopSequences []string        `json:"stop_sequences,omitempty"`
	Stream        bool            `json:"stream,omitempty"`
	System        string          `json:"system,omitempty"`
	Temperature   float32         `json:"temperature,omitempty"`
	ToolChoice    any             `json:"tool_choice,omitempty"`
	Tools         []AnthropicTool `json:"tools,omitempty"`
	TopK          int             `json:"top_k,omitempty"`
	TopP          float32         `json:"top_p,omitempty"`
}

type AnthropicChatCompletionRes struct {
	Id           string           `json:"id"`
	Type         string           `json:"type"`
	Role         string           `json:"role"`
	Content      AnthropicContent `json:"content"`
	Model        string           `json:"model"`
	StopReason   string           `json:"stop_reason"`
	StopSequence string           `json:"stop_sequence"`
	Usage        *Usage           `json:"usage"`
	Error        AnthropicError   `json:"error"`
}

type AnthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type AnthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type AnthropicError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type AnthropicErrorResponse struct {
	Error *AnthropicError `json:"error,omitempty"`
}

type AnthropicTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"input_schema"`
}
