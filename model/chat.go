package model

import "github.com/iimeta/go-openai"

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	Model    string                  `json:"model"`
	Messages []ChatCompletionMessage `json:"messages"`
	// MaxTokens The maximum number of tokens that can be generated in the chat completion.
	// This value can be used to control costs for text generated via API.
	// This value is now deprecated in favor of max_completion_tokens, and is not compatible with o1 series models.
	// refs: https://platform.openai.com/docs/api-reference/chat/create#chat-create-max_tokens
	MaxTokens int `json:"max_tokens,omitempty"`
	// MaxCompletionTokens An upper bound for the number of tokens that can be generated for a completion,
	// including visible output tokens and reasoning tokens https://platform.openai.com/docs/guides/reasoning
	MaxCompletionTokens int                                  `json:"max_completion_tokens,omitempty"`
	Temperature         float32                              `json:"temperature,omitempty"`
	TopP                float32                              `json:"top_p,omitempty"`
	TopK                int                                  `json:"top_k,omitempty"`
	N                   int                                  `json:"n,omitempty"`
	Stream              bool                                 `json:"stream,omitempty"`
	Stop                []string                             `json:"stop,omitempty"`
	PresencePenalty     float32                              `json:"presence_penalty,omitempty"`
	ResponseFormat      *openai.ChatCompletionResponseFormat `json:"response_format,omitempty"`
	Seed                *int                                 `json:"seed,omitempty"`
	FrequencyPenalty    float32                              `json:"frequency_penalty,omitempty"`
	// LogitBias is must be a token id string (specified by their token ID in the tokenizer), not a word string.
	// incorrect: `"logit_bias":{"You": 6}`, correct: `"logit_bias":{"1639": 6}`
	// refs: https://platform.openai.com/docs/api-reference/chat/create#chat/create-logit_bias
	LogitBias map[string]int `json:"logit_bias,omitempty"`
	// LogProbs indicates whether to return log probabilities of the output tokens or not.
	// If true, returns the log probabilities of each output token returned in the content of message.
	// This option is currently not available on the gpt-4-vision-preview model.
	LogProbs bool `json:"logprobs,omitempty"`
	// TopLogProbs is an integer between 0 and 5 specifying the number of most likely tokens to return at each
	// token position, each with an associated log probability.
	// logprobs must be set to true if this parameter is used.
	TopLogProbs int    `json:"top_logprobs,omitempty"`
	User        string `json:"user,omitempty"`
	// Deprecated: use Tools instead.
	Functions []openai.FunctionDefinition `json:"functions,omitempty"`
	// Deprecated: use ToolChoice instead.
	FunctionCall any `json:"function_call,omitempty"`
	Tools        any `json:"tools,omitempty"`
	// This can be either a string or an ToolChoice object.
	ToolChoice any `json:"tool_choice,omitempty"`
	// Options for streaming response. Only set this when you set stream: true.
	StreamOptions *openai.StreamOptions `json:"stream_options,omitempty"`
	// Disable the default behavior of parallel tool calls by setting it: false.
	ParallelToolCalls any `json:"parallel_tool_calls,omitempty"`
	// Store can be set to true to store the output of this completion request for use in distillations and evals.
	// https://platform.openai.com/docs/api-reference/chat/create#chat-create-store
	Store bool `json:"store,omitempty"`
	// Metadata to store with the completion.
	Metadata map[string]string `json:"metadata,omitempty"`
	// o1 models only
	ReasoningEffort string `json:"reasoning_effort,omitempty"`

	Modalities []string `json:"modalities,omitempty"`
	Audio      *struct {
		Voice  string `json:"voice,omitempty"`
		Format string `json:"format,omitempty"`
	} `json:"audio,omitempty"`
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	ID                string                    `json:"id"`
	Object            string                    `json:"object"`
	Created           int64                     `json:"created"`
	Model             string                    `json:"model"`
	Choices           []ChatCompletionChoice    `json:"choices"`
	Usage             *Usage                    `json:"usage"`
	SystemFingerprint string                    `json:"system_fingerprint,omitempty"`
	PromptAnnotations []openai.PromptAnnotation `json:"prompt_annotations,omitempty"`
	ResponseBytes     []byte                    `json:"-"`
	ConnTime          int64                     `json:"-"`
	Duration          int64                     `json:"-"`
	TotalTime         int64                     `json:"-"`
	Error             error                     `json:"-"`
}

type ChatCompletionMessage struct {
	Role             string                   `json:"role"`
	Content          any                      `json:"content"`
	ReasoningContent any                      `json:"reasoning_content,omitempty"`
	Refusal          string                   `json:"refusal,omitempty"`
	MultiContent     []openai.ChatMessagePart `json:"-"`

	// This property isn't in the official documentation, but it's in
	// the documentation for the official library for python:
	// - https://github.com/openai/openai-python/blob/main/chatml.md
	// - https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
	Name string `json:"name,omitempty"`

	FunctionCall *openai.FunctionCall `json:"function_call,omitempty"`

	// For Role=assistant prompts this may be set to the tool calls generated by the model, such as function calls.
	ToolCalls []openai.ToolCall `json:"tool_calls,omitempty"`

	// For Role=tool prompts this should be set to the ID given in the assistant's prior request to call a tool.
	ToolCallID string `json:"tool_call_id,omitempty"`

	Audio *openai.Audio `json:"audio,omitempty"`
}

type ChatCompletionChoice struct {
	Index        int                              `json:"index"`
	Message      *ChatCompletionMessage           `json:"message,omitempty"`
	Delta        *ChatCompletionStreamChoiceDelta `json:"delta,omitempty"`
	LogProbs     *openai.LogProbs                 `json:"logprobs,omitempty"`
	FinishReason openai.FinishReason              `json:"finish_reason"`
	//ContentFilterResults *openai.ContentFilterResults            `json:"content_filter_results,omitempty"`
}

// Usage Represents the total token usage per request to OpenAI.
type Usage struct {
	PromptTokens             int                             `json:"prompt_tokens"`
	CompletionTokens         int                             `json:"completion_tokens"`
	TotalTokens              int                             `json:"total_tokens"`
	PromptTokensDetails      *openai.PromptTokensDetails     `json:"prompt_tokens_details"`
	CompletionTokensDetails  *openai.CompletionTokensDetails `json:"completion_tokens_details"`
	SearchTokens             int                             `json:"search_tokens,omitempty"`
	CacheCreationInputTokens int                             `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens     int                             `json:"cache_read_input_tokens,omitempty"`
}

type ChatCompletionStreamChoiceDelta struct {
	Content          string               `json:"content"`
	ReasoningContent any                  `json:"reasoning_content,omitempty"`
	Role             string               `json:"role,omitempty"`
	FunctionCall     *openai.FunctionCall `json:"function_call,omitempty"`
	ToolCalls        []openai.ToolCall    `json:"tool_calls,omitempty"`
	Refusal          string               `json:"refusal,omitempty"`
	Audio            *openai.Audio        `json:"audio,omitempty"`
}
