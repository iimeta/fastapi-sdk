package xai

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (x *XAI) VideoCreate(ctx context.Context, request model.VideoCreateRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoCreate xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoCreate xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	var data any = request
	if !slices.Contains(x.ReqPassthroughParams, "req_data") {
		if data, err = x.ConvVideoCreateRequest(ctx, request); err != nil {
			logger.Errorf(ctx, "VideoCreate xAI ConvVideoCreateRequest error: %v", err)
			return response, err
		}
	}

	if x.Path == "" {
		x.Path = "/videos/generations"
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, x.BaseUrl+x.Path, x.header, data, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoCreate xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvVideoJobResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "VideoCreate xAI ConvVideoJobResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "VideoCreate xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) VideoRemix(ctx context.Context, request model.VideoRemixRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoRemix xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoRemix xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	if x.Path == "" {
		x.Path = "/videos/generations"
	}

	responseBytes, responseHeader, err := func() ([]byte, map[string][]string, error) {
		video := map[string]any{}
		if strings.HasPrefix(request.VideoId, "http") || strings.HasPrefix(request.VideoId, "data:") {
			video["url"] = request.VideoId
		} else {
			video["file_id"] = request.VideoId
		}
		body := map[string]any{"model": x.Model, "prompt": request.Prompt, "video": video}
		x.header["Content-Type"] = "application/json"
		return util.HttpPost(ctx, x.BaseUrl+x.Path, x.header, body, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	}()
	if err != nil {
		logger.Errorf(ctx, "VideoRemix xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvVideoJobResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "VideoRemix xAI ConvVideoJobResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "VideoRemix xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) VideoList(ctx context.Context, request model.VideoListRequest) (response model.VideoListResponse, err error) {

	logger.Infof(ctx, "VideoList xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoList xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	if x.Path == "" {
		x.Path = "/videos/generations"
	}

	bytes, _, err := util.HttpGet(ctx, x.BaseUrl+x.Path, x.header, request, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoList xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvVideoListResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoList xAI ConvVideoListResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoList xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) VideoRetrieve(ctx context.Context, request model.VideoRetrieveRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoRetrieve xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoRetrieve xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	if x.Path == "" {
		x.Path = fmt.Sprintf("/videos/generations/%s", request.VideoId)
	}

	bytes, _, err := util.HttpGet(ctx, x.BaseUrl+x.Path, x.header, nil, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoRetrieve xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoRetrieve xAI ConvVideoJobResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoRetrieve xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) VideoDelete(ctx context.Context, request model.VideoDeleteRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoDelete xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoDelete xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	if x.Path == "" {
		x.Path = fmt.Sprintf("/videos/generations/%s", request.VideoId)
	}

	bytes, _, err := util.HttpDelete(ctx, x.BaseUrl+x.Path, x.header, nil, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoDelete xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoDelete xAI ConvVideoJobResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoDelete xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) VideoContent(ctx context.Context, request model.VideoContentRequest) (response model.VideoContentResponse, err error) {

	logger.Infof(ctx, "VideoContent xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoContent xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	job, err := x.VideoRetrieve(ctx, model.VideoRetrieveRequest{VideoId: request.VideoId})
	if err != nil {
		logger.Errorf(ctx, "VideoContent xAI model: %s, error: %v", x.Model, err)
		return response, err
	}
	if job.VideoUrl == "" {
		return response, errors.New("video url is empty")
	}

	bytes, _, err := util.HttpGet(ctx, job.VideoUrl, nil, nil, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoContent xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvVideoContentResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoContent xAI ConvVideoContentResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoContent xAI model: %s finished", x.Model)

	return response, nil
}
