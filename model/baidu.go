package model

import "github.com/sashabaranov/go-openai"

type ErnieBotMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ErnieBotReq struct {
	Messages []openai.ChatCompletionMessage `json:"messages"`
}
type ErnieBotRes struct {
	Id               string `json:"id"`
	Object           string `json:"object"`
	Created          int64  `json:"created"`
	Result           string `json:"result"`
	IsTruncated      bool   `json:"is_truncated"`
	NeedClearHistory bool   `json:"need_clear_history"`
	Usage            struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

type GetAccessTokenRes struct {
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in"`
	SessionKey       string `json:"session_key"`
	AccessToken      string `json:"access_token"`
	Scope            string `json:"scope"`
	SessionSecret    string `json:"session_secret"`
	ErrorDescription string `json:"error_description"`
	Error            string `json:"error"`
}
