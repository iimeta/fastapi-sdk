package openai

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/sdkerr"
	"github.com/iimeta/go-openai"
	"net/http"
	"net/url"
)

type Client struct {
	client              *openai.Client
	model               string
	key                 string
	baseURL             string
	path                string
	proxyURL            string
	header              map[string]string
	isSupportSystemRole *bool
	isAzure             bool
}

func NewClient(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole *bool, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewClient OpenAI model: %s, key: %s", model, key)

	client := &Client{
		model:               model,
		key:                 key,
		baseURL:             "https://api.openai.com/v1",
		path:                "/chat/completions",
		isSupportSystemRole: isSupportSystemRole,
	}

	config := openai.DefaultConfig(key)

	if baseURL != "" {
		logger.Infof(ctx, "NewClient OpenAI model: %s, baseURL: %s", model, baseURL)
		config.BaseURL = baseURL
		client.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewClient OpenAI model: %s, path: %s", model, path)
		client.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewClient OpenAI model: %s, proxyURL: %s", model, proxyURL[0])

		proxyUrl, err := url.Parse(proxyURL[0])
		if err != nil {
			panic(err)
		}

		config.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}

		client.proxyURL = proxyURL[0]
	}

	header := make(map[string]string)
	header["Authorization"] = fmt.Sprintf("Bearer %s", key)

	client.header = header
	client.client = openai.NewClientWithConfig(config)

	return client
}

func NewAzureClient(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole *bool, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewAzureClient OpenAI model: %s, baseURL: %s, key: %s", model, baseURL, key)

	config := openai.DefaultAzureConfig(key, baseURL)

	if path != "" {
		logger.Infof(ctx, "NewAzureClient OpenAI model: %s, path: %s", model, path)

		split := gstr.Split(path, "?api-version=")

		if len(split) > 1 && split[1] != "" {
			config.APIVersion = split[1]
		}
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAzureClient OpenAI model: %s, proxyURL: %s", model, proxyURL[0])

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
		isAzure:             true,
	}
}

func (c *Client) requestErrorHandler(ctx context.Context, response *gclient.Response) (err error) {
	return sdkerr.NewRequestError(500, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, response.ReadAllString())))
}

func (c *Client) apiErrorHandler(err error) error {

	//apiError := &openai.APIError{}
	//if errors.As(err, &apiError) {
	//
	//	switch apiError.HTTPStatusCode {
	//	case 400:
	//		if apiError.Code == "context_length_exceeded" {
	//			return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
	//		}
	//	case 401:
	//		if apiError.Code == "invalid_api_key" {
	//			return sdkerr.ERR_INVALID_API_KEY
	//		}
	//	case 404:
	//		return sdkerr.ERR_MODEL_NOT_FOUND
	//	case 429:
	//		if apiError.Code == "insufficient_quota" {
	//			return sdkerr.ERR_INSUFFICIENT_QUOTA
	//		}
	//	}
	//
	//	return err
	//}
	//
	//reqError := &openai.RequestError{}
	//if errors.As(err, &reqError) {
	//	return sdkerr.NewRequestError(apiError.HTTPStatusCode, reqError.Err)
	//}

	return err
}
