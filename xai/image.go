package xai

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

func (x *XAI) ImageGenerations(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerations xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageGenerations xAI model: %s totalTime: %d ms", x.Model, gtime.TimestampMilli()-now)
	}()

	var request any
	if slices.Contains(x.ReqPassthroughParams, "req_data") {
		request = sanitizeImageGenerationsJSON(data)
	} else {
		var convReq model.ImageGenerationRequest
		if convReq, err = x.ConvImageGenerationsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "ImageGenerations xAI ConvImageGenerationsRequest error: %v", err)
			return response, err
		}
		request = convImageGenerationsOfficialJSON(convReq)
	}

	if x.Path == "" {
		x.Path = "/images/generations"
	}

	if x.Async {
		if strings.Contains(x.Path, "?") {
			x.Path += "&async=true"
		} else {
			x.Path += "?async=true"
		}
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, x.BaseUrl+x.Path, x.header, request, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageGenerations xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvImageGenerationsResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "ImageGenerations xAI ConvImageGenerationsResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	return response, nil
}

func (x *XAI) ImageGenerationsStream(ctx context.Context, data []byte) (responseChan chan *model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerationsStream xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ImageGenerationsStream xAI model: %s totalTime: %d ms", x.Model, gtime.TimestampMilli()-now)
		}
	}()

	var request any
	if slices.Contains(x.ReqPassthroughParams, "req_data") {
		request = sanitizeImageGenerationsJSON(data)
	} else {
		var convReq model.ImageGenerationRequest
		if convReq, err = x.ConvImageGenerationsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "ImageGenerationsStream xAI ConvImageGenerationsRequest error: %v", err)
			return nil, err
		}
		request = convImageGenerationsOfficialJSON(convReq)
	}

	if x.Path == "" {
		x.Path = "/images/generations"
	}

	stream, err := util.SSEClient(ctx, x.BaseUrl+x.Path, x.header, request, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageGenerationsStream xAI model: %s, error: %v", x.Model, err)
		return responseChan, err
	}

	streamResponseHeaders := stream.Response.Header

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ImageResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ImageGenerationsStream xAI model: %s, stream.Close error: %v", x.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ImageGenerationsStream xAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", x.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ImageGenerationsStream xAI model: %s finished", x.Model)
				} else {
					logger.Errorf(ctx, "ImageGenerationsStream xAI model: %s, error: %v", x.Model, err)
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
				logger.Errorf(ctx, "ImageGenerationsStream xAI ConvImageGenerationsStreamResponse error: %v", err)

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
		logger.Errorf(ctx, "ImageGenerationsStream xAI model: %s, error: %v", x.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (x *XAI) ImageEdits(ctx context.Context, request model.ImageEditRequest) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageEdits xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageEdits xAI model: %s totalTime: %d ms", x.Model, gtime.TimestampMilli()-now)
	}()

	data, err := x.ConvImageEditsRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "ImageEdits xAI ConvImageEditsRequest error: %v", err)
		return response, err
	}

	if x.Path == "" {
		x.Path = "/images/generations"
	}

	if x.Async {
		if strings.Contains(x.Path, "?") {
			x.Path += "&async=true"
		} else {
			x.Path += "?async=true"
		}
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, x.BaseUrl+x.Path, x.header, data, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageEdits xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvImageEditsResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "ImageEdits xAI ConvImageEditsResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	return response, nil
}

func (x *XAI) ImageEditsStream(ctx context.Context, request model.ImageEditRequest) (responseChan chan *model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageEditsStream xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ImageEditsStream xAI model: %s totalTime: %d ms", x.Model, gtime.TimestampMilli()-now)
		}
	}()

	data, err := x.ConvImageEditsRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "ImageEditsStream xAI ConvImageEditsRequest error: %v", err)
		return nil, err
	}

	if x.Path == "" {
		x.Path = "/images/generations"
	}

	stream, err := util.SSEClient(ctx, x.BaseUrl+x.Path, x.header, data, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageEditsStream xAI model: %s, error: %v", x.Model, err)
		return responseChan, err
	}

	streamResponseHeaders := stream.Response.Header

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ImageResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ImageEditsStream xAI model: %s, stream.Close error: %v", x.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ImageEditsStream xAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", x.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ImageEditsStream xAI model: %s finished", x.Model)
				} else {
					logger.Errorf(ctx, "ImageEditsStream xAI model: %s, error: %v", x.Model, err)
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
				logger.Errorf(ctx, "ImageEditsStream xAI ConvImageGenerationsStreamResponse error: %v", err)

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
		logger.Errorf(ctx, "ImageEditsStream xAI model: %s, error: %v", x.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
