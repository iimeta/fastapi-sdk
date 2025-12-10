package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/errors"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (g *Google) ChatCompletionsOfficial(ctx context.Context, data []byte) (response any, err error) {

	logger.Infof(ctx, "ChatCompletionsOfficial Google model: %s start", g.Model)

	var (
		now = gtime.TimestampMilli()
		res = &model.GoogleChatCompletionRes{}
	)

	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletionsOfficial Google model: %s totalTime: %d ms", g.Model, res.TotalTime)
	}()

	if g.Path == "" {
		g.Path = "/models/" + g.Model
	}

	if res.ResponseBytes, err = util.HttpPost(ctx, fmt.Sprintf("%s:generateContent?key=%s", g.BaseUrl+g.Path, g.Key), g.header, data, &res, g.Timeout, g.ProxyUrl, g.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "ChatCompletionsOfficial Google model: %s, error: %v", g.Model, err)
		return res, err
	}

	if res.Error.Code != 0 || res.Candidates[0].FinishReason != "STOP" {
		logger.Errorf(ctx, "ChatCompletionsOfficial Google model: %s, chatCompletionRes: %s", g.Model, gjson.MustEncodeString(res))

		err = g.apiErrorHandler(res)
		logger.Errorf(ctx, "ChatCompletionsOfficial Google model: %s, error: %v", g.Model, err)

		return res, err
	}

	return res, nil
}

func (g *Google) ChatCompletionsStreamOfficial(ctx context.Context, data []byte) (responseChan chan any, err error) {

	logger.Infof(ctx, "ChatCompletionsStreamOfficial Google model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStreamOfficial Google model: %s totalTime: %d ms", g.Model, gtime.TimestampMilli()-now)
		}
	}()

	if g.Path == "" {
		g.Path = "/models/" + g.Model
	}

	stream, err := util.SSEClient(ctx, fmt.Sprintf("%s:streamGenerateContent?alt=sse&key=%s", g.BaseUrl+g.Path, g.Key), nil, data, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStreamOfficial Google model: %s, error: %v", g.Model, err)
		return responseChan, err
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan any)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStreamOfficial Google model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", g.Model, duration-now, end-duration, end-now)

			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStreamOfficial Google model: %s, stream.Close error: %v", g.Model, err)
			}
		}()

		var (
			usageMetadata *model.UsageMetadata
		)

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ChatCompletionsStreamOfficial Google model: %s finished", g.Model)
				} else {
					logger.Errorf(ctx, "ChatCompletionsStreamOfficial Google model: %s, error: %v", g.Model, err)
				}

				end := gtime.TimestampMilli()
				responseChan <- &model.GoogleChatCompletionRes{
					UsageMetadata: usageMetadata,
					ConnTime:      duration - now,
					Duration:      end - duration,
					TotalTime:     end - now,
					Err:           err,
				}

				return
			}

			chatCompletionRes := model.GoogleChatCompletionRes{}
			if err := json.Unmarshal(responseBytes, &chatCompletionRes); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStreamOfficial Google model: %s, response: %s, error: %v", g.Model, responseBytes, err)

				end := gtime.TimestampMilli()
				responseChan <- &model.GoogleChatCompletionRes{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Err:       errors.New(fmt.Sprintf("response: %s, error: %v", responseBytes, err)),
				}

				return
			}

			if chatCompletionRes.Error.Code != 0 {
				logger.Errorf(ctx, "ChatCompletionsStreamOfficial Google model: %s, chatCompletionRes: %s", g.Model, gjson.MustEncodeString(chatCompletionRes))

				err = g.apiErrorHandler(&chatCompletionRes)
				logger.Errorf(ctx, "ChatCompletionsStreamOfficial Google model: %s, error: %v", g.Model, err)

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
				ResponseBytes: responseBytes,
				ConnTime:      duration - now,
			}

			end := gtime.TimestampMilli()
			response.Duration = end - duration
			response.TotalTime = end - now

			responseChan <- response
		}
	}, nil); err != nil {
		logger.Errorf(ctx, "ChatCompletionsStreamOfficial Google model: %s, error: %v", g.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
