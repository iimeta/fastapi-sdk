package model

// ---- 火山引擎视频生成 API 数据结构 ----

// VolcVideoContent 输入给模型的内容项
type VolcVideoContent struct {
	Type      string              `json:"type"`                 // text, image_url, video_url, audio_url, draft_task
	Text      string              `json:"text,omitempty"`       // type=text 时的文本提示词
	ImageUrl  *VolcVideoMediaUrl  `json:"image_url,omitempty"`  // type=image_url 时的图片对象
	VideoUrl  *VolcVideoMediaUrl  `json:"video_url,omitempty"`  // type=video_url 时的视频对象
	AudioUrl  *VolcVideoMediaUrl  `json:"audio_url,omitempty"`  // type=audio_url 时的音频对象
	DraftTask *VolcVideoDraftTask `json:"draft_task,omitempty"` // type=draft_task 时的样片任务
	Role      string              `json:"role,omitempty"`       // first_frame, last_frame, reference_image, reference_video, reference_audio
}

// VolcVideoMediaUrl 媒体 URL 对象（图片/视频/音频共用）
type VolcVideoMediaUrl struct {
	Url string `json:"url"` // URL / Base64 编码 / 素材 ID（asset://xxx）
}

// VolcVideoDraftTask 样片任务
type VolcVideoDraftTask struct {
	Id string `json:"id"` // 样片任务 ID
}

// VolcVideoTool 工具配置
type VolcVideoTool struct {
	Type string `json:"type"` // web_search
}

// VolcVideoCreateReq 创建视频生成任务请求
type VolcVideoCreateReq struct {
	Model                 string             `json:"model"`                             // 模型 ID
	Content               []VolcVideoContent `json:"content"`                           // 输入内容（文本/图片/视频/音频/样片）
	CallbackUrl           string             `json:"callback_url,omitempty"`            // 回调通知地址
	ReturnLastFrame       *bool              `json:"return_last_frame,omitempty"`       // 是否返回尾帧图像
	ServiceTier           string             `json:"service_tier,omitempty"`            // default / flex
	ExecutionExpiresAfter *int               `json:"execution_expires_after,omitempty"` // 任务超时阈值（秒）
	GenerateAudio         *bool              `json:"generate_audio,omitempty"`          // 是否生成有声视频
	Draft                 *bool              `json:"draft,omitempty"`                   // 是否开启样片模式
	Tools                 []VolcVideoTool    `json:"tools,omitempty"`                   // 工具配置
	SafetyIdentifier      string             `json:"safety_identifier,omitempty"`       // 终端用户标识
	Resolution            string             `json:"resolution,omitempty"`              // 480p / 720p / 1080p
	Ratio                 string             `json:"ratio,omitempty"`                   // 16:9 / 4:3 / 1:1 / 3:4 / 9:16 / 21:9 / adaptive
	Duration              *int               `json:"duration,omitempty"`                // 视频时长（秒）
	Frames                *int               `json:"frames,omitempty"`                  // 视频帧数
	Seed                  *int64             `json:"seed,omitempty"`                    // 随机种子
	CameraFixed           *bool              `json:"camera_fixed,omitempty"`            // 是否固定摄像头
	Watermark             *bool              `json:"watermark,omitempty"`               // 是否含水印
}

// VolcVideoContentResult 视频生成输出内容
type VolcVideoContentResult struct {
	VideoUrl     string `json:"video_url"`                // 生成视频 URL
	LastFrameUrl string `json:"last_frame_url,omitempty"` // 尾帧图像 URL
}

// VolcVideoError 错误信息
type VolcVideoError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// VolcVideoToolUsage 工具用量
type VolcVideoToolUsage struct {
	WebSearch int `json:"web_search,omitempty"` // 联网搜索次数
}

// VolcVideoUsage token 用量
type VolcVideoUsage struct {
	CompletionTokens int                 `json:"completion_tokens"`
	TotalTokens      int                 `json:"total_tokens"`
	ToolUsage        *VolcVideoToolUsage `json:"tool_usage,omitempty"`
}

// VolcVideoTaskRes 视频生成任务响应（查询/创建共用）
type VolcVideoTaskRes struct {
	Id                    string                  `json:"id"`
	Model                 string                  `json:"model"`
	Status                string                  `json:"status"`
	Error                 *VolcVideoError         `json:"error"`
	CreatedAt             int64                   `json:"created_at"`
	UpdatedAt             int64                   `json:"updated_at"`
	Content               *VolcVideoContentResult `json:"content"`
	Seed                  *int64                  `json:"seed,omitempty"`
	Resolution            string                  `json:"resolution,omitempty"`
	Ratio                 string                  `json:"ratio,omitempty"`
	Duration              *int                    `json:"duration,omitempty"`
	Frames                *int                    `json:"frames,omitempty"`
	FramesPerSecond       *int                    `json:"framespersecond,omitempty"`
	GenerateAudio         *bool                   `json:"generate_audio,omitempty"`
	Tools                 []VolcVideoTool         `json:"tools,omitempty"`
	SafetyIdentifier      string                  `json:"safety_identifier,omitempty"`
	Draft                 *bool                   `json:"draft,omitempty"`
	DraftTaskId           string                  `json:"draft_task_id,omitempty"`
	ServiceTier           string                  `json:"service_tier,omitempty"`
	ExecutionExpiresAfter int                     `json:"execution_expires_after,omitempty"`
	Usage                 *VolcVideoUsage         `json:"usage,omitempty"`
}

// VolcVideoListReq 查询视频任务列表请求参数（Query String）
type VolcVideoListReq struct {
	PageNum           *int   `json:"page_num,omitempty"`
	PageSize          *int   `json:"page_size,omitempty"`
	FilterStatus      string `json:"filter.status,omitempty"`
	FilterTaskIds     string `json:"filter.task_ids,omitempty"`
	FilterModel       string `json:"filter.model,omitempty"`
	FilterServiceTier string `json:"filter.service_tier,omitempty"`
}

// VolcVideoRetrieveReq 查询单个视频任务请求参数
type VolcVideoRetrieveReq struct {
	TaskId string `json:"task_id" in:"path"`
}

// VolcVideoDeleteReq 取消或删除视频任务请求参数
type VolcVideoDeleteReq struct {
	TaskId string `json:"task_id" in:"path"`
}

// VolcVideoListRes 查询视频任务列表响应
type VolcVideoListRes struct {
	Items []*VolcVideoTaskRes `json:"items"`
	Total int                 `json:"total"`
}
