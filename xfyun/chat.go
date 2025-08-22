package xfyun

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gorilla/websocket"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (x *Xfyun) ChatCompletions(ctx context.Context, data []byte) (res model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions Xfyun model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions Xfyun model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", x.Model, res.ConnTime, res.Duration, res.TotalTime)
	}()

	request, err := x.ConvChatCompletionsRequestOfficial(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions Xfyun ConvChatCompletionsRequestOfficial error: %v", err)
		return res, err
	}

	conn, err := util.WebSocketClient(ctx, x.getWebSocketUrl(ctx), nil, websocket.TextMessage, request, x.ProxyUrl)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions Xfyun model: %s, error: %v", x.Model, err)
		return res, err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Errorf(ctx, "ChatCompletions Xfyun model: %s, conn.Close error: %v", x.Model, err)
		}
	}()

	var (
		duration          = gtime.TimestampMilli()
		responseContent   = ""
		chatCompletionRes = model.XfyunChatCompletionRes{}
	)

	for {

		_, message, err := conn.ReadMessage(ctx)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletions Xfyun model: %s, error: %v", x.Model, err)
			return res, err
		}

		if err = json.Unmarshal(message, &chatCompletionRes); err != nil {
			logger.Errorf(ctx, "ChatCompletions Xfyun model: %s, message: %s, error: %v", x.Model, message, err)
			return res, errors.New(fmt.Sprintf("message: %s, error: %v", message, err))
		}

		if chatCompletionRes.Header.Code != 0 {
			logger.Errorf(ctx, "ChatCompletions Xfyun model: %s, chatCompletionRes: %s", x.Model, gjson.MustEncodeString(chatCompletionRes))

			err = x.apiErrorHandler(&chatCompletionRes)
			logger.Errorf(ctx, "ChatCompletions Xfyun model: %s, error: %v", x.Model, err)

			return res, err
		}

		responseContent += chatCompletionRes.Payload.Choices.Text[0].Content

		if chatCompletionRes.Header.Status == 2 {
			break
		}
	}

	res = model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Header.Sid,
		Object:  consts.COMPLETION_OBJECT,
		Created: gtime.Timestamp(),
		Model:   x.Model,
		Choices: []model.ChatCompletionChoice{{
			Index: chatCompletionRes.Payload.Choices.Seq,
			Message: &model.ChatCompletionMessage{
				Role:         chatCompletionRes.Payload.Choices.Text[0].Role,
				Content:      responseContent,
				FunctionCall: chatCompletionRes.Payload.Choices.Text[0].FunctionCall,
			},
		}},
		Usage: &model.Usage{
			PromptTokens:     chatCompletionRes.Payload.Usage.Text.PromptTokens,
			CompletionTokens: chatCompletionRes.Payload.Usage.Text.CompletionTokens,
			TotalTokens:      chatCompletionRes.Payload.Usage.Text.TotalTokens,
		},
		ConnTime: duration - now,
		Duration: gtime.TimestampMilli() - duration,
	}

	return res, nil
}

func (x *Xfyun) ChatCompletionsStream(ctx context.Context, data []byte) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream Xfyun model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream Xfyun model: %s totalTime: %d ms", x.Model, gtime.TimestampMilli()-now)
		}
	}()

	request, err := x.ConvChatCompletionsRequestOfficial(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream Xfyun ConvChatCompletionsRequestOfficial error: %v", err)
		return nil, err
	}

	conn, err := util.WebSocketClient(ctx, x.getWebSocketUrl(ctx), nil, websocket.TextMessage, request, x.ProxyUrl)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream Xfyun model: %s, error: %v", x.Model, err)
		return responseChan, err
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := conn.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream Xfyun model: %s, conn.Close error: %v", x.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream Xfyun model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", x.Model, duration-now, end-duration, end-now)
		}()

		var created = gtime.Timestamp()

		for {

			_, message, err := conn.ReadMessage(ctx)
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ChatCompletionsStream Xfyun model: %s finished", x.Model)
				} else {
					logger.Errorf(ctx, "ChatCompletionsStream Xfyun model: %s, error: %v", x.Model, err)
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

			chatCompletionRes := model.XfyunChatCompletionRes{}
			if err := json.Unmarshal(message, &chatCompletionRes); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream Xfyun model: %s, message: %s, error: %v", x.Model, message, err)

				end := gtime.TimestampMilli()
				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     errors.New(fmt.Sprintf("message: %s, error: %v", message, err)),
				}

				return
			}

			if chatCompletionRes.Header.Code != 0 {
				logger.Errorf(ctx, "ChatCompletionsStream Xfyun model: %s, chatCompletionRes: %s", x.Model, gjson.MustEncodeString(chatCompletionRes))

				err = x.apiErrorHandler(&chatCompletionRes)
				logger.Errorf(ctx, "ChatCompletionsStream Xfyun model: %s, error: %v", x.Model, err)

				end := gtime.TimestampMilli()
				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     err,
				}

				return
			}

			response := &model.ChatCompletionResponse{
				Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Header.Sid,
				Object:  consts.COMPLETION_STREAM_OBJECT,
				Created: created,
				Model:   x.Model,
				Choices: []model.ChatCompletionChoice{{
					Index: chatCompletionRes.Payload.Choices.Seq,
					Delta: &model.ChatCompletionStreamChoiceDelta{
						Role:         chatCompletionRes.Payload.Choices.Text[0].Role,
						Content:      chatCompletionRes.Payload.Choices.Text[0].Content,
						FunctionCall: chatCompletionRes.Payload.Choices.Text[0].FunctionCall,
					},
				}},
				ConnTime: duration - now,
			}

			if chatCompletionRes.Payload.Usage != nil {
				response.Usage = &model.Usage{
					PromptTokens:     chatCompletionRes.Payload.Usage.Text.PromptTokens,
					CompletionTokens: chatCompletionRes.Payload.Usage.Text.CompletionTokens,
					TotalTokens:      chatCompletionRes.Payload.Usage.Text.TotalTokens,
				}
			}

			if chatCompletionRes.Header.Status == 2 {

				logger.Infof(ctx, "ChatCompletionsStream Xfyun model: %s finished", x.Model)

				response.Choices[0].FinishReason = consts.FinishReasonStop

				end := gtime.TimestampMilli()
				response.Duration = end - duration
				response.TotalTime = end - now
				responseChan <- response

				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     io.EOF,
				}

				return
			}

			end := gtime.TimestampMilli()
			response.Duration = end - duration
			response.TotalTime = end - now

			responseChan <- response
		}

	}, nil); err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream Xfyun model: %s, error: %v", x.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
