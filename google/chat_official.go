package google

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
	"io"
)

func (c *Client) ChatCompletionOfficial(ctx context.Context, data []byte) (res model.GoogleChatCompletionRes, err error) {

	logger.Infof(ctx, "ChatCompletionOfficial Google model: %s start", c.model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletionOfficial Google model: %s totalTime: %d ms", c.model, res.TotalTime)
	}()

	if res.ResponseBytes, err = util.HttpPost(ctx, fmt.Sprintf("%s:generateContent?key=%s", c.baseURL+c.path, c.key), nil, data, &res, c.proxyURL); err != nil {
		logger.Errorf(ctx, "ChatCompletionOfficial Google model: %s, error: %v", c.model, err)
		return res, err
	}

	if res.Error.Code != 0 || res.Candidates[0].FinishReason != "STOP" {
		logger.Errorf(ctx, "ChatCompletionOfficial Google model: %s, chatCompletionRes: %s", c.model, gjson.MustEncodeString(res))

		err = c.apiErrorHandler(&res)
		logger.Errorf(ctx, "ChatCompletionOfficial Google model: %s, error: %v", c.model, err)

		return res, err
	}

	return res, nil
}

func (c *Client) ChatCompletionStreamOfficial(ctx context.Context, data []byte) (responseChan chan *model.GoogleChatCompletionRes, err error) {

	logger.Infof(ctx, "ChatCompletionStreamOfficial Google model: %s start", c.model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionStreamOfficial Google model: %s totalTime: %d ms", c.model, gtime.TimestampMilli()-now)
		}
	}()

	stream, err := util.SSEClient(ctx, fmt.Sprintf("%s:streamGenerateContent?alt=sse&key=%s", c.baseURL+c.path, c.key), nil, data, c.proxyURL, c.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStreamOfficial Google model: %s, error: %v", c.model, err)
		return responseChan, err
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.GoogleChatCompletionRes)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionStreamOfficial Google model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", c.model, duration-now, end-duration, end-now)

			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionStreamOfficial Google model: %s, stream.Close error: %v", c.model, err)
			}
		}()

		var (
			usageMetadata *model.UsageMetadata
		)

		for {

			streamResponse, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {

				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "ChatCompletionStreamOfficial Google model: %s, error: %v", c.model, err)
				}

				end := gtime.TimestampMilli()
				responseChan <- &model.GoogleChatCompletionRes{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Err:       err,
				}

				return
			}

			if errors.Is(err, io.EOF) {
				logger.Infof(ctx, "ChatCompletionStreamOfficial Google model: %s finished", c.model)

				end := gtime.TimestampMilli()
				responseChan <- &model.GoogleChatCompletionRes{
					UsageMetadata: usageMetadata,
					ConnTime:      duration - now,
					Duration:      end - duration,
					TotalTime:     end - now,
					Err:           io.EOF,
				}

				return
			}

			chatCompletionRes := new(model.GoogleChatCompletionRes)
			if err := gjson.Unmarshal(streamResponse, &chatCompletionRes); err != nil {
				logger.Errorf(ctx, "ChatCompletionStreamOfficial Google model: %s, streamResponse: %s, error: %v", c.model, streamResponse, err)

				end := gtime.TimestampMilli()
				responseChan <- &model.GoogleChatCompletionRes{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Err:       errors.New(fmt.Sprintf("streamResponse: %s, error: %v", streamResponse, err)),
				}

				return
			}

			if chatCompletionRes.Error.Code != 0 {
				logger.Errorf(ctx, "ChatCompletionStreamOfficial Google model: %s, chatCompletionRes: %s", c.model, gjson.MustEncodeString(chatCompletionRes))

				err = c.apiErrorHandler(chatCompletionRes)
				logger.Errorf(ctx, "ChatCompletionStreamOfficial Google model: %s, error: %v", c.model, err)

				end := gtime.TimestampMilli()
				responseChan <- &model.GoogleChatCompletionRes{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Err:       err,
				}

				return
			}

			if chatCompletionRes.UsageMetadata != nil {
				usageMetadata = chatCompletionRes.UsageMetadata
			}

			response := &model.GoogleChatCompletionRes{
				Candidates:    chatCompletionRes.Candidates,
				UsageMetadata: chatCompletionRes.UsageMetadata,
				Error:         chatCompletionRes.Error,
				ResponseBytes: streamResponse,
				ConnTime:      duration - now,
			}

			end := gtime.TimestampMilli()
			response.Duration = end - duration
			response.TotalTime = end - now

			responseChan <- response
		}
	}, nil); err != nil {
		logger.Errorf(ctx, "ChatCompletionStreamOfficial Google model: %s, error: %v", c.model, err)
		return responseChan, err
	}

	return responseChan, nil
}
