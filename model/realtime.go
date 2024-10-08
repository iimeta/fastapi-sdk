package model

type RealtimeRequest struct {
	MessageType int    `json:"message_type"`
	Message     []byte `json:"message"`
}

type RealtimeResponse struct {
	MessageType int    `json:"message_type"`
	Message     []byte `json:"message"`
	Usage       *Usage `json:"usage"`
	ConnTime    int64  `json:"-"`
	Duration    int64  `json:"-"`
	TotalTime   int64  `json:"-"`
	Error       error  `json:"-"`
}
