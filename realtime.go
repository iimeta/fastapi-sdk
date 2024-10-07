package sdk

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/sdkerr"
	"github.com/iimeta/fastapi-sdk/util"
	"github.com/iimeta/go-openai"
	"io"
	"net/http"
	"net/url"
)

type RealtimeClient struct {
	client   *openai.Client
	key      string
	proxyURL string
}

func NewRealtimeClient(ctx context.Context, model, key, baseURL, path string, proxyURL ...string) *RealtimeClient {

	logger.Infof(ctx, "NewClient OpenAI model: %s, key: %s", model, key)

	config := openai.DefaultConfig(key)

	if baseURL != "" {
		logger.Infof(ctx, "NewClient OpenAI model: %s, baseURL: %s", model, baseURL)
		config.BaseURL = baseURL
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewClient OpenAI model: %s, proxyURL: %s", model, proxyURL[0])

		proxyUrl, err := url.Parse(proxyURL[0])
		if err != nil {
			panic(err)
		}

		config.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	}

	return &RealtimeClient{
		client: openai.NewClientWithConfig(config),
		key:    key,
	}
}

func (c *RealtimeClient) Realtime(ctx context.Context, requestChan chan *model.RealtimeRequest) (responseChan chan *model.RealtimeResponse, err error) {

	responseChan = make(chan *model.RealtimeResponse)
	var requestModel string

	if err := grpool.Add(ctx, func(ctx context.Context) {

		now := gtime.Now().UnixMilli()
		defer func() {
			if err != nil {
				logger.Infof(ctx, "Realtime OpenAI model: %s totalTime: %d ms", requestModel, gtime.Now().UnixMilli()-now)
			}
		}()

		for {

			request := <-requestChan
			requestModel = request.Model
			logger.Infof(ctx, "Realtime OpenAI model: %s start", request.Model)

			requestHeader := http.Header{
				"Authorization": {"Bearer " + c.key},
				"OpenAI-Beta":   {"realtime=v1"},
			}

			conn, err := util.WebSocketClient(ctx, c.getWebSocketUrl(ctx, request.Model), requestHeader, request.MessageType, request.Message, c.proxyURL)
			if err != nil {
				logger.Errorf(ctx, "Realtime OpenAI model: %s, error: %v", request.Model, err)
				return
			}

			duration := gtime.Now().UnixMilli()

			if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

				defer func() {
					if err := conn.Close(); err != nil {
						logger.Errorf(ctx, "Realtime OpenAI model: %s, conn.Close error: %v", request.Model, err)
					}

					end := gtime.Now().UnixMilli()
					logger.Infof(ctx, "Realtime OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
				}()

				for {

					message, err := conn.ReadMessage(ctx)
					if err != nil && !errors.Is(err, io.EOF) {

						if !errors.Is(err, context.Canceled) {
							logger.Errorf(ctx, "Realtime OpenAI model: %s, error: %v", request.Model, err)
						}

						end := gtime.Now().UnixMilli()
						responseChan <- &model.RealtimeResponse{
							ConnTime:  duration - now,
							Duration:  end - duration,
							TotalTime: end - now,
							Error:     err,
						}

						return
					}

					response := &model.RealtimeResponse{
						Message:  message,
						ConnTime: duration - now,
					}

					end := gtime.Now().UnixMilli()
					response.Duration = end - duration
					response.TotalTime = end - now

					responseChan <- response
				}
			}, nil); err != nil {
				logger.Errorf(ctx, "Realtime OpenAI model: %s, error: %v", request.Model, err)
				return
			}
		}

	}); err != nil {
		logger.Errorf(ctx, "Realtime OpenAI model: %s, error: %v", requestModel, err)
		return nil, err
	}

	return responseChan, nil
}

func (c *RealtimeClient) getWebSocketUrl(ctx context.Context, model string) string {
	return fmt.Sprintf("wss://api.openai.com/v1/realtime?model=%s", model)
}

func (c *RealtimeClient) apiErrorHandler(response *model.XfyunChatCompletionRes) error {

	switch response.Header.Code {
	}

	return sdkerr.NewApiError(500, response.Header.Code, gjson.MustEncodeString(response), "api_error", "")
}
