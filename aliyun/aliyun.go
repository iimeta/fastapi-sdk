package aliyun

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	Key      string
	BaseURL  string
	Path     string
	ProxyURL string
}

func NewClient(ctx context.Context, model, key, baseURL, path string, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewClient Aliyun model: %s, key: %s", model, key)

	client := &Client{
		Key:     key,
		BaseURL: "https://aip.baidubce.com",
		Path:    path,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewClient Aliyun model: %s, baseURL: %s", model, baseURL)
		client.BaseURL = baseURL
	}

	return client
}

func (c *Client) ChatCompletion(ctx context.Context, request model.ChatCompletionRequest) (res model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletion Aliyun model: %s start", request.Model)

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "ChatCompletion Aliyun model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	qwenChatCompletionReq := model.QwenChatCompletionReq{
		Model: request.Model,
		Input: model.Input{
			Messages: request.Messages,
		},
	}

	header := make(map[string]string)
	header["Authorization"] = "Bearer " + c.Key

	qwenChatCompletionRes := new(model.QwenChatCompletionRes)
	err = util.HttpPostJson(ctx, c.BaseURL+c.Path, header, qwenChatCompletionReq, &qwenChatCompletionRes, c.ProxyURL)
	if err != nil {
		logger.Error(ctx, err)
		return
	}

	if qwenChatCompletionRes.Code != "" {
		err = errors.New(gjson.MustEncodeString(qwenChatCompletionRes))
		logger.Error(ctx)
		return
	}

	res = model.ChatCompletionResponse{
		ID:     qwenChatCompletionRes.RequestId,
		Object: qwenChatCompletionRes.Message,
		Model:  request.Model,
		Choices: []model.ChatCompletionChoice{{
			Message: &openai.ChatCompletionMessage{
				Role:    consts.ROLE_ASSISTANT,
				Content: qwenChatCompletionRes.Output.Text,
			},
		}},
		Usage: &openai.Usage{
			PromptTokens:     qwenChatCompletionRes.Usage.InputTokens,
			CompletionTokens: qwenChatCompletionRes.Usage.OutputTokens,
			TotalTokens:      qwenChatCompletionRes.Usage.InputTokens + qwenChatCompletionRes.Usage.OutputTokens,
		},
	}

	return res, nil
}

func (c *Client) ChatCompletionStream(ctx context.Context, request model.ChatCompletionRequest) (responseChan chan *model.ChatCompletionResponse, err error) {

	return
}

func (c *Client) Image(ctx context.Context, request model.ImageRequest) (res model.ImageResponse, err error) {

	return
}
