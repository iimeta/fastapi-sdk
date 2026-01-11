package google

import (
	"context"
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
	//TODO implement me
	panic("implement me")
}

func (g *Google) FileRetrieve(ctx context.Context, request model.FileRetrieveRequest) (response model.FileResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) FileDelete(ctx context.Context, request model.FileDeleteRequest) (response model.FileResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) FileContent(ctx context.Context, request model.FileContentRequest) (response model.FileContentResponse, err error) {
	//TODO implement me
	panic("implement me")
}
