package xai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/options"
)

type XAI struct {
	*options.AdapterOptions
	header map[string]string
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) *XAI {

	xai := &XAI{
		AdapterOptions: options,
		header: map[string]string{
			"Authorization": "Bearer " + options.Key,
		},
	}

	if xai.BaseUrl == "" {
		xai.BaseUrl = "https://api.x.ai/v1"
	}

	for k, v := range xai.PassthroughHeader {
		xai.header[k] = v
	}

	for k, v := range xai.Header {
		xai.header[k] = v
	}

	logger.Infof(ctx, "NewAdapter xAI model: %s, key: %s", xai.Model, xai.Key)

	return xai
}

func (x *XAI) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	errorResponse := errors.ErrorResponse{}
	if err := json.Unmarshal(bytes, &errorResponse); err != nil || errorResponse.Error == nil {
		return &errors.RequestError{
			HttpStatusCode: response.StatusCode,
			Err:            errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)),
		}
	}

	return errors.NewApiError(response.StatusCode, errorResponse.Error.Code, errorResponse.Error.Message, errorResponse.Error.Type, errorResponse.Error.Param)
}

func (x *XAI) apiErrorHandler(err error) error {
	return err
}
