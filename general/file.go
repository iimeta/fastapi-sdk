package general

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (g *General) FileUpload(ctx context.Context, request model.FileUploadRequest) (response model.FileResponse, err error) {

	logger.Infof(ctx, "FileUpload General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileUpload General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	data, err := g.ConvFileUploadRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "FileUpload General ConvFileUploadRequest error: %v", err)
		return response, err
	}

	bytes, err := util.HttpPost(ctx, g.BaseUrl+g.Path, g.header, data, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileUpload General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvFileResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileUpload General ConvFileResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "FileUpload General model: %s finished", g.Model)

	return response, nil
}

func (g *General) FileList(ctx context.Context, request model.FileListRequest) (response model.FileListResponse, err error) {

	logger.Infof(ctx, "FileList General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileList General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	bytes, err := util.HttpGet(ctx, g.BaseUrl+g.Path, g.header, request, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileList General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvFileListResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileList General ConvFileListResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "FileList General model: %s finished", g.Model)

	return response, nil
}

func (g *General) FileRetrieve(ctx context.Context, request model.FileRetrieveRequest) (response model.FileResponse, err error) {

	logger.Infof(ctx, "FileRetrieve General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileRetrieve General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	bytes, err := util.HttpGet(ctx, g.BaseUrl+g.Path, g.header, nil, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileRetrieve General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvFileResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileRetrieve General ConvFileResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "FileRetrieve General model: %s finished", g.Model)

	return response, nil
}

func (g *General) FileDelete(ctx context.Context, request model.FileDeleteRequest) (response model.FileResponse, err error) {

	logger.Infof(ctx, "FileDelete General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileDelete General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	bytes, err := util.HttpDelete(ctx, g.BaseUrl+g.Path, g.header, nil, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileDelete General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvFileResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileDelete General ConvFileResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "FileDelete General model: %s finished", g.Model)

	return response, nil
}

func (g *General) FileContent(ctx context.Context, request model.FileContentRequest) (response model.FileContentResponse, err error) {

	logger.Infof(ctx, "FileContent General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileContent General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	bytes, err := util.HttpGet(ctx, g.BaseUrl+g.Path, g.header, nil, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileContent General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvFileContentResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileContent General ConvFileContentResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "FileContent General model: %s finished", g.Model)

	return response, nil
}
