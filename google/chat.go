package google

import (
	"context"
	"fmt"
	"io"

	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (g *Google) ChatCompletions(ctx context.Context, data any) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions Google model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions Google model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	if !g.IsOfficialFormatRequest {

		request, err := g.ConvChatCompletionsRequest(ctx, data)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletions Google ConvChatCompletionsRequest error: %v", err)
			return response, err
		}

		if data, err = g.ConvChatCompletionsRequestOfficial(ctx, request); err != nil {
			logger.Errorf(ctx, "ChatCompletions Google ConvChatCompletionsRequestOfficial error: %v", err)
			return response, err
		}
	}

	if g.Path == "" {
		g.Path = "/models/" + g.Model
	}

	if g.Action == "" {
		g.Action = "generateContent"
	}

	var bytes []byte

	if g.isGcp {
		if bytes, err = util.HttpPost(ctx, fmt.Sprintf("%s%s:%s", g.BaseUrl, g.Path, g.Action), g.header, data, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler); err != nil {
			logger.Errorf(ctx, "ChatCompletions Google model: %s, error: %v", g.Model, err)
			return response, err
		}
	} else {
		if bytes, err = util.HttpPost(ctx, fmt.Sprintf("%s%s:%s?key=%s", g.BaseUrl, g.Path, g.Action, g.Key), g.header, data, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler); err != nil {
			logger.Errorf(ctx, "ChatCompletions Google model: %s, error: %v", g.Model, err)
			return response, err
		}
	}

	if response, err = g.ConvChatCompletionsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions Google ConvChatCompletionsResponse error: %v", err)
		return response, err
	}

	return response, nil
}

func (g *Google) ChatCompletionsStream(ctx context.Context, data any) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream Google model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream Google model: %s totalTime: %d ms", g.Model, gtime.TimestampMilli()-now)
		}
	}()

	if !g.IsOfficialFormatRequest {

		request, err := g.ConvChatCompletionsRequest(ctx, data)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream Google ConvChatCompletionsRequest error: %v", err)
			return responseChan, err
		}

		if data, err = g.ConvChatCompletionsRequestOfficial(ctx, request); err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream Google ConvChatCompletionsRequestOfficial error: %v", err)
			return responseChan, err
		}
	}

	if g.Path == "" {
		g.Path = "/models/" + g.Model
	}

	if g.Action == "" {
		g.Action = "streamGenerateContent"
	}

	var stream *util.StreamReader

	if g.isGcp {
		stream, err = util.SSEClient(ctx, fmt.Sprintf("%s%s:%s?alt=sse", g.BaseUrl, g.Path, g.Action), g.header, data, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream Google model: %s, error: %v", g.Model, err)
			return responseChan, err
		}
	} else {
		stream, err = util.SSEClient(ctx, fmt.Sprintf("%s%s:%s?alt=sse&key=%s", g.BaseUrl, g.Path, g.Action, g.Key), g.header, data, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream Google model: %s, error: %v", g.Model, err)
			return responseChan, err
		}
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream Google model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", g.Model, duration-now, end-duration, end-now)

			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream Google model: %s, stream.Close error: %v", g.Model, err)
			}
		}()

		var (
			usage *model.Usage
		)

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ChatCompletionsStream Google model: %s finished", g.Model)
				} else {
					logger.Errorf(ctx, "ChatCompletionsStream Google model: %s, error: %v", g.Model, err)
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
				logger.Errorf(ctx, "ChatCompletionsStream Google ConvChatCompletionsStreamResponse error: %v", err)

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
		logger.Errorf(ctx, "ChatCompletionsStream Google model: %s, error: %v", g.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
