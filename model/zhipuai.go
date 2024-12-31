package model

import "github.com/iimeta/go-openai"

type ZhipuAIChatCompletionReq struct {
	// 所要调用的模型编码
	Model string `json:"model"`
	// 调用语言模型时，将当前对话信息列表作为提示输入给模型
	// 按照 {"role": "user", "content": "你好"} 的json 数组形式进行传参
	// 可能的消息类型包括 System message、User message、Assistant message 和 Tool message。
	Messages []ChatCompletionMessage `json:"messages"`
	// 由用户端传参，需保证唯一性；用于区分每次请求的唯一标识，用户端不传时平台会默认生成。
	RequestId string `json:"request_id,omitempty"`
	// do_sample 为 true 时启用采样策略，do_sample 为 false 时采样策略 temperature、top_p 将不生效。默认值为 true。
	DoSample bool `json:"do_sample,omitempty"`
	// 使用同步调用时，此参数应当设置为 fasle 或者省略。表示模型生成完所有内容后一次性返回所有内容。默认值为 false。
	// 如果设置为 true，模型将通过标准 Event Stream ，逐块返回模型生成内容。Event Stream 结束时会返回一条data: [DONE]消息。
	// 注意：在模型流式输出生成内容的过程中，我们会分批对模型生成内容进行检测，当检测到违法及不良信息时，API会返回错误码（1301）。
	// 开发者识别到错误码（1301），应及时采取（清屏、重启对话）等措施删除生成内容，并确保不将含有违法及不良信息的内容传递给模型继续生成，避免其造成负面影响。
	Stream bool `json:"stream,omitempty"`
	// 采样温度，控制输出的随机性，必须为正数
	// 取值范围是：(0.0, 1.0)，不能等于 0，默认值为 0.95，值越大，会使输出更随机，更具创造性；值越小，输出会更加稳定或确定
	// 建议您根据应用场景调整 top_p 或 temperature 参数，但不要同时调整两个参数
	Temperature float32 `json:"temperature,omitempty"`
	// 用温度取样的另一种方法，称为核取样
	// 取值范围是：(0.0, 1.0) 开区间，不能等于 0 或 1，默认值为 0.7
	// 模型考虑具有 top_p 概率质量 tokens 的结果
	// 例如：0.1 意味着模型解码器只考虑从前 10% 的概率的候选集中取 tokens
	// 建议您根据应用场景调整 top_p 或 temperature 参数，但不要同时调整两个参数
	TopP float32 `json:"top_p,omitempty"`
	// 模型输出最大 tokens，最大输出为8192，默认值为1024
	MaxTokens int `json:"max_tokens,omitempty"`
	// 模型在遇到stop所制定的字符时将停止生成，目前仅支持单个停止词，格式为["stop_word1"]
	Stop []string `json:"stop,omitempty"`
	// 可供模型调用的工具列表,tools 字段会计算 tokens ，同样受到 tokens 长度的限制
	Tools any `json:"tools,omitempty"`
	// 用于控制模型是如何选择要调用的函数，仅当工具类型为function时补充。默认为auto，当前仅支持auto
	ToolChoice any `json:"tool_choice,omitempty"`
	// 终端用户的唯一ID，协助平台对终端用户的违规行为、生成违法及不良信息或其他滥用行为进行干预。ID长度要求：最少6个字符，最多128个字符。
	UserId string `json:"user_id,omitempty"`
}

type ZhipuAIChatCompletionRes struct {
	// 任务ID
	Id string `json:"id"`
	// 请求创建时间，是以秒为单位的 Unix 时间戳
	Created int64 `json:"created"`
	// 模型名称
	Model string `json:"model"`
	// 当前对话的模型输出内容
	Choices []Choice `json:"choices"`
	// 结束时返回本次模型调用的 tokens 数量统计。
	Usage *Usage `json:"usage"`
	// 当failed时会有错误信息
	Error ZhipuAIError `json:"error"`
}

type Choice struct {
	// 结果下标
	Index int `json:"index"`
	// 模型推理终止的原因。
	// stop代表推理自然结束或触发停止词。
	// tool_calls 代表模型命中函数。
	// length代表到达 tokens 长度上限。
	// sensitive 代表模型推理内容被安全审核接口拦截。请注意，针对此类内容，请用户自行判断并决定是否撤回已公开的内容。
	// network_error 代表模型推理异常。
	FinishReason openai.FinishReason `json:"finish_reason"`
	// 模型返回的文本信息
	Message *openai.ChatCompletionMessage `json:"message,omitempty"`
	// 模型返回的文本信息-流式
	Delta *openai.ChatCompletionStreamChoiceDelta `json:"delta,omitempty"`
}

type ZhipuAIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ZhipuAIErrorResponse struct {
	Error *ZhipuAIError `json:"error,omitempty"`
}
