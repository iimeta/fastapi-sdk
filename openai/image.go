package openai

import (
	"context"
	"slices"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (o *OpenAI) ImageGenerations(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerations OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageGenerations OpenAI model: %s totalTime: %d ms", o.Model, gtime.TimestampMilli()-now)
	}()

	var request any = data
	if !slices.Contains(o.ReqPassthroughParams, "req_data") {
		if request, err = o.ConvImageGenerationsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "ImageGenerations OpenAI ConvImageGenerationsRequest error: %v", err)
			return response, err
		}
	}

	if o.Path == "" {
		if o.isAzure {
			o.Path = "/images/generations?api-version=" + o.apiVersion
		} else {
			o.Path = "/images/generations"
		}
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, o.BaseUrl+o.Path, o.header, request, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageGenerations OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvImageGenerationsResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "ImageGenerations OpenAI ConvImageGenerationsResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

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

	if o.Path == "" {
		if o.isAzure {
			o.Path = "/images/edits?api-version=" + o.apiVersion
		} else {
			o.Path = "/images/edits"
		}
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, o.BaseUrl+o.Path, o.header, data, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageEdits OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvImageEditsResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "ImageEdits OpenAI ConvImageEditsResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	return response, nil
}
