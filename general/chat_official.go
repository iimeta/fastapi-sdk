package general

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (g *General) ChatCompletionsOfficial(ctx context.Context, data []byte) (response any, err error) {

	logger.Infof(ctx, "ChatCompletionsOfficial General model: %s start", g.Model)

	var (
		now = gtime.TimestampMilli()
		res = &model.ChatCompletionResponse{}
	)

	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletionsOfficial General model: %s totalTime: %d ms", g.Model, res.TotalTime)
	}()

	if res.ResponseBytes, err = util.HttpPost(ctx, g.BaseUrl+g.Path, g.header, data, &res, g.Timeout, g.ProxyUrl, g.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "ChatCompletionsOfficial General model: %s, error: %v", g.Model, err)
		return res, err
	}

	logger.Infof(ctx, "ChatCompletionsOfficial General model: %s finished", g.Model)

	return res, nil
}

func (g *General) ChatCompletionsStreamOfficial(ctx context.Context, data []byte) (responseChan chan any, err error) {

	logger.Infof(ctx, "ChatCompletionsStreamOfficial General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStreamOfficial General model: %s totalTime: %d ms", g.Model, gtime.TimestampMilli()-now)
		}
	}()

	stream, err := util.SSEClient(ctx, g.BaseUrl+g.Path, g.header, data, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStreamOfficial General model: %s, error: %v", g.Model, err)
		return responseChan, err
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan any)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStreamOfficial General model: %s, stream.Close error: %v", g.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStreamOfficial General model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", g.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ChatCompletionsStreamOfficial General model: %s finished", g.Model)
				} else {
					logger.Errorf(ctx, "ChatCompletionsStreamOfficial General model: %s, error: %v", g.Model, err)
				}

				end := gtime.TimestampMilli()
				responseChan <- &model.ChatCompletionResponse{
					ResponseBytes: responseBytes,
					ConnTime:      duration - now,
					Duration:      end - duration,
					TotalTime:     end - now,
					Error:         err,
				}

				return
			}

			response := model.ChatCompletionResponse{}
			if err := json.Unmarshal(responseBytes, &response); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStreamOfficial General model: %s, response: %s, error: %v", g.Model, responseBytes, err)

				end := gtime.TimestampMilli()
				responseChan <- &model.ChatCompletionResponse{
					ResponseBytes: responseBytes,
					ConnTime:      duration - now,
					Duration:      end - duration,
					TotalTime:     end - now,
					Error:         errors.New(fmt.Sprintf("response: %s, error: %v", responseBytes, err)),
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
		logger.Errorf(ctx, "ChatCompletionsStreamOfficial General model: %s, error: %v", g.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
