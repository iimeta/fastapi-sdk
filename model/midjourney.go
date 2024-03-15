package model

type MidjourneyProxy struct {
	ApiSecret              string `json:"api_secret"`
	ApiSecretHeader        string `json:"api_secret_header"`
	ImagineUrl             string `json:"imagine_url"`
	ChangeUrl              string `json:"change_url"`
	DescribeUrl            string `json:"describe_url"`
	BlendUrl               string `json:"blend_url"`
	SwapFaceUrl            string `json:"swap_face_url"`
	ActionUrl              string `json:"action_url"`
	ModalUrl               string `json:"modal_url"`
	ShortenUrl             string `json:"shorten_url"`
	UploadDiscordImagesUrl string `json:"upload_discord_images_url"`
	FetchUrl               string `json:"fetch_url"`
}

type MidjourneyProxyRequest struct {
	Prompt        string   `json:"prompt"`
	Base64        string   `json:"base64"`
	Base64Array   []string `json:"base64Array"`
	Action        string   `json:"action"`
	Index         int      `json:"index"`
	TaskId        string   `json:"taskId"`
	SourceBase64  string   `json:"sourceBase64"`
	TargetBase64  string   `json:"targetBase64"`
	NotifyHook    string   `json:"notifyHook"`
	State         string   `json:"state"`
	BotType       string   `json:"botType"`
	Dimensions    string   `json:"dimensions"`
	AccountFilter struct {
		ChannelId           string   `json:"channelId"`
		InstanceId          string   `json:"instanceId"`
		Modes               []string `json:"modes"`
		Remark              string   `json:"remark"`
		Remix               bool     `json:"remix"`
		RemixAutoConsidered bool     `json:"remixAutoConsidered"`
	} `json:"accountFilter"`
	MaskBase64 string `json:"maskBase64"`
	Filter     struct {
		ChannelId  string `json:"channelId"`
		InstanceId string `json:"instanceId"`
		Remark     string `json:"remark"`
	} `json:"filter"`
}

type MidjourneyProxyResponse struct {
	Code        int         `json:"code"`
	Description string      `json:"description"`
	Result      interface{} `json:"result"`
	Properties  struct {
		PromptEn   string `json:"promptEn"`
		BannedWord string `json:"bannedWord"`
	} `json:"properties"`
	TotalTime int64 `json:"-"`
}

type MidjourneyProxyFetchResponse struct {
	Action  string `json:"action"`
	Buttons []struct {
		CustomId string `json:"customId"`
		Emoji    string `json:"emoji"`
		Label    string `json:"label"`
		Style    int    `json:"style"`
		Type     int    `json:"type"`
	} `json:"buttons"`
	Description string `json:"description"`
	FailReason  string `json:"failReason"`
	FinishTime  int    `json:"finishTime"`
	Id          string `json:"id"`
	ImageUrl    string `json:"imageUrl"`
	Progress    string `json:"progress"`
	Prompt      string `json:"prompt"`
	PromptEn    string `json:"promptEn"`
	Properties  struct {
	} `json:"properties"`
	StartTime  int    `json:"startTime"`
	State      string `json:"state"`
	Status     string `json:"status"`
	SubmitTime int    `json:"submitTime"`
	TotalTime  int64  `json:"-"`
}
