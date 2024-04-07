package baidu

import (
	"context"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"net/url"
)

type Client struct {
	client *openai.Client
}

func NewClient(ctx context.Context, model, key string, baseURL ...string) *Client {

	logger.Infof(ctx, "NewClient Baidu model: %s, key: %s", model, key)

	config := openai.DefaultConfig(key)

	if len(baseURL) > 0 && baseURL[0] != "" {
		logger.Infof(ctx, "NewClient Baidu model: %s, baseURL: %s", model, baseURL[0])
		config.BaseURL = baseURL[0]
	}

	return &Client{
		client: openai.NewClientWithConfig(config),
	}
}

func NewProxyClient(ctx context.Context, model, key string, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewProxyClient Baidu model: %s, key: %s", model, key)

	config := openai.DefaultConfig(key)

	transport := &http.Transport{}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewProxyClient Baidu model: %s, proxyURL: %s", model, proxyURL[0])
		proxyUrl, err := url.Parse(proxyURL[0])
		if err != nil {
			panic(err)
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
	}

	config.HTTPClient = &http.Client{
		Transport: transport,
	}

	return &Client{
		client: openai.NewClientWithConfig(config),
	}
}

func (c *Client) ChatCompletion(ctx context.Context, request model.ChatCompletionRequest) (res model.ChatCompletionResponse, err error) {

	return
}

func (c *Client) ChatCompletionStream(ctx context.Context, request model.ChatCompletionRequest) (responseChan chan *model.ChatCompletionResponse, err error) {

	return
}

func (c *Client) Image(ctx context.Context, request model.ImageRequest) (res model.ImageResponse, err error) {

	return
}
