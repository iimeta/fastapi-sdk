package general

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (g *General) TextEmbeddings(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error) {

	logger.Infof(ctx, "TextEmbeddings General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "TextEmbeddings General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	request, err := g.ConvTextEmbeddingsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "TextEmbeddings General ConvTextEmbeddingsRequest error: %v", err)
		return response, err
	}

	bytes, err := util.HttpPost(ctx, g.BaseUrl+g.Path, g.header, request, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "TextEmbeddings General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvTextEmbeddingsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "TextEmbeddings General ConvTextEmbeddingsResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "TextEmbeddings General model: %s finished", g.Model)

	return response, nil
}
