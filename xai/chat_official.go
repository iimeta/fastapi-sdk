package xai

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

func (x *XAI) ChatCompletionsOfficial(ctx context.Context, data []byte) (response any, err error) {

	logger.Infof(ctx, "ChatCompletionsOfficial xAI model: %s start", x.Model)

	var (
		now = gtime.TimestampMilli()
		res = &model.ChatCompletionResponse{}
	)

	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletionsOfficial xAI model: %s totalTime: %d ms", x.Model, res.TotalTime)
	}()

	if x.Path == "" {
		x.Path = "/chat/completions"
	}

	responseBytes, responseHeader, postErr := util.HttpPost(ctx, x.BaseUrl+x.Path, x.header, data, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if postErr != nil {
		err = postErr
		logger.Errorf(ctx, "ChatCompletionsOfficial xAI model: %s, error: %v", x.Model, err)
		return res, err
	}

	converted, convErr := x.ConvChatCompletionsResponse(ctx, responseBytes)
	if convErr != nil {
		err = convErr
		logger.Errorf(ctx, "ChatCompletionsOfficial xAI ConvChatCompletionsResponse error: %v", err)
		return res, err
	}
	converted.ResponseHeaders = responseHeader
	*res = converted

	logger.Infof(ctx, "ChatCompletionsOfficial xAI model: %s finished", x.Model)

	return res, nil
}

func (x *XAI) ChatCompletionsStreamOfficial(ctx context.Context, data []byte) (responseChan chan any, err error) {

	logger.Infof(ctx, "ChatCompletionsStreamOfficial xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStreamOfficial xAI model: %s totalTime: %d ms", x.Model, gtime.TimestampMilli()-now)
		}
	}()

	if x.Path == "" {
		x.Path = "/chat/completions"
	}

	stream, err := util.SSEClient(ctx, x.BaseUrl+x.Path, x.header, data, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStreamOfficial xAI model: %s, error: %v", x.Model, err)
		return responseChan, err
	}

	streamResponseHeaders := stream.Response.Header

	duration := gtime.TimestampMilli()

	responseChan = make(chan any)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStreamOfficial xAI model: %s, stream.Close error: %v", x.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStreamOfficial xAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", x.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ChatCompletionsStreamOfficial xAI model: %s finished", x.Model)
				} else {
					logger.Errorf(ctx, "ChatCompletionsStreamOfficial xAI model: %s, error: %v", x.Model, err)
				}

				end := gtime.TimestampMilli()
				responseChan <- &model.ChatCompletionResponse{
					ResponseBytes:   responseBytes,
					ResponseHeaders: streamResponseHeaders,
					ConnTime:        duration - now,
					Duration:        end - duration,
					TotalTime:       end - now,
					Error:           err,
				}

				return
			}

			response, convErr := x.ConvChatCompletionsStreamResponse(ctx, responseBytes)
			if convErr != nil {
				logger.Errorf(ctx, "ChatCompletionsStreamOfficial xAI model: %s, response: %s, error: %v", x.Model, responseBytes, convErr)

				end := gtime.TimestampMilli()
				responseChan <- &model.ChatCompletionResponse{
					ResponseBytes:   responseBytes,
					ResponseHeaders: streamResponseHeaders,
					ConnTime:        duration - now,
					Duration:        end - duration,
					TotalTime:       end - now,
					Error:           errors.New(fmt.Sprintf("response: %s, error: %v", responseBytes, convErr)),
				}

				return
			}

			end := gtime.TimestampMilli()

			response.ResponseBytes = responseBytes
			response.ResponseHeaders = streamResponseHeaders
			response.ConnTime = duration - now
			response.Duration = end - duration
			response.TotalTime = end - now

			responseChan <- &response
		}

	}, nil); err != nil {
		logger.Errorf(ctx, "ChatCompletionsStreamOfficial xAI model: %s, error: %v", x.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
