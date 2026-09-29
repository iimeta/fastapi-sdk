package xai

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (x *XAI) BatchCreate(ctx context.Context, request model.BatchCreateRequest) (response model.BatchResponse, err error) {

	logger.Infof(ctx, "BatchCreate xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "BatchCreate xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	if x.Path == "" {
		x.Path = "/batches"
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, x.BaseUrl+x.Path, x.header, request, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "BatchCreate xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvBatchResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "BatchCreate xAI ConvBatchResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "BatchCreate xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) BatchList(ctx context.Context, request model.BatchListRequest) (response model.BatchListResponse, err error) {

	logger.Infof(ctx, "BatchList xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "BatchList xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	if x.Path == "" {
		x.Path = "/batches"
	}

	bytes, _, err := util.HttpGet(ctx, x.BaseUrl+x.Path, x.header, request, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "BatchList xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvBatchListResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "BatchList xAI ConvBatchListResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "BatchList xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) BatchRetrieve(ctx context.Context, request model.BatchRetrieveRequest) (response model.BatchResponse, err error) {

	logger.Infof(ctx, "BatchRetrieve xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "BatchRetrieve xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	if x.Path == "" {
		x.Path = fmt.Sprintf("/batches/%s", request.BatchId)
	}

	bytes, _, err := util.HttpGet(ctx, x.BaseUrl+x.Path, x.header, nil, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "BatchRetrieve xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvBatchResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "BatchRetrieve xAI ConvBatchResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "BatchRetrieve xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) BatchCancel(ctx context.Context, request model.BatchCancelRequest) (response model.BatchResponse, err error) {

	logger.Infof(ctx, "BatchCancel xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "BatchCancel xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	if x.Path == "" {
		x.Path = fmt.Sprintf("/batches/%s/cancel", request.BatchId)
	}

	bytes, _, err := util.HttpPost(ctx, x.BaseUrl+x.Path, x.header, nil, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "BatchCancel xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvBatchResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "BatchCancel xAI ConvBatchResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "BatchCancel xAI model: %s finished", x.Model)

	return response, nil
}
