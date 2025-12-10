package consts

const (
	PROVIDER_OPENAI         = "OpenAI"
	PROVIDER_ANTHROPIC      = "Anthropic"
	PROVIDER_GOOGLE         = "Google"
	PROVIDER_AZURE          = "Azure"
	PROVIDER_DEEPSEEK       = "DeepSeek"
	PROVIDER_DEEPSEEK_BAIDU = "DeepSeek-Baidu"
	PROVIDER_BAIDU          = "Baidu"
	PROVIDER_ALIYUN         = "Aliyun"
	PROVIDER_XFYUN          = "Xfyun"
	PROVIDER_ZHIPUAI        = "ZhipuAI"
	PROVIDER_VOLC_ENGINE    = "VolcEngine"
	PROVIDER_AWS_CLAUDE     = "AWSClaude"
	PROVIDER_GCP_CLAUDE     = "GCPClaude"
	PROVIDER_GCP_GEMINI     = "GCPGemini"
	PROVIDER_MIDJOURNEY     = "Midjourney"
)

const (
	ROLE_SYSTEM    = "system"
	ROLE_DEVELOPER = "developer"
	ROLE_USER      = "user"
	ROLE_ASSISTANT = "assistant"
	ROLE_FUNCTION  = "function"
	ROLE_TOOL      = "tool"
	ROLE_MODEL     = "model"
)

const (
	DELTA_TYPE_TEXT       = "text_delta"
	DELTA_TYPE_INPUT_JSON = "input_json_delta"
)

const (
	COMPLETION_ID_PREFIX     = "chatcmpl-"
	COMPLETION_OBJECT        = "chat.completion"
	COMPLETION_STREAM_OBJECT = "chat.completion.chunk"
)

var MIME_TYPE_MAP = map[string]string{
	"pdf":  "application/pdf",
	"js":   "application/x-javascript",
	"py":   "application/x-python",
	"jpg":  "image/jpeg",
	"png":  "image/png",
	"webp": "image/webp",
	"gif":  "image/gif",
	"txt":  "text/plain",
	"html": "text/html",
	"css":  "text/css",
	"md":   "text/md",
	"csv":  "text/csv",
	"xml":  "text/xml",
	"rtf":  "text/rtf",
}

const (
	FinishReasonStop          = "stop"
	FinishReasonLength        = "length"
	FinishReasonFunctionCall  = "function_call"
	FinishReasonToolCalls     = "tool_calls"
	FinishReasonContentFilter = "content_filter"
	FinishReasonNull          = "null"
)
