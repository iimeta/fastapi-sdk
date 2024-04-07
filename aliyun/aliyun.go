package aliyun

import (
	"context"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
}

func NewClient(ctx context.Context, model, key, baseURL, path string, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewClient Aliyun model: %s, key: %s", model, key)

	config := openai.DefaultConfig(key)

	if baseURL != "" {
		logger.Infof(ctx, "NewClient Aliyun model: %s, baseURL: %s", model, baseURL)
		config.BaseURL = baseURL
	}

	return &Client{
		client: openai.NewClientWithConfig(config),
	}
}

func (s *Client) ChatCompletion(ctx context.Context, request model.ChatCompletionRequest) (res model.ChatCompletionResponse, err error) {

	return
}

func (s *Client) ChatCompletionStream(ctx context.Context, request model.ChatCompletionRequest) (responseChan chan *model.ChatCompletionResponse, err error) {

	return
}

func (s *Client) Image(ctx context.Context, request model.ImageRequest) (res model.ImageResponse, err error) {

	return
}
