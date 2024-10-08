package xfyun

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/gorilla/websocket"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
	"github.com/iimeta/go-openai"
	"io"
)

func (c *Client) ChatCompletion(ctx context.Context, request model.ChatCompletionRequest) (res model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletion Xfyun model: %s start", request.Model)

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "ChatCompletion Xfyun model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, res.ConnTime, res.Duration, res.TotalTime)
	}()

	var messages []model.ChatCompletionMessage
	if c.isSupportSystemRole != nil {
		messages = common.HandleMessages(request.Messages, *c.isSupportSystemRole)
	} else {
		messages = common.HandleMessages(request.Messages, true)
	}

	if len(messages) == 1 && messages[0].Role == consts.ROLE_SYSTEM {
		messages[0].Role = consts.ROLE_USER
	}

	maxTokens := request.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	chatCompletionReq := model.XfyunChatCompletionReq{
		Header: model.Header{
			AppId: c.appId,
			Uid:   grand.Digits(10),
		},
		Parameter: model.Parameter{
			Chat: &model.Chat{
				Domain:      c.domain,
				MaxTokens:   maxTokens,
				Temperature: request.Temperature,
				TopK:        request.N,
				ChatId:      request.User,
			},
		},
		Payload: model.Payload{
			Message: &model.Message{
				Text: messages,
			},
		},
	}

	if request.Functions != nil && len(request.Functions) > 0 {
		chatCompletionReq.Payload.Functions = new(model.Functions)
		chatCompletionReq.Payload.Functions.Text = append(chatCompletionReq.Payload.Functions.Text, request.Functions...)
	}

	data, err := gjson.Marshal(chatCompletionReq)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletion Xfyun model: %s, error: %v", request.Model, err)
		return res, err
	}

	conn, err := util.WebSocketClient(ctx, c.getWebSocketUrl(ctx), nil, websocket.TextMessage, data, c.proxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletion Xfyun model: %s, error: %v", request.Model, err)
		return res, err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Errorf(ctx, "ChatCompletion Xfyun model: %s, conn.Close error: %v", request.Model, err)
		}
	}()

	duration := gtime.Now().UnixMilli()

	responseContent := ""
	chatCompletionRes := new(model.XfyunChatCompletionRes)

	for {

		_, message, err := conn.ReadMessage(ctx)
		if err != nil && !errors.Is(err, io.EOF) {
			logger.Errorf(ctx, "ChatCompletion Xfyun model: %s, error: %v", request.Model, err)
			return res, err
		}

		if err = gjson.Unmarshal(message, &chatCompletionRes); err != nil {
			logger.Errorf(ctx, "ChatCompletion Xfyun model: %s, message: %s, error: %v", request.Model, message, err)
			return res, errors.New(fmt.Sprintf("message: %s, error: %v", message, err))
		}

		if chatCompletionRes.Header.Code != 0 {
			logger.Errorf(ctx, "ChatCompletion Xfyun model: %s, chatCompletionRes: %s", request.Model, gjson.MustEncodeString(chatCompletionRes))

			err = c.apiErrorHandler(chatCompletionRes)
			logger.Errorf(ctx, "ChatCompletion Xfyun model: %s, error: %v", request.Model, err)

			return res, err
		}

		responseContent += chatCompletionRes.Payload.Choices.Text[0].Content

		if chatCompletionRes.Header.Status == 2 {
			break
		}
	}

	res = model.ChatCompletionResponse{
		ID:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Header.Sid,
		Object:  consts.COMPLETION_OBJECT,
		Created: gtime.Now().Unix(),
		Model:   request.Model,
		Choices: []model.ChatCompletionChoice{{
			Index: chatCompletionRes.Payload.Choices.Seq,
			Message: &openai.ChatCompletionMessage{
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
		Duration: gtime.Now().UnixMilli() - duration,
	}

	return res, nil
}

func (c *Client) ChatCompletionStream(ctx context.Context, request model.ChatCompletionRequest) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionStream Xfyun model: %s start", request.Model)

	now := gtime.Now().UnixMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionStream Xfyun model: %s totalTime: %d ms", request.Model, gtime.Now().UnixMilli()-now)
		}
	}()

	var messages []model.ChatCompletionMessage
	if c.isSupportSystemRole != nil {
		messages = common.HandleMessages(request.Messages, *c.isSupportSystemRole)
	} else {
		messages = common.HandleMessages(request.Messages, true)
	}

	if len(messages) == 1 && messages[0].Role == consts.ROLE_SYSTEM {
		messages[0].Role = consts.ROLE_USER
	}

	maxTokens := request.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	chatCompletionReq := model.XfyunChatCompletionReq{
		Header: model.Header{
			AppId: c.appId,
			Uid:   grand.Digits(10),
		},
		Parameter: model.Parameter{
			Chat: &model.Chat{
				Domain:      c.domain,
				MaxTokens:   maxTokens,
				Temperature: request.Temperature,
				TopK:        request.N,
				ChatId:      request.User,
			},
		},
		Payload: model.Payload{
			Message: &model.Message{
				Text: messages,
			},
		},
	}

	if request.Functions != nil && len(request.Functions) > 0 {
		chatCompletionReq.Payload.Functions = new(model.Functions)
		chatCompletionReq.Payload.Functions.Text = append(chatCompletionReq.Payload.Functions.Text, request.Functions...)
	}

	data, err := gjson.Marshal(chatCompletionReq)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStream Xfyun model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	conn, err := util.WebSocketClient(ctx, c.getWebSocketUrl(ctx), nil, websocket.TextMessage, data, c.proxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStream Xfyun model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	duration := gtime.Now().UnixMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := conn.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream Xfyun model: %s, conn.Close error: %v", request.Model, err)
			}

			end := gtime.Now().UnixMilli()
			logger.Infof(ctx, "ChatCompletionStream Xfyun model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		var created = gtime.Now().Unix()

		for {

			_, message, err := conn.ReadMessage(ctx)
			if err != nil && !errors.Is(err, io.EOF) {

				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "ChatCompletionStream Xfyun model: %s, error: %v", request.Model, err)
				}

				end := gtime.Now().UnixMilli()
				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     err,
				}

				return
			}

			chatCompletionRes := new(model.XfyunChatCompletionRes)
			if err := gjson.Unmarshal(message, &chatCompletionRes); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream Xfyun model: %s, message: %s, error: %v", request.Model, message, err)

				end := gtime.Now().UnixMilli()
				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     errors.New(fmt.Sprintf("message: %s, error: %v", message, err)),
				}

				return
			}

			if chatCompletionRes.Header.Code != 0 {
				logger.Errorf(ctx, "ChatCompletionStream Xfyun model: %s, chatCompletionRes: %s", request.Model, gjson.MustEncodeString(chatCompletionRes))

				err = c.apiErrorHandler(chatCompletionRes)
				logger.Errorf(ctx, "ChatCompletionStream Xfyun model: %s, error: %v", request.Model, err)

				end := gtime.Now().UnixMilli()
				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     err,
				}

				return
			}

			response := &model.ChatCompletionResponse{
				ID:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Header.Sid,
				Object:  consts.COMPLETION_STREAM_OBJECT,
				Created: created,
				Model:   request.Model,
				Choices: []model.ChatCompletionChoice{{
					Index: chatCompletionRes.Payload.Choices.Seq,
					Delta: &openai.ChatCompletionStreamChoiceDelta{
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

				logger.Infof(ctx, "ChatCompletionStream Xfyun model: %s finished", request.Model)

				response.Choices[0].FinishReason = openai.FinishReasonStop

				end := gtime.Now().UnixMilli()
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

			end := gtime.Now().UnixMilli()
			response.Duration = end - duration
			response.TotalTime = end - now

			responseChan <- response
		}
	}, nil); err != nil {
		logger.Errorf(ctx, "ChatCompletionStream Xfyun model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
