package baidu

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"time"
)

type Client struct {
	AppId    string
	Key      string
	Secret   string
	BaseURL  string
	Path     string
	ProxyURL string
}

func NewClient(ctx context.Context, model, key, baseURL, path string, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewClient Baidu model: %s, key: %s", model, key)

	result := gstr.Split(key, "|")

	client := &Client{
		AppId:   result[0],
		Key:     result[1],
		Secret:  result[2],
		BaseURL: "https://aip.baidubce.com",
		Path:    path,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewClient Baidu model: %s, baseURL: %s", model, baseURL)
		client.BaseURL = baseURL
	}

	return client
}

func (c *Client) ChatCompletion(ctx context.Context, request model.ChatCompletionRequest) (res model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletion Baidu model: %s start", request.Model)

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "ChatCompletion Baidu model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	req := model.ErnieBotReq{
		Messages: request.Messages,
	}

	ernieBotRes := new(model.ErnieBotRes)
	err = util.HttpPostJson(ctx, fmt.Sprintf("%s?access_token=%s", c.BaseURL+c.Path, c.GetAccessToken(ctx)), nil, req, &ernieBotRes, c.ProxyURL)
	if err != nil {
		logger.Error(ctx, err)
		return
	}

	if ernieBotRes.ErrorCode != 0 {
		err = errors.New(gjson.MustEncodeString(ernieBotRes))
		logger.Error(ctx, err)
		return
	}

	res = model.ChatCompletionResponse{
		ID:      ernieBotRes.Id,
		Object:  ernieBotRes.Object,
		Created: ernieBotRes.Created,
		Model:   request.Model,
		Choices: []model.ChatCompletionChoice{{
			Message: &openai.ChatCompletionMessage{
				Role:    consts.ROLE_ASSISTANT,
				Content: ernieBotRes.Result,
			},
		}},
		Usage: &openai.Usage{
			PromptTokens:     ernieBotRes.Usage.PromptTokens,
			CompletionTokens: ernieBotRes.Usage.CompletionTokens,
			TotalTokens:      ernieBotRes.Usage.TotalTokens,
		},
	}

	return res, nil
}

func (c *Client) ChatCompletionStream(ctx context.Context, request model.ChatCompletionRequest) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionStream Baidu model: %s start", request.Model)

	now := gtime.Now().UnixMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionStream Baidu model: %s totalTime: %d ms", request.Model, gtime.Now().UnixMilli()-now)
		}
	}()

	req := model.ErnieBotReq{
		Messages: request.Messages,
		Stream:   request.Stream,
	}

	stream, err := util.SSEClient(ctx, http.MethodPost, fmt.Sprintf("%s?access_token=%s", c.BaseURL+c.Path, c.GetAccessToken(ctx)), nil, req)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStream Baidu model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	duration := gtime.Now().UnixMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.Now().UnixMilli()
			logger.Infof(ctx, "ChatCompletionStream Baidu model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		for {

			streamResponse, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {

				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "ChatCompletionStream Baidu model: %s, error: %v", request.Model, err)
				}

				responseChan <- nil
				time.Sleep(time.Millisecond)
				close(responseChan)

				return
			}

			ernieBotRes := new(model.ErnieBotRes)
			if err = gjson.Unmarshal(streamResponse, &ernieBotRes); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream Baidu model: %s, error: %v", request.Model, err)

				responseChan <- nil
				time.Sleep(time.Millisecond)
				close(responseChan)

				return
			}

			if ernieBotRes.ErrorCode != 0 {

				err = errors.New(gjson.MustEncodeString(ernieBotRes))
				logger.Errorf(ctx, "ChatCompletionStream Baidu model: %s, error: %v", request.Model, err)

				if err = stream.Close(); err != nil {
					logger.Errorf(ctx, "ChatCompletionStream Baidu model: %s, stream.Close error: %v", request.Model, err)
				}

				end := gtime.Now().UnixMilli()

				responseChan <- &model.ChatCompletionResponse{
					ID:      ernieBotRes.Id,
					Object:  ernieBotRes.Object,
					Created: ernieBotRes.Created,
					Model:   request.Model,
					Choices: []model.ChatCompletionChoice{{
						Index: 0,
						Delta: &openai.ChatCompletionStreamChoiceDelta{
							Role:    consts.ROLE_ASSISTANT,
							Content: ernieBotRes.ErrorMsg,
						},
						FinishReason: "stop",
					}},
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
				}

				return
			}

			response := &model.ChatCompletionResponse{
				ID:      ernieBotRes.Id,
				Object:  ernieBotRes.Object,
				Created: ernieBotRes.Created,
				Model:   request.Model,
				Choices: []model.ChatCompletionChoice{{
					Index: ernieBotRes.SentenceId,
					Delta: &openai.ChatCompletionStreamChoiceDelta{
						Role:    consts.ROLE_ASSISTANT,
						Content: ernieBotRes.Result,
					},
				}},
				Usage:    ernieBotRes.Usage,
				ConnTime: duration - now,
			}

			if ernieBotRes.IsEnd {

				logger.Infof(ctx, "ChatCompletionStream Baidu model: %s finished", request.Model)

				if err = stream.Close(); err != nil {
					logger.Errorf(ctx, "ChatCompletionStream Baidu model: %s, stream.Close error: %v", request.Model, err)
				}

				response.Choices[0].FinishReason = "stop"

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

func (c *Client) Image(ctx context.Context, request model.ImageRequest) (res model.ImageResponse, err error) {

	return
}

func (c *Client) GetAccessToken(ctx context.Context) string {

	reply, err := g.Redis().Get(ctx, fmt.Sprintf(consts.ACCESS_TOKEN_KEY, c.AppId))
	if err == nil && reply.String() != "" {
		return reply.String()
	}

	data := g.Map{
		"grant_type":    "client_credentials",
		"client_id":     c.Key,
		"client_secret": c.Secret,
	}

	getAccessTokenRes := new(model.GetAccessTokenRes)
	err = util.HttpPost(ctx, c.BaseURL+"/oauth/2.0/token", nil, data, &getAccessTokenRes, c.ProxyURL)
	if err != nil {
		logger.Error(ctx, err)
		return ""
	}

	if getAccessTokenRes.Error != "" {
		logger.Error(ctx, getAccessTokenRes.Error)
		return ""
	}

	_ = g.Redis().SetEX(ctx, fmt.Sprintf(consts.ACCESS_TOKEN_KEY, c.AppId), getAccessTokenRes.AccessToken, getAccessTokenRes.ExpiresIn)

	return getAccessTokenRes.AccessToken
}
