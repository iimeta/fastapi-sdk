package consts

const (
	CORP_OPENAI         = "OpenAI"
	CORP_AZURE          = "Azure"
	CORP_BAIDU          = "Baidu"
	CORP_XFYUN          = "Xfyun"
	CORP_ALIYUN         = "Aliyun"
	CORP_ZHIPUAI        = "ZhipuAI"
	CORP_GOOGLE         = "Google"
	CORP_DEEPSEEK       = "DeepSeek"
	CORP_DEEPSEEK_BAIDU = "DeepSeek-Baidu"
	CORP_360AI          = "360AI"
	CORP_MIDJOURNEY     = "Midjourney"
	CORP_ANTHROPIC      = "Anthropic"
	CORP_GCP_CLAUDE     = "GCPClaude"
	CORP_AWS_CLAUDE     = "AWSClaude"
)

const (
	ROLE_SYSTEM    = "system"
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
