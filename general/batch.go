package general

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (g *General) BatchCreate(ctx context.Context, request model.BatchCreateRequest) (response model.BatchResponse, err error) {

	logger.Infof(ctx, "BatchCreate General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "BatchCreate General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	bytes, err := util.HttpPost(ctx, g.BaseUrl+g.Path, g.header, request, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "BatchCreate General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvBatchResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "BatchCreate General ConvBatchResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "BatchCreate General model: %s finished", g.Model)

	return response, nil
}

func (g *General) BatchList(ctx context.Context, request model.BatchListRequest) (response model.BatchListResponse, err error) {

	logger.Infof(ctx, "BatchList General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "BatchList General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	bytes, err := util.HttpGet(ctx, g.BaseUrl+g.Path, g.header, request, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "BatchList General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvBatchListResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "BatchList General ConvBatchListResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "BatchList General model: %s finished", g.Model)

	return response, nil
}

func (g *General) BatchRetrieve(ctx context.Context, request model.BatchRetrieveRequest) (response model.BatchResponse, err error) {

	logger.Infof(ctx, "BatchRetrieve General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "BatchRetrieve General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	bytes, err := util.HttpGet(ctx, g.BaseUrl+g.Path, g.header, nil, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "BatchRetrieve General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvBatchResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "BatchRetrieve General ConvBatchResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "BatchRetrieve General model: %s finished", g.Model)

	return response, nil
}

func (g *General) BatchCancel(ctx context.Context, request model.BatchCancelRequest) (response model.BatchResponse, err error) {

	logger.Infof(ctx, "BatchCancel General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "BatchCancel General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	bytes, err := util.HttpPost(ctx, g.BaseUrl+g.Path, g.header, nil, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "BatchCancel General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvBatchResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "BatchCancel General ConvBatchResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "BatchCancel General model: %s finished", g.Model)

	return response, nil
}
