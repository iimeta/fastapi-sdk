package model

type OpenAIResponsesReq struct {
	Model  string `json:"model,omitempty"`
	Input  any    `json:"input"`
	Stream bool   `json:"stream,omitempty"`
}

type OpenAIResponsesRes struct {
	Model string          `json:"model"`
	Usage *ResponsesUsage `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
	ResponseBytes []byte `json:"-"`
	ConnTime      int64  `json:"-"`
	Duration      int64  `json:"-"`
	TotalTime     int64  `json:"-"`
	Err           error  `json:"-"`
}

type ResponsesUsage struct {
	InputTokens        int `json:"input_tokens"`
	InputTokensDetails struct {
		CachedTokens int `json:"cached_tokens"`
	} `json:"input_tokens_details"`
	OutputTokens        int `json:"output_tokens"`
	OutputTokensDetails struct {
		ReasoningTokens int `json:"reasoning_tokens"`
	} `json:"output_tokens_details"`
	TotalTokens int `json:"total_tokens"`
}
