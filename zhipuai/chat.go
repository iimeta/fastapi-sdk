package zhipuai

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

func (z *ZhipuAI) ChatCompletions(ctx context.Context, data any) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions ZhipuAI model: %s start", z.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions ZhipuAI model: %s totalTime: %d ms", z.Model, response.TotalTime)
	}()

	if !z.IsOfficialFormatRequest {

		request, err := z.ConvChatCompletionsRequest(ctx, data)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletions ZhipuAI ConvChatCompletionsRequest error: %v", err)
			return response, err
		}

		if data, err = z.ConvChatCompletionsRequestOfficial(ctx, request); err != nil {
			logger.Errorf(ctx, "ChatCompletions ZhipuAI ConvChatCompletionsRequestOfficial error: %v", err)
			return response, err
		}
	}

	bytes, err := util.HttpPost(ctx, z.BaseUrl+z.Path, z.header, data, nil, z.Timeout, z.ProxyUrl, z.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions ZhipuAI model: %s, error: %v", z.Model, err)
		return response, err
	}

	if response, err = z.ConvChatCompletionsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions ZhipuAI ConvChatCompletionsResponse error: %v", err)
		return response, err
	}

	return response, nil
}

func (z *ZhipuAI) ChatCompletionsStream(ctx context.Context, data any) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream ZhipuAI model: %s start", z.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream ZhipuAI model: %s totalTime: %d ms", z.Model, gtime.TimestampMilli()-now)
		}
	}()

	if !z.IsOfficialFormatRequest {

		request, err := z.ConvChatCompletionsRequest(ctx, data)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream ZhipuAI ConvChatCompletionsRequest error: %v", err)
			return responseChan, err
		}

		if data, err = z.ConvChatCompletionsRequestOfficial(ctx, request); err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream ZhipuAI ConvChatCompletionsRequestOfficial error: %v", err)
			return responseChan, err
		}
	}

	stream, err := util.SSEClient(ctx, z.BaseUrl+z.Path, z.header, data, z.Timeout, z.ProxyUrl, z.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream ZhipuAI model: %s, error: %v", z.Model, err)
		return responseChan, err
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream ZhipuAI model: %s, stream.Close error: %v", z.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream ZhipuAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", z.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ChatCompletionsStream ZhipuAI model: %s finished", z.Model)
				} else {
					logger.Errorf(ctx, "ChatCompletionsStream ZhipuAI model: %s, error: %v", z.Model, err)
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

			response, err := z.ConvChatCompletionsStreamResponse(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream ZhipuAI ConvChatCompletionsStreamResponse error: %v", err)

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
		logger.Errorf(ctx, "ChatCompletionsStream ZhipuAI model: %s, error: %v", z.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
