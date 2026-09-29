package model

import "net/http"

type XAISearchParameters struct {
	Mode             string            `json:"mode,omitempty"` // auto / on / off
	ReturnCitations  *bool             `json:"return_citations,omitempty"`
	FromDate         string            `json:"from_date,omitempty"` // YYYY-MM-DD
	ToDate           string            `json:"to_date,omitempty"`
	MaxSearchResults int               `json:"max_search_results,omitempty"`
	Sources          []XAISearchSource `json:"sources,omitempty"`
}

type XAISearchSource struct {
	Type              string   `json:"type"` // web / x / news / rss
	Country           string   `json:"country,omitempty"`
	ExcludedWebsites  []string `json:"excluded_websites,omitempty"`
	AllowedWebsites   []string `json:"allowed_websites,omitempty"`
	SafeSearch        *bool    `json:"safe_search,omitempty"`
	IncludedXHandles  []string `json:"included_x_handles,omitempty"`
	ExcludedXHandles  []string `json:"excluded_x_handles,omitempty"`
	XHandles          []string `json:"x_handles,omitempty"` // 兼容旧字段
	PostFavoriteCount int      `json:"post_favorite_count,omitempty"`
	PostViewCount     int      `json:"post_view_count,omitempty"`
	Links             []string `json:"links,omitempty"` // rss
}

type XAIWebSearchOptions struct {
	Filters           any    `json:"filters,omitempty"`
	SearchContextSize string `json:"search_context_size,omitempty"`
	UserLocation      any    `json:"user_location,omitempty"`
}

type XAIStreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

type XAIUsage struct {
	PromptTokens               int                            `json:"prompt_tokens,omitempty"`
	CompletionTokens           int                            `json:"completion_tokens,omitempty"`
	InputTokens                int                            `json:"input_tokens,omitempty"`
	OutputTokens               int                            `json:"output_tokens,omitempty"`
	TotalTokens                int                            `json:"total_tokens,omitempty"`
	PromptTokensDetails        *XAIPromptTokensDetails        `json:"prompt_tokens_details,omitempty"`
	CompletionTokensDetails    *XAICompletionTokensDetails    `json:"completion_tokens_details,omitempty"`
	InputTokensDetails         *XAIInputTokensDetails         `json:"input_tokens_details,omitempty"`
	OutputTokensDetails        *XAIOutputTokensDetails        `json:"output_tokens_details,omitempty"`
	CostInUsdTicks             int64                          `json:"cost_in_usd_ticks,omitempty"`
	CostInNanoUsd              *int64                         `json:"cost_in_nano_usd,omitempty"`
	NumSourcesUsed             int                            `json:"num_sources_used,omitempty"`
	NumServerSideToolsUsed     int                            `json:"num_server_side_tools_used,omitempty"`
	ServerSideToolUsageDetails *XAIServerSideToolUsageDetails `json:"server_side_tool_usage_details,omitempty"`
	ContextDetails             *XAIUsageContextDetails        `json:"context_details,omitempty"`
	DroppedMessageCount        int                            `json:"dropped_message_count,omitempty"`
	PromptTextTokens           int                            `json:"prompt_text_tokens,omitempty"`
	CachedPromptTextTokens     int                            `json:"cached_prompt_text_tokens,omitempty"`
	PromptImageTokens          int                            `json:"prompt_image_tokens,omitempty"`
	ReasoningTokens            int                            `json:"reasoning_tokens,omitempty"`
}

type XAIPromptTokensDetails struct {
	AudioTokens  int `json:"audio_tokens,omitempty"`
	CachedTokens int `json:"cached_tokens,omitempty"`
	ImageTokens  int `json:"image_tokens,omitempty"`
	TextTokens   int `json:"text_tokens,omitempty"`
}

type XAICompletionTokensDetails struct {
	AcceptedPredictionTokens int `json:"accepted_prediction_tokens,omitempty"`
	AudioTokens              int `json:"audio_tokens,omitempty"`
	ReasoningTokens          int `json:"reasoning_tokens,omitempty"`
	RejectedPredictionTokens int `json:"rejected_prediction_tokens,omitempty"`
	TextTokens               int `json:"text_tokens,omitempty"`
	ImageTokens              int `json:"image_tokens,omitempty"`
}

type XAIInputTokensDetails struct {
	AudioTokens  int `json:"audio_tokens,omitempty"`
	CachedTokens int `json:"cached_tokens,omitempty"`
	ImageTokens  int `json:"image_tokens,omitempty"`
	TextTokens   int `json:"text_tokens,omitempty"`
}

type XAIOutputTokensDetails struct {
	ImageTokens     int `json:"image_tokens,omitempty"`
	ReasoningTokens int `json:"reasoning_tokens,omitempty"`
	TextTokens      int `json:"text_tokens,omitempty"`
}

type XAIServerSideToolUsageDetails struct {
	CodeInterpreterCalls int `json:"code_interpreter_calls,omitempty"`
	DocumentSearchCalls  int `json:"document_search_calls,omitempty"`
	FileSearchCalls      int `json:"file_search_calls,omitempty"`
	ImageGenerationCalls int `json:"image_generation_calls,omitempty"`
	McpCalls             int `json:"mcp_calls,omitempty"`
	WebSearchCalls       int `json:"web_search_calls,omitempty"`
	XPostsFetched        int `json:"x_posts_fetched,omitempty"`
	XSearchCalls         int `json:"x_search_calls,omitempty"`
	XUsersFetched        int `json:"x_users_fetched,omitempty"`
}

type XAIUsageContextDetails struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
}

type XAIChatCompletionReq struct {
	Model               string                  `json:"model,omitempty"`
	Messages            []ChatCompletionMessage `json:"messages,omitempty"`
	Deferred            *bool                   `json:"deferred,omitempty"`
	FrequencyPenalty    float32                 `json:"frequency_penalty,omitempty"`
	LogitBias           map[string]int          `json:"logit_bias,omitempty"`
	Logprobs            bool                    `json:"logprobs,omitempty"`
	MaxCompletionTokens int                     `json:"max_completion_tokens,omitempty"`
	MaxTokens           int                     `json:"max_tokens,omitempty"`
	N                   int                     `json:"n,omitempty"`
	ParallelToolCalls   any                     `json:"parallel_tool_calls,omitempty"`
	PresencePenalty     float32                 `json:"presence_penalty,omitempty"`
	PromptCacheKey      string                  `json:"prompt_cache_key,omitempty"`
	ReasoningEffort     string                  `json:"reasoning_effort,omitempty"`
	ResponseFormat      any                     `json:"response_format,omitempty"`
	SafetyIdentifier    string                  `json:"safety_identifier,omitempty"`
	SearchParameters    any                     `json:"search_parameters,omitempty"`
	Seed                *int                    `json:"seed,omitempty"`
	ServiceTier         string                  `json:"service_tier,omitempty"`
	Stop                []string                `json:"stop,omitempty"`
	Stream              bool                    `json:"stream,omitempty"`
	StreamOptions       *XAIStreamOptions       `json:"stream_options,omitempty"`
	Temperature         float32                 `json:"temperature,omitempty"`
	ToolChoice          any                     `json:"tool_choice,omitempty"`
	Tools               any                     `json:"tools,omitempty"`
	TopLogprobs         int                     `json:"top_logprobs,omitempty"`
	TopP                float32                 `json:"top_p,omitempty"`
	User                string                  `json:"user,omitempty"`
	WebSearchOptions    any                     `json:"web_search_options,omitempty"`
}

type XAIChatCompletionRes struct {
	Id                string          `json:"id,omitempty"`
	Object            string          `json:"object,omitempty"`
	Created           int64           `json:"created,omitempty"`
	Model             string          `json:"model,omitempty"`
	Choices           []XAIChatChoice `json:"choices,omitempty"`
	Usage             *XAIUsage       `json:"usage,omitempty"`
	ServiceTier       string          `json:"service_tier,omitempty"`
	SystemFingerprint string          `json:"system_fingerprint,omitempty"`
	Citations         []any           `json:"citations,omitempty"`
	OutputFiles       any             `json:"output_files,omitempty"`
	RequestId         string          `json:"request_id,omitempty"` // deferred
}

type XAIChatChoice struct {
	Index        int             `json:"index"`
	Message      *XAIChatMessage `json:"message,omitempty"`
	Delta        *XAIChatMessage `json:"delta,omitempty"`
	LogProbs     any             `json:"logprobs"`
	FinishReason string          `json:"finish_reason"`
}

type XAIChatMessage struct {
	Role             string  `json:"role,omitempty"`
	Content          any     `json:"content,omitempty"`
	ReasoningContent any     `json:"reasoning_content,omitempty"`
	Refusal          *string `json:"refusal,omitempty"`
	ToolCalls        any     `json:"tool_calls,omitempty"`
}

type XAIDeferredCompletionRes struct {
	RequestId string `json:"request_id"`
}

type XAIResponsesReq struct {
	Background         *bool                  `json:"background,omitempty"`
	ContextManagement  any                    `json:"context_management,omitempty"`
	Include            any                    `json:"include,omitempty"`
	Input              any                    `json:"input"`
	Instructions       *string                `json:"instructions,omitempty"`
	Logprobs           *bool                  `json:"logprobs,omitempty"`
	MaxOutputTokens    *int                   `json:"max_output_tokens,omitempty"`
	MaxTurns           *int                   `json:"max_turns,omitempty"`
	Metadata           map[string]string      `json:"metadata,omitempty"`
	MinP               *float32               `json:"min_p,omitempty"`
	Model              string                 `json:"model,omitempty"`
	ParallelToolCalls  *bool                  `json:"parallel_tool_calls,omitempty"`
	PreviousResponseId *string                `json:"previous_response_id,omitempty"`
	PromptCacheKey     *string                `json:"prompt_cache_key,omitempty"`
	Reasoning          *XAIResponsesReasoning `json:"reasoning,omitempty"`
	ReasoningEffort    string                 `json:"reasoning_effort,omitempty"`
	SafetyIdentifier   *string                `json:"safety_identifier,omitempty"`
	SearchParameters   any                    `json:"search_parameters,omitempty"`
	ServiceTier        string                 `json:"service_tier,omitempty"`
	Store              *bool                  `json:"store,omitempty"`
	Stream             *bool                  `json:"stream,omitempty"`
	Temperature        *float32               `json:"temperature,omitempty"`
	Text               any                    `json:"text,omitempty"`
	ToolChoice         any                    `json:"tool_choice,omitempty"`
	Tools              any                    `json:"tools,omitempty"`
	TopK               *int                   `json:"top_k,omitempty"`
	TopLogprobs        *int                   `json:"top_logprobs,omitempty"`
	TopP               *float32               `json:"top_p,omitempty"`
	Truncation         *string                `json:"truncation,omitempty"`
	User               *string                `json:"user,omitempty"`
}

type XAIResponsesRes struct {
	Background         bool                   `json:"background"`
	CompletedAt        *int64                 `json:"completed_at"`
	CreatedAt          int64                  `json:"created_at"`
	Error              *XAIResponsesError     `json:"error"`
	FrequencyPenalty   float32                `json:"frequency_penalty"`
	Id                 string                 `json:"id"`
	IncompleteDetails  any                    `json:"incomplete_details"`
	Instructions       any                    `json:"instructions"`
	MaxOutputTokens    *int                   `json:"max_output_tokens"`
	MaxToolCalls       *int                   `json:"max_tool_calls"`
	Metadata           map[string]string      `json:"metadata"`
	Model              string                 `json:"model"`
	Object             string                 `json:"object"`
	Output             []XAIResponsesOutput   `json:"output"`
	ParallelToolCalls  bool                   `json:"parallel_tool_calls"`
	PresencePenalty    float32                `json:"presence_penalty"`
	PreviousResponseId *string                `json:"previous_response_id"`
	PromptCacheKey     *string                `json:"prompt_cache_key"`
	Reasoning          *XAIResponsesReasoning `json:"reasoning"`
	SafetyIdentifier   *string                `json:"safety_identifier"`
	ServiceTier        string                 `json:"service_tier"`
	Status             string                 `json:"status"`
	Store              bool                   `json:"store"`
	Temperature        *float32               `json:"temperature"`
	Text               XAIResponsesText       `json:"text"`
	ToolChoice         any                    `json:"tool_choice"`
	Tools              any                    `json:"tools"`
	TopLogprobs        int                    `json:"top_logprobs"`
	TopP               *float32               `json:"top_p"`
	Truncation         string                 `json:"truncation"`
	Usage              *XAIUsage              `json:"usage"`
	User               *string                `json:"user"`
	ResponseBytes      []byte                 `json:"-"`
	ResponseHeaders    http.Header            `json:"-"`
	ConnTime           int64                  `json:"-"`
	Duration           int64                  `json:"-"`
	TotalTime          int64                  `json:"-"`
	Err                error                  `json:"-"`
}

type XAIResponsesStreamRes struct {
	Type           string              `json:"type"`
	SequenceNumber int                 `json:"sequence_number"`
	Response       XAIResponsesRes     `json:"response"`
	OutputIndex    int                 `json:"output_index"`
	ContentIndex   int                 `json:"content_index"`
	ItemId         string              `json:"item_id"`
	Item           XAIResponsesOutput  `json:"item"`
	Delta          string              `json:"delta"`
	Part           XAIResponsesContent `json:"part"`
	Arguments      any                 `json:"arguments"`
	Error          *XAIResponsesError  `json:"error"`
}

type XAIResponsesReasoning struct {
	Effort          string `json:"effort,omitempty"`
	GenerateSummary string `json:"generate_summary,omitempty"`
	Summary         string `json:"summary,omitempty"`
}

type XAIResponsesText struct {
	Format any `json:"format"`
}

type XAIResponsesOutput struct {
	Type      string                `json:"type"`
	Id        string                `json:"id"`
	Status    string                `json:"status,omitempty"`
	Role      string                `json:"role,omitempty"`
	Content   []XAIResponsesContent `json:"content,omitempty"`
	Summary   []XAIResponsesSummary `json:"summary,omitempty"`
	Arguments any                   `json:"arguments,omitempty"`
	CallId    string                `json:"call_id,omitempty"`
	Name      string                `json:"name,omitempty"`
}

type XAIResponsesContent struct {
	Type        string `json:"type"`
	Text        string `json:"text,omitempty"`
	ImageUrl    string `json:"image_url,omitempty"`
	Annotations []any  `json:"annotations,omitempty"`
	Logprobs    any    `json:"logprobs,omitempty"`
}

type XAIResponsesSummary struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type XAIResponsesError struct {
	Code           string `json:"code"`
	Message        string `json:"message"`
	Param          string `json:"param"`
	SequenceNumber int    `json:"sequence_number"`
	Type           string `json:"type"`
}

type XAIResponsesCompactReq struct {
	Input any    `json:"input"`
	Model string `json:"model"`
}

type XAIResponsesCompactRes struct {
	CreatedAt int64                `json:"created_at"`
	Id        string               `json:"id"`
	Model     string               `json:"model"`
	Object    string               `json:"object"`
	Output    []XAIResponsesOutput `json:"output"`
	Usage     *XAIUsage            `json:"usage"`
	Error     *XAIResponsesError   `json:"error"`
}

type XAIResponsesDeleteRes struct {
	Deleted bool   `json:"deleted"`
	Id      string `json:"id"`
	Object  string `json:"object"`
}

type XAIMedia struct {
	Url      string `json:"url,omitempty"`
	ImageUrl string `json:"image_url,omitempty"`
	FileId   string `json:"file_id,omitempty"`
	VoiceId  string `json:"voice_id,omitempty"`
}

type XAIStorageOptions struct {
	ExpiresAfter *int   `json:"expires_after,omitempty"`
	Filename     string `json:"filename,omitempty"`
	PublicUrl    any    `json:"public_url,omitempty"`
}

type XAIFileOutput struct {
	FileId             string `json:"file_id"`
	Filename           string `json:"filename"`
	PublicUrl          string `json:"public_url,omitempty"`
	PublicUrlError     string `json:"public_url_error,omitempty"`
	ExpiresAt          *int64 `json:"expires_at,omitempty"`
	PublicUrlExpiresAt *int64 `json:"public_url_expires_at,omitempty"`
}

type XAIImageReq struct {
	Prompt         string             `json:"prompt,omitempty"`
	Model          string             `json:"model,omitempty"`
	N              int                `json:"n,omitempty"`
	AspectRatio    string             `json:"aspect_ratio,omitempty"`
	Resolution     string             `json:"resolution,omitempty"` // 1k / 2k / 1.5k
	ResponseFormat string             `json:"response_format,omitempty"`
	Quality        string             `json:"quality,omitempty"` // low / medium
	User           string             `json:"user,omitempty"`
	StorageOptions *XAIStorageOptions `json:"storage_options,omitempty"`
	Image          any                `json:"image,omitempty"`
	Images         any                `json:"images,omitempty"`
}

type XAIImageRes struct {
	Data  []XAIImageData `json:"data"`
	Usage *XAIUsage      `json:"usage,omitempty"`
}

type XAIImageData struct {
	B64Json      string         `json:"b64_json,omitempty"`
	Url          string         `json:"url,omitempty"`
	MimeType     string         `json:"mime_type,omitempty"`
	StorageError string         `json:"storage_error,omitempty"`
	FileOutput   *XAIFileOutput `json:"file_output,omitempty"`
}

type XAIVideoMedia struct {
	Url    string `json:"url,omitempty"`
	FileId string `json:"file_id,omitempty"`
}

type XAIVideoKeyframe struct {
	Image       *XAIMedia `json:"image,omitempty"`
	ImageUrl    string    `json:"image_url,omitempty"`
	ImageFileId string    `json:"image_file_id,omitempty"`
	Timestamp   *float64  `json:"timestamp,omitempty"`
	TimestampS  *float64  `json:"timestamp_s,omitempty"`
}

type XAIVideoCreateReq struct {
	Model           string             `json:"model,omitempty"`
	Prompt          string             `json:"prompt,omitempty"`
	Duration        int                `json:"duration,omitempty"`
	AspectRatio     string             `json:"aspect_ratio,omitempty"`
	Resolution      string             `json:"resolution,omitempty"` // 480p / 720p / 1080p
	Image           *XAIVideoMedia     `json:"image,omitempty"`
	Video           *XAIVideoMedia     `json:"video,omitempty"`
	LastFrame       *XAIVideoMedia     `json:"last_frame,omitempty"`
	Keyframes       any                `json:"keyframes,omitempty"`
	ReferenceImages any                `json:"reference_images,omitempty"`
	ReferenceAudios any                `json:"reference_audios,omitempty"`
	GenerateAudio   *bool              `json:"generate_audio,omitempty"`
	StorageOptions  *XAIStorageOptions `json:"storage_options,omitempty"`
	User            string             `json:"user,omitempty"`
}

type XAIVideoCreateRes struct {
	RequestId string `json:"request_id"`
}

type XAIVideoRetrieveReq struct {
	RequestId string `json:"request_id" in:"path"`
}

type XAIVideoGenerated struct {
	Url               string         `json:"url,omitempty"`
	Duration          int            `json:"duration,omitempty"`
	RespectModeration *bool          `json:"respect_moderation,omitempty"`
	FileOutput        *XAIFileOutput `json:"file_output,omitempty"`
	StorageError      string         `json:"storage_error,omitempty"`
}

type XAIVideoJobRes struct {
	RequestId   string             `json:"request_id,omitempty"`
	Id          string             `json:"id,omitempty"`
	Object      string             `json:"object,omitempty"`
	Model       string             `json:"model,omitempty"`
	Status      string             `json:"status,omitempty"`
	Prompt      string             `json:"prompt,omitempty"`
	Duration    int                `json:"duration,omitempty"`
	AspectRatio string             `json:"aspect_ratio,omitempty"`
	Resolution  string             `json:"resolution,omitempty"`
	Video       *XAIVideoGenerated `json:"video,omitempty"`
	Usage       *XAIUsage          `json:"usage,omitempty"`
	Error       *XAIVideoError     `json:"error,omitempty"`
}

type XAIVideoError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
