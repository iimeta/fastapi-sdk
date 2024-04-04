package sdk

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"net/url"
	"time"
)

func NewClient(ctx context.Context, model, apiKey string, baseURL ...string) *openai.Client {

	logger.Infof(ctx, "NewClient OpenAI model: %s, apiKey: %s", model, apiKey)

	config := openai.DefaultConfig(apiKey)

	if len(baseURL) > 0 && baseURL[0] != "" {
		logger.Infof(ctx, "NewClient OpenAI model: %s, baseURL: %s", model, baseURL[0])
		config.BaseURL = baseURL[0]
	}

	return openai.NewClientWithConfig(config)
}

func NewProxyClient(ctx context.Context, model, apiKey string, proxyURL ...string) *openai.Client {

	logger.Infof(ctx, "NewProxyClient OpenAI model: %s, apiKey: %s", model, apiKey)

	config := openai.DefaultConfig(apiKey)

	transport := &http.Transport{}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewProxyClient OpenAI model: %s, proxyURL: %s", model, proxyURL[0])
		proxyUrl, err := url.Parse(proxyURL[0])
		if err != nil {
			panic(err)
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
	}

	config.HTTPClient = &http.Client{
		Transport: transport,
	}

	return openai.NewClientWithConfig(config)
}

func ChatCompletion(ctx context.Context, client *openai.Client, request openai.ChatCompletionRequest) (res model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletion OpenAI model: %s start", request.Model)

	now := gtime.Now().UnixMilli()

	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "ChatCompletion OpenAI model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletion OpenAI model: %s, error: %v", request.Model, err)
		return res, err
	}

	logger.Infof(ctx, "ChatCompletion OpenAI model: %s finished", request.Model)

	res = model.ChatCompletionResponse{
		ID:                response.ID,
		Object:            response.Object,
		Created:           response.Created,
		Model:             response.Model,
		Choices:           response.Choices,
		Usage:             response.Usage,
		SystemFingerprint: response.SystemFingerprint,
	}

	return res, nil
}

func ChatCompletionStream(ctx context.Context, client *openai.Client, request openai.ChatCompletionRequest) (responseChan chan *model.ChatCompletionStreamResponse, err error) {

	logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s start", request.Model)

	now := gtime.Now().UnixMilli()

	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s totalTime: %d ms", request.Model, gtime.Now().UnixMilli()-now)
		}
	}()

	stream, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStream OpenAI model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	duration := gtime.Now().UnixMilli()

	responseChan = make(chan *model.ChatCompletionStreamResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.Now().UnixMilli()
			logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		for {

			streamResponse, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {
				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "ChatCompletionStream OpenAI model: %s, error: %v", request.Model, err)
				}
				responseChan <- nil
				time.Sleep(time.Millisecond)
				close(responseChan)
				return
			}

			response := &model.ChatCompletionStreamResponse{
				ID:                streamResponse.ID,
				Object:            streamResponse.Object,
				Created:           streamResponse.Created,
				Model:             streamResponse.Model,
				PromptAnnotations: streamResponse.PromptAnnotations,
				ConnTime:          duration - now,
			}

			for _, choice := range streamResponse.Choices {
				response.Choices = append(response.Choices, model.ChatCompletionStreamChoice{
					Index:                choice.Index,
					Delta:                choice.Delta,
					FinishReason:         choice.FinishReason,
					ContentFilterResults: choice.ContentFilterResults,
				})
			}

			if errors.Is(err, io.EOF) || response.Choices[0].FinishReason == "stop" {

				logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s finished", request.Model)

				if err = stream.Close(); err != nil {
					logger.Error(ctx, err)
				}

				end := gtime.Now().UnixMilli()
				response.Duration = end - duration
				response.TotalTime = end - now
				responseChan <- response

				return
			}

			end := gtime.Now().UnixMilli()
			response.Duration = end - duration
			response.TotalTime = end - now

			responseChan <- response
		}
	}, nil); err != nil {
		logger.Error(ctx, err)
		return responseChan, err
	}

	return responseChan, nil
}

func Image(ctx context.Context, client *openai.Client, request openai.ImageRequest) (res model.ImageResponse, err error) {

	logger.Infof(ctx, "Image OpenAI model: %s start", request.Model)

	now := gtime.Now().UnixMilli()

	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Image OpenAI model: %s totalTime: %d ms", request.Model, gtime.Now().UnixMilli()-now)
	}()

	response, err := client.CreateImage(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "Image OpenAI model: %s, error: %v", request.Model, err)
		return res, err
	}

	res = model.ImageResponse{
		Created: response.Created,
		Data:    response.Data,
	}

	return res, nil
}
