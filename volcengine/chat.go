package volcengine

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

func (v *VolcEngine) ChatCompletions(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions VolcEngine model: %s start", v.model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions VolcEngine model: %s totalTime: %d ms", v.model, response.TotalTime)
	}()

	request, err := v.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions VolcEngine ConvChatCompletionsRequest error: %v", err)
		return response, err
	}

	bytes, err := util.HttpPost(ctx, v.baseURL+v.path, v.header, request, nil, v.proxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions VolcEngine model: %s, error: %v", v.model, err)
		return response, v.apiErrorHandler(err)
	}

	if response, err = v.ConvChatCompletionsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions VolcEngine ConvChatCompletionsResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "ChatCompletions VolcEngine model: %s finished", v.model)

	return response, nil
}

func (v *VolcEngine) ChatCompletionsStream(ctx context.Context, data []byte) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream VolcEngine model: %s start", v.model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream VolcEngine model: %s totalTime: %d ms", v.model, gtime.TimestampMilli()-now)
		}
	}()

	request, err := v.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream VolcEngine ConvChatCompletionsRequest error: %v", err)
		return nil, err
	}

	stream, err := util.SSEClient(ctx, v.baseURL+v.path, v.header, gjson.MustEncode(request), v.proxyURL, nil)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream VolcEngine model: %s, error: %v", v.model, err)
		return responseChan, v.apiErrorHandler(err)
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream VolcEngine model: %s, stream.Close error: %v", v.model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream VolcEngine model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", v.model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
					logger.Errorf(ctx, "ChatCompletionsStream VolcEngine model: %s, error: %v", v.model, err)
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

			response, err := v.ConvChatCompletionsStreamResponse(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream VolcEngine ConvChatCompletionsStreamResponse error: %v", err)

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
		logger.Errorf(ctx, "ChatCompletionsStream VolcEngine model: %s, error: %v", v.model, err)
		return responseChan, err
	}

	return responseChan, nil
}
