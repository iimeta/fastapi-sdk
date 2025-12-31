package volcengine

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

func (v *VolcEngine) ChatCompletions(ctx context.Context, data any) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions VolcEngine model: %s start", v.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions VolcEngine model: %s totalTime: %d ms", v.Model, response.TotalTime)
	}()

	if !v.IsOfficialFormatRequest {
		if data, err = v.ConvChatCompletionsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "ChatCompletions VolcEngine ConvChatCompletionsRequest error: %v", err)
			return response, err
		}
	}

	if v.Path == "" {
		v.Path = "/chat/completions"
	}

	bytes, err := util.HttpPost(ctx, v.BaseUrl+v.Path, v.header, data, nil, v.Timeout, v.ProxyUrl, v.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions VolcEngine model: %s, error: %v", v.Model, err)
		return response, err
	}

	if response, err = v.ConvChatCompletionsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions VolcEngine ConvChatCompletionsResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "ChatCompletions VolcEngine model: %s finished", v.Model)

	return response, nil
}

func (v *VolcEngine) ChatCompletionsStream(ctx context.Context, data any) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream VolcEngine model: %s start", v.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream VolcEngine model: %s totalTime: %d ms", v.Model, gtime.TimestampMilli()-now)
		}
	}()

	if !v.IsOfficialFormatRequest {
		if data, err = v.ConvChatCompletionsRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream VolcEngine ConvChatCompletionsRequest error: %v", err)
			return nil, err
		}
	}

	if v.Path == "" {
		v.Path = "/chat/completions"
	}

	stream, err := util.SSEClient(ctx, v.BaseUrl+v.Path, v.header, data, v.Timeout, v.ProxyUrl, v.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream VolcEngine model: %s, error: %v", v.Model, err)
		return responseChan, err
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream VolcEngine model: %s, stream.Close error: %v", v.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream VolcEngine model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", v.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ChatCompletionsStream VolcEngine model: %s finished", v.Model)
				} else {
					logger.Errorf(ctx, "ChatCompletionsStream VolcEngine model: %s, error: %v", v.Model, err)
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
		logger.Errorf(ctx, "ChatCompletionsStream VolcEngine model: %s, error: %v", v.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
