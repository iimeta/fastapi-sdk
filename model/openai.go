package model

type OpenAIResponsesReq struct {
	Model              string                    `json:"model,omitempty"`
	Input              any                       `json:"input"`
	Stream             bool                      `json:"stream,omitempty"`
	Background         bool                      `json:"background,omitempty"`
	Include            any                       `json:"include,omitempty"`
	Instructions       string                    `json:"instructions,omitempty"`
	MaxOutputTokens    int                       `json:"max_output_tokens,omitempty"`
	Metadata           map[string]string         `json:"metadata,omitempty"`
	ParallelToolCalls  bool                      `json:"parallel_tool_calls,omitempty"`
	PreviousResponseId string                    `json:"previous_response_id,omitempty"`
	Reasoning          *OpenAIResponsesReasoning `json:"reasoning,omitempty"`
	ServiceTier        string                    `json:"service_tier,omitempty"`
	Store              bool                      `json:"store,omitempty"`
	Temperature        float32                   `json:"temperature,omitempty"`
	Text               any                       `json:"text,omitempty"`
	Tools              any                       `json:"tools,omitempty"`
	ToolChoice         any                       `json:"tool_choice,omitempty"`
	TopP               float32                   `json:"top_p,omitempty"`
	Truncation         string                    `json:"truncation,omitempty"`
	User               string                    `json:"user,omitempty"`
}

type OpenAIResponsesRes struct {
	Id                 string                   `json:"id"`
	Object             string                   `json:"object"`
	Model              string                   `json:"model"`
	CreatedAt          int64                    `json:"created_at"`
	Status             string                   `json:"status"`
	Background         bool                     `json:"background"`
	IncompleteDetails  any                      `json:"incomplete_details"`
	Instructions       any                      `json:"instructions"`
	MaxOutputTokens    int                      `json:"max_output_tokens"`
	Metadata           map[string]string        `json:"metadata"`
	Output             []OpenAIResponsesOutput  `json:"output"`
	ParallelToolCalls  bool                     `json:"parallel_tool_calls"`
	PreviousResponseId string                   `json:"previous_response_id"`
	Reasoning          OpenAIResponsesReasoning `json:"reasoning"`
	ServiceTier        string                   `json:"service_tier"`
	Store              bool                     `json:"store"`
	Temperature        float32                  `json:"temperature"`
	Text               OpenAIResponsesText      `json:"text"`
	Tools              any                      `json:"tools"`
	ToolChoice         string                   `json:"tool_choice"`
	TopP               float32                  `json:"top_p"`
	Truncation         string                   `json:"truncation"`
	User               string                   `json:"user"`
	Usage              *Usage                   `json:"usage"`
	Error              *OpenAIResponsesError    `json:"error"`
	ResponseBytes      []byte                   `json:"-"`
	ConnTime           int64                    `json:"-"`
	Duration           int64                    `json:"-"`
	TotalTime          int64                    `json:"-"`
	Err                error                    `json:"-"`
}

type OpenAIResponsesStreamRes struct {
	Type           string                  `json:"type"`
	SequenceNumber int                     `json:"sequence_number"`
	Response       OpenAIResponsesResponse `json:"response"`
	OutputIndex    int                     `json:"output_index"`
	ContentIndex   int                     `json:"content_index"`
	ItemId         string                  `json:"item_id"`
	Item           OpenAIResponsesItem     `json:"item"`
	Delta          string                  `json:"delta"`
	Part           OpenAIResponsesPart     `json:"part"`
	Arguments      string                  `json:"arguments"`
	ResponseBytes  []byte                  `json:"-"`
	ConnTime       int64                   `json:"-"`
	Duration       int64                   `json:"-"`
	TotalTime      int64                   `json:"-"`
	Err            error                   `json:"-"`
}

type OpenAIResponsesResponse struct {
	Id                 string                   `json:"id"`
	Object             string                   `json:"object"`
	CreatedAt          int64                    `json:"created_at"`
	Status             string                   `json:"status"`
	Background         bool                     `json:"background"`
	IncompleteDetails  any                      `json:"incomplete_details"`
	Instructions       any                      `json:"instructions"`
	MaxOutputTokens    int                      `json:"max_output_tokens"`
	Model              string                   `json:"model"`
	Output             []OpenAIResponsesOutput  `json:"output"`
	ParallelToolCalls  bool                     `json:"parallel_tool_calls"`
	PreviousResponseId string                   `json:"previous_response_id"`
	Reasoning          OpenAIResponsesReasoning `json:"reasoning"`
	ServiceTier        string                   `json:"service_tier"`
	Store              bool                     `json:"store"`
	Temperature        float32                  `json:"temperature"`
	Text               OpenAIResponsesText      `json:"text"`
	ToolChoice         string                   `json:"tool_choice"`
	Tools              any                      `json:"tools"`
	TopP               float32                  `json:"top_p"`
	Truncation         string                   `json:"truncation"`
	User               string                   `json:"user"`
	Metadata           map[string]string        `json:"metadata"`
	Usage              *Usage                   `json:"usage"`
	Error              *OpenAIResponsesError    `json:"error"`
}

type OpenAIResponsesInput struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type OpenAIResponsesReasoning struct {
	Effort  string `json:"effort"`
	Summary string `json:"summary"`
}

type OpenAIResponsesOutput struct {
	Type      string                   `json:"type"`
	Id        string                   `json:"id"`
	Status    string                   `json:"status,omitempty"`
	Role      string                   `json:"role,omitempty"`
	Content   []OpenAIResponsesContent `json:"content,omitempty"`
	Summary   []OpenAIResponsesSummary `json:"summary,omitempty"`
	Arguments string                   `json:"arguments,omitempty"`
	CallId    string                   `json:"call_id,omitempty"`
	Name      string                   `json:"name,omitempty"`
}

type OpenAIResponsesContent struct {
	Type        string `json:"type"`
	Text        string `json:"text,omitempty"`
	ImageUrl    string `json:"image_url,omitempty"`
	Annotations []any  `json:"annotations,omitempty"`
}

type OpenAIResponsesItem struct {
	Id        string                   `json:"id"`
	Type      string                   `json:"type"`
	Status    string                   `json:"status"`
	Content   []OpenAIResponsesContent `json:"content"`
	Role      string                   `json:"role"`
	Arguments string                   `json:"arguments"`
	CallId    string                   `json:"call_id"`
	Name      string                   `json:"name"`
	Summary   []OpenAIResponsesSummary `json:"summary"`
}

type OpenAIResponsesPart struct {
	Type        string `json:"type"`
	Annotations []any  `json:"annotations"`
	Text        string `json:"text"`
}

type OpenAIResponsesText struct {
	Format OpenAIResponsesFormat `json:"format"`
}

type OpenAIResponsesFormat struct {
	Type string `json:"type"`
}

type OpenAIResponsesSummary struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type OpenAIResponsesError struct {
	Code           string `json:"code"`
	Message        string `json:"message"`
	Param          string `json:"param"`
	SequenceNumber int    `json:"sequence_number"`
	Type           string `json:"type"`
}
