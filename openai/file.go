package openai

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (o *OpenAI) FileUpload(ctx context.Context, request model.FileUploadRequest) (response model.FileResponse, err error) {

	logger.Infof(ctx, "FileUpload OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileUpload OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	data, err := o.ConvFileUploadRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "FileUpload OpenAI ConvFileUploadRequest error: %v", err)
		return response, err
	}

	if o.Path == "" {
		o.Path = "/files"
	}

	bytes, err := util.HttpPost(ctx, o.BaseUrl+o.Path, o.header, data, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileUpload OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvFileResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileUpload OpenAI ConvFileResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "FileUpload OpenAI model: %s finished", o.Model)

	return response, nil
}

func (o *OpenAI) FileList(ctx context.Context, request model.FileListRequest) (response model.FileListResponse, err error) {

	logger.Infof(ctx, "FileList OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileList OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	if o.Path == "" {
		o.Path = "/files"
	}

	bytes, err := util.HttpGet(ctx, o.BaseUrl+o.Path, o.header, request, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileList OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvFileListResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileList OpenAI ConvFileListResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "FileList OpenAI model: %s finished", o.Model)

	return response, nil
}

func (o *OpenAI) FileRetrieve(ctx context.Context, request model.FileRetrieveRequest) (response model.FileResponse, err error) {

	logger.Infof(ctx, "FileRetrieve OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileRetrieve OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	if o.Path == "" {
		o.Path = fmt.Sprintf("/files/%s", request.FileId)
	}

	bytes, err := util.HttpGet(ctx, o.BaseUrl+o.Path, o.header, nil, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileRetrieve OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvFileResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileRetrieve OpenAI ConvFileResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "FileRetrieve OpenAI model: %s finished", o.Model)

	return response, nil
}

func (o *OpenAI) FileDelete(ctx context.Context, request model.FileDeleteRequest) (response model.FileResponse, err error) {

	logger.Infof(ctx, "FileDelete OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileDelete OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	if o.Path == "" {
		o.Path = fmt.Sprintf("/files/%s", request.FileId)
	}

	bytes, err := util.HttpDelete(ctx, o.BaseUrl+o.Path, o.header, nil, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileDelete OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvFileResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileDelete OpenAI ConvFileResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "FileDelete OpenAI model: %s finished", o.Model)

	return response, nil
}

func (o *OpenAI) FileContent(ctx context.Context, request model.FileContentRequest) (response model.FileContentResponse, err error) {

	logger.Infof(ctx, "FileContent OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "FileContent OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	if o.Path == "" {
		o.Path = fmt.Sprintf("/files/%s/content", request.FileId)
	}

	bytes, err := util.HttpGet(ctx, o.BaseUrl+o.Path, o.header, nil, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "FileContent OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvFileContentResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "FileContent OpenAI ConvFileContentResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "FileContent OpenAI model: %s finished", o.Model)

	return response, nil
}
