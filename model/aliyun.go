package model

import "github.com/iimeta/go-openai"

type AliyunChatCompletionReq struct {
	// 指定用于对话的通义千问模型名
	// 目前可选择qwen-turbo、qwen-plus、qwen-max、qwen-max-0403、qwen-max-0107、qwen-max-1201和qwen-max-longcontext。
	Model      string     `json:"model"`
	Input      Input      `json:"input"`
	Parameters Parameters `json:"parameters"`
}

type Input struct {
	// 用户与模型的对话历史，list中的每个元素形式为{"role":角色, "content": 内容}。
	// 角色当前可选值：system、user、assistant和tool。未来可以扩展到更多role。
	Messages []ChatCompletionMessage `json:"messages"`
}

type Parameters struct {
	// "text"表示旧版本的text
	// "message"表示兼容openai的message
	ResultFormat string `json:"resultFormat,omitempty"`
	// 生成时使用的随机数种子，用户控制模型生成内容的随机性。
	// seed支持无符号64位整数，默认值为1234。
	// 在使用seed时，模型将尽可能生成相同或相似的结果，但目前不保证每次生成的结果完全相同。
	Seed *int `json:"seed,omitempty"`
	// 用于限制模型生成token的数量，max_tokens设置的是生成上限，并不表示一定会生成这么多的token数量。
	// 其中qwen-turbo最大值和默认值为1500，qwen-max、qwen-max-1201 、qwen-max-longcontext 和 qwen-plus最大值和默认值均为2000。
	MaxTokens int `json:"max_tokens,omitempty"`
	// 生成时，核采样方法的概率阈值。
	// 例如，取值为0.8时，仅保留累计概率之和大于等于0.8的概率分布中的token，作为随机采样的候选集。
	// 取值范围为（0,1.0)，取值越大，生成的随机性越高；取值越低，生成的随机性越低。
	// 默认值为0.8。注意，取值不要大于等于1
	TopP float32 `json:"top_p,omitempty"`
	// 生成时，采样候选集的大小。
	// 例如，取值为50时，仅将单次生成中得分最高的50个token组成随机采样的候选集。
	// 取值越大，生成的随机性越高；取值越小，生成的确定性越高。
	// 注意：如果top_k参数为空或者top_k的值大于100，表示不启用top_k策略，此时仅有top_p策略生效，默认是空。
	TopK int `json:"top_k,omitempty"`
	// 用于控制模型生成时的重复度。
	// 提高repetition_penalty时可以降低模型生成的重复度。
	// 1.0表示不做惩罚。默认为1.1。
	RepetitionPenalty float32 `json:"repetition_penalty,omitempty"`
	// 用于控制随机性和多样性的程度。
	// 具体来说，temperature值控制了生成文本时对每个候选词的概率分布进行平滑的程度。
	// 较高的temperature值会降低概率分布的峰值，使得更多的低概率词被选择，生成结果更加多样化；
	// 而较低的temperature值则会增强概率分布的峰值，使得高概率词更容易被选择，生成结果更加确定。
	// 取值范围：[0, 2)，系统默认值0.85。不建议取值为0，无意义。
	Temperature float32 `json:"temperature,omitempty"`
	// stop参数用于实现内容生成过程的精确控制，在生成内容即将包含指定的字符串或token_ids时自动停止，生成内容不包含指定的内容。
	// 例如，如果指定stop为"你好"，表示将要生成"你好"时停止；如果指定stop为[37763, 367]，表示将要生成"Observation"时停止。
	// stop参数支持以list方式传入字符串数组或者token_ids数组，支持使用多个stop的场景。
	Stop []string `json:"stop,omitempty"` // String/List[String]用于指定字符串；List[Integer]/List[List[Integer]]用于指定token_ids；注意: list模式下不支持字符串和token_ids混用，元素类型要相同。
	// 模型内置了互联网搜索服务，该参数控制模型在生成文本时是否参考使用互联网搜索结果。取值如下：
	// true：启用互联网搜索，模型会将搜索结果作为文本生成过程中的参考信息，但模型会基于其内部逻辑“自行判断”是否使用互联网搜索结果。
	// false（默认）：关闭互联网搜索。
	EnableSearch bool `json:"enable_search,omitempty"`
	// 用于控制流式输出模式，默认false，即后面内容会包含已经输出的内容；
	// 设置为true，将开启增量输出模式，后面输出不会包含已经输出的内容，您需要自行拼接整体输出，参考流式输出示例代码。
	// 该参数只能与stream输出模式配合使用。
	// 注意: incremental_output暂时无法和tools参数同时使用。
	IncrementalOutput bool `json:"incremental_output,omitempty"`
	// 模型可选调用的工具列表。目前仅支持function，并且即使输入多个function，模型仅会选择其中一个生成结果。
	// 模型根据tools参数内容可以生产函数调用的参数，tools中每一个tool的结构如下：
	// type，类型为string，表示tools的类型，当前仅支持function。
	// function，类型为dict，包括name，description和parameters：
	// name，类型为string，表示function的名称，必须是字母、数字，或包含下划线和短划线，最大长度为64。
	// description，类型为string，表示function的描述，供模型选择何时以及如何调用function。
	// parameters，类型为dict，表示function的参数描述，需要是一个合法的json schema。json schema的描述可以见链接。参考代码中给出了一个参数描述的示例。如果parameters参数缺省了，表示function没有入参。
	// 使用tools功能时需要指定result_format为message。
	// 在多轮对话中，无论是发起function_call的轮次，还是向模型提交function的执行结果，均请设置tools参数。
	// 当前支持qwen-turbo、qwen-plus、qwen-max和qwen-max-longcontext。
	// 注意: tools暂时无法和incremental_output参数同时使用。
	Tools any `json:"tools,omitempty"`
}

type AliyunChatCompletionRes struct {
	// 入参result_format=text时候的返回值
	Output Output `json:"output"`
	Usage  struct {
		// 本次请求输入内容的 token 数目。
		// 在打开了搜索的情况下，输入的 token 数目因为还需要添加搜索相关内容支持，所以会超出客户在请求中的输入。
		InputTokens int `json:"input_tokens"`
		// 本次请求算法输出内容的 token 数目。
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	// 本次请求的系统唯一码。
	RequestId string `json:"request_id"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

type Output struct {
	// 包含本次请求的算法输出内容。
	Text string `json:"text"`
	// 有三种情况：正在生成时为null，生成结束时如果由于停止token导致则为stop，生成结束时如果因为生成长度过长导致则为length。
	FinishReason openai.FinishReason `json:"finish_reason"`
	// 入参result_format=message时候的返回值
	Choices []ChatCompletionChoice `json:"choices"`
}
