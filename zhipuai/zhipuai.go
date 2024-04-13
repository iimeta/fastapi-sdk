package zhipuai

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	Key      string
	BaseURL  string
	Path     string
	ProxyURL string
}

func NewClient(ctx context.Context, model, key, baseURL, path string, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewClient ZhipuAI model: %s, key: %s", model, key)

	client := &Client{
		Key:     key,
		BaseURL: "https://open.bigmodel.cn/api/paas/v4",
		Path:    path,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewClient ZhipuAI model: %s, baseURL: %s", model, baseURL)
		client.BaseURL = baseURL
	}

	return client
}

func (c *Client) ChatCompletion(ctx context.Context, request model.ChatCompletionRequest) (res model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletion ZhipuAI model: %s start", request.Model)

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "ChatCompletion ZhipuAI model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	req := model.ZhipuAIChatCompletionReq{
		Model:       request.Model,
		Messages:    request.Messages,
		MaxTokens:   request.MaxTokens,
		Temperature: request.Temperature,
		TopP:        request.TopP,
		Stream:      request.Stream,
		Stop:        request.Stop,
		Tools:       request.Tools,
		ToolChoice:  request.ToolChoice,
		UserId:      request.User,
	}

	header := make(map[string]string)
	header["Authorization"] = "Bearer " + c.generateToken(ctx)

	chatCompletionRes := new(model.ZhipuAIChatCompletionRes)
	err = util.HttpPostJson(ctx, c.BaseURL+c.Path, header, req, &chatCompletionRes, c.ProxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletion ZhipuAI model: %s, error: %v", request.Model, err)
		return
	}

	res = model.ChatCompletionResponse{
		ID:      chatCompletionRes.Id,
		Created: chatCompletionRes.Created,
		Model:   request.Model,
		Usage:   chatCompletionRes.Usage,
	}

	for _, choice := range chatCompletionRes.Choices {
		res.Choices = append(res.Choices, model.ChatCompletionChoice{
			Index:        choice.Index,
			Message:      choice.Message,
			FinishReason: choice.FinishReason,
		})
	}

	return res, nil
}

func (c *Client) ChatCompletionStream(ctx context.Context, request model.ChatCompletionRequest) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionStream ZhipuAI model: %s start", request.Model)

	now := gtime.Now().UnixMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionStream ZhipuAI model: %s totalTime: %d ms", request.Model, gtime.Now().UnixMilli()-now)
		}
	}()

	req := model.ZhipuAIChatCompletionReq{
		Model:       request.Model,
		Messages:    request.Messages,
		MaxTokens:   request.MaxTokens,
		Temperature: request.Temperature,
		TopP:        request.TopP,
		Stream:      request.Stream,
		Stop:        request.Stop,
		Tools:       request.Tools,
		ToolChoice:  request.ToolChoice,
		UserId:      request.User,
	}

	header := make(map[string]string)
	header["Authorization"] = "Bearer " + c.generateToken(ctx)

	stream, err := util.SSEClient(ctx, http.MethodPost, c.BaseURL+c.Path, header, req)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStream ZhipuAI model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	duration := gtime.Now().UnixMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.Now().UnixMilli()
			logger.Infof(ctx, "ChatCompletionStream ZhipuAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		for {

			streamResponse, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {

				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "ChatCompletionStream ZhipuAI model: %s, error: %v", request.Model, err)
				}

				responseChan <- nil
				time.Sleep(time.Millisecond)
				close(responseChan)

				return
			}

			chatCompletionRes := new(model.ZhipuAIChatCompletionRes)
			if err := gjson.Unmarshal(streamResponse, &chatCompletionRes); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream ZhipuAI model: %s, error: %v", request.Model, err)

				responseChan <- nil
				time.Sleep(time.Millisecond)
				close(responseChan)

				return
			}

			response := &model.ChatCompletionResponse{
				ID:       chatCompletionRes.Id,
				Created:  chatCompletionRes.Created,
				Model:    request.Model,
				Usage:    chatCompletionRes.Usage,
				ConnTime: duration - now,
			}

			for _, choice := range chatCompletionRes.Choices {
				response.Choices = append(response.Choices, model.ChatCompletionChoice{
					Index:        choice.Index,
					Message:      choice.Message,
					FinishReason: choice.FinishReason,
				})
			}

			if errors.Is(err, io.EOF) || response.Choices[0].FinishReason == openai.FinishReasonStop {

				logger.Infof(ctx, "ChatCompletionStream ZhipuAI model: %s finished", request.Model)

				if err = stream.Close(); err != nil {
					logger.Errorf(ctx, "ChatCompletionStream ZhipuAI model: %s, stream.Close error: %v", request.Model, err)
				}

				if len(response.Choices) == 0 {
					response.Choices = append(response.Choices, model.ChatCompletionChoice{
						Delta:        new(openai.ChatCompletionStreamChoiceDelta),
						FinishReason: openai.FinishReasonStop,
					})
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
		logger.Errorf(ctx, "ChatCompletionStream ZhipuAI model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (c *Client) Image(ctx context.Context, request model.ImageRequest) (res model.ImageResponse, err error) {

	return
}

func (c *Client) generateToken(ctx context.Context) string {

	split := strings.Split(c.Key, ".")
	if len(split) != 2 {
		return c.Key
	}

	now := gtime.Now()

	claims := jwt.MapClaims{
		"api_key":   split[0],
		"exp":       now.Add(time.Minute * 10).UnixMilli(),
		"timestamp": now.UnixMilli(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token.Header["alg"] = "HS256"
	token.Header["sign_type"] = "SIGN"

	sign, err := token.SignedString([]byte(split[1]))
	if err != nil {
		logger.Error(ctx, err)
	}

	return sign
}
