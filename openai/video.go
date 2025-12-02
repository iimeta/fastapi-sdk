package openai

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (o *OpenAI) VideoCreate(ctx context.Context, request model.VideoRequest) (response model.VideoResponse, err error) {

	logger.Infof(ctx, "VideoCreate OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoCreate OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	data, err := o.ConvVideoRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "VideoCreate OpenAI ConvVideoRequest error: %v", err)
		return response, err
	}

	bytes, err := util.HttpPost(ctx, o.BaseUrl+o.Path, o.header, data, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoCreate OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvVideoResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoCreate OpenAI ConvVideoResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoCreate OpenAI model: %s finished", o.Model)

	return response, nil
}
