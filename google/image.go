package google

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (g *Google) ImageGenerations(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerations Google model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageGenerations Google model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	if !g.IsOfficialFormatRequest {

		request, err := g.ConvImageGenerationsRequest(ctx, data)
		if err != nil {
			logger.Errorf(ctx, "ImageGenerations Google ConvImageGenerationsRequest error: %v", err)
			return response, err
		}

		if data, err = g.ConvImageGenerationsRequestOfficial(ctx, request); err != nil {
			logger.Errorf(ctx, "ImageGenerations Google ConvImageGenerationsRequestOfficial error: %v", err)
			return response, err
		}
	}

	if g.Path == "" {
		g.Path = "/models/" + g.Model
	}

	if g.Action == "" {
		g.Action = "generateContent"
	}

	var bytes []byte

	if g.isGcp {
		if bytes, err = util.HttpPost(ctx, fmt.Sprintf("%s%s:%s", g.BaseUrl, g.Path, g.Action), g.header, data, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler); err != nil {
			logger.Errorf(ctx, "ImageGenerations Google model: %s, error: %v", g.Model, err)
			return response, err
		}
	} else {
		if bytes, err = util.HttpPost(ctx, fmt.Sprintf("%s%s:%s?key=%s", g.BaseUrl, g.Path, g.Action, g.Key), g.header, data, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler); err != nil {
			logger.Errorf(ctx, "ImageGenerations Google model: %s, error: %v", g.Model, err)
			return response, err
		}
	}

	if response, err = g.ConvImageGenerationsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ImageGenerations Google ConvImageGenerationsResponse error: %v", err)
		return response, err
	}

	return response, nil
}

func (g *Google) ImageEdits(ctx context.Context, request model.ImageEditRequest) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageEdits Google model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageEdits Google model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	data, err := g.ConvImageEditsRequestOfficial(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "ImageEdits Google ConvImageEditsRequestOfficial error: %v", err)
		return response, err
	}

	if g.Path == "" {
		g.Path = "/models/" + g.Model
	}

	if g.Action == "" {
		g.Action = "generateContent"
	}

	var bytes []byte

	if g.isGcp {
		if bytes, err = util.HttpPost(ctx, fmt.Sprintf("%s%s:%s", g.BaseUrl, g.Path, g.Action), g.header, data, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler); err != nil {
			logger.Errorf(ctx, "ImageEdits Google model: %s, error: %v", g.Model, err)
			return response, err
		}
	} else {
		if bytes, err = util.HttpPost(ctx, fmt.Sprintf("%s%s:%s?key=%s", g.BaseUrl, g.Path, g.Action, g.Key), g.header, data, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler); err != nil {
			logger.Errorf(ctx, "ImageEdits Google model: %s, error: %v", g.Model, err)
			return response, err
		}
	}

	if response, err = g.ConvImageGenerationsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ImageEdits Google ConvImageGenerationsResponse error: %v", err)
		return response, err
	}

	return response, nil
}
