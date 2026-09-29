package xai

import (
	"context"
	"slices"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (x *XAI) TextEmbeddings(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error) {

	logger.Infof(ctx, "TextEmbeddings xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "TextEmbeddings xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	var request any = data
	if !slices.Contains(x.ReqPassthroughParams, "req_data") {
		if request, err = x.ConvTextEmbeddingsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "TextEmbeddings xAI ConvTextEmbeddingsRequest error: %v", err)
			return response, err
		}
	}

	if x.Path == "" {
		x.Path = "/embeddings"
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, x.BaseUrl+x.Path, x.header, request, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "TextEmbeddings xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvTextEmbeddingsResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "TextEmbeddings xAI ConvTextEmbeddingsResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "TextEmbeddings xAI model: %s finished", x.Model)

	return response, nil
}
