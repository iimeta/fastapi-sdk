package sdk

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"net/url"
)

func NewClient(ctx context.Context, model, apiKey string, baseURL ...string) *openai.Client {

	Infof(ctx, "NewClient OpenAI model: %s, apiKey: %s", model, apiKey)

	config := openai.DefaultConfig(apiKey)

	if len(baseURL) > 0 {
		Infof(ctx, "NewClient OpenAI model: %s, baseURL: %s", model, baseURL[0])
		config.BaseURL = baseURL[0]
	}

	return openai.NewClientWithConfig(config)
}

func NewProxyClient(ctx context.Context, model, apiKey string, proxyURL ...string) *openai.Client {

	Infof(ctx, "NewProxyClient OpenAI model: %s, apiKey: %s", model, apiKey)

	config := openai.DefaultConfig(apiKey)

	transport := &http.Transport{}

	if len(proxyURL) > 0 {
		Infof(ctx, "NewProxyClient OpenAI model: %s, proxyURL: %s", model, proxyURL[0])
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

func ChatCompletion(ctx context.Context, client *openai.Client, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {

	Infof(ctx, "ChatCompletion OpenAI model: %s", request.Model)

	now := gtime.Now().UnixMilli()

	defer func() {
		Infof(ctx, "ChatCompletion OpenAI model: %s, totalTime: %d ms", request.Model, gtime.Now().UnixMilli()-now)
	}()

	response, err := client.CreateChatCompletion(ctx, request)

	if err != nil {
		Errorf(ctx, "ChatCompletion OpenAI model: %s, error: %v", request.Model, err)
		return openai.ChatCompletionResponse{}, err
	}

	Infof(ctx, "ChatCompletion OpenAI model: %s, response: %s", request.Model, gjson.MustEncodeString(response))

	return response, nil
}

func ChatCompletionStream(ctx context.Context, client *openai.Client, request openai.ChatCompletionRequest) (responseChan chan openai.ChatCompletionStreamResponse, err error) {

	Infof(ctx, "ChatCompletionStream OpenAI model: %s", request.Model)

	now := gtime.Now().UnixMilli()

	defer func() {
		if err != nil {
			Infof(ctx, "ChatCompletionStream OpenAI model: %s, totalTime: %d ms", request.Model, gtime.Now().UnixMilli()-now)
		}
	}()

	stream, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		Errorf(ctx, "ChatCompletionStream OpenAI model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	Infof(ctx, "ChatCompletionStream OpenAI model: %s, start", request.Model)

	duration := gtime.Now().UnixMilli()

	responseChan = make(chan openai.ChatCompletionStreamResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.Now().UnixMilli()
			Infof(ctx, "ChatCompletionStream OpenAI model: %s, connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		for {

			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				Infof(ctx, "ChatCompletionStream OpenAI model: %s, finished", request.Model)
				stream.Close()
				responseChan <- response
				return
			}

			if err != nil {
				Errorf(ctx, "ChatCompletionStream OpenAI model: %s, error: %v", request.Model, err)
				close(responseChan)
				return
			}

			responseChan <- response
		}
	}, nil); err != nil {
		Error(ctx, err)
		return responseChan, err
	}

	return responseChan, nil
}

func GenImage(ctx context.Context, client *openai.Client, model, prompt string) (url string, err error) {

	Infof(ctx, "GenImage OpenAI model: %s, prompt: %s", model, prompt)

	now := gtime.Now().UnixMilli()

	defer func() {
		Infof(ctx, "GenImage OpenAI model: %s, url: %s", model, url)
		Infof(ctx, "GenImage OpenAI model: %s, totalTime: %d ms", model, gtime.Now().UnixMilli()-now)
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
		Errorf(ctx, "GenImage OpenAI creation error: %v", err)
		return "", err
	}

	url = respUrl.Data[0].URL

	return url, nil
}

func GenImageBase64(ctx context.Context, client *openai.Client, model, prompt string) (string, error) {

	Infof(ctx, "GenImageBase64 OpenAI model: %s, prompt: %s", model, prompt)

	now := gtime.Now().UnixMilli()

	imgBase64 := ""

	defer func() {
		Infof(ctx, "GenImageBase64 OpenAI model: %s, len: %d", model, len(imgBase64))
		Infof(ctx, "GenImageBase64 OpenAI model: %s, totalTime: %d ms", model, gtime.Now().UnixMilli()-now)
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
		Errorf(ctx, "GenImageBase64 OpenAI model: %s, creation error: %v", model, err)
		return "", err
	}

	imgBase64 = respBase64.Data[0].B64JSON

	return imgBase64, nil
}
