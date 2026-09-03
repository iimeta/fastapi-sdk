package model

// ---- 阿里云百炼 万相图像生成 API 数据结构 ----

// 输入消息内容项(text 与 image 二选一)
type BailianImageContent struct {
	Text  string `json:"text,omitempty"`  // 提示词, 不超过5000字符
	Image string `json:"image,omitempty"` // 输入图像 URL 或 Base64(data:{MIME};base64,xxx)
}

// 输入消息(仅支持单轮, role 固定为 user)
type BailianImageMessage struct {
	Role    string                `json:"role"`
	Content []BailianImageContent `json:"content"`
}

// 输入信息
type BailianImageInput struct {
	Messages []BailianImageMessage `json:"messages"`
}

// 自定义颜色主题项
type BailianColorPalette struct {
	Hex   string `json:"hex"`   // 十六进制色值
	Ratio string `json:"ratio"` // 占比, 精确到小数点后两位, 如 "25.00%"
}

// 模型参数配置
type BailianImageParameters struct {
	BboxList         [][][]int             `json:"bbox_list,omitempty"`         // 交互式编辑框选区域, 长度与输入图片数量一致
	EnableSequential *bool                 `json:"enable_sequential,omitempty"` // 是否启用组图模式
	Size             string                `json:"size,omitempty"`              // 1K / 2K / 4K 或 宽*高
	N                int                   `json:"n,omitempty"`                 // 生成数量, 非组图 1-4, 组图 1-12
	ThinkingMode     *bool                 `json:"thinking_mode,omitempty"`     // 是否开启思考模式
	ColorPalette     []BailianColorPalette `json:"color_palette,omitempty"`     // 自定义颜色主题, 3-10 种
	Watermark        *bool                 `json:"watermark,omitempty"`         // 是否添加水印
	Seed             *int64                `json:"seed,omitempty"`              // 随机种子 [0, 2147483647]
}

// 图像生成请求
type BailianImageGenerationReq struct {
	Model      string                  `json:"model"`
	Input      BailianImageInput       `json:"input"`
	Parameters *BailianImageParameters `json:"parameters,omitempty"`
}

// 模型输出内容项
type BailianImageOutputContent struct {
	Type  string `json:"type"`            // image / text
	Image string `json:"image,omitempty"` // 生成图像 URL(PNG), 有效期24小时
	Text  string `json:"text,omitempty"`  // 生成的文字
}

// 模型返回的消息
type BailianImageOutputMessage struct {
	Role    string                      `json:"role"`
	Content []BailianImageOutputContent `json:"content"`
}

// 模型生成的输出内容
type BailianImageChoice struct {
	FinishReason string                    `json:"finish_reason"`
	Message      BailianImageOutputMessage `json:"message"`
}

// 任务输出信息(同步/异步查询共用)
type BailianImageOutput struct {
	Choices       []BailianImageChoice `json:"choices,omitempty"`
	Finished      *bool                `json:"finished,omitempty"`
	TaskId        string               `json:"task_id,omitempty"`
	TaskStatus    string               `json:"task_status,omitempty"` // PENDING / RUNNING / SUCCEEDED / FAILED / CANCELED / UNKNOWN
	SubmitTime    string               `json:"submit_time,omitempty"`
	ScheduledTime string               `json:"scheduled_time,omitempty"`
	EndTime       string               `json:"end_time,omitempty"`
}

// 输出信息统计, 只对成功的结果计数
type BailianImageUsage struct {
	ImageCount   int    `json:"image_count"`
	Size         string `json:"size"` // 生成的图像分辩率, 如 1376*768
	InputTokens  int    `json:"input_tokens"`
	OutputTokens int    `json:"output_tokens"`
	TotalTokens  int    `json:"total_tokens"`
}

// 图像生成响应(成功与失败共用)
type BailianImageGenerationRes struct {
	Output    *BailianImageOutput `json:"output,omitempty"`
	Usage     *BailianImageUsage  `json:"usage,omitempty"`
	RequestId string              `json:"request_id,omitempty"`
	Code      string              `json:"code,omitempty"`    // 失败时返回
	Message   string              `json:"message,omitempty"` // 失败时返回
}

// ---- 阿里云百炼 万相视频生成 API 数据结构 ----

// 媒体素材
type BailianVideoMedia struct {
	Type string `json:"type"` // first_frame / last_frame / reference_image / reference_video / reference_audio / file / link
	Url  string `json:"url"`  // 公网 URL / oss:// 临时 URL / Base64(data:{MIME};base64,xxx)
}

// 输入信息, prompt 与 media 必填其一
type BailianVideoInput struct {
	Prompt string              `json:"prompt,omitempty"`
	Media  []BailianVideoMedia `json:"media,omitempty"`
}

// 视频处理参数
type BailianVideoParameters struct {
	Resolution   string `json:"resolution,omitempty"`    // 1080P(默认) / 720P / 480P
	Ratio        string `json:"ratio,omitempty"`         // adaptive(默认) / 16:9 / 4:3 / 1:1 / 3:4 / 9:16
	Duration     *int   `json:"duration,omitempty"`      // 时长(秒), 默认5, [2,30], -1 为智能时长
	Audio        *bool  `json:"audio,omitempty"`         // 是否包含音频, 默认 true
	Seed         *int64 `json:"seed,omitempty"`          // 随机种子 [0, 2147483647]
	PromptExtend *bool  `json:"prompt_extend,omitempty"` // 是否开启 prompt 智能改写, 默认 true
	Watermark    *bool  `json:"watermark,omitempty"`     // 是否添加水印, 默认 false
}

// 创建视频生成任务请求
type BailianVideoCreateReq struct {
	Model      string                  `json:"model"`
	Input      BailianVideoInput       `json:"input"`
	Parameters *BailianVideoParameters `json:"parameters,omitempty"`
}

// 任务输出信息(创建/查询共用)
type BailianVideoOutput struct {
	TaskId        string `json:"task_id"`
	TaskStatus    string `json:"task_status"` // PENDING / RUNNING / SUCCEEDED / FAILED / CANCELED / UNKNOWN
	SubmitTime    string `json:"submit_time,omitempty"`
	ScheduledTime string `json:"scheduled_time,omitempty"`
	EndTime       string `json:"end_time,omitempty"`
	OrigPrompt    string `json:"orig_prompt,omitempty"`
	VideoUrl      string `json:"video_url,omitempty"` // 任务成功时返回, 有效期24小时
	Code          string `json:"code,omitempty"`      // 任务失败时返回
	Message       string `json:"message,omitempty"`   // 任务失败时返回
}

// 输出信息统计, 只对成功的结果计数
type BailianVideoUsage struct {
	VideoCount          int     `json:"video_count"`
	Duration            float64 `json:"duration"`
	InputVideoDuration  float64 `json:"input_video_duration"`
	OutputVideoDuration float64 `json:"output_video_duration"`
	Fps                 int     `json:"fps"`
	SR                  int     `json:"SR"`
	Ratio               string  `json:"ratio"`
}

// 视频任务响应(创建/查询共用, 成功与失败共用)
type BailianVideoTaskRes struct {
	Output    *BailianVideoOutput `json:"output,omitempty"`
	Usage     *BailianVideoUsage  `json:"usage,omitempty"`
	RequestId string              `json:"request_id,omitempty"`
	Code      string              `json:"code,omitempty"`    // 请求失败时返回
	Message   string              `json:"message,omitempty"` // 请求失败时返回
}

// 查询任务结果请求参数
type BailianVideoRetrieveReq struct {
	TaskId string `json:"task_id" in:"path"`
}
