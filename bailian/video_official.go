package bailian

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

// 官方格式请求体与响应体原样透传, 仅替换鉴权并追加异步请求头
func (b *Bailian) VideoCreateOfficial(ctx context.Context, data []byte) (responseBytes []byte, responseHeader http.Header, err error) {

	logger.Infof(ctx, "VideoCreateOfficial Bailian model: %s start", b.Model)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "VideoCreateOfficial Bailian model: %s totalTime: %d ms", b.Model, gtime.TimestampMilli()-now)
	}()

	if b.Path == "" {
		b.Path = "/services/aigc/video-generation/video-synthesis"
	}

	if responseBytes, responseHeader, err = util.HttpPost(ctx, b.BaseUrl+b.Path, b.asyncHeader(), data, nil, b.Timeout, b.ProxyUrl, b.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "VideoCreateOfficial Bailian model: %s, error: %v", b.Model, err)
		return nil, nil, err
	}

	logger.Infof(ctx, "VideoCreateOfficial Bailian model: %s finished", b.Model)

	return responseBytes, responseHeader, nil
}

func (b *Bailian) VideoListOfficial(ctx context.Context, params model.VolcVideoListReq) (responseBytes []byte, responseHeader http.Header, err error) {
	return nil, nil, errors.New("Bailian video does not support list")
}

func (b *Bailian) VideoRetrieveOfficial(ctx context.Context, taskId string) (responseBytes []byte, responseHeader http.Header, err error) {

	logger.Infof(ctx, "VideoRetrieveOfficial Bailian model: %s, taskId: %s start", b.Model, taskId)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "VideoRetrieveOfficial Bailian model: %s totalTime: %d ms", b.Model, gtime.TimestampMilli()-now)
	}()

	if responseBytes, responseHeader, err = util.HttpGet(ctx, b.BaseUrl+fmt.Sprintf("/tasks/%s", taskId), b.header, nil, nil, b.Timeout, b.ProxyUrl, b.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "VideoRetrieveOfficial Bailian model: %s, error: %v", b.Model, err)
		return nil, nil, err
	}

	logger.Infof(ctx, "VideoRetrieveOfficial Bailian model: %s, taskId: %s finished", b.Model, taskId)

	return responseBytes, responseHeader, nil
}

func (b *Bailian) VideoDeleteOfficial(ctx context.Context, taskId string) (err error) {
	return errors.New("Bailian video does not support delete")
}
