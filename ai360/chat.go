package ai360

import (
	"context"
	"errors"
	"io"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (a *AI360) ChatCompletions(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	request, err := a.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions 360AI ConvChatCompletionsRequest error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "ChatCompletions 360AI model: %s start", a.model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions 360AI model: %s totalTime: %d ms", a.model, response.TotalTime)
	}()

	bytes, err := util.HttpPost(ctx, a.baseURL+a.path, a.header, gjson.MustEncode(request), nil, a.proxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions 360AI model: %s, error: %v", a.model, err)
		return response, err
	}

	if response, err = a.ConvChatCompletionsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions 360AI ConvChatCompletionsResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "ChatCompletions 360AI model: %s finished", a.model)

	return response, nil
}

func (a *AI360) ChatCompletionsStream(ctx context.Context, data []byte) (responseChan chan *model.ChatCompletionResponse, err error) {

	request, err := a.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream 360AI ConvChatCompletionsRequest error: %v", err)
		return nil, err
	}

	logger.Infof(ctx, "ChatCompletionsStream 360AI model: %s start", a.model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream 360AI model: %s totalTime: %d ms", a.model, gtime.TimestampMilli()-now)
		}
	}()

	stream, err := util.SSEClient(ctx, a.baseURL+a.path, a.header, gjson.MustEncode(request), a.proxyURL, nil)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream 360AI model: %s, error: %v", a.model, err)
		return responseChan, a.apiErrorHandler(err)
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream 360AI model: %s, stream.Close error: %v", a.model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream 360AI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", a.model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
					logger.Errorf(ctx, "ChatCompletionsStream 360AI model: %s, error: %v", a.model, err)
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

			response, err := a.ConvChatCompletionsStreamResponse(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream 360AI ConvChatCompletionsStreamResponse error: %v", err)

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
		logger.Errorf(ctx, "ChatCompletionsStream 360AI model: %s, error: %v", a.model, err)
		return responseChan, err
	}

	return responseChan, nil
}
