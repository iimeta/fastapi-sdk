package model

type RealtimeRequest struct {
	Model       string `json:"model"`
	MessageType int    `json:"message_type"`
	Message     []byte `json:"message"`
}

type RealtimeResponse struct {
	Message   []byte `json:"-"`
	Usage     *Usage `json:"usage"`
	ConnTime  int64  `json:"-"`
	Duration  int64  `json:"-"`
	TotalTime int64  `json:"-"`
	Error     error  `json:"-"`
}
