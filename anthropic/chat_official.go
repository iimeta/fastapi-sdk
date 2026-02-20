package anthropic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/v2/anthropic/aws"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (a *Anthropic) ChatCompletionsOfficial(ctx context.Context, data []byte) (response any, err error) {

	logger.Infof(ctx, "ChatCompletionsOfficial Anthropic model: %s start", a.Model)

	var (
		now = gtime.TimestampMilli()
		res = &model.AnthropicChatCompletionRes{}
	)

	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletionsOfficial Anthropic model: %s totalTime: %d ms", a.Model, res.TotalTime)
	}()

	if a.isGcp || a.isAws {

		request := make(map[string]any)
		if err = json.Unmarshal(data, &request); err != nil {
			logger.Errorf(ctx, "ChatCompletionsOfficial Anthropic model: %s, data: %s, json.Unmarshal error: %v", a.Model, data, err)
			return res, err
		}

		if a.isGcp {
			delete(request, "model")
			data = gjson.MustEncode(request)
		}

		if a.isAws {

			request["anthropic_version"] = "bedrock-2023-05-31"
			delete(request, "metadata")
			delete(request, "model")
			delete(request, "stream")

			data = gjson.MustEncode(request)

			a.header = aws.SignHeader(a.Path, a.region, a.accessKey, a.secretKey, data)
		}
	}

	if a.Path == "" {
		a.Path = "/messages"
	}

	if res.ResponseBytes, err = util.HttpPost(ctx, a.BaseUrl+a.Path, a.header, data, &res, a.Timeout, a.ProxyUrl, a.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "ChatCompletionsOfficial Anthropic model: %s, error: %v", a.Model, err)
		return res, err
	}

	if res.Error != nil && res.Error.Type != "" {
		logger.Errorf(ctx, "ChatCompletionsOfficial Anthropic model: %s, chatCompletionRes: %s", a.Model, gjson.MustEncodeString(res))

		err = a.apiErrorHandler(res)
		logger.Errorf(ctx, "ChatCompletionsOfficial Anthropic model: %s, error: %v", a.Model, err)

		return res, err
	}

	return res, nil
}

func (a *Anthropic) ChatCompletionsStreamOfficial(ctx context.Context, data []byte) (responseChan chan any, err error) {

	logger.Infof(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s start", a.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s totalTime: %d ms", a.Model, gtime.TimestampMilli()-now)
		}
	}()

	request := make(map[string]any)
	if err = json.Unmarshal(data, &request); err != nil {
		logger.Errorf(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s, data: %s, json.Unmarshal error: %v", a.Model, data, err)
		return responseChan, err
	}

	if a.isGcp {
		delete(request, "model")
	}

	if a.isAws {

		request["anthropic_version"] = "bedrock-2023-05-31"
		delete(request, "stream")
		delete(request, "model")

		data = gjson.MustEncode(request)

		a.header = aws.SignHeader(a.Path, a.region, a.accessKey, a.secretKey, data)

		stream, err := util.SSEClient(ctx, a.BaseUrl+a.Path, a.header, data, a.Timeout, a.ProxyUrl, a.requestErrorHandler)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s, error: %v", a.Model, err)
			return responseChan, err
		}

		payloadBuf := make([]byte, 10*1024)

		duration := gtime.TimestampMilli()

		responseChan = make(chan any)

		if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

			defer func() {
				if err := stream.Close(); err != nil {
					logger.Errorf(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s, stream.Close error: %v", a.Model, err)
				}

				end := gtime.TimestampMilli()
				logger.Infof(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", a.Model, duration-now, end-duration, end-now)
			}()

			for {

				payloadBuf = payloadBuf[0:0]

				decodedMessage, err := aws.DecodeMessage(stream.Response.Body, payloadBuf)
				if err != nil {

					if errors.Is(err, io.EOF) {
						logger.Infof(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s finished", a.Model)
					} else {
						logger.Errorf(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s, error: %v", a.Model, err)
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

				payload := make(map[string]any)
				if err := json.Unmarshal(decodedMessage.Payload, &payload); err != nil {
					logger.Errorf(ctx, "ChatCompletionsStreamOfficial Anthropic json.Unmarshal(decodedMessage.Payload, &payload), payload: %s, error: %v", decodedMessage.Payload, err)

					end := gtime.TimestampMilli()
					responseChan <- &model.ChatCompletionResponse{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Error:     err,
					}
				}

				bytes, err := base64.StdEncoding.DecodeString(gconv.String(payload["bytes"]))
				if err != nil {
					logger.Errorf(ctx, `ChatCompletionsStreamOfficial Anthropic base64.StdEncoding.DecodeString(gconv.String(payload["bytes"])), bytes: %s, error: %v`, payload["bytes"], err)

					end := gtime.TimestampMilli()
					responseChan <- &model.ChatCompletionResponse{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Error:     err,
					}
				}

				chatCompletionRes := model.AnthropicChatCompletionRes{}
				if err := json.Unmarshal(bytes, &chatCompletionRes); err != nil {
					logger.Errorf(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s, bytes: %s, error: %v", a.Model, bytes, err)

					end := gtime.TimestampMilli()
					responseChan <- &model.AnthropicChatCompletionRes{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Err:       errors.New(fmt.Sprintf("bytes: %s, error: %v", bytes, err)),
					}

					return
				}

				if chatCompletionRes.Error != nil && chatCompletionRes.Error.Type != "" {
					logger.Errorf(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s, chatCompletionRes: %s", a.Model, gjson.MustEncodeString(chatCompletionRes))

					err = a.apiErrorHandler(&chatCompletionRes)
					logger.Errorf(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s, error: %v", a.Model, err)

					end := gtime.TimestampMilli()
					responseChan <- &model.AnthropicChatCompletionRes{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Err:       err,
					}

					return
				}

				response := &model.AnthropicChatCompletionRes{
					Id:            chatCompletionRes.Id,
					Type:          chatCompletionRes.Type,
					Role:          chatCompletionRes.Role,
					Content:       chatCompletionRes.Content,
					Model:         chatCompletionRes.Model,
					StopReason:    chatCompletionRes.StopReason,
					StopSequence:  chatCompletionRes.StopSequence,
					Message:       chatCompletionRes.Message,
					Index:         chatCompletionRes.Index,
					Delta:         chatCompletionRes.Delta,
					Usage:         chatCompletionRes.Usage,
					Error:         chatCompletionRes.Error,
					ResponseBytes: bytes,
					ConnTime:      duration - now,
				}

				if chatCompletionRes.Delta.StopReason != "" {
					logger.Infof(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s finished", a.Model)

					end := gtime.TimestampMilli()
					response.Duration = end - duration
					response.TotalTime = end - now
					responseChan <- response

					responseChan <- &model.AnthropicChatCompletionRes{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Err:       io.EOF,
					}

					return
				}

				end := gtime.TimestampMilli()
				response.Duration = end - duration
				response.TotalTime = end - now

				responseChan <- response
			}
		}, nil); err != nil {
			logger.Errorf(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s, error: %v", a.Model, err)
			return responseChan, err
		}

	} else {

		if a.Path == "" {
			a.Path = "/messages"
		}

		stream, err := util.SSEClient(ctx, a.BaseUrl+a.Path, a.header, gjson.MustEncode(request), a.Timeout, a.ProxyUrl, a.requestErrorHandler)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s, error: %v", a.Model, err)
			return responseChan, err
		}

		duration := gtime.TimestampMilli()

		responseChan = make(chan any)

		if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

			defer func() {
				if err := stream.Close(); err != nil {
					logger.Errorf(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s, stream.Close error: %v", a.Model, err)
				}

				end := gtime.TimestampMilli()
				logger.Infof(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", a.Model, duration-now, end-duration, end-now)
			}()

			for {

				responseBytes, err := stream.Recv()
				if err != nil {

					if errors.Is(err, io.EOF) {
						logger.Infof(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s finished", a.Model)
					} else {
						logger.Errorf(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s, error: %v", a.Model, err)
					}

					end := gtime.TimestampMilli()
					responseChan <- &model.AnthropicChatCompletionRes{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Err:       err,
					}

					return
				}

				chatCompletionRes := model.AnthropicChatCompletionRes{}
				if err := json.Unmarshal(responseBytes, &chatCompletionRes); err != nil {
					logger.Errorf(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s, response: %s, error: %v", a.Model, responseBytes, err)

					end := gtime.TimestampMilli()
					responseChan <- &model.AnthropicChatCompletionRes{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Err:       errors.New(fmt.Sprintf("response: %s, error: %v", responseBytes, err)),
					}

					return
				}

				if chatCompletionRes.Error != nil && chatCompletionRes.Error.Type != "" {
					logger.Errorf(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s, chatCompletionRes: %s", a.Model, gjson.MustEncodeString(chatCompletionRes))

					err = a.apiErrorHandler(&chatCompletionRes)
					logger.Errorf(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s, error: %v", a.Model, err)

					end := gtime.TimestampMilli()
					responseChan <- &model.AnthropicChatCompletionRes{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Err:       err,
					}

					return
				}

				response := &model.AnthropicChatCompletionRes{
					Id:            chatCompletionRes.Id,
					Type:          chatCompletionRes.Type,
					Role:          chatCompletionRes.Role,
					Content:       chatCompletionRes.Content,
					Model:         chatCompletionRes.Model,
					StopReason:    chatCompletionRes.StopReason,
					StopSequence:  chatCompletionRes.StopSequence,
					Message:       chatCompletionRes.Message,
					Index:         chatCompletionRes.Index,
					Delta:         chatCompletionRes.Delta,
					Usage:         chatCompletionRes.Usage,
					Error:         chatCompletionRes.Error,
					ResponseBytes: responseBytes,
					ConnTime:      duration - now,
				}

				if chatCompletionRes.Delta.StopReason != "" {
					logger.Infof(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s finished", a.Model)

					end := gtime.TimestampMilli()
					response.Duration = end - duration
					response.TotalTime = end - now
					responseChan <- response

					responseChan <- &model.AnthropicChatCompletionRes{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Err:       io.EOF,
					}

					return
				}

				end := gtime.TimestampMilli()
				response.Duration = end - duration
				response.TotalTime = end - now

				responseChan <- response
			}

		}, nil); err != nil {
			logger.Errorf(ctx, "ChatCompletionsStreamOfficial Anthropic model: %s, error: %v", a.Model, err)
			return responseChan, err
		}
	}

	return responseChan, nil
}
