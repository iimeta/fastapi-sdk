package deepseek

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/options"
)

type DeepSeek struct {
	*options.AdapterOptions
	header map[string]string
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) *DeepSeek {

	deepseek := &DeepSeek{
		AdapterOptions: options,
		header: map[string]string{
			"Authorization": "Bearer " + options.Key,
		},
	}

	if deepseek.BaseUrl == "" {
		deepseek.BaseUrl = "https://api.deepseek.com/v1"
	}

	for k, v := range deepseek.Header {
		deepseek.header[k] = v
	}

	logger.Infof(ctx, "NewAdapter DeepSeek model: %s, key: %s", deepseek.Model, deepseek.Key)

	return deepseek
}

func NewAdapterBaidu(ctx context.Context, options *options.AdapterOptions) *DeepSeek {

	split := gstr.Split(options.Key, "|")

	baidu := &DeepSeek{
		AdapterOptions: options,
		header: map[string]string{
			"appid": split[0],
		},
	}

	baidu.Key = split[1]

	if baidu.BaseUrl == "" {
		baidu.BaseUrl = "https://qianfan.baidubce.com/v2"
	}

	for k, v := range baidu.Header {
		baidu.header[k] = v
	}

	logger.Infof(ctx, "NewAdapterBaidu DeepSeek model: %s, key: %s", baidu.Model, baidu.Key)

	return baidu
}

func (d *DeepSeek) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return errors.NewRequestError(response.StatusCode, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)))
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
