package model

type MidjourneyProxyRequest struct {
	Prompt        string         `json:"prompt,omitempty"`
	Base64        string         `json:"base64,omitempty"`
	Base64Array   []string       `json:"base64Array,omitempty"`
	Action        string         `json:"action,omitempty"`
	Index         int            `json:"index,omitempty"`
	TaskId        string         `json:"taskId,omitempty"`
	SourceBase64  string         `json:"sourceBase64,omitempty"`
	TargetBase64  string         `json:"targetBase64,omitempty"`
	NotifyHook    string         `json:"notifyHook,omitempty"`
	State         string         `json:"state,omitempty"`
	BotType       string         `json:"botType,omitempty"`
	Dimensions    string         `json:"dimensions,omitempty"`
	AccountFilter *AccountFilter `json:"accountFilter,omitempty"`
	MaskBase64    string         `json:"maskBase64,omitempty"`
	Filter        *Filter        `json:"filter,omitempty"`
}

type AccountFilter struct {
	ChannelId           string   `json:"channelId,omitempty"`
	InstanceId          string   `json:"instanceId,omitempty"`
	Modes               []string `json:"modes,omitempty"`
	Remark              string   `json:"remark,omitempty"`
	Remix               bool     `json:"remix,omitempty"`
	RemixAutoConsidered bool     `json:"remixAutoConsidered,omitempty"`
}

type Filter struct {
	ChannelId  string `json:"channelId,omitempty"`
	InstanceId string `json:"instanceId,omitempty"`
	Remark     string `json:"remark,omitempty"`
}

type MidjourneyProxyResponse struct {
	Code        int         `json:"code,omitempty"`
	Description string      `json:"description,omitempty"`
	Result      string      `json:"result,omitempty"`
	Properties  *Properties `json:"properties,omitempty"`
	TotalTime   int64       `json:"-"`
}

type MidjourneyResponse struct {
	Response  []byte `json:"response,omitempty"`
	TotalTime int64  `json:"-"`
}

type MidjourneyProxyFetchResponse struct {
	Id          string      `json:"id,omitempty"`
	Action      string      `json:"action,omitempty"`
	Buttons     []*Button   `json:"buttons,omitempty"`
	Description string      `json:"description,omitempty"`
	FailReason  string      `json:"failReason,omitempty"`
	ImageUrl    string      `json:"imageUrl,omitempty"`
	Progress    string      `json:"progress,omitempty"`
	Prompt      string      `json:"prompt,omitempty"`
	PromptEn    string      `json:"promptEn,omitempty"`
	Properties  *Properties `json:"properties,omitempty"`
	SubmitTime  int         `json:"submitTime,omitempty"`
	StartTime   int         `json:"startTime,omitempty"`
	FinishTime  int         `json:"finishTime,omitempty"`
	State       string      `json:"state,omitempty"`
	Status      string      `json:"status,omitempty"`
	TotalTime   int64       `json:"-"`
}

type Properties struct {
	NotifyHook        string `json:"notifyHook,omitempty"`
	FinalPrompt       string `json:"finalPrompt,omitempty"`
	MessageId         string `json:"messageId,omitempty"`
	MessageHash       string `json:"messageHash,omitempty"`
	ProgressMessageId string `json:"progressMessageId,omitempty"`
	Flags             int    `json:"flags,omitempty"`
	Nonce             string `json:"nonce,omitempty"`
	DiscordInstanceId string `json:"discordInstanceId,omitempty"`
	PromptEn          string `json:"promptEn,omitempty"`
	BannedWord        string `json:"bannedWord,omitempty"`
}

type Button struct {
	CustomId string `json:"customId,omitempty"`
	Emoji    string `json:"emoji,omitempty"`
	Label    string `json:"label,omitempty"`
	Style    int    `json:"style,omitempty"`
	Type     int    `json:"type,omitempty"`
}
