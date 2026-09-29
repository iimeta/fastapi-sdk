package xai

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (x *XAI) FileUpload(ctx context.Context, request model.FileUploadRequest) (response model.FileResponse, err error) {

	logger.Infof(ctx, "FileUpload xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileUpload xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	data, err := x.ConvFileUploadRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "FileUpload xAI ConvFileUploadRequest error: %v", err)
		return response, err
	}

	if x.Path == "" {
		x.Path = "/files"
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, x.BaseUrl+x.Path, x.header, data, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileUpload xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvFileResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "FileUpload xAI ConvFileResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "FileUpload xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) FileList(ctx context.Context, request model.FileListRequest) (response model.FileListResponse, err error) {

	logger.Infof(ctx, "FileList xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileList xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	if x.Path == "" {
		x.Path = "/files"
	}

	bytes, _, err := util.HttpGet(ctx, x.BaseUrl+x.Path, x.header, request, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileList xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvFileListResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileList xAI ConvFileListResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "FileList xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) FileRetrieve(ctx context.Context, request model.FileRetrieveRequest) (response model.FileResponse, err error) {

	logger.Infof(ctx, "FileRetrieve xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileRetrieve xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	if x.Path == "" {
		x.Path = fmt.Sprintf("/files/%s", request.FileId)
	}

	bytes, _, err := util.HttpGet(ctx, x.BaseUrl+x.Path, x.header, nil, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileRetrieve xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvFileResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileRetrieve xAI ConvFileResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "FileRetrieve xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) FileDelete(ctx context.Context, request model.FileDeleteRequest) (response model.FileResponse, err error) {

	logger.Infof(ctx, "FileDelete xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileDelete xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	if x.Path == "" {
		x.Path = fmt.Sprintf("/files/%s", request.FileId)
	}

	bytes, _, err := util.HttpDelete(ctx, x.BaseUrl+x.Path, x.header, nil, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileDelete xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvFileResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileDelete xAI ConvFileResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "FileDelete xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) FileContent(ctx context.Context, request model.FileContentRequest) (response model.FileContentResponse, err error) {

	logger.Infof(ctx, "FileContent xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileContent xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	if x.Path == "" {
		x.Path = fmt.Sprintf("/files/%s/content", request.FileId)
	}

	bytes, _, err := util.HttpGet(ctx, x.BaseUrl+x.Path, x.header, nil, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileContent xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvFileContentResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileContent xAI ConvFileContentResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "FileContent xAI model: %s finished", x.Model)

	return response, nil
}
