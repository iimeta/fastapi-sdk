package deepseek

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/sdkerr"
	"github.com/iimeta/go-openai"
)

type DeepSeek struct {
	client              *openai.Client
	isSupportSystemRole *bool
}

func NewAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *DeepSeek {

	logger.Infof(ctx, "NewAdapter DeepSeek model: %s, key: %s", model, key)

	config := openai.DefaultConfig(key)

	if baseURL != "" {
		logger.Infof(ctx, "NewAdapter DeepSeek model: %s, baseURL: %s", model, baseURL)
		config.BaseURL = baseURL
	} else {
		config.BaseURL = "https://api.deepseek.com/v1"
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAdapter DeepSeek model: %s, proxyURL: %s", model, proxyURL[0])

		proxyUrl, err := url.Parse(proxyURL[0])
		if err != nil {
			panic(err)
		}

		config.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	}

	return &DeepSeek{
		client:              openai.NewClientWithConfig(config),
		isSupportSystemRole: isSupportSystemRole,
	}
}

func NewAdapterBaidu(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *DeepSeek {

	logger.Infof(ctx, "NewAdapter DeepSeek model: %s, key: %s", model, key)

	split := gstr.Split(key, "|")

	config := openai.DefaultConfig(split[1])

	if baseURL != "" {
		logger.Infof(ctx, "NewAdapter DeepSeek model: %s, baseURL: %s", model, baseURL)
		config.BaseURL = baseURL
	} else {
		config.BaseURL = "https://qianfan.baidubce.com/v2"
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAdapter DeepSeek model: %s, proxyURL: %s", model, proxyURL[0])

		proxyUrl, err := url.Parse(proxyURL[0])
		if err != nil {
			panic(err)
		}

		config.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	}

	client := openai.NewClientWithConfig(config)
	client.Header = http.Header{
		"appid": []string{split[0]},
	}

	return &DeepSeek{
		client:              client,
		isSupportSystemRole: isSupportSystemRole,
	}
}

func (d *DeepSeek) apiErrorHandler(err error) error {

	apiError := &openai.APIError{}
	if errors.As(err, &apiError) {

		switch apiError.HTTPStatusCode {
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

	reqError := &openai.RequestError{}
	if errors.As(err, &reqError) {
		return sdkerr.NewRequestError(apiError.HTTPStatusCode, reqError.Err)
	}

	return err
}
