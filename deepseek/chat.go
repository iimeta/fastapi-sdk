package deepseek

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

func (d *DeepSeek) ChatCompletions(ctx context.Context, data any) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions DeepSeek model: %s start", d.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions DeepSeek model: %s totalTime: %d ms", d.Model, response.TotalTime)
	}()

	if !d.IsOfficialFormatRequest {
		if data, err = d.ConvChatCompletionsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "ChatCompletions DeepSeek ConvChatCompletionsRequest error: %v", err)
			return response, err
		}
	}

	if d.Path == "" {
		d.Path = "/chat/completions"
	}

	bytes, err := util.HttpPost(ctx, d.BaseUrl+d.Path, d.header, data, nil, d.Timeout, d.ProxyUrl, d.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions DeepSeek model: %s, error: %v", d.Model, err)
		return response, err
	}

	if response, err = d.ConvChatCompletionsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions DeepSeek ConvChatCompletionsResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "ChatCompletions DeepSeek model: %s finished", d.Model)

	return response, nil
}

func (d *DeepSeek) ChatCompletionsStream(ctx context.Context, data any) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream DeepSeek model: %s start", d.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream DeepSeek model: %s totalTime: %d ms", d.Model, gtime.TimestampMilli()-now)
		}
	}()

	if !d.IsOfficialFormatRequest {
		if data, err = d.ConvChatCompletionsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream DeepSeek ConvChatCompletionsRequest error: %v", err)
			return nil, err
		}
	}

	if d.Path == "" {
		d.Path = "/chat/completions"
	}

	stream, err := util.SSEClient(ctx, d.BaseUrl+d.Path, d.header, data, d.Timeout, d.ProxyUrl, d.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream DeepSeek model: %s, error: %v", d.Model, err)
		return responseChan, err
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream DeepSeek model: %s, stream.Close error: %v", d.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream DeepSeek model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", d.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ChatCompletionsStream DeepSeek model: %s finished", d.Model)
				} else {
					logger.Errorf(ctx, "ChatCompletionsStream DeepSeek model: %s, error: %v", d.Model, err)
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

			response, err := d.ConvChatCompletionsStreamResponse(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream DeepSeek ConvChatCompletionsStreamResponse error: %v", err)

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
		logger.Errorf(ctx, "ChatCompletionsStream DeepSeek model: %s, error: %v", d.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
