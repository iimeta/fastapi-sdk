package model

type BatchCreateRequest struct {
	Model              string             `json:"model,omitempty"`
	InputFileId        string             `json:"input_file_id"`
	Endpoint           string             `json:"endpoint"`
	CompletionWindow   string             `json:"completion_window"`
	Metadata           any                `json:"metadata"`
	OutputExpiresAfter OutputExpiresAfter `json:"output_expires_after"`
}

type OutputExpiresAfter struct {
	Anchor  string `json:"anchor"`
	Seconds int    `json:"seconds"`
}

type BatchListRequest struct {
	After string `json:"after"`
	Limit int64  `json:"limit"`
}

type BatchListResponse struct {
	Object    string          `json:"object"`
	Data      []BatchResponse `json:"data"`
	FirstId   *string         `json:"first_id"`
	LastId    *string         `json:"last_id"`
	HasMore   bool            `json:"has_more"`
	TotalTime int64           `json:"-"`
}

type BatchRetrieveRequest struct {
	BatchId string `json:"batch_id"`
}

type BatchCancelRequest struct {
	BatchId string `json:"batch_id"`
}

type BatchResponse struct {
	Id               string        `json:"id"`
	Object           string        `json:"object"`
	Endpoint         string        `json:"endpoint"`
	Model            string        `json:"model"`
	InputFileId      string        `json:"input_file_id"`
	CompletionWindow string        `json:"completion_window"`
	Status           string        `json:"status"`
	OutputFileId     string        `json:"output_file_id"`
	ErrorFileId      string        `json:"error_file_id"`
	CreatedAt        int64         `json:"created_at"`
	InProgressAt     int64         `json:"in_progress_at"`
	ExpiresAt        int64         `json:"expires_at"`
	FinalizingAt     int64         `json:"finalizing_at"`
	CompletedAt      int64         `json:"completed_at"`
	FailedAt         int64         `json:"failed_at"`
	ExpiredAt        int64         `json:"expired_at"`
	CancellingAt     int64         `json:"cancelling_at"`
	CancelledAt      int64         `json:"cancelled_at"`
	RequestCounts    RequestCounts `json:"request_counts"`
	Metadata         any           `json:"metadata"`
	Usage            Usage         `json:"usage"`
	Errors           *BatchError   `json:"errors"`
	TotalTime        int64         `json:"-"`
}

type RequestCounts struct {
	Completed int `json:"completed"`
	Total     int `json:"total"`
	Failed    int `json:"failed"`
}

type BatchRequestInput struct {
	CustomId string `json:"custom_id"`
	Method   string `json:"method"`
	Url      string `json:"url"`
	Body     struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	} `json:"body"`
}

type BatchRequestOutput struct {
	Id       string `json:"id"`
	CustomId string `json:"custom_id"`
	Response struct {
		StatusCode int    `json:"status_code"`
		RequestId  string `json:"request_id"`
		Body       struct {
			Id      string `json:"id"`
			Object  string `json:"object"`
			Created int    `json:"created"`
			Model   string `json:"model"`
			Choices []struct {
				Index   int `json:"index"`
				Message struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				} `json:"message"`
				FinishReason string `json:"finish_reason"`
			} `json:"choices"`
			Usage struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				TotalTokens      int `json:"total_tokens"`
			} `json:"usage"`
			SystemFingerprint interface{} `json:"system_fingerprint"`
		} `json:"body"`
	} `json:"response"`
	Error interface{} `json:"error"`
}

type BatchError struct {
	Object string `json:"object"`
	Data   []struct {
		Code    string `json:"code"`
		Line    int    `json:"line"`
		Message string `json:"message"`
		Param   string `json:"param"`
	}
}
