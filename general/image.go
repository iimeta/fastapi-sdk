package general

import (
	"context"
	"slices"
	"strings"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (g *General) ImageGenerations(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerations General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageGenerations General model: %s totalTime: %d ms", g.Model, gtime.TimestampMilli()-now)
	}()

	path := g.Path
	if g.Async {
		if strings.Contains(path, "?") {
			path += "&async=true"
		} else {
			path += "?async=true"
		}
	}

	bytes, responseHeader, err := util.HttpPost(ctx, g.BaseUrl+path, g.header, data, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageGenerations General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvImageGenerationsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ImageGenerations General ConvImageGenerationsResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = bytes
	response.ResponseHeaders = responseHeader

	return response, nil
}

func (g *General) ImageEdits(ctx context.Context, request model.ImageEditRequest) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageEdits General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageEdits General model: %s totalTime: %d ms", g.Model, gtime.TimestampMilli()-now)
	}()

	var data any = request
	if !slices.Contains(g.ReqPassthroughParams, "req_data") {
		data, err = g.ConvImageEditsRequest(ctx, request)
		if err != nil {
			logger.Errorf(ctx, "ImageEdits General ConvImageEditsRequest error: %v", err)
			return response, err
		}
	}

	path := g.Path
	if g.Async {
		if strings.Contains(path, "?") {
			path += "&async=true"
		} else {
			path += "?async=true"
		}
	}

	bytes, responseHeader, err := util.HttpPost(ctx, g.BaseUrl+path, g.header, data, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageEdits General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvImageEditsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ImageEdits General ConvImageEditsResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = bytes
	response.ResponseHeaders = responseHeader

	return response, nil
}

func (g *General) ImageGenerationsStream(ctx context.Context, data []byte) (responseChan chan *model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) ImageEditsStream(ctx context.Context, request model.ImageEditRequest) (responseChan chan *model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}
