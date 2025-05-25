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
	Service_tier       string                    `json:"service_tier,omitempty"`
	Store              bool                      `json:"store,omitempty"`
	Temperature        float32                   `json:"temperature,omitempty"`
	Text               any                       `json:"text,omitempty"`
	ToolChoice         string                    `json:"tool_choice,omitempty"`
	Tools              any                       `json:"tools,omitempty"`
	TopP               float32                   `json:"top_p,omitempty"`
	Truncation         string                    `json:"truncation,omitempty"`
	User               string                    `json:"user,omitempty"`
}

type OpenAIResponsesRes struct {
	Id                 string                  `json:"id"`
	Object             string                  `json:"object"`
	Model              string                  `json:"model"`
	CreatedAt          int64                   `json:"created_at"`
	Status             string                  `json:"status"`
	Background         bool                    `json:"background"`
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
	Temperature        float32                 `json:"temperature"`
	Text               any                     `json:"text"`
	ToolChoice         string                  `json:"tool_choice"`
	Tools              any                     `json:"tools"`
	TopP               float32                 `json:"top_p"`
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
	ResponseBytes  []byte `json:"-"`
	ConnTime       int64  `json:"-"`
	Duration       int64  `json:"-"`
	TotalTime      int64  `json:"-"`
	Err            error  `json:"-"`
	Type           string `json:"type"`
	SequenceNumber int    `json:"sequence_number"`
	Response       struct {
		Id                string      `json:"id"`
		Object            string      `json:"object"`
		CreatedAt         int         `json:"created_at"`
		Status            string      `json:"status"`
		Background        bool        `json:"background"`
		Error             interface{} `json:"error"`
		IncompleteDetails interface{} `json:"incomplete_details"`
		Instructions      interface{} `json:"instructions"`
		MaxOutputTokens   interface{} `json:"max_output_tokens"`
		Model             string      `json:"model"`
		Output            []struct {
			Id      string `json:"id"`
			Type    string `json:"type"`
			Status  string `json:"status"`
			Content []struct {
				Type        string        `json:"type"`
				Annotations []interface{} `json:"annotations"`
				Text        string        `json:"text"`
			} `json:"content"`
			Role string `json:"role"`
		} `json:"output"`
		ParallelToolCalls  bool        `json:"parallel_tool_calls"`
		PreviousResponseId interface{} `json:"previous_response_id"`
		Reasoning          struct {
			Effort  interface{} `json:"effort"`
			Summary interface{} `json:"summary"`
		} `json:"reasoning"`
		ServiceTier string  `json:"service_tier"`
		Store       bool    `json:"store"`
		Temperature float64 `json:"temperature"`
		Text        struct {
			Format struct {
				Type string `json:"type"`
			} `json:"format"`
		} `json:"text"`
		ToolChoice string        `json:"tool_choice"`
		Tools      []interface{} `json:"tools"`
		TopP       float64       `json:"top_p"`
		Truncation string        `json:"truncation"`
		Usage      *Usage        `json:"usage"`
		User       interface{}   `json:"user"`
		Metadata   struct {
		} `json:"metadata"`
	} `json:"response"`
	OutputIndex int `json:"output_index"`
	Item        struct {
		Id      string `json:"id"`
		Type    string `json:"type"`
		Status  string `json:"status"`
		Content []struct {
			Type        string        `json:"type"`
			Annotations []interface{} `json:"annotations"`
			Text        string        `json:"text"`
		} `json:"content"`
		Role string `json:"role"`
	} `json:"item"`
	ItemId       string `json:"item_id"`
	ContentIndex int    `json:"content_index"`
	Delta        string `json:"delta"`
	Part         struct {
		Type        string        `json:"type"`
		Annotations []interface{} `json:"annotations"`
		Text        string        `json:"text"`
	} `json:"part"`
}

type OpenAIResponsesInput struct {
	Role    string `json:"role"`
	Content []struct {
		Type     string `json:"type"`
		Text     string `json:"text,omitempty"`
		ImageUrl string `json:"image_url,omitempty"`
	} `json:"content"`
}

type OpenAIResponsesReasoning struct {
	Effort  string `json:"effort"`
	Summary string `json:"summary"`
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
