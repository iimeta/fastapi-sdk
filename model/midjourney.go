package model

type MidjourneyProxy struct {
	ApiSecret       string `json:"api_secret"`
	ApiSecretHeader string `json:"api_secret_header"`
	ImagineUrl      string `json:"imagine_url"`
	ChangeUrl       string `json:"change_url"`
	DescribeUrl     string `json:"describe_url"`
	BlendUrl        string `json:"blend_url"`
	FetchUrl        string `json:"fetch_url"`
}

type MidjourneyProxyRequest struct {
	Prompt      string   `json:"prompt"`
	Base64      string   `json:"base64"`
	Base64Array []string `json:"base64Array"`
	Action      string   `json:"action"`
	Index       int      `json:"index"`
	TaskId      string   `json:"taskId"`
}

type MidjourneyProxyResponse struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	Result      string `json:"result"`
	Properties  struct {
		PromptEn   string `json:"promptEn"`
		BannedWord string `json:"bannedWord"`
	} `json:"properties"`
	TotalTime int64 `json:"-"`
}

type MidjourneyProxyFetchResponse struct {
	Action      string      `json:"action"`
	Id          string      `json:"id"`
	Prompt      string      `json:"prompt"`
	PromptEn    string      `json:"promptEn"`
	Description string      `json:"description"`
	State       interface{} `json:"state"`
	SubmitTime  int64       `json:"submitTime"`
	StartTime   int64       `json:"startTime"`
	FinishTime  int64       `json:"finishTime"`
	ImageUrl    string      `json:"imageUrl"`
	Status      string      `json:"status"`
	Progress    string      `json:"progress"`
	FailReason  string      `json:"failReason"`
	TotalTime   int64       `json:"-"`
}
