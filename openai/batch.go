package openai

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (o *OpenAI) BatchCreate(ctx context.Context, request model.BatchCreateRequest) (response model.BatchResponse, err error) {

	logger.Infof(ctx, "BatchCreate OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "BatchCreate OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	if o.Path == "" {
		o.Path = "/batches"
	}

	bytes, err := util.HttpPost(ctx, o.BaseUrl+o.Path, o.header, request, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "BatchCreate OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvBatchResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "BatchCreate OpenAI ConvBatchResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "BatchCreate OpenAI model: %s finished", o.Model)

	return response, nil
}

func (o *OpenAI) BatchList(ctx context.Context, request model.BatchListRequest) (response model.BatchListResponse, err error) {

	logger.Infof(ctx, "BatchList OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "BatchList OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	if o.Path == "" {
		o.Path = "/batches"
	}

	bytes, err := util.HttpGet(ctx, o.BaseUrl+o.Path, o.header, request, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "BatchList OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvBatchListResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "BatchList OpenAI ConvBatchListResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "BatchList OpenAI model: %s finished", o.Model)

	return response, nil
}

func (o *OpenAI) BatchRetrieve(ctx context.Context, request model.BatchRetrieveRequest) (response model.BatchResponse, err error) {

	logger.Infof(ctx, "BatchRetrieve OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "BatchRetrieve OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	if o.Path == "" {
		o.Path = fmt.Sprintf("/batches/%s", request.BatchId)
	}

	bytes, err := util.HttpGet(ctx, o.BaseUrl+o.Path, o.header, nil, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "BatchRetrieve OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvBatchResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "BatchRetrieve OpenAI ConvBatchResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "BatchRetrieve OpenAI model: %s finished", o.Model)

	return response, nil
}

func (o *OpenAI) BatchCancel(ctx context.Context, request model.BatchCancelRequest) (response model.BatchResponse, err error) {

	logger.Infof(ctx, "BatchCancel OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "BatchCancel OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	if o.Path == "" {
		o.Path = fmt.Sprintf("/batches/%s/cancel", request.BatchId)
	}

	bytes, err := util.HttpPost(ctx, o.BaseUrl+o.Path, o.header, nil, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "BatchCancel OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvBatchResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "BatchCancel OpenAI ConvBatchResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "BatchCancel OpenAI model: %s finished", o.Model)

	return response, nil
}
