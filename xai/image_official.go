package xai

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

func (x *XAI) ImageGenerationsOfficial(ctx context.Context, data []byte) (responseBytes []byte, responseHeader http.Header, err error) {

	logger.Infof(ctx, "ImageGenerationsOfficial xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "ImageGenerationsOfficial xAI model: %s totalTime: %d ms", x.Model, gtime.TimestampMilli()-now)
	}()

	if x.Path == "" {
		x.Path = "/images/generations"
	}

	data = sanitizeImageGenerationsJSON(data)

	if responseBytes, responseHeader, err = util.HttpPost(ctx, x.BaseUrl+x.Path, x.header, data, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "ImageGenerationsOfficial xAI model: %s, error: %v", x.Model, err)
		return nil, nil, err
	}

	logger.Infof(ctx, "ImageGenerationsOfficial xAI model: %s finished", x.Model)

	return responseBytes, responseHeader, nil
}

func (x *XAI) ImageGenerationsStreamOfficial(ctx context.Context, data []byte) (responseChan chan *model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerationsStreamOfficial xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ImageGenerationsStreamOfficial xAI model: %s totalTime: %d ms", x.Model, gtime.TimestampMilli()-now)
		}
	}()

	if x.Path == "" {
		x.Path = "/images/generations"
	}

	stream, err := util.SSEClient(ctx, x.BaseUrl+x.Path, x.header, data, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageGenerationsStreamOfficial xAI model: %s, error: %v", x.Model, err)
		return responseChan, err
	}

	streamResponseHeaders := stream.Response.Header

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ImageResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ImageGenerationsStreamOfficial xAI model: %s, stream.Close error: %v", x.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ImageGenerationsStreamOfficial xAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", x.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ImageGenerationsStreamOfficial xAI model: %s finished", x.Model)
				} else {
					logger.Errorf(ctx, "ImageGenerationsStreamOfficial xAI model: %s, error: %v", x.Model, err)
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

			response, err := x.ConvImageGenerationsStreamResponse(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ImageGenerationsStreamOfficial xAI ConvImageGenerationsStreamResponse error: %v", err)

				end := gtime.TimestampMilli()
				responseChan <- &model.ImageResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     err,
				}

				return
			}

			end := gtime.TimestampMilli()

			response.ConnTime = duration - now
			response.Duration = end - duration
			response.TotalTime = end - now
			response.ResponseHeaders = streamResponseHeaders
			response.Event = stream.Event()

			responseChan <- &response
		}

	}, nil); err != nil {
		logger.Errorf(ctx, "ImageGenerationsStreamOfficial xAI model: %s, error: %v", x.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
