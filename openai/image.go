package openai

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (o *OpenAI) ImageGenerations(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerations OpenAI model: %s start", o.Model)

	request, err := o.ConvImageGenerationsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ImageGenerations OpenAI ConvImageGenerationsRequest error: %v", err)
		return response, err
	}

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageGenerations OpenAI model: %s totalTime: %d ms", o.Model, gtime.TimestampMilli()-now)
	}()

	bytes, err := util.HttpPost(ctx, o.BaseUrl+"/images/generations", o.header, request, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageGenerations OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvImageGenerationsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ImageGenerations OpenAI ConvImageGenerationsResponse error: %v", err)
		return response, err
	}

	return response, nil
}

func (o *OpenAI) ImageEdits(ctx context.Context, request model.ImageEditRequest) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageEdits OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageEdits OpenAI model: %s totalTime: %d ms", o.Model, gtime.TimestampMilli()-now)
	}()

	data, err := o.ConvImageEditsRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "ImageEdits OpenAI ConvImageEditsRequest error: %v", err)
		return response, err
	}

	bytes, err := util.HttpPost(ctx, o.BaseUrl+"/images/edits", o.header, data, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageEdits OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvImageEditsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ImageEdits OpenAI ConvImageEditsResponse error: %v", err)
		return response, err
	}

	return response, nil
}
