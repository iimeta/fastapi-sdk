package ai360

import (
	"context"
	"errors"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/sdkerr"
	"github.com/iimeta/go-openai"
	"net/http"
	"net/url"
)

type Client struct {
	client              *openai.Client
	isSupportSystemRole *bool
}

func NewClient(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewClient 360AI model: %s, key: %s", model, key)

	config := openai.DefaultConfig(key)

	if baseURL != "" {
		logger.Infof(ctx, "NewClient 360AI model: %s, baseURL: %s", model, baseURL)
		config.BaseURL = baseURL
	} else {
		config.BaseURL = "https://api.360.cn/v1"
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewClient 360AI model: %s, proxyURL: %s", model, proxyURL[0])

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

	return &Client{
		client:              openai.NewClientWithConfig(config),
		isSupportSystemRole: isSupportSystemRole,
	}
}

func (c *Client) apiErrorHandler(err error) error {

	apiError := &openai.APIError{}
	if errors.As(err, &apiError) {

		switch apiError.HTTPStatusCode {
		case 400:
			if apiError.Code == "1001" {
				return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
			}
		case 401:

			if apiError.Code == "1002" {
				return sdkerr.ERR_INVALID_API_KEY
			}

			if apiError.Code == "1004" || apiError.Code == "1006" {
				return sdkerr.ERR_INSUFFICIENT_QUOTA
			}

		case 404:
			return sdkerr.ERR_MODEL_NOT_FOUND
		case 429:
			if apiError.Code == "1005" {
				return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
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
