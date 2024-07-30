package model

import "github.com/iimeta/go-openai"

type XfyunChatCompletionReq struct {
	Header    Header    `json:"header"`
	Parameter Parameter `json:"parameter"`
	Payload   Payload   `json:"payload"`
}

type XfyunChatCompletionRes struct {
	Header  Header  `json:"header"`
	Payload Payload `json:"payload"`
}

type Header struct {
	// req
	AppId string `json:"app_id"` // 应用appid，从开放平台控制台创建的应用中获取
	Uid   string `json:"uid"`    // 每个用户的id，用于区分不同用户，最大长度32
	// res
	Code    int    `json:"code,omitempty"`    // 错误码，0表示正常，非0表示出错；详细释义可在接口说明文档最后的错误码说明了解
	Message string `json:"message,omitempty"` // 会话是否成功的描述信息
	Sid     string `json:"sid,omitempty"`     // 会话的唯一id，用于讯飞技术人员查询服务端会话日志使用,出现调用错误时建议留存该字段
	Status  int    `json:"status,omitempty"`  // 会话状态，取值为[0,1,2]；0代表首次结果；1代表中间结果；2代表最后一个结果
}

type Parameter struct {
	// req
	Chat *Chat `json:"chat"`
}

type Chat struct {
	// 指定访问的领域:
	// general指向V1.5版本;
	// generalv2指向V2版本;
	// generalv3指向V3版本;
	// generalv3.5指向V3.5版本;
	// 注意：不同的取值对应的url也不一样！
	Domain string `json:"domain"`
	// 核采样阈值。用于决定结果随机性，取值越高随机性越强即相同的问题得到的不同答案的可能性越高
	// 取值范围 (0，1] ，默认值0.5
	Temperature float32 `json:"temperature,omitempty"`
	// 模型回答的tokens的最大长度
	// V1.5取值为[1,4096]
	// V2.0、V3.0和V3.5取值为[1,8192]，默认为2048。
	MaxTokens int `json:"max_tokens,omitempty"`
	// 从k个候选中随机选择⼀个（⾮等概率）
	// 取值为[1，6],默认为4
	TopK int `json:"top_k,omitempty"`
	// 用于关联用户会话
	// 需要保障用户下的唯一性
	ChatId string `json:"chat_id,omitempty"`
	// 图片的宽度
	Width int `json:"width,omitempty"`
	// 图片的宽度
	Height int `json:"height,omitempty"`
}

type Payload struct {
	// req
	Message   *Message   `json:"message,omitempty"`
	Functions *Functions `json:"functions,omitempty"`
	// res
	Choices *Choices    `json:"choices,omitempty"`
	Usage   *XfyunUsage `json:"usage,omitempty"`
}

type Message struct {
	// req
	Text []ChatCompletionMessage `json:"text"`
}

type Functions struct {
	// req
	Text []openai.FunctionDefinition `json:"text"`
}

type Text struct {
	// 角色标识，固定为assistant，标识角色为AI
	Role string `json:"role,omitempty"`
	// AI的回答内容
	Content string `json:"content,omitempty"`
	// 结果序号，取值为[0,10]; 当前为保留字段，开发者可忽略
	Index int `json:"index,omitempty"`
	// 内容类型
	ContentType string `json:"content_type,omitempty"`
	// function call 返回结果
	FunctionCall *openai.FunctionCall `json:"function_call,omitempty"`
	// 保留字段，可忽略
	QuestionTokens int `json:"question_tokens,omitempty"`
	// 包含历史问题的总tokens大小
	PromptTokens int `json:"prompt_tokens,omitempty"`
	// 回答的tokens大小
	CompletionTokens int `json:"completion_tokens,omitempty"`
	// prompt_tokens和completion_tokens的和，也是本次交互计费的tokens大小
	TotalTokens int `json:"total_tokens,omitempty"`
}

type Choices struct {
	// 文本响应状态，取值为[0,1,2]; 0代表首个文本结果；1代表中间文本结果；2代表最后一个文本结果
	Status int `json:"status,omitempty"`
	// 返回的数据序号，取值为[0,9999999]
	Seq  int    `json:"seq,omitempty"`
	Text []Text `json:"text,omitempty"`
}

type XfyunUsage struct {
	// res
	Text *Text `json:"text,omitempty"`
}
