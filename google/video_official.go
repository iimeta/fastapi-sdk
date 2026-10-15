package google

import (
	"context"
	"net/http"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (g *Google) VideoCreateOfficial(ctx context.Context, data []byte) (responseBytes []byte, responseHeader http.Header, err error) {

	logger.Infof(ctx, "VideoCreateOfficial Google model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "VideoCreateOfficial Google model: %s totalTime: %d ms", g.Model, gtime.TimestampMilli()-now)
	}()

	data = forceInteractionBackground(data)

	responseBytes, responseHeader, err = util.HttpPost(ctx, g.interactionsURL("", false), g.header, data, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoCreateOfficial Google model: %s, error: %v", g.Model, err)
		return responseBytes, responseHeader, err
	}

	return responseBytes, responseHeader, nil
}

func (g *Google) VideoListOfficial(ctx context.Context, params model.VolcVideoListReq) (responseBytes []byte, responseHeader http.Header, err error) {
	return nil, nil, errors.New("Google video does not support list")
}

func (g *Google) VideoRetrieveOfficial(ctx context.Context, taskId string) (responseBytes []byte, responseHeader http.Header, err error) {

	logger.Infof(ctx, "VideoRetrieveOfficial Google model: %s, taskId: %s start", g.Model, taskId)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "VideoRetrieveOfficial Google model: %s totalTime: %d ms", g.Model, gtime.TimestampMilli()-now)
	}()

	responseBytes, responseHeader, err = util.HttpGet(ctx, g.interactionsURL(taskId, false), g.header, nil, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoRetrieveOfficial Google model: %s, error: %v", g.Model, err)
		return responseBytes, responseHeader, err
	}

	return responseBytes, responseHeader, nil
}

func (g *Google) VideoDeleteOfficial(ctx context.Context, taskId string) (err error) {
	return errors.New("Google video does not support delete")
}
