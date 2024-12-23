package sdk

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
	"io"
	"net/http"
)

type RealtimeClient struct {
	model    string
	key      string
	baseURL  string
	path     string
	proxyURL string
}

func NewRealtimeClient(ctx context.Context, model, key, baseURL, path string, proxyURL ...string) *RealtimeClient {

	logger.Infof(ctx, "NewRealtimeClient OpenAI model: %s, key: %s", model, key)

	realtimeClient := &RealtimeClient{
		model:   model,
		key:     key,
		baseURL: "wss://api.openai.com/v1",
		path:    "/realtime",
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewRealtimeClient OpenAI model: %s, baseURL: %s", model, baseURL)
		realtimeClient.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewRealtimeClient OpenAI model: %s, path: %s", model, path)
		realtimeClient.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewRealtimeClient OpenAI model: %s, proxyURL: %s", model, proxyURL[0])
		realtimeClient.proxyURL = proxyURL[0]
	}

	return realtimeClient
}

func (c *RealtimeClient) Realtime(ctx context.Context, requestChan chan *model.RealtimeRequest) (responseChan chan *model.RealtimeResponse, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "Realtime OpenAI model: %s totalTime: %d ms", c.model, gtime.TimestampMilli()-now)
	}()

	logger.Infof(ctx, "Realtime OpenAI model: %s start", c.model)

	requestHeader := http.Header{
		"Authorization": {"Bearer " + c.key},
		"OpenAI-Beta":   {"realtime=v1"},
	}

	conn, err := util.WebSocketClient(ctx, c.getWebSocketUrl(ctx), requestHeader, 0, nil, c.proxyURL)
	if err != nil {
		logger.Errorf(ctx, "Realtime OpenAI model: %s, error: %v", c.model, err)
		return
	}

	duration := gtime.TimestampMilli()
	responseChan = make(chan *model.RealtimeResponse)

	// WriteMessage
	if err := grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			logger.Infof(ctx, "Realtime OpenAI WriteMessage model: %s totalTime: %d ms", c.model, gtime.TimestampMilli()-now)
		}()

		for {

			request := <-requestChan

			if request == nil || request.MessageType == -1 {

				if err := conn.Close(); err != nil {
					logger.Errorf(ctx, "Realtime OpenAI WriteMessage model: %s, conn.Close error: %v", c.model, err)
				}

				responseChan <- nil

				return
			}

			if err := conn.WriteMessage(ctx, request.MessageType, request.Message); err != nil {
				logger.Errorf(ctx, "Realtime OpenAI WriteMessage model: %s, error: %v", c.model, err)
				return
			}
		}

	}, nil); err != nil {
		logger.Errorf(ctx, "Realtime OpenAI WriteMessage model: %s, error: %v", c.model, err)
		return nil, err
	}

	// ReadMessage
	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.TimestampMilli()
			logger.Infof(ctx, "Realtime OpenAI ReadMessage model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", c.model, duration-now, end-duration, end-now)
		}()

		for {

			messageType, message, err := conn.ReadMessage(ctx)
			if err != nil && !errors.Is(err, io.EOF) {

				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "Realtime OpenAI ReadMessage model: %s, error: %v", c.model, err)
				}

				end := gtime.TimestampMilli()
				responseChan <- &model.RealtimeResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     err,
				}

				return
			}

			if messageType == -1 {
				return
			}

			response := &model.RealtimeResponse{
				MessageType: messageType,
				Message:     message,
				ConnTime:    duration - now,
			}

			end := gtime.TimestampMilli()
			response.Duration = end - duration
			response.TotalTime = end - now

			responseChan <- response
		}

	}, nil); err != nil {
		logger.Errorf(ctx, "Realtime OpenAI ReadMessage model: %s, error: %v", c.model, err)
		return
	}

	return responseChan, nil
}

func (c *RealtimeClient) getWebSocketUrl(ctx context.Context) string {
	webSocketUrl := gstr.Replace(gstr.Replace(fmt.Sprintf("%s%s?model=%s", c.baseURL, c.path, c.model), "https://", "wss://"), "http://", "ws://")
	logger.Infof(ctx, "Realtime OpenAI model: %s, webSocketUrl: %s", c.model, webSocketUrl)
	return webSocketUrl
}
