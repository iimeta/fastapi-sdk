package baidu

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (b *Baidu) ChatCompletions(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions Baidu model: %s start", b.model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions Baidu model: %s totalTime: %d ms", b.model, response.TotalTime)
	}()

	request, err := b.ConvChatCompletionsRequestOfficial(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions Baidu ConvChatCompletionsRequestOfficial error: %v", err)
		return response, err
	}

	bytes, err := util.HttpPost(ctx, fmt.Sprintf("%s?access_token=%s", b.baseURL+b.path, b.accessToken), nil, request, nil, b.proxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions Baidu model: %s, error: %v", b.model, err)
		return response, err
	}

	if response, err = b.ConvChatCompletionsResponseOfficial(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions Baidu ConvChatCompletionsResponseOfficial error: %v", err)
		return response, err
	}

	return response, nil
}

func (b *Baidu) ChatCompletionsStream(ctx context.Context, data []byte) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream Baidu model: %s start", b.model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream Baidu model: %s totalTime: %d ms", b.model, gtime.TimestampMilli()-now)
		}
	}()

	request, err := b.ConvChatCompletionsRequestOfficial(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream Baidu ConvChatCompletionsRequestOfficial error: %v", err)
		return nil, err
	}

	stream, err := util.SSEClient(ctx, fmt.Sprintf("%s?access_token=%s", b.baseURL+b.path, b.accessToken), nil, request, b.proxyURL, b.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream Baidu model: %s, error: %v", b.model, err)
		return responseChan, err
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream Baidu model: %s, stream.Close error: %v", b.model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream Baidu model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", b.model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
					logger.Errorf(ctx, "ChatCompletionsStream Baidu model: %s, error: %v", b.model, err)
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

			response, err := b.ConvChatCompletionsStreamResponseOfficial(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream Baidu ConvChatCompletionsStreamResponseOfficial error: %v", err)

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

			responseChan <- &response
		}

	}, nil); err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream Baidu model: %s, error: %v", b.model, err)
		return responseChan, err
	}

	return responseChan, nil
}
