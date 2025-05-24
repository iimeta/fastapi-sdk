package model

type OpenAIResponsesReq struct {
	Model              string `json:"model,omitempty"`
	Input              any    `json:"input"`
	Stream             bool   `json:"stream,omitempty"`
	Background         bool   `json:"background,omitempty"`
	Include            any    `json:"include,omitempty"`
	Instructions       string `json:"instructions,omitempty"`
	MaxOutputTokens    int    `json:"max_output_tokens,omitempty"`
	Metadata           any    `json:"metadata,omitempty"`
	ParallelToolCalls  bool   `json:"parallel_tool_calls,omitempty"`
	PreviousResponseId string `json:"previous_response_id,omitempty"`
	Reasoning          any    `json:"reasoning,omitempty"`
	Service_tier       string `json:"service_tier,omitempty"`
	Store              bool   `json:"store,omitempty"`
	Temperature        int    `json:"temperature,omitempty"`
	Text               any    `json:"text,omitempty"`
	ToolChoice         string `json:"tool_choice,omitempty"`
	Tools              any    `json:"tools,omitempty"`
	TopP               int    `json:"top_p,omitempty"`
	Truncation         string `json:"truncation,omitempty"`
	User               string `json:"user,omitempty"`
}

type OpenAIResponsesRes struct {
	Id                 string                  `json:"id"`
	Object             string                  `json:"object"`
	Model              string                  `json:"model"`
	CreatedAt          int                     `json:"created_at"`
	Status             string                  `json:"status"`
	Background         string                  `json:"background"`
	IncompleteDetails  any                     `json:"incomplete_details"`
	Instructions       any                     `json:"instructions"`
	MaxOutputTokens    int                     `json:"max_output_tokens"`
	Metadata           any                     `json:"metadata"`
	Output             []OpenAIResponsesOutput `json:"output"`
	ParallelToolCalls  bool                    `json:"parallel_tool_calls"`
	PreviousResponseId string                  `json:"previous_response_id"`
	Reasoning          any                     `json:"reasoning"`
	ServiceTier        string                  `json:"service_tier"`
	Store              bool                    `json:"store"`
	Temperature        int                     `json:"temperature"`
	Text               any                     `json:"text"`
	ToolChoice         string                  `json:"tool_choice"`
	Tools              any                     `json:"tools"`
	TopP               int                     `json:"top_p"`
	Truncation         string                  `json:"truncation"`
	User               string                  `json:"user"`
	Usage              *Usage                  `json:"usage"`
	Error              *struct {
		Code           string `json:"code"`
		Message        string `json:"message"`
		Param          string `json:"param"`
		SequenceNumber int    `json:"sequence_number"`
		Type           string `json:"type"`
	} `json:"error"`
	ResponseBytes []byte `json:"-"`
	ConnTime      int64  `json:"-"`
	Duration      int64  `json:"-"`
	TotalTime     int64  `json:"-"`
	Err           error  `json:"-"`
}

type OpenAIResponsesOutput struct {
	Type    string                   `json:"type"`
	Id      string                   `json:"id"`
	Status  string                   `json:"status"`
	Role    string                   `json:"role"`
	Content []OpenAIResponsesContent `json:"content"`
}

type OpenAIResponsesContent struct {
	Type        string `json:"type"`
	Text        string `json:"text"`
	Annotations []any  `json:"annotations"`
}
