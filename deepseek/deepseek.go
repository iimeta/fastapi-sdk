package deepseek

import (
	"context"
	"errors"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/sdkerr"
)

type DeepSeek struct {
	model               string
	key                 string
	baseURL             string
	path                string
	proxyURL            string
	header              map[string]string
	isSupportSystemRole *bool
	isSupportStream     *bool
}

func NewAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *DeepSeek {

	logger.Infof(ctx, "NewAdapter DeepSeek model: %s, key: %s", model, key)

	deepseek := &DeepSeek{
		model:   model,
		key:     key,
		baseURL: "https://api.deepseek.com/v1",
		path:    "/chat/completions",
		header: g.MapStrStr{
			"Authorization": "Bearer " + key,
		},
		isSupportSystemRole: isSupportSystemRole,
		isSupportStream:     isSupportStream,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewAdapter DeepSeek model: %s, baseURL: %s", model, baseURL)
		deepseek.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewAdapter DeepSeek model: %s, path: %s", model, path)
		deepseek.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAdapter DeepSeek model: %s, proxyURL: %s", model, proxyURL[0])
		deepseek.proxyURL = proxyURL[0]
	}

	return deepseek
}

func NewAdapterBaidu(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *DeepSeek {

	logger.Infof(ctx, "NewAdapterBaidu DeepSeek model: %s, key: %s", model, key)

	split := gstr.Split(key, "|")

	baidu := &DeepSeek{
		model:   model,
		key:     split[1],
		baseURL: "https://qianfan.baidubce.com/v2",
		path:    "/chat/completions",
		header: g.MapStrStr{
			"appid": split[0],
		},
		isSupportSystemRole: isSupportSystemRole,
		isSupportStream:     isSupportStream,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewAdapterBaidu DeepSeek model: %s, baseURL: %s", model, baseURL)
		baidu.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewAdapterBaidu DeepSeek model: %s, path: %s", model, path)
		baidu.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAdapterBaidu DeepSeek model: %s, proxyURL: %s", model, proxyURL[0])
		baidu.proxyURL = proxyURL[0]
	}

	return baidu
}

func (d *DeepSeek) apiErrorHandler(err error) error {

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
