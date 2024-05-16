package zhipuai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/sdkerr"
	"github.com/iimeta/fastapi-sdk/util"
	"github.com/sashabaranov/go-openai"
	"io"
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
		Path:    "/chat/completions",
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewClient ZhipuAI model: %s, baseURL: %s", model, baseURL)
		client.BaseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewClient ZhipuAI model: %s, path: %s", model, path)
		client.Path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewClient ZhipuAI model: %s, proxyURL: %s", model, proxyURL[0])
		client.ProxyURL = proxyURL[0]
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

	chatCompletionReq := model.ZhipuAIChatCompletionReq{
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

	if chatCompletionReq.TopP == 1 {
		chatCompletionReq.TopP -= 0.01
	} else if chatCompletionReq.TopP == 0 {
		chatCompletionReq.TopP += 0.01
	}

	if chatCompletionReq.Temperature == 1 {
		chatCompletionReq.Temperature -= 0.01
	} else if chatCompletionReq.Temperature == 0 {
		chatCompletionReq.Temperature += 0.01
	}

	if chatCompletionReq.MaxTokens == 1 {
		chatCompletionReq.MaxTokens = 2
	}

	if chatCompletionReq.Messages[0].Role == openai.ChatMessageRoleSystem && chatCompletionReq.Messages[0].Content == "" && len(chatCompletionReq.Messages[0].ToolCalls) == 0 {
		chatCompletionReq.Messages = chatCompletionReq.Messages[1:]
	}

	header := make(map[string]string)
	header["Authorization"] = "Bearer " + c.generateToken(ctx)

	chatCompletionRes := new(model.ZhipuAIChatCompletionRes)
	if err = util.HttpPost(ctx, c.BaseURL+c.Path, header, chatCompletionReq, &chatCompletionRes, c.ProxyURL); err != nil {
		logger.Errorf(ctx, "ChatCompletion ZhipuAI model: %s, error: %v", request.Model, err)
		return
	}

	if chatCompletionRes.Error.Code != "" && chatCompletionRes.Error.Code != "200" {
		logger.Errorf(ctx, "ChatCompletion ZhipuAI model: %s, chatCompletionRes: %s", request.Model, gjson.MustEncodeString(chatCompletionRes))

		err = c.apiErrorHandler(chatCompletionRes)
		logger.Errorf(ctx, "ChatCompletion ZhipuAI model: %s, error: %v", request.Model, err)

		return
	}

	res = model.ChatCompletionResponse{
		ID:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Id,
		Object:  consts.COMPLETION_OBJECT,
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

	chatCompletionReq := model.ZhipuAIChatCompletionReq{
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

	if chatCompletionReq.TopP == 1 {
		chatCompletionReq.TopP -= 0.01
	} else if chatCompletionReq.TopP == 0 {
		chatCompletionReq.TopP += 0.01
	}

	if chatCompletionReq.Temperature == 1 {
		chatCompletionReq.Temperature -= 0.01
	} else if chatCompletionReq.Temperature == 0 {
		chatCompletionReq.Temperature += 0.01
	}

	if chatCompletionReq.MaxTokens == 1 {
		chatCompletionReq.MaxTokens = 2
	}

	if chatCompletionReq.Messages[0].Role == openai.ChatMessageRoleSystem && chatCompletionReq.Messages[0].Content == "" && len(chatCompletionReq.Messages[0].ToolCalls) == 0 {
		chatCompletionReq.Messages = chatCompletionReq.Messages[1:]
	}

	header := make(map[string]string)
	header["Authorization"] = "Bearer " + c.generateToken(ctx)

	stream, err := util.SSEClient(ctx, c.BaseURL+c.Path, header, chatCompletionReq, c.ProxyURL, c.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStream ZhipuAI model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	duration := gtime.Now().UnixMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream ZhipuAI model: %s, stream.Close error: %v", request.Model, err)
			}

			end := gtime.Now().UnixMilli()
			logger.Infof(ctx, "ChatCompletionStream ZhipuAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		for {

			streamResponse, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {

				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "ChatCompletionStream ZhipuAI model: %s, error: %v", request.Model, err)
				}

				end := gtime.Now().UnixMilli()
				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     err,
				}

				time.Sleep(time.Millisecond)
				close(responseChan)

				return
			}

			chatCompletionRes := new(model.ZhipuAIChatCompletionRes)
			if err := gjson.Unmarshal(streamResponse, &chatCompletionRes); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream ZhipuAI model: %s, error: %v", request.Model, err)

				end := gtime.Now().UnixMilli()
				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     err,
				}

				time.Sleep(time.Millisecond)
				close(responseChan)

				return
			}

			if chatCompletionRes.Error.Code != "" && chatCompletionRes.Error.Code != "200" {
				logger.Errorf(ctx, "ChatCompletionStream ZhipuAI model: %s, chatCompletionRes: %s", request.Model, gjson.MustEncodeString(chatCompletionRes))

				err = c.apiErrorHandler(chatCompletionRes)
				logger.Errorf(ctx, "ChatCompletionStream ZhipuAI model: %s, error: %v", request.Model, err)

				end := gtime.Now().UnixMilli()
				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     err,
				}

				time.Sleep(time.Millisecond)
				close(responseChan)

				return
			}

			response := &model.ChatCompletionResponse{
				ID:       consts.COMPLETION_ID_PREFIX + chatCompletionRes.Id,
				Object:   consts.COMPLETION_STREAM_OBJECT,
				Created:  chatCompletionRes.Created,
				Model:    request.Model,
				Usage:    chatCompletionRes.Usage,
				ConnTime: duration - now,
			}

			for _, choice := range chatCompletionRes.Choices {
				response.Choices = append(response.Choices, model.ChatCompletionChoice{
					Index:        choice.Index,
					Delta:        choice.Delta,
					FinishReason: choice.FinishReason,
				})
			}

			if errors.Is(err, io.EOF) || response.Choices[0].FinishReason != "" {
				logger.Infof(ctx, "ChatCompletionStream ZhipuAI model: %s finished", request.Model)

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

func (c *Client) requestErrorHandler(ctx context.Context, response *gclient.Response) error {

	errRes := model.ZhipuAIErrorResponse{}
	if err := json.NewDecoder(response.Body).Decode(&errRes); err != nil || errRes.Error == nil {

		reqErr := &sdkerr.RequestError{
			HttpStatusCode: response.StatusCode,
			Err:            err,
		}

		if errRes.Error != nil {
			reqErr.Err = errors.New(gjson.MustEncodeString(errRes.Error))
		}

		return reqErr
	}

	switch errRes.Error.Code {
	case "1261":
		return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
	case "1113":
		return sdkerr.ERR_INSUFFICIENT_QUOTA
	}

	return errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, gjson.MustEncodeString(errRes.Error)))
}

func (c *Client) apiErrorHandler(response *model.ZhipuAIChatCompletionRes) error {

	switch response.Error.Code {
	case "1261":
		return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
	case "1113":
		return sdkerr.ERR_INSUFFICIENT_QUOTA
	}

	return errors.New(gjson.MustEncodeString(response))
}
