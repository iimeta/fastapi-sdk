package model

type OpenAIResponsesReq struct {
	Model  string `json:"model,omitempty"`
	Input  any    `json:"input"`
	Stream bool   `json:"stream,omitempty"`
}

type OpenAIResponsesRes struct {
	Model string `json:"model"`
	Usage *Usage `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
	ResponseBytes []byte `json:"-"`
	ConnTime      int64  `json:"-"`
	Duration      int64  `json:"-"`
	TotalTime     int64  `json:"-"`
	Err           error  `json:"-"`
}
