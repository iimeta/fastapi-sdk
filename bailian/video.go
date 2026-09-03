package bailian

import (
	"context"
	"fmt"
	"slices"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

// 视频生成为异步接口: 创建任务返回 task_id, 再通过 VideoRetrieve 轮询结果
func (b *Bailian) VideoCreate(ctx context.Context, request model.VideoCreateRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoCreate Bailian model: %s start", b.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoCreate Bailian model: %s totalTime: %d ms", b.Model, response.TotalTime)
	}()

	var data any = request
	if !slices.Contains(b.ReqPassthroughParams, "req_data") {
		if data, err = b.ConvVideoCreateRequest(ctx, request); err != nil {
			logger.Errorf(ctx, "VideoCreate Bailian ConvVideoCreateRequest error: %v", err)
			return response, err
		}
	}

	if b.Path == "" {
		b.Path = "/services/aigc/video-generation/video-synthesis"
	}

	bytes, responseHeader, err := util.HttpPost(ctx, b.BaseUrl+b.Path, b.asyncHeader(), data, nil, b.Timeout, b.ProxyUrl, b.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoCreate Bailian model: %s, error: %v", b.Model, err)
		return response, err
	}

	if response, err = b.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoCreate Bailian ConvVideoJobResponse error: %v", err)
		return response, err
	}

	response.Model = b.Model
	response.Prompt = request.Prompt
	response.ResponseBytes = bytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "VideoCreate Bailian model: %s finished, id: %s", b.Model, response.Id)

	return response, nil
}

func (b *Bailian) VideoRemix(ctx context.Context, request model.VideoRemixRequest) (response model.VideoJobResponse, err error) {
	return response, errors.New("Bailian video does not support remix")
}

func (b *Bailian) VideoList(ctx context.Context, request model.VideoListRequest) (response model.VideoListResponse, err error) {
	return response, errors.New("Bailian video does not support list")
}

func (b *Bailian) VideoRetrieve(ctx context.Context, request model.VideoRetrieveRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoRetrieve Bailian model: %s, videoId: %s start", b.Model, request.VideoId)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoRetrieve Bailian model: %s totalTime: %d ms", b.Model, response.TotalTime)
	}()

	bytes, _, err := util.HttpGet(ctx, b.BaseUrl+fmt.Sprintf("/tasks/%s", request.VideoId), b.header, nil, nil, b.Timeout, b.ProxyUrl, b.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoRetrieve Bailian model: %s, error: %v", b.Model, err)
		return response, err
	}

	if response, err = b.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoRetrieve Bailian ConvVideoJobResponse error: %v", err)
		return response, err
	}

	response.Model = b.Model

	logger.Infof(ctx, "VideoRetrieve Bailian model: %s, videoId: %s, status: %s finished", b.Model, request.VideoId, response.Status)

	return response, nil
}

func (b *Bailian) VideoDelete(ctx context.Context, request model.VideoDeleteRequest) (response model.VideoJobResponse, err error) {
	return response, errors.New("Bailian video does not support delete")
}

// 先查询任务拿到 video_url, 再下载视频字节
func (b *Bailian) VideoContent(ctx context.Context, request model.VideoContentRequest) (response model.VideoContentResponse, err error) {

	logger.Infof(ctx, "VideoContent Bailian model: %s, videoId: %s start", b.Model, request.VideoId)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoContent Bailian model: %s totalTime: %d ms", b.Model, response.TotalTime)
	}()

	retrieve, err := b.VideoRetrieve(ctx, model.VideoRetrieveRequest{VideoId: request.VideoId})
	if err != nil {
		logger.Errorf(ctx, "VideoContent Bailian VideoRetrieve error: %v", err)
		return response, err
	}

	if retrieve.VideoUrl == "" {
		return response, fmt.Errorf("VideoContent Bailian: video_url is empty for videoId %s", request.VideoId)
	}

	data, _, err := util.HttpGet(ctx, retrieve.VideoUrl, nil, nil, nil, b.Timeout, b.ProxyUrl, nil)
	if err != nil {
		logger.Errorf(ctx, "VideoContent Bailian download error: %v", err)
		return response, err
	}

	response = model.VideoContentResponse{Data: data}

	logger.Infof(ctx, "VideoContent Bailian model: %s, videoId: %s finished, size: %d bytes", b.Model, request.VideoId, len(data))

	return response, nil
}
