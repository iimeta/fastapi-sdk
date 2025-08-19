package zhipuai

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

func (z *ZhipuAI) ChatCompletions(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions ZhipuAI model: %s start", z.model)

	request, err := z.ConvChatCompletionsRequestOfficial(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions ZhipuAI ConvChatCompletionsRequestOfficial error: %v", err)
		return response, err
	}

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions ZhipuAI model: %s totalTime: %d ms", z.model, response.TotalTime)
	}()

	bytes, err := util.HttpPost(ctx, z.baseURL+z.path, z.header, request, nil, z.proxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions ZhipuAI model: %s, error: %v", z.model, err)
		return response, err
	}

	if response, err = z.ConvChatCompletionsResponseOfficial(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions ZhipuAI ConvChatCompletionsResponseOfficial error: %v", err)
		return response, err
	}

	return response, nil
}

func (z *ZhipuAI) ChatCompletionsStream(ctx context.Context, data []byte) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream ZhipuAI model: %s start", z.model)

	request, err := z.ConvChatCompletionsRequestOfficial(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream ZhipuAI ConvChatCompletionsRequestOfficial error: %v", err)
		return nil, err
	}

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream ZhipuAI model: %s totalTime: %d ms", z.model, gtime.TimestampMilli()-now)
		}
	}()

	stream, err := util.SSEClient(ctx, z.baseURL+z.path, z.header, request, z.proxyURL, z.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream ZhipuAI model: %s, error: %v", z.model, err)
		return responseChan, err
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream ZhipuAI model: %s, stream.Close error: %v", z.model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream ZhipuAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", z.model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
					logger.Errorf(ctx, "ChatCompletionsStream ZhipuAI model: %s, error: %v", z.model, err)
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

			response, err := z.ConvChatCompletionsStreamResponseOfficial(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream ZhipuAI ConvChatCompletionsStreamResponseOfficial error: %v", err)

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
		logger.Errorf(ctx, "ChatCompletionsStream ZhipuAI model: %s, error: %v", z.model, err)
		return responseChan, err
	}

	return responseChan, nil
}
