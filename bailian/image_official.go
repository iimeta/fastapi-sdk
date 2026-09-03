package bailian

import (
	"context"
	"net/http"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

// 官方格式请求体与响应体原样透传, 仅替换鉴权
func (b *Bailian) ImageGenerationsOfficial(ctx context.Context, data []byte) (responseBytes []byte, responseHeader http.Header, err error) {

	logger.Infof(ctx, "ImageGenerationsOfficial Bailian model: %s start", b.Model)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "ImageGenerationsOfficial Bailian model: %s totalTime: %d ms", b.Model, gtime.TimestampMilli()-now)
	}()

	if b.Path == "" {
		b.Path = "/services/aigc/multimodal-generation/generation"
	}

	if responseBytes, responseHeader, err = util.HttpPost(ctx, b.BaseUrl+b.Path, b.header, data, nil, b.Timeout, b.ProxyUrl, b.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "ImageGenerationsOfficial Bailian model: %s, error: %v", b.Model, err)
		return nil, nil, err
	}

	logger.Infof(ctx, "ImageGenerationsOfficial Bailian model: %s finished", b.Model)

	return responseBytes, responseHeader, nil
}

func (b *Bailian) ImageGenerationsStreamOfficial(ctx context.Context, data []byte) (responseChan chan *model.ImageResponse, err error) {
	return nil, errors.New("Bailian image generations does not support stream")
}
