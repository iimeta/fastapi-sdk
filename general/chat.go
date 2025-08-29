package general

import (
	"context"
	"io"

	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/errors"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (g *General) ChatCompletions(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	bytes, err := util.HttpPost(ctx, g.BaseUrl+g.Path, g.header, data, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvChatCompletionsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions General ConvChatCompletionsResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "ChatCompletions General model: %s finished", g.Model)

	return response, nil
}

func (g *General) ChatCompletionsStream(ctx context.Context, data []byte) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream General model: %s totalTime: %d ms", g.Model, gtime.TimestampMilli()-now)
		}
	}()

	stream, err := util.SSEClient(ctx, g.BaseUrl+g.Path, g.header, data, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream General model: %s, error: %v", g.Model, err)
		return responseChan, err
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream General model: %s, stream.Close error: %v", g.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream General model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", g.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ChatCompletionsStream General model: %s finished", g.Model)
				} else {
					logger.Errorf(ctx, "ChatCompletionsStream General model: %s, error: %v", g.Model, err)
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

			response, err := g.ConvChatCompletionsStreamResponse(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream General ConvChatCompletionsStreamResponse error: %v", err)

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
		logger.Errorf(ctx, "ChatCompletionsStream General model: %s, error: %v", g.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
