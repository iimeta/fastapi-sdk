package model

type ModerationRequest struct {
	Model string `json:"model"`
	Input any    `json:"input"`
}

type ModerationResponse struct {
	Id        string `json:"id,omitempty"`
	Model     string `json:"model,omitempty"`
	Results   any    `json:"results,omitempty"`
	Error     any    `json:"error,omitempty"`
	Usage     *Usage `json:"-"`
	TotalTime int64  `json:"-"`
}
