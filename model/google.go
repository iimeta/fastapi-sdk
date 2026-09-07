package model

import (
	"net/http"
	"time"
)

type GoogleChatCompletionReq struct {
	Contents         []Content        `json:"contents"`
	GenerationConfig GenerationConfig `json:"generationConfig,omitempty"`
	Tools            any              `json:"tools,omitempty"`
}

type GoogleChatCompletionRes struct {
	Candidates    []Candidate    `json:"candidates"`
	UsageMetadata *UsageMetadata `json:"usageMetadata"`
	ModelVersion  string         `json:"modelVersion"`
	Error         struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
		Details []struct {
			Type     string `json:"@type"`
			Reason   string `json:"reason"`
			Domain   string `json:"domain"`
			Metadata struct {
				Service string `json:"service"`
			} `json:"metadata"`
		} `json:"details"`
	} `json:"error"`
	ResponseBytes   []byte      `json:"-"`
	ResponseHeaders http.Header `json:"-"`
	ConnTime        int64       `json:"-"`
	Duration        int64       `json:"-"`
	TotalTime       int64       `json:"-"`
	Err             error       `json:"-"`
}

type Content struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

type Part struct {
	Text             string      `json:"text,omitempty"`
	InlineData       *InlineData `json:"inlineData,omitempty"`
	FileData         *FileData   `json:"fileData,omitempty"`
	FunctionCall     any         `json:"functionCall,omitempty"`
	FunctionResponse any         `json:"functionResponse,omitempty"`
	ThoughtSignature any         `json:"thoughtSignature,omitempty"`
}

type InlineData struct {
	MimeType string `json:"mime_type,omitempty"`
	Data     string `json:"data,omitempty"`
}

type FileData struct {
	FileUri  string `json:"file_uri,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
}

type Candidate struct {
	Content       Content        `json:"content"`
	FinishReason  string         `json:"finishReason"`
	FinishMessage string         `json:"finishMessage"`
	Index         int            `json:"index"`
	SafetyRatings []SafetyRating `json:"safetyRatings"`
}

type SafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

type UsageMetadata struct {
	PromptTokenCount           int                  `json:"promptTokenCount"`           // 提示中的 token 数量。如果设置了 cachedContent，这仍然是有效提示的总大小，这意味着它包含缓存内容中的词元数。
	CachedContentTokenCount    int                  `json:"cachedContentTokenCount"`    // 提示的缓存部分（缓存内容）中的 token 数量
	CandidatesTokenCount       int                  `json:"candidatesTokenCount"`       // 所有生成的回答候选项中的 token 总数。
	ToolUsePromptTokenCount    int                  `json:"toolUsePromptTokenCount"`    // 仅限输出。工具使用提示中的 token 数量。
	ThoughtsTokenCount         int                  `json:"thoughtsTokenCount"`         // 仅限输出。思考模型用于思考的 token 数量。
	TotalTokenCount            int                  `json:"totalTokenCount"`            // 生成请求（提示 + 想法 + 回答候选）的总 token 数量。
	PromptTokensDetails        []ModalityTokenCount `json:"promptTokensDetails"`        // 仅限输出。请求输入中处理的模态列表。
	CacheTokensDetails         []ModalityTokenCount `json:"cacheTokensDetails"`         // 仅限输出。请求输入中缓存内容的模态列表。
	CandidatesTokensDetails    []ModalityTokenCount `json:"candidatesTokensDetails"`    // 仅限输出。响应中返回的模态列表。
	ToolUsePromptTokensDetails []ModalityTokenCount `json:"toolUsePromptTokensDetails"` // 仅限输出。为工具使用请求输入处理的模态列表。
	ServiceTier                string               `json:"serviceTier"`                // 仅限输出。请求的服务等级。[unspecified:默认服务层级（标准）, standard:标准服务层级, flex:灵活服务层级, priority:优先服务层级]
}

type ModalityTokenCount struct {
	Modality   string `json:"modality"`   // 与此 token 数量关联的模态。[MODALITY_UNSPECIFIED:未指定模态, TEXT:纯文本, IMAGE:图片, VIDEO:视频, AUDIO:音频, DOCUMENT:文档，例如 PDF]
	TokenCount int    `json:"tokenCount"` // 令牌数量。
}

type GenerationConfig struct {
	StopSequences      []string              `json:"stopSequences,omitempty"`
	CandidateCount     int                   `json:"candidateCount,omitempty"`
	MaxOutputTokens    int                   `json:"maxOutputTokens,omitempty"`
	Temperature        float32               `json:"temperature,omitempty"`
	TopP               float32               `json:"topP,omitempty"`
	TopK               int                   `json:"topK,omitempty"`
	ResponseModalities []string              `json:"responseModalities,omitempty"`
	ImageConfig        *ImageConfig          `json:"imageConfig,omitempty"`
	ResponseFormat     *GoogleResponseFormat `json:"responseFormat,omitempty"`
}

type ImageConfig struct {
	AspectRatio string `json:"aspectRatio,omitempty"`
	ImageSize   string `json:"imageSize,omitempty"`
}

type GoogleResponseFormat struct {
	Image *GoogleImage `json:"image,omitempty"`
}

type GoogleImage struct {
	AspectRatio string `json:"aspectRatio,omitempty"`
	ImageSize   string `json:"imageSize,omitempty"`
}

type GoogleFileResponse struct {
	File GoogleFile `json:"file"`
}

type GoogleFileListResponse struct {
	Files     []GoogleFile `json:"files"`
	TotalTime int64        `json:"-"`
}

type GoogleFile struct {
	Name           string    `json:"name"`
	MimeType       string    `json:"mimeType"`
	SizeBytes      string    `json:"sizeBytes"`
	CreateTime     time.Time `json:"createTime"`
	UpdateTime     time.Time `json:"updateTime"`
	ExpirationTime time.Time `json:"expirationTime"`
	Sha256Hash     string    `json:"sha256Hash"`
	Uri            string    `json:"uri"`
	State          string    `json:"state"`
	VideoMetadata  struct {
		VideoDuration string `json:"videoDuration"`
	} `json:"videoMetadata"`
	Source string `json:"source"`
}

type GoogleImageGenerationReq struct {
	Contents         []Content        `json:"contents"`
	Tools            any              `json:"tools,omitempty"`
	GenerationConfig GenerationConfig `json:"generationConfig,omitempty"`
}
