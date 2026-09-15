package minimax

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

func (m *MiniMax) VideoCreate(ctx context.Context, request model.VideoCreateRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoCreate MiniMax model: %s start", m.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoCreate MiniMax model: %s totalTime: %d ms", m.Model, response.TotalTime)
	}()

	var data any = request
	if !slices.Contains(m.ReqPassthroughParams, "req_data") {
		if data, err = m.ConvVideoCreateRequest(ctx, request); err != nil {
			logger.Errorf(ctx, "VideoCreate MiniMax ConvVideoCreateRequest error: %v", err)
			return response, err
		}
	}

	if m.Path == "" {
		m.Path = "/v2/video_generation"
	}

	bytes, responseHeader, err := util.HttpPost(ctx, m.BaseUrl+m.Path, m.header, data, nil, m.Timeout, m.ProxyUrl, m.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoCreate MiniMax model: %s, error: %v", m.Model, err)
		return response, err
	}

	if response, err = m.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoCreate MiniMax ConvVideoJobResponse error: %v", err)
		return response, err
	}

	response.Model = m.Model
	response.Prompt = request.Prompt
	response.ResponseBytes = bytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "VideoCreate MiniMax model: %s finished, id: %s", m.Model, response.Id)

	return response, nil
}

func (m *MiniMax) VideoRemix(ctx context.Context, request model.VideoRemixRequest) (response model.VideoJobResponse, err error) {
	return response, errors.New("MiniMax video does not support remix")
}

func (m *MiniMax) VideoList(ctx context.Context, request model.VideoListRequest) (response model.VideoListResponse, err error) {
	return response, errors.New("MiniMax video does not support list")
}

func (m *MiniMax) VideoRetrieve(ctx context.Context, request model.VideoRetrieveRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoRetrieve MiniMax model: %s, videoId: %s start", m.Model, request.VideoId)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoRetrieve MiniMax model: %s totalTime: %d ms", m.Model, response.TotalTime)
	}()

	bytes, _, err := util.HttpGet(ctx, m.BaseUrl+fmt.Sprintf("/v2/query/video_generation/%s", request.VideoId), m.header, nil, nil, m.Timeout, m.ProxyUrl, m.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoRetrieve MiniMax model: %s, error: %v", m.Model, err)
		return response, err
	}

	if response, err = m.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoRetrieve MiniMax ConvVideoJobResponse error: %v", err)
		return response, err
	}

	if response.Model == "" {
		response.Model = m.Model
	}

	logger.Infof(ctx, "VideoRetrieve MiniMax model: %s, videoId: %s, status: %s finished", m.Model, request.VideoId, response.Status)

	return response, nil
}

func (m *MiniMax) VideoDelete(ctx context.Context, request model.VideoDeleteRequest) (response model.VideoJobResponse, err error) {
	return response, errors.New("MiniMax video does not support delete")
}

// 先查询任务拿到 video url, 再下载视频字节
func (m *MiniMax) VideoContent(ctx context.Context, request model.VideoContentRequest) (response model.VideoContentResponse, err error) {

	logger.Infof(ctx, "VideoContent MiniMax model: %s, videoId: %s start", m.Model, request.VideoId)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoContent MiniMax model: %s totalTime: %d ms", m.Model, response.TotalTime)
	}()

	retrieve, err := m.VideoRetrieve(ctx, model.VideoRetrieveRequest{VideoId: request.VideoId})
	if err != nil {
		logger.Errorf(ctx, "VideoContent MiniMax VideoRetrieve error: %v", err)
		return response, err
	}

	if retrieve.VideoUrl == "" {
		return response, fmt.Errorf("VideoContent MiniMax: video_url is empty for videoId %s", request.VideoId)
	}

	data, _, err := util.HttpGet(ctx, retrieve.VideoUrl, nil, nil, nil, m.Timeout, m.ProxyUrl, nil)
	if err != nil {
		logger.Errorf(ctx, "VideoContent MiniMax download error: %v", err)
		return response, err
	}

	response = model.VideoContentResponse{Data: data}

	logger.Infof(ctx, "VideoContent MiniMax model: %s, videoId: %s finished, size: %d bytes", m.Model, request.VideoId, len(data))

	return response, nil
}
