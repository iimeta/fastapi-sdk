package google

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

func (g *Google) ChatCompletions(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions Google model: %s start", g.model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions Google model: %s totalTime: %d ms", g.model, response.TotalTime)
	}()

	request, err := g.ConvChatCompletionsRequestOfficial(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions Google ConvChatCompletionsRequest error: %v", err)
		return response, err
	}

	var bytes []byte

	if g.isGcp {
		if bytes, err = util.HttpPost(ctx, fmt.Sprintf("%s:generateContent", g.baseURL+g.path), g.header, request, nil, g.proxyURL); err != nil {
			logger.Errorf(ctx, "ChatCompletions Google model: %s, error: %v", g.model, err)
			return response, err
		}
	} else {
		if bytes, err = util.HttpPost(ctx, fmt.Sprintf("%s:generateContent?key=%s", g.baseURL+g.path, g.key), nil, request, nil, g.proxyURL); err != nil {
			logger.Errorf(ctx, "ChatCompletions Google model: %s, error: %v", g.model, err)
			return response, err
		}
	}

	if response, err = g.ConvChatCompletionsResponseOfficial(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions Google ConvChatCompletionsResponseOfficial error: %v", err)
		return response, err
	}

	return response, nil
}

func (g *Google) ChatCompletionsStream(ctx context.Context, data []byte) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream Google model: %s start", g.model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream Google model: %s totalTime: %d ms", g.model, gtime.TimestampMilli()-now)
		}
	}()

	request, err := g.ConvChatCompletionsRequestOfficial(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream Google ConvChatCompletionsRequestOfficial error: %v", err)
		return nil, err
	}

	var stream *util.StreamReader

	if g.isGcp {
		stream, err = util.SSEClient(ctx, fmt.Sprintf("%s:streamGenerateContent?alt=sse", g.baseURL+g.path), g.header, request, g.proxyURL, g.requestErrorHandler)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream Google model: %s, error: %v", g.model, err)
			return responseChan, err
		}
	} else {
		stream, err = util.SSEClient(ctx, fmt.Sprintf("%s:streamGenerateContent?alt=sse&key=%s", g.baseURL+g.path, g.key), nil, request, g.proxyURL, g.requestErrorHandler)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream Google model: %s, error: %v", g.model, err)
			return responseChan, err
		}
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream Google model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", g.model, duration-now, end-duration, end-now)

			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream Google model: %s, stream.Close error: %v", g.model, err)
			}
		}()

		var (
			usage *model.Usage
		)

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
					logger.Errorf(ctx, "ChatCompletionsStream Google model: %s, error: %v", g.model, err)
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

			response, err := g.ConvChatCompletionsStreamResponseOfficial(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream Google ConvChatCompletionsStreamResponseOfficial error: %v", err)

				end := gtime.TimestampMilli()
				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     err,
				}

				return
			}

			if response.Usage != nil {
				usage = response.Usage
			} else {
				response.Usage = usage
			}

			end := gtime.TimestampMilli()

			response.ConnTime = duration - now
			response.Duration = end - duration
			response.TotalTime = end - now

			responseChan <- &response
		}

	}, nil); err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream Google model: %s, error: %v", g.model, err)
		return responseChan, err
	}

	return responseChan, nil
}
