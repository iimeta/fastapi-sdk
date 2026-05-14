package openai

import (
	"context"
	"slices"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (o *OpenAI) TextEmbeddings(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error) {

	logger.Infof(ctx, "TextEmbeddings OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "TextEmbeddings OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	var request any = data
	if !slices.Contains(o.ReqPassthroughParams, "req_data") {
		if request, err = o.ConvTextEmbeddingsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "TextEmbeddings OpenAI ConvTextEmbeddingsRequest error: %v", err)
			return response, err
		}
	}

	if o.Path == "" {
		if o.isAzure {
			o.Path = "/embeddings?api-version=" + o.apiVersion
		} else {
			o.Path = "/embeddings"
		}
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, o.BaseUrl+o.Path, o.header, request, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "TextEmbeddings OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvTextEmbeddingsResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "TextEmbeddings OpenAI ConvTextEmbeddingsResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "TextEmbeddings OpenAI model: %s finished", o.Model)

	return response, nil
}
