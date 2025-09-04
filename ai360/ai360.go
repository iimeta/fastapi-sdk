package ai360

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/iimeta/fastapi-sdk/errors"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/options"
)

type AI360 struct {
	*options.AdapterOptions
	header map[string]string
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) *AI360 {

	ai360 := &AI360{
		AdapterOptions: options,
		header: g.MapStrStr{
			"Authorization": "Bearer " + options.Key,
		},
	}

	if ai360.BaseUrl == "" {
		ai360.BaseUrl = "https://api.360.cn/v1"
	}

	if ai360.Path == "" {
		ai360.Path = "/chat/completions"
	}

	logger.Infof(ctx, "NewAdapter 360AI model: %s, key: %s", ai360.Model, ai360.Key)

	return ai360
}

func (a *AI360) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return errors.NewRequestError(response.StatusCode, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)))
}

func (a *AI360) apiErrorHandler(err error) error {

	apiError := &errors.ApiError{}
	if errors.As(err, &apiError) {

		switch apiError.HttpStatusCode {
		case 400:
			if apiError.Code == "1001" {
				return errors.ERR_CONTEXT_LENGTH_EXCEEDED
			}
		case 401:

			if apiError.Code == "1002" {
				return errors.ERR_INVALID_API_KEY
			}

			if apiError.Code == "1004" || apiError.Code == "1006" {
				return errors.ERR_INSUFFICIENT_QUOTA
			}

		case 404:
			return errors.ERR_MODEL_NOT_FOUND
		case 429:
			if apiError.Code == "1005" {
				return errors.ERR_CONTEXT_LENGTH_EXCEEDED
			}
		}

		return err
	}

	reqError := &errors.RequestError{}
	if errors.As(err, &reqError) {
		return errors.NewRequestError(apiError.HttpStatusCode, reqError.Err)
	}

	return err
}
