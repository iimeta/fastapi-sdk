package openai

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/go-openai"
)

func (c *Client) Embeddings(ctx context.Context, request model.EmbeddingRequest) (res model.EmbeddingResponse, err error) {

	logger.Infof(ctx, "Embeddings OpenAI model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "Embeddings OpenAI model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	response, err := c.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input:          request.Input,
		Model:          request.Model,
		User:           request.User,
		EncodingFormat: request.EncodingFormat,
		Dimensions:     request.Dimensions,
	})
	if err != nil {
		logger.Errorf(ctx, "Embeddings OpenAI model: %s, error: %v", request.Model, err)
		return res, err
	}

	logger.Infof(ctx, "Embeddings OpenAI model: %s finished", request.Model)

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
