package openai

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

func (c *Client) Responses(ctx context.Context, data []byte) (res model.OpenAIResponsesRes, err error) {

	logger.Infof(ctx, "Responses OpenAI model: %s start", c.model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "Responses OpenAI model: %s totalTime: %d ms", c.model, res.TotalTime)
	}()

	if res.ResponseBytes, err = util.HttpPost(ctx, fmt.Sprintf("%s%s", c.baseURL, c.path), c.header, data, &res, c.proxyURL); err != nil {
		logger.Errorf(ctx, "Responses OpenAI model: %s, error: %v", c.model, err)
		return res, err
	}

	if res.Error != nil {
		logger.Errorf(ctx, "Responses OpenAI model: %s, responsesRes: %s", c.model, gjson.MustEncodeString(res))

		// todo
		//err = c.apiErrorHandler(&res)
		logger.Errorf(ctx, "Responses OpenAI model: %s, error: %v", c.model, err)

		return res, err
	}

	return res, nil
}

func (c *Client) ResponsesStream(ctx context.Context, data []byte) (responseChan chan *model.OpenAIResponsesRes, err error) {

	logger.Infof(ctx, "ResponsesStream OpenAI model: %s start", c.model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ResponsesStream OpenAI model: %s totalTime: %d ms", c.model, gtime.TimestampMilli()-now)
		}
	}()

	stream, err := util.SSEClient(ctx, fmt.Sprintf("%s%s", c.baseURL, c.path), c.header, data, c.proxyURL, c.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ResponsesStream OpenAI model: %s, error: %v", c.model, err)
		return responseChan, err
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.OpenAIResponsesRes)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ResponsesStream OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", c.model, duration-now, end-duration, end-now)

			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ResponsesStream OpenAI model: %s, stream.Close error: %v", c.model, err)
			}
		}()

		for {

			streamResponse, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {

				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "ResponsesStream OpenAI model: %s, error: %v", c.model, err)
				}

				end := gtime.TimestampMilli()
				responseChan <- &model.OpenAIResponsesRes{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Err:       err,
				}

				return
			}

			if errors.Is(err, io.EOF) {
				logger.Infof(ctx, "ResponsesStream OpenAI model: %s finished", c.model)

				end := gtime.TimestampMilli()
				responseChan <- &model.OpenAIResponsesRes{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Err:       io.EOF,
				}

				return
			}

			responsesRes := new(model.OpenAIResponsesRes)
			if err := gjson.Unmarshal(streamResponse, &responsesRes); err != nil {
				logger.Errorf(ctx, "ResponsesStream OpenAI model: %s, streamResponse: %s, error: %v", c.model, streamResponse, err)

				end := gtime.TimestampMilli()
				responseChan <- &model.OpenAIResponsesRes{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Err:       errors.New(fmt.Sprintf("streamResponse: %s, error: %v", streamResponse, err)),
				}

				return
			}

			if responsesRes.Error != nil {
				logger.Errorf(ctx, "ResponsesStream OpenAI model: %s, responsesRes: %s", c.model, gjson.MustEncodeString(responsesRes))

				// todo
				//err = c.apiErrorHandler(responsesRes.Error)
				logger.Errorf(ctx, "ResponsesStream OpenAI model: %s, error: %v", c.model, err)

				end := gtime.TimestampMilli()
				responseChan <- &model.OpenAIResponsesRes{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Err:       err,
				}

				return
			}

			response := &model.OpenAIResponsesRes{
				Error:         responsesRes.Error,
				ResponseBytes: streamResponse,
				ConnTime:      duration - now,
			}

			end := gtime.TimestampMilli()
			response.Duration = end - duration
			response.TotalTime = end - now

			responseChan <- response
		}
	}, nil); err != nil {
		logger.Errorf(ctx, "ResponsesStream OpenAI model: %s, error: %v", c.model, err)
		return responseChan, err
	}

	return responseChan, nil
}
