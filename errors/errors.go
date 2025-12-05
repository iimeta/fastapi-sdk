package errors

import (
	"encoding/json"
	"errors"
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

type ApiError struct {
	HttpStatusCode int    `json:"-"`
	Code           any    `json:"code"`
	Message        string `json:"message"`
	Type           string `json:"type"`
	Param          any    `json:"param"`
}

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

	if err = json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	if _, ok := rawMap["code"]; ok {
		if err = json.Unmarshal(rawMap["code"], &e.Code); err != nil {
			return err
		}
	}

	if err = json.Unmarshal(rawMap["message"], &e.Message); err != nil {

		var messages []string
		if err = json.Unmarshal(rawMap["message"], &messages); err != nil {
			return err
		}

		e.Message = strings.Join(messages, "; ")
	}

	if _, ok := rawMap["type"]; ok {
		if err = json.Unmarshal(rawMap["type"], &e.Type); err != nil {
			return err
		}
	}

	if _, ok := rawMap["param"]; ok {
		if err = json.Unmarshal(rawMap["param"], &e.Param); err != nil {
			return err
		}
	}

	return nil
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("error, status code: %d, response: %s", e.HttpStatusCode, e.Err)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}

func NewApiError(httpStatusCode int, code any, message, typ string, param any) error {
	return &ApiError{
		HttpStatusCode: httpStatusCode,
		Code:           code,
		Message:        message,
		Type:           typ,
		Param:          param,
	}
}

func NewRequestError(httpStatusCode int, err error) error {
	return &RequestError{
		HttpStatusCode: httpStatusCode,
		Err:            err,
	}
}

func New(text string) error {
	return errors.New(text)
}

func Newf(format string, args ...any) error {
	return errors.New(fmt.Sprintf(format, args...))
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}
