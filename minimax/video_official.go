package minimax

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

// 官方格式请求体与响应体原样透传, 仅替换鉴权
func (m *MiniMax) VideoCreateOfficial(ctx context.Context, data []byte) (responseBytes []byte, responseHeader http.Header, err error) {

	logger.Infof(ctx, "VideoCreateOfficial MiniMax model: %s start", m.Model)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "VideoCreateOfficial MiniMax model: %s totalTime: %d ms", m.Model, gtime.TimestampMilli()-now)
	}()

	if m.Path == "" {
		m.Path = "/v2/video_generation"
	}

	if responseBytes, responseHeader, err = util.HttpPost(ctx, m.BaseUrl+m.Path, m.header, data, nil, m.Timeout, m.ProxyUrl, m.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "VideoCreateOfficial MiniMax model: %s, error: %v", m.Model, err)
		return nil, nil, err
	}

	logger.Infof(ctx, "VideoCreateOfficial MiniMax model: %s finished", m.Model)

	return responseBytes, responseHeader, nil
}

func (m *MiniMax) VideoListOfficial(ctx context.Context, params model.VolcVideoListReq) (responseBytes []byte, responseHeader http.Header, err error) {
	return nil, nil, errors.New("MiniMax video does not support list")
}

func (m *MiniMax) VideoRetrieveOfficial(ctx context.Context, taskId string) (responseBytes []byte, responseHeader http.Header, err error) {

	logger.Infof(ctx, "VideoRetrieveOfficial MiniMax model: %s, taskId: %s start", m.Model, taskId)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "VideoRetrieveOfficial MiniMax model: %s totalTime: %d ms", m.Model, gtime.TimestampMilli()-now)
	}()

	if responseBytes, responseHeader, err = util.HttpGet(ctx, m.BaseUrl+fmt.Sprintf("/v2/query/video_generation/%s", taskId), m.header, nil, nil, m.Timeout, m.ProxyUrl, m.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "VideoRetrieveOfficial MiniMax model: %s, error: %v", m.Model, err)
		return nil, nil, err
	}

	logger.Infof(ctx, "VideoRetrieveOfficial MiniMax model: %s, taskId: %s finished", m.Model, taskId)

	return responseBytes, responseHeader, nil
}

func (m *MiniMax) VideoDeleteOfficial(ctx context.Context, taskId string) (err error) {
	return errors.New("MiniMax video does not support delete")
}
