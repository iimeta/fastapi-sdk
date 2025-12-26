package openai

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (o *OpenAI) VideoCreate(ctx context.Context, request model.VideoCreateRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoCreate OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoCreate OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	data, err := o.ConvVideoCreateRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "VideoCreate OpenAI ConvVideoCreateRequest error: %v", err)
		return response, err
	}

	if o.Path == "" {
		o.Path = "/videos"
	}

	bytes, err := util.HttpPost(ctx, o.BaseUrl+o.Path, o.header, data, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoCreate OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoCreate OpenAI ConvVideoJobResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoCreate OpenAI model: %s finished", o.Model)

	return response, nil
}

func (o *OpenAI) VideoRemix(ctx context.Context, request model.VideoRemixRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoRemix OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoRemix OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	if o.Path == "" {
		o.Path = fmt.Sprintf("/videos/%s/remix", request.VideoId)
	}

	bytes, err := util.HttpPost(ctx, o.BaseUrl+o.Path, o.header, request, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoRemix OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoRemix OpenAI ConvVideoJobResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoRemix OpenAI model: %s finished", o.Model)

	return response, nil
}

func (o *OpenAI) VideoList(ctx context.Context, request model.VideoListRequest) (response model.VideoListResponse, err error) {

	logger.Infof(ctx, "VideoList OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoList OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	if o.Path == "" {
		o.Path = "/videos"
	}

	bytes, err := util.HttpGet(ctx, o.BaseUrl+o.Path, o.header, request, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoList OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvVideoListResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoList OpenAI ConvVideoListResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoList OpenAI model: %s finished", o.Model)

	return response, nil
}

func (o *OpenAI) VideoRetrieve(ctx context.Context, request model.VideoRetrieveRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoRetrieve OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoRetrieve OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	if o.Path == "" {
		o.Path = fmt.Sprintf("/videos/%s", request.VideoId)
	}

	bytes, err := util.HttpGet(ctx, o.BaseUrl+o.Path, o.header, nil, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoRetrieve OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoRetrieve OpenAI ConvVideoJobResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoRetrieve OpenAI model: %s finished", o.Model)

	return response, nil
}

func (o *OpenAI) VideoDelete(ctx context.Context, request model.VideoDeleteRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoDelete OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoDelete OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	if o.Path == "" {
		o.Path = fmt.Sprintf("/videos/%s", request.VideoId)
	}

	bytes, err := util.HttpDelete(ctx, o.BaseUrl+o.Path, o.header, nil, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoDelete OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoDelete OpenAI ConvVideoJobResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoDelete OpenAI model: %s finished", o.Model)

	return response, nil
}

func (o *OpenAI) VideoContent(ctx context.Context, request model.VideoContentRequest) (response model.VideoContentResponse, err error) {

	logger.Infof(ctx, "VideoContent OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoContent OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	if o.Path == "" {
		o.Path = fmt.Sprintf("/videos/%s/content", request.VideoId)
	}

	bytes, err := util.HttpGet(ctx, o.BaseUrl+o.Path, o.header, nil, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoContent OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvVideoContentResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoContent OpenAI ConvVideoContentResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoContent OpenAI model: %s finished", o.Model)

	return response, nil
}
