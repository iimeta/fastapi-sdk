package openai

import (
	"context"
	"io"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/errors"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (o *OpenAI) ChatCompletions(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions OpenAI model: %s totalTime: %d ms", o.Model, response.TotalTime)
	}()

	request, err := o.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions OpenAI ConvChatCompletionsRequest error: %v", err)
		return response, err
	}

	bytes, err := util.HttpPost(ctx, o.BaseUrl+o.Path, o.header, request, nil, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions OpenAI model: %s, error: %v", o.Model, err)
		return response, err
	}

	if response, err = o.ConvChatCompletionsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions OpenAI ConvChatCompletionsResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "ChatCompletions OpenAI model: %s finished", o.Model)

	return response, nil
}

func (o *OpenAI) ChatCompletionsStream(ctx context.Context, data []byte) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream OpenAI model: %s totalTime: %d ms", o.Model, gtime.TimestampMilli()-now)
		}
	}()

	request, err := o.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream OpenAI ConvChatCompletionsRequest error: %v", err)
		return nil, err
	}

	if (o.IsSupportStream != nil && !*o.IsSupportStream) || (gstr.HasPrefix(o.Model, "o") && o.isAzure) {
		return o.ChatCompletionStreamToNonStream(ctx, data)
	}

	stream, err := util.SSEClient(ctx, o.BaseUrl+o.Path, o.header, gjson.MustEncode(request), o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream OpenAI model: %s, error: %v", o.Model, err)
		return responseChan, err
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream OpenAI model: %s, stream.Close error: %v", o.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", o.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ChatCompletionsStream OpenAI model: %s finished", o.Model)
				} else {
					logger.Errorf(ctx, "ChatCompletionsStream OpenAI model: %s, error: %v", o.Model, err)
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

			response, err := o.ConvChatCompletionsStreamResponse(ctx, responseBytes)
			if err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream OpenAI ConvChatCompletionsStreamResponse error: %v", err)

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
		logger.Errorf(ctx, "ChatCompletionsStream OpenAI model: %s, error: %v", o.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (o *OpenAI) ChatCompletionStreamToNonStream(ctx context.Context, data []byte) (responseChan chan *model.ChatCompletionResponse, err error) {

	request, err := o.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStreamToNonStream OpenAI ConvChatCompletionsRequest error: %v", err)
		return nil, err
	}

	responseChan = make(chan *model.ChatCompletionResponse)

	now := gtime.TimestampMilli()
	duration := now

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionStreamToNonStream OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", o.Model, duration-now, end-duration, end-now)
		}()

		request.Stream = false

		response, err := o.ChatCompletions(ctx, gjson.MustEncode(request))
		if err != nil {

			if errors.Is(err, io.EOF) {
				logger.Infof(ctx, "ChatCompletionStreamToNonStream OpenAI model: %s finished", o.Model)
			} else {
				logger.Errorf(ctx, "ChatCompletionStreamToNonStream OpenAI model: %s, error: %v", o.Model, err)
			}

			end := gtime.TimestampMilli()
			responseChan <- &model.ChatCompletionResponse{
				ConnTime:  gtime.TimestampMilli() - now,
				Duration:  end - gtime.TimestampMilli(),
				TotalTime: end - now,
				Error:     err,
			}

			return
		}

		duration = gtime.TimestampMilli()
		response.ConnTime = duration - now

		end := gtime.TimestampMilli()
		response.Duration = end - duration
		response.TotalTime = end - now

		responseChan <- &response

		end = gtime.TimestampMilli()
		responseChan <- &model.ChatCompletionResponse{
			Duration:  end - duration,
			TotalTime: end - now,
			Error:     io.EOF,
		}

	}, nil); err != nil {
		logger.Errorf(ctx, "ChatCompletionStreamToNonStream OpenAI model: %s, error: %v", o.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
