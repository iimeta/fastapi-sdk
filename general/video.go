package general

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (g *General) VideoCreate(ctx context.Context, request model.VideoCreateRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoCreate General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoCreate General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	data, err := g.ConvVideoCreateRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "VideoCreate General ConvVideoCreateRequest error: %v", err)
		return response, err
	}

	bytes, err := util.HttpPost(ctx, g.BaseUrl+g.Path, g.header, data, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoCreate General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoCreate General ConvVideoJobResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoCreate General model: %s finished", g.Model)

	return response, nil
}

func (g *General) VideoRemix(ctx context.Context, request model.VideoRemixRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoRemix General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoRemix General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	bytes, err := util.HttpPost(ctx, g.BaseUrl+g.Path, g.header, request, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoRemix General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoRemix General ConvVideoJobResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoRemix General model: %s finished", g.Model)

	return response, nil
}

func (g *General) VideoList(ctx context.Context, request model.VideoListRequest) (response model.VideoListResponse, err error) {

	logger.Infof(ctx, "VideoList General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoList General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	bytes, err := util.HttpGet(ctx, g.BaseUrl+g.Path, g.header, request, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoList General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvVideoListResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoList General ConvVideoListResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoList General model: %s finished", g.Model)

	return response, nil
}

func (g *General) VideoRetrieve(ctx context.Context, request model.VideoRetrieveRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoRetrieve General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoRetrieve General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	bytes, err := util.HttpGet(ctx, g.BaseUrl+g.Path, g.header, nil, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoRetrieve General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoRetrieve General ConvVideoJobResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoRetrieve General model: %s finished", g.Model)

	return response, nil
}

func (g *General) VideoDelete(ctx context.Context, request model.VideoDeleteRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoDelete General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoDelete General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	bytes, err := util.HttpDelete(ctx, g.BaseUrl+g.Path, g.header, nil, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoDelete General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoDelete General ConvVideoJobResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoDelete General model: %s finished", g.Model)

	return response, nil
}

func (g *General) VideoContent(ctx context.Context, request model.VideoContentRequest) (response model.VideoContentResponse, err error) {

	logger.Infof(ctx, "VideoContent General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoContent General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	bytes, err := util.HttpGet(ctx, g.BaseUrl+g.Path, g.header, nil, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoContent General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvVideoContentResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoContent General ConvVideoContentResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoContent General model: %s finished", g.Model)

	return response, nil
}
