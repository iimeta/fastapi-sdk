package openai

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (o *OpenAI) TextEmbeddings(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error) {

	logger.Infof(ctx, "TextEmbeddings OpenAI model: %s start", o.model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "TextEmbeddings OpenAI model: %s totalTime: %d ms", o.model, response.TotalTime)
	}()

	request, err := o.ConvTextEmbeddingsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "TextEmbeddings OpenAI ConvTextEmbeddingsRequest error: %v", err)
		return response, err
	}

	bytes, err := util.HttpPost(ctx, o.baseURL+"/embeddings", o.header, request, nil, o.proxyURL, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "TextEmbeddings OpenAI model: %s, error: %v", o.model, err)
		return response, err
	}

	if response, err = o.ConvTextEmbeddingsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "TextEmbeddings OpenAI ConvTextEmbeddingsResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "TextEmbeddings OpenAI model: %s finished", o.model)

	return response, nil
}
