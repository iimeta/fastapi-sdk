package general

import (
	"context"
	"io"
	"slices"
	"strings"

	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (g *General) ImageGenerations(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerations General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageGenerations General model: %s totalTime: %d ms", g.Model, gtime.TimestampMilli()-now)
	}()

	var request any = data
	if !slices.Contains(g.ReqPassthroughParams, "req_data") {
		if request, err = g.ConvImageGenerationsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "ImageGenerations General ConvImageGenerationsRequest error: %v", err)
			return response, err
		}
	}

	path := g.Path
	if g.Async {
		if strings.Contains(path, "?") {
			path += "&async=true"
		} else {
			path += "?async=true"
		}
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, g.BaseUrl+path, g.header, request, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageGenerations General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvImageGenerationsResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "ImageGenerations General ConvImageGenerationsResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	return response, nil
}

func (g *General) ImageGenerationsStream(ctx context.Context, data []byte) (responseChan chan *model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerationsStream General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ImageGenerationsStream General model: %s totalTime: %d ms", g.Model, gtime.TimestampMilli()-now)
		}
	}()

	var request any = data
	if !slices.Contains(g.ReqPassthroughParams, "req_data") {
		if request, err = g.ConvImageGenerationsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "ImageGenerationsStream General ConvImageGenerationsRequest error: %v", err)
			return nil, err
		}
	}

	path := g.Path
	if g.Async {
		if strings.Contains(path, "?") {
			path += "&async=true"
		} else {
			path += "?async=true"
		}
	}

	stream, err := util.SSEClient(ctx, g.BaseUrl+path, g.header, request, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageGenerationsStream General model: %s, error: %v", g.Model, err)
		return responseChan, err
	}

	streamResponseHeaders := stream.Response.Header

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ImageResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ImageGenerationsStream General model: %s, stream.Close error: %v", g.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ImageGenerationsStream General model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", g.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ImageGenerationsStream General model: %s finished", g.Model)
				} else {
					logger.Errorf(ctx, "ImageGenerationsStream General model: %s, error: %v", g.Model, err)
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

			response, err := g.ConvImageGenerationsStreamResponse(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ImageGenerationsStream General ConvImageGenerationsStreamResponse error: %v", err)

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
		logger.Errorf(ctx, "ImageGenerationsStream General model: %s, error: %v", g.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (g *General) ImageEdits(ctx context.Context, request model.ImageEditRequest) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageEdits General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageEdits General model: %s totalTime: %d ms", g.Model, gtime.TimestampMilli()-now)
	}()

	var data any = request
	if !slices.Contains(g.ReqPassthroughParams, "req_data") {
		data, err = g.ConvImageEditsRequest(ctx, request)
		if err != nil {
			logger.Errorf(ctx, "ImageEdits General ConvImageEditsRequest error: %v", err)
			return response, err
		}
	}

	path := g.Path
	if g.Async {
		if strings.Contains(path, "?") {
			path += "&async=true"
		} else {
			path += "?async=true"
		}
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, g.BaseUrl+path, g.header, data, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageEdits General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvImageEditsResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "ImageEdits General ConvImageEditsResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	return response, nil
}

func (g *General) ImageEditsStream(ctx context.Context, request model.ImageEditRequest) (responseChan chan *model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageEditsStream General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ImageEditsStream General model: %s totalTime: %d ms", g.Model, gtime.TimestampMilli()-now)
		}
	}()

	var data any = request
	if !slices.Contains(g.ReqPassthroughParams, "req_data") {
		data, err = g.ConvImageEditsRequest(ctx, request)
		if err != nil {
			logger.Errorf(ctx, "ImageEditsStream General ConvImageEditsRequest error: %v", err)
			return nil, err
		}
	}

	path := g.Path
	if g.Async {
		if strings.Contains(path, "?") {
			path += "&async=true"
		} else {
			path += "?async=true"
		}
	}

	stream, err := util.SSEClient(ctx, g.BaseUrl+path, g.header, data, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageEditsStream General model: %s, error: %v", g.Model, err)
		return responseChan, err
	}

	streamResponseHeaders := stream.Response.Header

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ImageResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ImageEditsStream General model: %s, stream.Close error: %v", g.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ImageEditsStream General model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", g.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ImageEditsStream General model: %s finished", g.Model)
				} else {
					logger.Errorf(ctx, "ImageEditsStream General model: %s, error: %v", g.Model, err)
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

			response, err := g.ConvImageGenerationsStreamResponse(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ImageEditsStream General ConvImageGenerationsStreamResponse error: %v", err)

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
		logger.Errorf(ctx, "ImageEditsStream General model: %s, error: %v", g.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
