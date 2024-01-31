package sdk

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/sdk"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"net/url"
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

	logger.Infof(ctx, "ChatCompletion OpenAI model: %s", request.Model)

	now := gtime.Now().UnixMilli()

	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "ChatCompletion OpenAI model: %s, totalTime: %d ms", request.Model, res.TotalTime)
	}()

	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletion OpenAI model: %s, error: %v", request.Model, err)
		return model.ChatCompletionResponse{}, err
	}

	logger.Infof(ctx, "ChatCompletion OpenAI model: %s, response: %s", request.Model, gjson.MustEncodeString(response))

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

func ChatCompletionStream(ctx context.Context, client *openai.Client, request openai.ChatCompletionRequest) (responseChan chan model.ChatCompletionStreamResponse, err error) {

	logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s", request.Model)

	now := gtime.Now().UnixMilli()

	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s, totalTime: %d ms", request.Model, gtime.Now().UnixMilli()-now)
		}
	}()

	stream, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStream OpenAI model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s, start", request.Model)

	duration := gtime.Now().UnixMilli()

	responseChan = make(chan model.ChatCompletionStreamResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.Now().UnixMilli()
			logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s, connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		promptTokens, err := sdk.NumTokensFromMessages(request.Messages, request.Model)
		if err != nil {
			logger.Error(ctx, err)
			return
		}

		completionTokens := 0

		for {

			streamResponse, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {
				logger.Errorf(ctx, "ChatCompletionStream OpenAI model: %s, error: %v", request.Model, err)
				close(responseChan)
				return
			}

			response := model.ChatCompletionStreamResponse{
				ID:                streamResponse.ID,
				Object:            streamResponse.Object,
				Created:           streamResponse.Created,
				Model:             streamResponse.Model,
				Choices:           streamResponse.Choices,
				PromptAnnotations: streamResponse.PromptAnnotations,
				ConnTime:          duration - now,
			}

			if response.Usage.CompletionTokens, err = sdk.NumTokensFromString(response.Choices[0].Delta.Content, request.Model); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream OpenAI model: %s, NumTokensFromString error: %v", request.Model, err)
				return
			}

			completionTokens += response.Usage.CompletionTokens
			response.Usage.PromptTokens = promptTokens
			response.Usage.CompletionTokens = completionTokens
			response.Usage.TotalTokens = response.Usage.PromptTokens + response.Usage.CompletionTokens

			if errors.Is(err, io.EOF) || response.Choices[0].FinishReason == "stop" {
				logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s, finished", request.Model)
				stream.Close()
				end := gtime.Now().UnixMilli()
				response.Duration = end - duration
				response.TotalTime = end - now
				responseChan <- response
				return
			}

			responseChan <- response
		}
	}, nil); err != nil {
		logger.Error(ctx, err)
		return responseChan, err
	}

	return responseChan, nil
}

func GenImage(ctx context.Context, client *openai.Client, model, prompt string) (url string, err error) {

	logger.Infof(ctx, "GenImage OpenAI model: %s, prompt: %s", model, prompt)

	now := gtime.Now().UnixMilli()

	defer func() {
		logger.Infof(ctx, "GenImage OpenAI model: %s, url: %s", model, url)
		logger.Infof(ctx, "GenImage OpenAI model: %s, totalTime: %d ms", model, gtime.Now().UnixMilli()-now)
	}()

	reqUrl := openai.ImageRequest{
		Model:          model,
		Prompt:         prompt,
		Size:           openai.CreateImageSize1024x1024,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		N:              1,
	}

	respUrl, err := client.CreateImage(ctx, reqUrl)
	if err != nil {
		logger.Errorf(ctx, "GenImage OpenAI creation error: %v", err)
		return "", err
	}

	url = respUrl.Data[0].URL

	return url, nil
}

func GenImageBase64(ctx context.Context, client *openai.Client, model, prompt string) (string, error) {

	logger.Infof(ctx, "GenImageBase64 OpenAI model: %s, prompt: %s", model, prompt)

	now := gtime.Now().UnixMilli()

	imgBase64 := ""

	defer func() {
		logger.Infof(ctx, "GenImageBase64 OpenAI model: %s, len: %d", model, len(imgBase64))
		logger.Infof(ctx, "GenImageBase64 OpenAI model: %s, totalTime: %d ms", model, gtime.Now().UnixMilli()-now)
	}()

	reqBase64 := openai.ImageRequest{
		Model:          model,
		Prompt:         prompt,
		Size:           openai.CreateImageSize1024x1024,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}

	respBase64, err := client.CreateImage(ctx, reqBase64)
	if err != nil {
		logger.Errorf(ctx, "GenImageBase64 OpenAI model: %s, creation error: %v", model, err)
		return "", err
	}

	imgBase64 = respBase64.Data[0].B64JSON

	return imgBase64, nil
}
