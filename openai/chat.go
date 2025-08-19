package openai

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (o *OpenAI) ChatCompletions(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions OpenAI model: %s start", o.model)

	request, err := o.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions OpenAI ConvChatCompletionsRequest error: %v", err)
		return response, err
	}

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions OpenAI model: %s totalTime: %d ms", o.model, response.TotalTime)
	}()

	bytes, err := util.HttpPost(ctx, o.baseURL+o.path, o.header, gjson.MustEncode(request), nil, o.proxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions OpenAI model: %s, error: %v", o.model, err)
		return response, err
	}

	if response, err = o.ConvChatCompletionsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions OpenAI ConvChatCompletionsResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "ChatCompletions OpenAI model: %s finished", o.model)

	return response, nil
}

func (o *OpenAI) ChatCompletionsStream(ctx context.Context, data []byte) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream OpenAI model: %s start", o.model)

	request, err := o.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream OpenAI ConvChatCompletionsRequest error: %v", err)
		return nil, err
	}

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream OpenAI model: %s totalTime: %d ms", o.model, gtime.TimestampMilli()-now)
		}
	}()

	if (o.isSupportStream != nil && !*o.isSupportStream) || (gstr.HasPrefix(o.model, "o") && o.isAzure) {
		return o.ChatCompletionStreamToNonStream(ctx, data)
	}

	stream, err := util.SSEClient(ctx, o.baseURL+o.path, o.header, gjson.MustEncode(request), o.proxyURL, nil)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream OpenAI model: %s, error: %v", o.model, err)
		return responseChan, err
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream OpenAI model: %s, stream.Close error: %v", o.model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", o.model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
					logger.Errorf(ctx, "ChatCompletionsStream OpenAI model: %s, error: %v", o.model, err)
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

			// todo
			if response.Usage != nil {
				fmt.Println(stream.ReqTime, response.ResTime, end, end-gconv.Int64(response.ResTime), end-gconv.Int64(stream.ReqTime)-response.ResTotalTime, "end")
			}

			response.ResponseBytes = responseBytes
			response.ConnTime = duration - now
			response.Duration = end - duration
			response.TotalTime = end - now

			responseChan <- &response
		}

	}, nil); err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream OpenAI model: %s, error: %v", o.model, err)
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
			logger.Infof(ctx, "ChatCompletionStreamToNonStream OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", o.model, duration-now, end-duration, end-now)
		}()

		request.Stream = false

		response, err := o.ChatCompletions(ctx, gjson.MustEncode(request))
		if err != nil {

			if !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
				logger.Errorf(ctx, "ChatCompletionStreamToNonStream OpenAI model: %s, error: %v", o.model, err)
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
		logger.Errorf(ctx, "ChatCompletionStreamToNonStream OpenAI model: %s, error: %v", o.model, err)
		return responseChan, err
	}

	return responseChan, nil
}
