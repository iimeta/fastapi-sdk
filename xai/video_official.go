package xai

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (x *XAI) VideoCreateOfficial(ctx context.Context, data []byte) (responseBytes []byte, responseHeader http.Header, err error) {

	logger.Infof(ctx, "VideoCreateOfficial xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "VideoCreateOfficial xAI model: %s totalTime: %d ms", x.Model, gtime.TimestampMilli()-now)
	}()

	if x.Path == "" {
		x.Path = "/videos/generations"
	}

	if responseBytes, responseHeader, err = util.HttpPost(ctx, x.BaseUrl+x.Path, x.header, data, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "VideoCreateOfficial xAI model: %s, error: %v", x.Model, err)
		return nil, nil, err
	}

	logger.Infof(ctx, "VideoCreateOfficial xAI model: %s finished", x.Model)

	return responseBytes, responseHeader, nil
}

func (x *XAI) VideoListOfficial(ctx context.Context, params model.VolcVideoListReq) (responseBytes []byte, responseHeader http.Header, err error) {

	logger.Infof(ctx, "VideoListOfficial xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "VideoListOfficial xAI model: %s totalTime: %d ms", x.Model, gtime.TimestampMilli()-now)
	}()

	if x.Path == "" {
		x.Path = "/videos/generations"
	}

	if responseBytes, responseHeader, err = util.HttpGet(ctx, x.BaseUrl+x.Path, x.header, params, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "VideoListOfficial xAI model: %s, error: %v", x.Model, err)
		return nil, nil, err
	}

	logger.Infof(ctx, "VideoListOfficial xAI model: %s finished", x.Model)

	return responseBytes, responseHeader, nil
}

func (x *XAI) VideoRetrieveOfficial(ctx context.Context, taskId string) (responseBytes []byte, responseHeader http.Header, err error) {

	logger.Infof(ctx, "VideoRetrieveOfficial xAI model: %s, taskId: %s start", x.Model, taskId)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "VideoRetrieveOfficial xAI model: %s totalTime: %d ms", x.Model, gtime.TimestampMilli()-now)
	}()

	if responseBytes, responseHeader, err = util.HttpGet(ctx, x.BaseUrl+fmt.Sprintf("/videos/generations/%s", taskId), x.header, nil, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "VideoRetrieveOfficial xAI model: %s, error: %v", x.Model, err)
		return nil, nil, err
	}

	logger.Infof(ctx, "VideoRetrieveOfficial xAI model: %s, taskId: %s finished", x.Model, taskId)

	return responseBytes, responseHeader, nil
}

func (x *XAI) VideoDeleteOfficial(ctx context.Context, taskId string) (err error) {

	logger.Infof(ctx, "VideoDeleteOfficial xAI model: %s, taskId: %s start", x.Model, taskId)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "VideoDeleteOfficial xAI model: %s totalTime: %d ms", x.Model, gtime.TimestampMilli()-now)
	}()

	if x.Path == "" {
		x.Path = fmt.Sprintf("/videos/generations/%s", taskId)
	}

	if _, _, err = util.HttpDelete(ctx, x.BaseUrl+x.Path, x.header, nil, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "VideoDeleteOfficial xAI model: %s, error: %v", x.Model, err)
		return err
	}

	logger.Infof(ctx, "VideoDeleteOfficial xAI model: %s, taskId: %s finished", x.Model, taskId)

	return nil
}
