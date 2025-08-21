package volcengine

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/options"
	"github.com/iimeta/fastapi-sdk/sdkerr"
)

type VolcEngine struct {
	*options.AdapterOptions
	header map[string]string
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) *VolcEngine {

	split := gstr.Split(options.Key, "|")

	volcengine := &VolcEngine{
		AdapterOptions: options,
		header: g.MapStrStr{
			"Authorization": "Bearer " + split[1],
		},
	}

	volcengine.Model = split[0]
	volcengine.Key = split[1]

	if volcengine.BaseUrl == "" {
		volcengine.BaseUrl = "https://ark.cn-beijing.volces.com/api/v3"
	}

	if volcengine.Path == "" {
		volcengine.Path = "/chat/completions"
	}

	logger.Infof(ctx, "NewAdapter VolcEngine model: %s, key: %s", volcengine.Model, volcengine.Key)

	return volcengine
}

func (v *VolcEngine) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return sdkerr.NewRequestError(500, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)))
}

func (v *VolcEngine) apiErrorHandler(err error) error {

	apiError := &sdkerr.ApiError{}
	if errors.As(err, &apiError) {

		switch apiError.HttpStatusCode {
		case 400:
			if apiError.Code == "context_length_exceeded" {
				return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
			}
		case 401:
			if apiError.Code == "invalid_api_key" {
				return sdkerr.ERR_INVALID_API_KEY
			}
		case 404:
			return sdkerr.ERR_MODEL_NOT_FOUND
		case 429:
			if apiError.Code == "insufficient_quota" {
				return sdkerr.ERR_INSUFFICIENT_QUOTA
			}
		}

		return err
	}

	reqError := &sdkerr.RequestError{}
	if errors.As(err, &reqError) {
		return sdkerr.NewRequestError(apiError.HttpStatusCode, reqError.Err)
	}

	return err
}
