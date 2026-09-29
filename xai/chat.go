package xai

import (
	"context"
	"io"
	"slices"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (x *XAI) ChatCompletions(ctx context.Context, data any) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	if !slices.Contains(x.ReqPassthroughParams, "req_data") {
		var convReq model.ChatCompletionRequest
		if convReq, err = x.ConvChatCompletionsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "ChatCompletions xAI ConvChatCompletionsRequest error: %v", err)
			return response, err
		}
		data = convChatCompletionsOfficialJSON(convReq)
	}

	if x.Path == "" {
		x.Path = "/chat/completions"
	}

	bytes, responseHeader, err := util.HttpPost(ctx, x.BaseUrl+x.Path, x.header, data, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvChatCompletionsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions xAI ConvChatCompletionsResponse error: %v", err)
		return response, err
	}

	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "ChatCompletions xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) ChatCompletionsStream(ctx context.Context, data any) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream xAI model: %s totalTime: %d ms", x.Model, gtime.TimestampMilli()-now)
		}
	}()

	if !slices.Contains(x.ReqPassthroughParams, "req_data") {
		var convReq model.ChatCompletionRequest
		if convReq, err = x.ConvChatCompletionsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream xAI ConvChatCompletionsRequest error: %v", err)
			return nil, err
		}
		data = convChatCompletionsOfficialJSON(convReq)
	}

	if x.IsSupportStream != nil && !*x.IsSupportStream {
		return x.ChatCompletionStreamToNonStream(ctx, data)
	}

	if x.Path == "" {
		x.Path = "/chat/completions"
	}

	stream, err := util.SSEClient(ctx, x.BaseUrl+x.Path, x.header, data, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream xAI model: %s, error: %v", x.Model, err)
		return responseChan, err
	}

	streamResponseHeaders := stream.Response.Header

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream xAI model: %s, stream.Close error: %v", x.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream xAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", x.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ChatCompletionsStream xAI model: %s finished", x.Model)
				} else {
					logger.Errorf(ctx, "ChatCompletionsStream xAI model: %s, error: %v", x.Model, err)
				}

				end := gtime.TimestampMilli()
				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     err,
				}

				return
			}

			response, err := x.ConvChatCompletionsStreamResponse(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream xAI ConvChatCompletionsStreamResponse error: %v", err)

				end := gtime.TimestampMilli()
				responseChan <- &model.ChatCompletionResponse{
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

			responseChan <- &response
		}

	}, nil); err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream xAI model: %s, error: %v", x.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (x *XAI) ChatCompletionStreamToNonStream(ctx context.Context, data any) (responseChan chan *model.ChatCompletionResponse, err error) {

	request, err := x.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStreamToNonStream xAI ConvChatCompletionsRequest error: %v", err)
		return nil, err
	}

	responseChan = make(chan *model.ChatCompletionResponse)

	now := gtime.TimestampMilli()
	duration := now

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionStreamToNonStream xAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", x.Model, duration-now, end-duration, end-now)
		}()

		request.Stream = false

		response, err := x.ChatCompletions(ctx, gjson.MustEncode(request))
		if err != nil {

			if errors.Is(err, io.EOF) {
				logger.Infof(ctx, "ChatCompletionStreamToNonStream xAI model: %s finished", x.Model)
			} else {
				logger.Errorf(ctx, "ChatCompletionStreamToNonStream xAI model: %s, error: %v", x.Model, err)
			}

			end := gtime.TimestampMilli()
			responseChan <- &model.ChatCompletionResponse{
				ConnTime:  gtime.TimestampMilli() - now,
				Duration:  end - gtime.TimestampMilli(),
				TotalTime: end - now,
				Error:     err,
			}

			return
		}

		duration = gtime.TimestampMilli()
		response.ConnTime = duration - now

		end := gtime.TimestampMilli()
		response.Duration = end - duration
		response.TotalTime = end - now

		responseChan <- &response

		end = gtime.TimestampMilli()
		responseChan <- &model.ChatCompletionResponse{
			Duration:  end - duration,
			TotalTime: end - now,
			Error:     io.EOF,
		}

	}, nil); err != nil {
		logger.Errorf(ctx, "ChatCompletionStreamToNonStream xAI model: %s, error: %v", x.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (x *XAI) DeferredChatCompletion(ctx context.Context, requestId string) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "DeferredChatCompletion xAI model: %s, requestId: %s start", x.Model, requestId)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "DeferredChatCompletion xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	path := x.Path
	if path == "" {
		path = "/chat/deferred-completion/" + requestId
	}

	bytes, responseHeader, err := util.HttpGet(ctx, x.BaseUrl+path, x.header, nil, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "DeferredChatCompletion xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvChatCompletionsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "DeferredChatCompletion xAI ConvChatCompletionsResponse error: %v", err)
		return response, err
	}

	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "DeferredChatCompletion xAI model: %s finished", x.Model)

	return response, nil
}
