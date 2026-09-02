package volcengine

import (
	"context"
	"io"
	"net/http"

	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (v *VolcEngine) ImageGenerationsOfficial(ctx context.Context, data []byte) (responseBytes []byte, responseHeader http.Header, err error) {

	logger.Infof(ctx, "ImageGenerationsOfficial VolcEngine model: %s start", v.Model)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "ImageGenerationsOfficial VolcEngine model: %s totalTime: %d ms", v.Model, gtime.TimestampMilli()-now)
	}()

	if v.Path == "" {
		v.Path = "/images/generations"
	}

	if responseBytes, responseHeader, err = util.HttpPost(ctx, v.BaseUrl+v.Path, v.header, data, nil, v.Timeout, v.ProxyUrl, v.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "ImageGenerationsOfficial VolcEngine model: %s, error: %v", v.Model, err)
		return nil, nil, err
	}

	logger.Infof(ctx, "ImageGenerationsOfficial VolcEngine model: %s finished", v.Model)

	return responseBytes, responseHeader, nil
}

func (v *VolcEngine) ImageGenerationsStreamOfficial(ctx context.Context, data []byte) (responseChan chan *model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerationsStreamOfficial VolcEngine model: %s start", v.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ImageGenerationsStreamOfficial VolcEngine model: %s totalTime: %d ms", v.Model, gtime.TimestampMilli()-now)
		}
	}()

	if v.Path == "" {
		v.Path = "/images/generations"
	}

	stream, err := util.SSEClient(ctx, v.BaseUrl+v.Path, v.header, data, v.Timeout, v.ProxyUrl, v.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageGenerationsStreamOfficial VolcEngine model: %s, error: %v", v.Model, err)
		return responseChan, err
	}

	streamResponseHeaders := stream.Response.Header
	duration := gtime.TimestampMilli()
	responseChan = make(chan *model.ImageResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ImageGenerationsStreamOfficial VolcEngine model: %s, stream.Close error: %v", v.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ImageGenerationsStreamOfficial VolcEngine model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", v.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ImageGenerationsStreamOfficial VolcEngine model: %s finished", v.Model)
				} else {
					logger.Errorf(ctx, "ImageGenerationsStreamOfficial VolcEngine model: %s, error: %v", v.Model, err)
				}

				end := gtime.TimestampMilli()
				responseChan <- &model.ImageResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     err,
				}
				return
			}

			response, convErr := v.ConvImageGenerationsStreamResponse(ctx, responseBytes)
			if convErr != nil {
				response = model.ImageResponse{ResponseBytes: responseBytes}
			}

			end := gtime.TimestampMilli()
			response.ConnTime = duration - now
			response.Duration = end - duration
			response.TotalTime = end - now
			response.ResponseBytes = responseBytes
			response.ResponseHeaders = streamResponseHeaders
			response.Event = stream.Event()
			responseChan <- &response
		}

	}, nil); err != nil {
		logger.Errorf(ctx, "ImageGenerationsStreamOfficial VolcEngine model: %s, error: %v", v.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
