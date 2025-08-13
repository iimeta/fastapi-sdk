package sdk

import (
	"context"
	"net/http"
	"net/url"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/go-openai"
)

type EmbeddingClient struct {
	client *openai.Client
}

func NewEmbeddingClient(ctx context.Context, model, key, baseURL, path string, proxyURL ...string) *EmbeddingClient {

	logger.Infof(ctx, "NewAdapter OpenAI model: %s, key: %s", model, key)

	config := openai.DefaultConfig(key)

	if baseURL != "" {
		logger.Infof(ctx, "NewAdapter OpenAI model: %s, baseURL: %s", model, baseURL)
		config.BaseURL = baseURL
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAdapter OpenAI model: %s, proxyURL: %s", model, proxyURL[0])

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

	return &EmbeddingClient{
		client: openai.NewClientWithConfig(config),
	}
}

func (c *EmbeddingClient) Embeddings(ctx context.Context, request model.EmbeddingRequest) (res model.EmbeddingResponse, err error) {

	logger.Infof(ctx, "TextEmbeddings OpenAI model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "TextEmbeddings OpenAI model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	response, err := c.client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Input:          request.Input,
		Model:          request.Model,
		User:           request.User,
		EncodingFormat: request.EncodingFormat,
		Dimensions:     request.Dimensions,
	})
	if err != nil {
		logger.Errorf(ctx, "TextEmbeddings OpenAI model: %s, error: %v", request.Model, err)
		return res, err
	}

	logger.Infof(ctx, "TextEmbeddings OpenAI model: %s finished", request.Model)

	res = model.EmbeddingResponse{
		Object: response.Object,
		Data:   response.Data,
		Model:  response.Model,
		Usage: &model.Usage{
			PromptTokens:     response.Usage.PromptTokens,
			CompletionTokens: response.Usage.CompletionTokens,
			TotalTokens:      response.Usage.TotalTokens,
		},
	}

	return res, nil
}
