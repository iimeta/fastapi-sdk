package sdkerr

import (
	"encoding/json"
	"fmt"
	"strings"
)

var (
	ERR_CONTEXT_LENGTH_EXCEEDED = NewApiError(400, "context_length_exceeded", "Please reduce the length of the messages.", "invalid_request_error", "messages")
	ERR_INVALID_API_KEY         = NewApiError(401, "invalid_api_key", "Incorrect API key provided or has been disabled.", "invalid_request_error", "")
	ERR_MODEL_NOT_FOUND         = NewApiError(404, "model_not_found", "The model does not exist or you do not have access to it.", "invalid_request_error", "")
	ERR_INSUFFICIENT_QUOTA      = NewApiError(429, "insufficient_quota", "You exceeded your current quota.", "insufficient_quota", "")
	ERR_RATE_LIMIT_EXCEEDED     = NewApiError(429, "rate_limit_exceeded", "Rate limit reached, Please try again later.", "requests", "")
)

// ApiError provides error information returned by the OpenAI API.
// InnerError struct is only valid for Azure OpenAI Service.
type ApiError struct {
	HttpStatusCode int     `json:"-"`
	Code           any     `json:"code,omitempty"`
	Message        string  `json:"message"`
	Type           string  `json:"type"`
	Param          *string `json:"param,omitempty"`
}

// RequestError provides information about generic request sdkerr.
type RequestError struct {
	HttpStatusCode int
	Err            error
}

type ErrorResponse struct {
	Error *ApiError `json:"error,omitempty"`
}

func (e *ApiError) Error() string {
	if e.HttpStatusCode > 0 {
		return fmt.Sprintf("error, status code: %d, response: %s", e.HttpStatusCode, e.Message)
	}

	return e.Message
}

func (e *ApiError) UnmarshalJSON(data []byte) (err error) {
	var rawMap map[string]json.RawMessage
	err = json.Unmarshal(data, &rawMap)
	if err != nil {
		return
	}

	err = json.Unmarshal(rawMap["message"], &e.Message)
	if err != nil {
		// If the parameter field of a function call is invalid as a JSON schema
		// refs: https://github.com/iimeta/go-openai/issues/381
		var messages []string
		err = json.Unmarshal(rawMap["message"], &messages)
		if err != nil {
			return
		}
		e.Message = strings.Join(messages, ", ")
	}

	// optional fields for azure openai
	// refs: https://github.com/iimeta/go-openai/issues/343
	if _, ok := rawMap["type"]; ok {
		err = json.Unmarshal(rawMap["type"], &e.Type)
		if err != nil {
			return
		}
	}

	// optional fields
	if _, ok := rawMap["param"]; ok {
		err = json.Unmarshal(rawMap["param"], &e.Param)
		if err != nil {
			return
		}
	}

	if _, ok := rawMap["code"]; !ok {
		return nil
	}

	// if the api returned a number, we need to force an integer
	// since the json package defaults to float64
	var intCode int
	err = json.Unmarshal(rawMap["code"], &intCode)
	if err == nil {
		e.Code = intCode
		return nil
	}

	return json.Unmarshal(rawMap["code"], &e.Code)
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("error, status code: %d, response: %s", e.HttpStatusCode, e.Err)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}

func NewApiError(httpStatusCode int, code any, message, typ, param string) error {
	return &ApiError{
		HttpStatusCode: httpStatusCode,
		Code:           code,
		Message:        message,
		Type:           typ,
		Param:          &param,
	}
}

func NewRequestError(httpStatusCode int, err error) error {
	return &RequestError{
		HttpStatusCode: httpStatusCode,
		Err:            err,
	}
}
