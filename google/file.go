package google

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (g *Google) FileUpload(ctx context.Context, request model.FileUploadRequest) (response model.FileResponse, err error) {

	logger.Infof(ctx, "FileUpload Google model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileUpload Google model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	data, err := g.ConvFileUploadRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "FileUpload Google ConvFileUploadRequest error: %v", err)
		return response, err
	}

	if g.Path == "" {
		g.Path = "/upload/v1beta/files"
	}

	if strings.HasSuffix(g.BaseUrl, "/v1beta") && strings.HasPrefix(g.Path, "/upload/v1beta") {
		g.BaseUrl = strings.TrimSuffix(g.BaseUrl, "/v1beta")
	}

	bytes, err := util.HttpPost(ctx, g.BaseUrl+g.Path, g.header, data, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileUpload Google model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvFileResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileUpload Google ConvFileResponse error: %v", err)
		return response, err
	}

	response.Filename = request.File.Filename

	if request.Purpose != "" {
		response.Purpose = request.Purpose
	}

	logger.Infof(ctx, "FileUpload Google model: %s finished", g.Model)

	return response, nil
}

func (g *Google) FileList(ctx context.Context, request model.FileListRequest) (response model.FileListResponse, err error) {

	logger.Infof(ctx, "FileList Google model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileList Google model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	if g.Path == "" {
		g.Path = "/files"
	}

	bytes, err := util.HttpGet(ctx, g.BaseUrl+g.Path, g.header, request, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileList Google model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvFileListResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileList Google ConvFileListResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "FileList Google model: %s finished", g.Model)

	return response, nil
}

func (g *Google) FileRetrieve(ctx context.Context, request model.FileRetrieveRequest) (response model.FileResponse, err error) {

	logger.Infof(ctx, "FileRetrieve Google model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileRetrieve Google model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	if g.Path == "" {
		g.Path = fmt.Sprintf("/files/%s", request.FileId)
	}

	bytes, err := util.HttpGet(ctx, g.BaseUrl+g.Path, g.header, nil, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileRetrieve Google model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvFileResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileRetrieve Google ConvFileResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "FileRetrieve Google model: %s finished", g.Model)

	return response, nil
}

func (g *Google) FileDelete(ctx context.Context, request model.FileDeleteRequest) (response model.FileResponse, err error) {

	logger.Infof(ctx, "FileDelete Google model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileDelete Google model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	if g.Path == "" {
		g.Path = fmt.Sprintf("/files/%s", request.FileId)
	}

	_, err = util.HttpDelete(ctx, g.BaseUrl+g.Path, g.header, nil, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileDelete Google model: %s, error: %v", g.Model, err)
		return response, err
	}

	response.Deleted = true

	logger.Infof(ctx, "FileDelete Google model: %s finished", g.Model)

	return response, nil
}

func (g *Google) FileContent(ctx context.Context, request model.FileContentRequest) (response model.FileContentResponse, err error) {
	//TODO implement me
	return response, err
}
