package deepseek

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/errors"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/options"
)

type DeepSeek struct {
	*options.AdapterOptions
	header map[string]string
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) *DeepSeek {

	deepseek := &DeepSeek{
		AdapterOptions: options,
		header: g.MapStrStr{
			"Authorization": "Bearer " + options.Key,
		},
	}

	if deepseek.BaseUrl == "" {
		deepseek.BaseUrl = "https://api.deepseek.com/v1"
	}

	if deepseek.Path == "" {
		deepseek.Path = "/chat/completions"
	}

	logger.Infof(ctx, "NewAdapter DeepSeek model: %s, key: %s", deepseek.Model, deepseek.Key)

	return deepseek
}

func NewAdapterBaidu(ctx context.Context, options *options.AdapterOptions) *DeepSeek {

	split := gstr.Split(options.Key, "|")

	baidu := &DeepSeek{
		AdapterOptions: options,
		header: g.MapStrStr{
			"appid": split[0],
		},
	}

	baidu.Key = split[1]

	if baidu.BaseUrl == "" {
		baidu.BaseUrl = "https://qianfan.baidubce.com/v2"
	}

	if baidu.Path == "" {
		baidu.Path = "/chat/completions"
	}

	logger.Infof(ctx, "NewAdapterBaidu DeepSeek model: %s, key: %s", baidu.Model, baidu.Key)

	return baidu
}

func (d *DeepSeek) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return errors.NewRequestError(500, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)))
}

func (d *DeepSeek) apiErrorHandler(err error) error {

	apiError := &errors.ApiError{}
	if errors.As(err, &apiError) {

		switch apiError.HttpStatusCode {
		case 400:
			if apiError.Code == "context_length_exceeded" {
				return errors.ERR_CONTEXT_LENGTH_EXCEEDED
			}
		case 401:
			if apiError.Code == "invalid_api_key" {
				return errors.ERR_INVALID_API_KEY
			}
		case 404:
			return errors.ERR_MODEL_NOT_FOUND
		case 429:
			if apiError.Code == "insufficient_quota" {
				return errors.ERR_INSUFFICIENT_QUOTA
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
