package model

type ChatCompletionRequest struct {
	Model               string                        `json:"model"`
	Messages            []ChatCompletionMessage       `json:"messages"`
	MaxTokens           int                           `json:"max_tokens,omitempty"`
	MaxCompletionTokens int                           `json:"max_completion_tokens,omitempty"`
	Temperature         float32                       `json:"temperature,omitempty"`
	TopP                float32                       `json:"top_p,omitempty"`
	TopK                int                           `json:"top_k,omitempty"`
	N                   int                           `json:"n,omitempty"`
	Stream              bool                          `json:"stream,omitempty"`
	Stop                []string                      `json:"stop,omitempty"`
	PresencePenalty     float32                       `json:"presence_penalty,omitempty"`
	ResponseFormat      *ChatCompletionResponseFormat `json:"response_format,omitempty"`
	Seed                *int                          `json:"seed,omitempty"`
	FrequencyPenalty    float32                       `json:"frequency_penalty,omitempty"`
	LogitBias           map[string]int                `json:"logit_bias,omitempty"`
	LogProbs            bool                          `json:"logprobs,omitempty"`
	TopLogProbs         int                           `json:"top_logprobs,omitempty"`
	User                string                        `json:"user,omitempty"`
	Functions           []FunctionDefinition          `json:"functions,omitempty"`
	FunctionCall        any                           `json:"function_call,omitempty"`
	Tools               any                           `json:"tools,omitempty"`
	ToolChoice          any                           `json:"tool_choice,omitempty"`
	StreamOptions       *StreamOptions                `json:"stream_options,omitempty"`
	ParallelToolCalls   any                           `json:"parallel_tool_calls,omitempty"`
	Store               bool                          `json:"store,omitempty"`
	Metadata            map[string]string             `json:"metadata,omitempty"`
	ReasoningEffort     string                        `json:"reasoning_effort,omitempty"`
	Modalities          []string                      `json:"modalities,omitempty"`
	Audio               *Audio                        `json:"audio,omitempty"`
	WebSearchOptions    any                           `json:"web_search_options,omitempty"`
}

type ChatCompletionResponse struct {
	Id                string                 `json:"id"`
	Object            string                 `json:"object"`
	Created           int64                  `json:"created"`
	Model             string                 `json:"model"`
	Choices           []ChatCompletionChoice `json:"choices"`
	Usage             *Usage                 `json:"usage"`
	ServiceTier       string                 `json:"service_tier,omitempty"`
	SystemFingerprint string                 `json:"system_fingerprint,omitempty"`
	Obfuscation       string                 `json:"obfuscation,omitempty"`
	PromptAnnotations []PromptAnnotation     `json:"prompt_annotations,omitempty"`
	ResponseBytes     []byte                 `json:"-"`
	ConnTime          int64                  `json:"-"`
	Duration          int64                  `json:"-"`
	TotalTime         int64                  `json:"-"`
	Error             error                  `json:"-"`
}

type ChatCompletionMessage struct {
	Role             string        `json:"role"`
	Content          any           `json:"content"`
	ReasoningContent any           `json:"reasoning_content,omitempty"`
	Refusal          *string       `json:"refusal"`
	Name             string        `json:"name,omitempty"`
	FunctionCall     *FunctionCall `json:"function_call,omitempty"`
	ToolCalls        any           `json:"tool_calls,omitempty"`
	ToolCallID       string        `json:"tool_call_id,omitempty"`
	Audio            *Audio        `json:"audio,omitempty"`
	Annotations      []any         `json:"annotations"`
}

type ChatCompletionChoice struct {
	Index        int                              `json:"index"`
	Message      *ChatCompletionMessage           `json:"message,omitempty"`
	Delta        *ChatCompletionStreamChoiceDelta `json:"delta,omitempty"`
	LogProbs     *LogProbs                        `json:"logprobs"`
	FinishReason string                           `json:"finish_reason"`
}

type Usage struct {
	PromptTokens             int                     `json:"prompt_tokens"`
	CompletionTokens         int                     `json:"completion_tokens"`
	TotalTokens              int                     `json:"total_tokens"`
	PromptTokensDetails      PromptTokensDetails     `json:"prompt_tokens_details"`
	CompletionTokensDetails  CompletionTokensDetails `json:"completion_tokens_details"`
	SearchTokens             int                     `json:"search_tokens,omitempty"`
	CacheCreationInputTokens int                     `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens     int                     `json:"cache_read_input_tokens,omitempty"`
	InputTokens              int                     `json:"input_tokens"`
	OutputTokens             int                     `json:"output_tokens"`
	InputTokensDetails       InputTokensDetails      `json:"input_tokens_details"`
	OutputTokensDetails      OutputTokensDetails     `json:"output_tokens_details"`
}
type PromptTokensDetails struct {
	AudioTokens     int `json:"audio_tokens"`
	CachedTokens    int `json:"cached_tokens"`
	ReasoningTokens int `json:"reasoning_tokens"`
	TextTokens      int `json:"text_tokens"`
}

type CompletionTokensDetails struct {
	AudioTokens              int `json:"audio_tokens"`
	ReasoningTokens          int `json:"reasoning_tokens"`
	CachedTokens             int `json:"cached_tokens"`
	CachedTokensInternal     int `json:"cached_tokens_internal"`
	TextTokens               int `json:"text_tokens"`
	ImageTokens              int `json:"image_tokens"`
	AcceptedPredictionTokens int `json:"accepted_prediction_tokens"`
	RejectedPredictionTokens int `json:"rejected_prediction_tokens"`
}

type InputTokensDetails struct {
	TextTokens   int `json:"text_tokens"`
	ImageTokens  int `json:"image_tokens"`
	CachedTokens int `json:"cached_tokens"`
}

type OutputTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}

type ChatCompletionStreamChoiceDelta struct {
	Content          string        `json:"content"`
	ReasoningContent any           `json:"reasoning_content,omitempty"`
	Role             string        `json:"role,omitempty"`
	FunctionCall     *FunctionCall `json:"function_call,omitempty"`
	ToolCalls        any           `json:"tool_calls,omitempty"`
	Refusal          *string       `json:"refusal,omitempty"`
	Audio            *Audio        `json:"audio,omitempty"`
	Annotations      any           `json:"annotations,omitempty"`
}

type ChatCompletionResponseFormat struct {
	Type       string `json:"type,omitempty"`
	JSONSchema any    `json:"json_schema,omitempty"`
}

type FunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Strict      bool   `json:"strict,omitempty"`
	Parameters  any    `json:"parameters"`
}

type StreamOptions struct {
	IncludeUsage bool `json:"include_usage,omitempty"`
}

type Audio struct {
	Voice      string `json:"voice,omitempty"`
	Format     string `json:"format,omitempty"`
	Id         string `json:"id,omitempty"`
	Data       string `json:"data,omitempty"`
	ExpiresAt  int    `json:"expires_at,omitempty"`
	Transcript string `json:"transcript,omitempty"`
}

type PromptAnnotation struct {
	PromptIndex          int                  `json:"prompt_index,omitempty"`
	ContentFilterResults ContentFilterResults `json:"content_filter_results,omitempty"`
}

type Hate struct {
	Filtered bool   `json:"filtered"`
	Severity string `json:"severity,omitempty"`
}
type SelfHarm struct {
	Filtered bool   `json:"filtered"`
	Severity string `json:"severity,omitempty"`
}
type Sexual struct {
	Filtered bool   `json:"filtered"`
	Severity string `json:"severity,omitempty"`
}
type Violence struct {
	Filtered bool   `json:"filtered"`
	Severity string `json:"severity,omitempty"`
}

type JailBreak struct {
	Filtered bool `json:"filtered"`
	Detected bool `json:"detected"`
}

type Profanity struct {
	Filtered bool `json:"filtered"`
	Detected bool `json:"detected"`
}

type ContentFilterResults struct {
	Hate      Hate      `json:"hate,omitempty"`
	SelfHarm  SelfHarm  `json:"self_harm,omitempty"`
	Sexual    Sexual    `json:"sexual,omitempty"`
	Violence  Violence  `json:"violence,omitempty"`
	JailBreak JailBreak `json:"jailbreak,omitempty"`
	Profanity Profanity `json:"profanity,omitempty"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	Index    *int         `json:"index,omitempty"`
	ID       string       `json:"id,omitempty"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type TopLogProbs struct {
	Token   string  `json:"token"`
	LogProb float64 `json:"logprob"`
	Bytes   []byte  `json:"bytes,omitempty"`
}

type LogProb struct {
	Token       string        `json:"token"`
	LogProb     float64       `json:"logprob"`
	Bytes       []byte        `json:"bytes,omitempty"`
	TopLogProbs []TopLogProbs `json:"top_logprobs"`
}

type LogProbs struct {
	Content []LogProb `json:"content"`
}
