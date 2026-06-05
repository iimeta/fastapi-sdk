package openai

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

func (o *OpenAI) ImageGenerations(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerations OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageGenerations OpenAI model: %s totalTime: %d ms", o.Model, gtime.TimestampMilli()-now)
	}()

	var request any = data
	if !slices.Contains(o.ReqPassthroughParams, "req_data") {
		if request, err = o.ConvImageGenerationsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "ImageGenerations OpenAI ConvImageGenerationsRequest error: %v", err)
			return response, err
		}
	}

	if o.Path == "" {
		if o.isAzure {
			o.Path = "/images/generations?api-version=" + o.apiVersion
		} else {
			o.Path = "/images/generations"
		}
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, o.BaseUrl+o.Path, o.header, request, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageGenerations OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvImageGenerationsResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "ImageGenerations OpenAI ConvImageGenerationsResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	return response, nil
}

func (o *OpenAI) ImageGenerationsStream(ctx context.Context, data []byte) (responseChan chan *model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerationsStream OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ImageGenerationsStream OpenAI model: %s totalTime: %d ms", o.Model, gtime.TimestampMilli()-now)
		}
	}()

	var request any = data
	if !slices.Contains(o.ReqPassthroughParams, "req_data") {
		if request, err = o.ConvImageGenerationsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "ImageGenerationsStream OpenAI ConvImageGenerationsRequest error: %v", err)
			return nil, err
		}
	}

	if o.Path == "" {
		if o.isAzure {
			o.Path = "/images/generations?api-version=" + o.apiVersion
		} else {
			o.Path = "/images/generations"
		}
	}

	stream, err := util.SSEClient(ctx, o.BaseUrl+o.Path, o.header, request, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageGenerationsStream OpenAI model: %s, error: %v", o.Model, err)
		return responseChan, err
	}

	streamResponseHeaders := stream.Response.Header

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ImageResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ImageGenerationsStream OpenAI model: %s, stream.Close error: %v", o.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ImageGenerationsStream OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", o.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ImageGenerationsStream OpenAI model: %s finished", o.Model)
				} else {
					logger.Errorf(ctx, "ImageGenerationsStream OpenAI model: %s, error: %v", o.Model, err)
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

			response, err := o.ConvImageGenerationsStreamResponse(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ImageGenerationsStream OpenAI ConvImageGenerationsStreamResponse error: %v", err)

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
		logger.Errorf(ctx, "ImageGenerationsStream OpenAI model: %s, error: %v", o.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (o *OpenAI) ImageEdits(ctx context.Context, request model.ImageEditRequest) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageEdits OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageEdits OpenAI model: %s totalTime: %d ms", o.Model, gtime.TimestampMilli()-now)
	}()

	data, err := o.ConvImageEditsRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "ImageEdits OpenAI ConvImageEditsRequest error: %v", err)
		return response, err
	}

	if o.Path == "" {
		if o.isAzure {
			o.Path = "/images/edits?api-version=" + o.apiVersion
		} else {
			o.Path = "/images/edits"
		}
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, o.BaseUrl+o.Path, o.header, data, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageEdits OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvImageEditsResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "ImageEdits OpenAI ConvImageEditsResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	return response, nil
}

func (o *OpenAI) ImageEditsStream(ctx context.Context, request model.ImageEditRequest) (responseChan chan *model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageEditsStream OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ImageEditsStream OpenAI model: %s totalTime: %d ms", o.Model, gtime.TimestampMilli()-now)
		}
	}()

	data, err := o.ConvImageEditsRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "ImageEditsStream OpenAI ConvImageEditsRequest error: %v", err)
		return nil, err
	}

	if o.Path == "" {
		if o.isAzure {
			o.Path = "/images/edits?api-version=" + o.apiVersion
		} else {
			o.Path = "/images/edits"
		}
	}

	stream, err := util.SSEClient(ctx, o.BaseUrl+o.Path, o.header, data, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageEditsStream OpenAI model: %s, error: %v", o.Model, err)
		return responseChan, err
	}

	streamResponseHeaders := stream.Response.Header

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ImageResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ImageEditsStream OpenAI model: %s, stream.Close error: %v", o.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ImageEditsStream OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", o.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ImageEditsStream OpenAI model: %s finished", o.Model)
				} else {
					logger.Errorf(ctx, "ImageEditsStream OpenAI model: %s, error: %v", o.Model, err)
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

			response, err := o.ConvImageGenerationsStreamResponse(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ImageEditsStream OpenAI ConvImageGenerationsStreamResponse error: %v", err)

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
		logger.Errorf(ctx, "ImageEditsStream OpenAI model: %s, error: %v", o.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
