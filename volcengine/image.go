package volcengine

import (
	"context"
	"io"
	"slices"

	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (v *VolcEngine) ImageGenerations(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerations VolcEngine model: %s start", v.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageGenerations VolcEngine model: %s totalTime: %d ms", v.Model, response.TotalTime)
	}()

	var request any = data
	if !slices.Contains(v.ReqPassthroughParams, "req_data") {
		if request, err = v.ConvImageGenerationsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "ImageGenerations VolcEngine ConvImageGenerationsRequest error: %v", err)
			return response, err
		}
	}

	if v.Path == "" {
		v.Path = "/images/generations"
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, v.BaseUrl+v.Path, v.header, request, nil, v.Timeout, v.ProxyUrl, v.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageGenerations VolcEngine model: %s, error: %v", v.Model, err)
		return response, err
	}

	if response, err = v.ConvImageGenerationsResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "ImageGenerations VolcEngine ConvImageGenerationsResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "ImageGenerations VolcEngine model: %s finished", v.Model)

	return response, nil
}

func (v *VolcEngine) ImageEdits(ctx context.Context, request model.ImageEditRequest) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ImageGenerationsStream(ctx context.Context, data []byte) (responseChan chan *model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerationsStream VolcEngine model: %s start", v.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ImageGenerationsStream VolcEngine model: %s totalTime: %d ms", v.Model, gtime.TimestampMilli()-now)
		}
	}()

	var request any = data
	if !slices.Contains(v.ReqPassthroughParams, "req_data") {
		if request, err = v.ConvImageGenerationsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "ImageGenerationsStream VolcEngine ConvImageGenerationsRequest error: %v", err)
			return nil, err
		}
	}

	if v.Path == "" {
		v.Path = "/images/generations"
	}

	stream, err := util.SSEClient(ctx, v.BaseUrl+v.Path, v.header, request, v.Timeout, v.ProxyUrl, v.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageGenerationsStream VolcEngine model: %s, error: %v", v.Model, err)
		return responseChan, err
	}

	streamResponseHeaders := stream.Response.Header
	duration := gtime.TimestampMilli()
	responseChan = make(chan *model.ImageResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ImageGenerationsStream VolcEngine model: %s, stream.Close error: %v", v.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ImageGenerationsStream VolcEngine model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", v.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ImageGenerationsStream VolcEngine model: %s finished", v.Model)
				} else {
					logger.Errorf(ctx, "ImageGenerationsStream VolcEngine model: %s, error: %v", v.Model, err)
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

			response, err := v.ConvImageGenerationsStreamResponse(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ImageGenerationsStream VolcEngine ConvImageGenerationsStreamResponse error: %v", err)

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
		logger.Errorf(ctx, "ImageGenerationsStream VolcEngine model: %s, error: %v", v.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (v *VolcEngine) ImageEditsStream(ctx context.Context, request model.ImageEditRequest) (responseChan chan *model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}
