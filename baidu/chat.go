package baidu

import (
	"context"
	"io"

	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (b *Baidu) ChatCompletions(ctx context.Context, data any) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions Baidu model: %s start", b.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions Baidu model: %s totalTime: %d ms", b.Model, response.TotalTime)
	}()

	if !b.IsOfficialFormatRequest {

		request, err := b.ConvChatCompletionsRequest(ctx, data)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletions Baidu ConvChatCompletionsRequest error: %v", err)
			return response, err
		}

		if data, err = b.ConvChatCompletionsRequestOfficial(ctx, request); err != nil {
			logger.Errorf(ctx, "ChatCompletions Baidu ConvChatCompletionsRequestOfficial error: %v", err)
			return response, err
		}
	}

	bytes, err := util.HttpPost(ctx, b.BaseUrl+b.Path, b.header, data, nil, b.Timeout, b.ProxyUrl, b.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions Baidu model: %s, error: %v", b.Model, err)
		return response, err
	}

	if response, err = b.ConvChatCompletionsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions Baidu ConvChatCompletionsResponse error: %v", err)
		return response, err
	}

	return response, nil
}

func (b *Baidu) ChatCompletionsStream(ctx context.Context, data any) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream Baidu model: %s start", b.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream Baidu model: %s totalTime: %d ms", b.Model, gtime.TimestampMilli()-now)
		}
	}()

	if !b.IsOfficialFormatRequest {

		request, err := b.ConvChatCompletionsRequest(ctx, data)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream Baidu ConvChatCompletionsRequest error: %v", err)
			return responseChan, err
		}

		if data, err = b.ConvChatCompletionsRequestOfficial(ctx, request); err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream Baidu ConvChatCompletionsRequestOfficial error: %v", err)
			return responseChan, err
		}
	}

	stream, err := util.SSEClient(ctx, b.BaseUrl+b.Path, b.header, data, b.Timeout, b.ProxyUrl, b.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream Baidu model: %s, error: %v", b.Model, err)
		return responseChan, err
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream Baidu model: %s, stream.Close error: %v", b.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream Baidu model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", b.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ChatCompletionsStream Baidu model: %s finished", b.Model)
				} else {
					logger.Errorf(ctx, "ChatCompletionsStream Baidu model: %s, error: %v", b.Model, err)
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

			response, err := b.ConvChatCompletionsStreamResponse(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream Baidu ConvChatCompletionsStreamResponse error: %v", err)

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
		logger.Errorf(ctx, "ChatCompletionsStream Baidu model: %s, error: %v", b.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
