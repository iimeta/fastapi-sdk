package model

type BaiduChatCompletionReq struct {
	// 聊天上下文信息。说明：
	//（1）messages成员不能为空，1个成员表示单轮对话，多个成员表示多轮对话，例如：
	//  1个成员示例，"messages": [ {"role": "user","content": "你好"}]
	//  3个成员示例，"messages": [ {"role": "user","content": "你好"},{"role":"assistant","content":"需要什么帮助"},{"role":"user","content":"自我介绍下"}]
	//（2）最后一个message为当前请求的信息，前面的message为历史对话信息
	//（3）成员数目必须为奇数，成员中message的role值说明如下：奇数位message的role值必须为user，偶数位message的role值为assistant。例如：
	// 示例中message中的role值分别为user、assistant、user、assistant、user；奇数位（红框）message中的role值为user，即第1、3、5个message中的role值为user；偶数位（蓝框）值为assistant，即第2、4个message中的role值为assistant
	Messages []ChatCompletionMessage `json:"messages"`
	//（1）较高的数值会使输出更加随机，而较低的数值会使其更加集中和确定
	//（2）默认0.8，范围 (0, 1.0]，不能为0
	Temperature float32 `json:"temperature,omitempty"`
	//（1）影响输出文本的多样性，取值越大，生成文本的多样性越强
	//（2）默认0.8，取值范围 [0, 1.0]
	TopP float32 `json:"top_p,omitempty"`
	// 通过对已生成的token增加惩罚，减少重复生成的现象。说明：
	//（1）值越大表示惩罚越大
	//（2）默认1.0，取值范围：[1.0, 2.0]
	PenaltyScore float32 `json:"penalty_score,omitempty"`
	// 是否以流式接口的形式返回数据，默认false
	Stream bool `json:"stream,omitempty"`
	// 模型人设，主要用于人设设定，例如，你是xxx公司制作的AI助手，说明：
	//（1）长度限制，最后一个message的content长度（即此轮对话的问题）和system字段总内容不能超过20000个字符，且不能超过5120 tokens
	System string `json:"system,omitempty"`
	// 生成停止标识，当模型生成结果以stop中某个元素结尾时，停止文本生成。说明：
	//（1）每个元素长度不超过20字符
	//（2）最多4个元素
	Stop []string `json:"stop,omitempty"`
	// 是否强制关闭实时搜索功能，默认false，表示不关闭
	DisableSearch bool `json:"disable_search,omitempty"`
	// 是否开启上角标返回，说明：
	//（1）开启后，有概率触发搜索溯源信息search_info，search_info内容见响应参数介绍
	//（2）默认false，不开启
	EnableCitation bool `json:"enable_citation,omitempty"`
	// 指定模型最大输出token数，说明：
	//（1）如果设置此参数，范围[2, 2048]
	//（2）如果不设置此参数，最大输出token数为1024
	MaxOutputTokens int `json:"max_output_tokens,omitempty"`
	// 指定响应内容的格式，说明：
	//（1）可选值：
	//  json_object：以json格式返回，可能出现不满足效果情况
	//  text：以文本格式返回
	//（2）如果不填写参数response_format值，默认为text
	ResponseFormat string `json:"response_format,omitempty"`
	// 表示最终用户的唯一标识符
	UserId string `json:"user_id,omitempty"`
}

type BaiduChatCompletionRes struct {
	// 本轮对话的id
	Id string `json:"id"`
	// 回包类型
	// chat.completion：多轮对话返回
	Object string `json:"object"`
	// 时间戳
	Created int64 `json:"created"`
	// 表示当前子句的序号。只有在流式接口模式下会返回该字段
	SentenceId int `json:"sentence_id"`
	// 表示当前子句是否是最后一句。只有在流式接口模式下会返回该字段
	IsEnd bool `json:"is_end"`
	// 当前生成的结果是否被截断
	IsTruncated bool `json:"is_truncated"`
	// 输出内容标识，说明：
	//  normal：输出内容完全由大模型生成，未触发截断、替换
	//  stop：输出结果命中入参stop中指定的字段后被截断
	//  length：达到了最大的token数，根据EB返回结果is_truncated来截断
	//  content_filter：输出内容被截断、兜底、替换为**等
	FinishReason string `json:"finish_reason"`
	// 搜索数据，当请求参数enable_citation为true并且触发搜索时，会返回该字段
	SearchInfo *SearchInfo
	// 对话返回结果
	Result string `json:"result"`
	// 表示用户输入是否存在安全风险，是否关闭当前会话，清理历史会话信息
	// true：是，表示用户输入存在安全风险，建议关闭当前会话，清理历史会话信息
	// false：否，表示用户输入无安全风险
	NeedClearHistory bool `json:"need_clear_history"`
	// 0：正常返回
	// 其他：非正常
	Flag int `json:"flag"`
	// 当need_clear_history为true时，此字段会告知第几轮对话有敏感信息，如果是当前问题，ban_round=-1
	BanRound int `json:"ban_round"`
	// token统计信息
	Usage     *Usage `json:"usage,omitempty"`
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

type SearchInfo struct {
	SearchResults []SearchResult `json:"search_results,omitempty"` // 搜索结果列表
}

type SearchResult struct {
	Index int    `json:"index,omitempty"` // 序号
	Url   string `json:"url,omitempty"`   // 搜索结果URL
	Title string `json:"title,omitempty"` // 搜索结果标题
}
