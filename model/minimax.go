package model

// ---- MiniMax 视频生成 V2 API 数据结构 ----

// MiniMaxMediaUrl 媒体 URL 对象（图片/视频/音频共用）
type MiniMaxMediaUrl struct {
	Url string `json:"url"` // 公网 URL / mm_file://{file_id} / data URI
}

// MiniMaxContentItem 多模态输入内容项
type MiniMaxContentItem struct {
	Type     string           `json:"type"`                // text, image_url, video_url, audio_url
	Text     string           `json:"text,omitempty"`      // type=text 时的文本提示词
	ImageUrl *MiniMaxMediaUrl `json:"image_url,omitempty"` // type=image_url 时的图片对象
	VideoUrl *MiniMaxMediaUrl `json:"video_url,omitempty"` // type=video_url 时的视频对象
	AudioUrl *MiniMaxMediaUrl `json:"audio_url,omitempty"` // type=audio_url 时的音频对象
	Role     string           `json:"role,omitempty"`      // first_frame, last_frame, reference_image, reference_video, reference_audio
}

// MiniMaxVideoExtra MiniMax-H3-Max 额外生成选项
type MiniMaxVideoExtra struct {
	PromptExpansionMode string `json:"prompt_expansion_mode,omitempty"` // disabled / balanced / quality
}

// MiniMaxVideoCreateReq 创建视频生成任务请求
type MiniMaxVideoCreateReq struct {
	Model         string               `json:"model"`                    // MiniMax-H3 / MiniMax-H3-Max
	Content       []MiniMaxContentItem `json:"content"`                  // 多模态输入
	Resolution    string               `json:"resolution"`               // 480P / 768P / 2K
	Duration      int                  `json:"duration"`                 // 视频时长（秒）
	Ratio         string               `json:"ratio,omitempty"`          // adaptive / 21:9 / 16:9 / 4:3 / 1:1 / 3:4 / 9:16
	Extra         *MiniMaxVideoExtra   `json:"extra,omitempty"`          // MiniMax-H3-Max 额外选项
	CallbackUrl   string               `json:"callback_url,omitempty"`   // 任务状态回调地址
	AigcWatermark *bool                `json:"aigc_watermark,omitempty"` // 是否添加 AIGC 水印
}

// MiniMaxVideoCreateRes 创建视频生成任务响应
type MiniMaxVideoCreateRes struct {
	TaskId string `json:"task_id"`
}

// MiniMaxVideoRetrieveReq 查询单个视频任务请求参数
type MiniMaxVideoRetrieveReq struct {
	TaskId string `json:"task_id" in:"path"`
}

// MiniMaxVideoQueryRes 查询任务响应
type MiniMaxVideoQueryRes struct {
	Task *MiniMaxVideoTask `json:"task"`
}

// MiniMaxVideoTask 任务对象（查询接口返回）
type MiniMaxVideoTask struct {
	Id         string               `json:"id"`
	Model      string               `json:"model"`
	Status     string               `json:"status"` // queued / running / succeeded / failed / cancelled
	Error      *MiniMaxVideoError   `json:"error,omitempty"`
	CreatedAt  int64                `json:"created_at"`
	UpdatedAt  int64                `json:"updated_at"`
	Content    *MiniMaxVideoContent `json:"content,omitempty"`
	Resolution string               `json:"resolution,omitempty"`
	Duration   int                  `json:"duration,omitempty"`
	Usage      *MiniMaxVideoUsage   `json:"usage,omitempty"`
	Ratio      string               `json:"ratio,omitempty"`
	TaskType   string               `json:"task_type,omitempty"` // generation / h3_context_ir / regeneration
	Modality   string               `json:"modality,omitempty"`  // video / text
}

// MiniMaxVideoError 任务错误信息
type MiniMaxVideoError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// MiniMaxVideoContent 任务产物
type MiniMaxVideoContent struct {
	Url    string `json:"url,omitempty"`    // 视频限时下载 URL
	Prompt string `json:"prompt,omitempty"` // H3-Context-IR 增强提示词
}

// MiniMaxVideoUsage 本次请求用量
type MiniMaxVideoUsage struct {
	TotalSeconds      int `json:"total_seconds,omitempty"`
	InputSeconds      int `json:"input_seconds,omitempty"`
	OutputSeconds     int `json:"output_seconds,omitempty"`
	InputImageCount   int `json:"input_image_count,omitempty"`
	InputAudioSeconds int `json:"input_audio_seconds,omitempty"`
	TotalTokens       int `json:"total_tokens,omitempty"`
	PromptTokens      int `json:"prompt_tokens,omitempty"`
	CompletionTokens  int `json:"completion_tokens,omitempty"`
}

// MiniMaxErrorRes OpenAI 风格错误响应
type MiniMaxErrorRes struct {
	Type      string              `json:"type"`
	Error     *MiniMaxErrorDetail `json:"error,omitempty"`
	RequestId string              `json:"request_id,omitempty"`
}

// MiniMaxErrorDetail 错误详情
type MiniMaxErrorDetail struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	HttpCode string `json:"http_code,omitempty"`
}
