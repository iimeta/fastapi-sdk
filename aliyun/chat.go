package aliyun

import (
	"context"
	"errors"
	"io"

	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (a *Aliyun) ChatCompletions(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions Aliyun model: %s start", a.model)

	request, err := a.ConvChatCompletionsRequestOfficial(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions Aliyun ConvChatCompletionsRequestOfficial error: %v", err)
		return response, err
	}

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions Aliyun model: %s totalTime: %d ms", a.model, response.TotalTime)
	}()

	bytes, err := util.HttpPost(ctx, a.baseURL+a.path, a.header, request, nil, a.proxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions Aliyun model: %s, error: %v", a.model, err)
		return response, err
	}

	if response, err = a.ConvChatCompletionsResponseOfficial(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions Aliyun ConvChatCompletionsResponseOfficial error: %v", err)
		return response, err
	}

	return response, nil
}

func (a *Aliyun) ChatCompletionsStream(ctx context.Context, data []byte) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream Aliyun model: %s start", a.model)

	request, err := a.ConvChatCompletionsRequestOfficial(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream Aliyun ConvChatCompletionsRequestOfficial error: %v", err)
		return nil, err
	}

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream Aliyun model: %s totalTime: %d ms", a.model, gtime.TimestampMilli()-now)
		}
	}()

	stream, err := util.SSEClient(ctx, a.baseURL+a.path, a.header, request, a.proxyURL, a.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream Aliyun model: %s, error: %v", a.model, err)
		return responseChan, err
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream Aliyun model: %s, stream.Close error: %v", a.model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream Aliyun model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", a.model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
					logger.Errorf(ctx, "ChatCompletionsStream Aliyun model: %s, error: %v", a.model, err)
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

			response, err := a.ConvChatCompletionsStreamResponseOfficial(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream Aliyun ConvChatCompletionsStreamResponseOfficial error: %v", err)

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

			response.ResponseBytes = responseBytes
			response.ConnTime = duration - now
			response.Duration = end - duration
			response.TotalTime = end - now

			responseChan <- &response
		}

	}, nil); err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream Aliyun model: %s, error: %v", a.model, err)
		return responseChan, err
	}

	return responseChan, nil
}
