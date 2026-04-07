package volcengine

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (v *VolcEngine) VideoCreate(ctx context.Context, request model.VideoCreateRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoCreate VolcEngine model: %s start", v.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoCreate VolcEngine model: %s totalTime: %d ms", v.Model, response.TotalTime)
	}()

	data, err := v.ConvVideoCreateRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "VideoCreate VolcEngine ConvVideoCreateRequest error: %v", err)
		return response, err
	}

	if v.Path == "" {
		v.Path = "/contents/generations/tasks"
	}

	bytes, err := util.HttpPost(ctx, v.BaseUrl+v.Path, v.header, data, nil, v.Timeout, v.ProxyUrl, v.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoCreate VolcEngine model: %s, error: %v", v.Model, err)
		return response, err
	}

	if response, err = v.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoCreate VolcEngine ConvVideoJobResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoCreate VolcEngine model: %s finished, id: %s", v.Model, response.Id)

	return response, nil
}

func (v *VolcEngine) VideoRemix(ctx context.Context, request model.VideoRemixRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoRemix VolcEngine model: %s start", v.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoRemix VolcEngine model: %s totalTime: %d ms", v.Model, response.TotalTime)
	}()

	createReq := model.VideoCreateRequest{
		Model:  v.Model,
		Prompt: request.Prompt,
	}

	response, err = v.VideoCreate(ctx, createReq)
	if err != nil {
		return response, err
	}

	response.RemixedFromVideoId = &request.VideoId

	return response, nil
}

func (v *VolcEngine) VideoList(ctx context.Context, request model.VideoListRequest) (response model.VideoListResponse, err error) {

	logger.Infof(ctx, "VideoList VolcEngine model: %s start", v.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoList VolcEngine model: %s totalTime: %d ms", v.Model, response.TotalTime)
	}()

	if v.Path == "" {
		v.Path = "/contents/generations/tasks"
	}

	bytes, err := util.HttpGet(ctx, v.BaseUrl+v.Path, v.header, request, nil, v.Timeout, v.ProxyUrl, v.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoList VolcEngine model: %s, error: %v", v.Model, err)
		return response, err
	}

	if response, err = v.ConvVideoListResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoList VolcEngine ConvVideoListResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoList VolcEngine model: %s finished", v.Model)

	return response, nil
}

func (v *VolcEngine) VideoRetrieve(ctx context.Context, request model.VideoRetrieveRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoRetrieve VolcEngine model: %s, videoId: %s start", v.Model, request.VideoId)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoRetrieve VolcEngine model: %s totalTime: %d ms", v.Model, response.TotalTime)
	}()

	if v.Path == "" {
		v.Path = fmt.Sprintf("/contents/generations/tasks/%s", request.VideoId)
	}

	bytes, err := util.HttpGet(ctx, v.BaseUrl+v.Path, v.header, nil, nil, v.Timeout, v.ProxyUrl, v.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoRetrieve VolcEngine model: %s, error: %v", v.Model, err)
		return response, err
	}

	if response, err = v.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoRetrieve VolcEngine ConvVideoJobResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoRetrieve VolcEngine model: %s, videoId: %s, status: %s finished", v.Model, request.VideoId, response.Status)

	return response, nil
}

func (v *VolcEngine) VideoDelete(ctx context.Context, request model.VideoDeleteRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoDelete VolcEngine model: %s, videoId: %s start", v.Model, request.VideoId)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoDelete VolcEngine model: %s totalTime: %d ms", v.Model, response.TotalTime)
	}()

	if v.Path == "" {
		v.Path = fmt.Sprintf("/contents/generations/tasks/%s", request.VideoId)
	}

	bytes, err := util.HttpDelete(ctx, v.BaseUrl+v.Path, v.header, nil, nil, v.Timeout, v.ProxyUrl, v.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoDelete VolcEngine model: %s, error: %v", v.Model, err)
		return response, err
	}

	if len(bytes) > 0 {
		if response, err = v.ConvVideoJobResponse(ctx, bytes); err != nil {
			logger.Errorf(ctx, "VideoDelete VolcEngine ConvVideoJobResponse error: %v", err)
			return response, err
		}
	} else {
		response = model.VideoJobResponse{
			Id:      request.VideoId,
			Object:  "video.deleted",
			Status:  "deleted",
			Deleted: true,
		}
	}

	logger.Infof(ctx, "VideoDelete VolcEngine model: %s, videoId: %s finished", v.Model, request.VideoId)

	return response, nil
}

func (v *VolcEngine) VideoContent(ctx context.Context, request model.VideoContentRequest) (response model.VideoContentResponse, err error) {

	logger.Infof(ctx, "VideoContent VolcEngine model: %s, videoId: %s start", v.Model, request.VideoId)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoContent VolcEngine model: %s totalTime: %d ms", v.Model, response.TotalTime)
	}()

	retrieve, err := v.VideoRetrieve(ctx, model.VideoRetrieveRequest{VideoId: request.VideoId})
	if err != nil {
		logger.Errorf(ctx, "VideoContent VolcEngine VideoRetrieve error: %v", err)
		return response, err
	}

	if retrieve.VideoUrl == "" {
		return response, fmt.Errorf("VideoContent VolcEngine: video_url is empty for videoId %s", request.VideoId)
	}

	data, err := util.HttpGet(ctx, retrieve.VideoUrl, nil, nil, nil, v.Timeout, v.ProxyUrl, nil)
	if err != nil {
		logger.Errorf(ctx, "VideoContent VolcEngine download error: %v", err)
		return response, err
	}

	response = model.VideoContentResponse{Data: data}

	logger.Infof(ctx, "VideoContent VolcEngine model: %s, videoId: %s finished, size: %d bytes", v.Model, request.VideoId, len(data))

	return response, nil
}
