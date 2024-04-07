package xfyun

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/encoding/gurl"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
	"github.com/sashabaranov/go-openai"
	"net/url"
	"time"
)

type Client struct {
	AppId       string `json:"appid"`
	Key         string `json:"key"`
	Secret      string `json:"secret"`
	OriginalURL string `json:"original_url"`
	BaseURL     string `json:"base_url"`
	Path        string `json:"path"`
	ProxyURL    string `json:"proxy_url"`
	MaxTokens   int    `json:"max_tokens"`
	Domain      string `json:"domain"`
}

func NewClient(ctx context.Context, model, key string, baseURL ...string) *Client {

	logger.Infof(ctx, "NewClient Xfyun model: %s, key: %s", model, key)

	client := &Client{
		AppId:       "",
		Key:         "",
		Secret:      "",
		OriginalURL: "https://spark-api.xf-yun.com",
		BaseURL:     "https://spark-api.xf-yun.com",
		Path:        "/v3.5/chat",
		MaxTokens:   8192,
		Domain:      "generalv3.5",
	}

	//if len(baseURL) > 0 && baseURL[0] != "" {
	//	logger.Infof(ctx, "NewClient Xfyun model: %s, baseURL: %s", model, baseURL[0])
	//	client.BaseURL = baseURL[0]
	//}

	return client
}

func NewProxyClient(ctx context.Context, model, key string, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewProxyClient Xfyun model: %s, key: %s", model, key)

	client := &Client{
		AppId:       "",
		Key:         "",
		Secret:      "",
		OriginalURL: "https://spark-api.xf-yun.com",
		BaseURL:     "https://spark-api.xf-yun.com",
		Path:        "/v3.5/chat",
		MaxTokens:   8192,
		Domain:      "generalv3.5",
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewProxyClient Xfyun model: %s, proxyURL: %s", model, proxyURL[0])
		client.ProxyURL = proxyURL[0]
	}

	return client
}

func (c *Client) ChatCompletion(ctx context.Context, request model.ChatCompletionRequest) (res model.ChatCompletionResponse, err error) {

	now := gtime.Now().UnixMilli()

	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "ChatCompletion Xfyun model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	sparkReq := model.SparkReq{
		Header: model.Header{
			AppId: c.AppId,
			Uid:   uuid.NewString(),
		},
		Parameter: model.Parameter{
			Chat: &model.Chat{
				Domain:          c.Domain,
				RandomThreshold: 0,
				MaxTokens:       c.MaxTokens,
			},
		},
		Payload: model.Payload{
			Message: &model.Message{
				Text: request.Messages,
			},
		},
	}

	data, err := gjson.Marshal(sparkReq)
	if err != nil {
		logger.Error(ctx, err)
		return
	}

	authorizationUrl := c.getAuthorizationUrl(ctx)

	logger.Infof(ctx, "ChatCompletion Xfyun model: %s, appid: %s, authorizationUrl: %s", request.Model, c.AppId, authorizationUrl)

	result := make(chan []byte)
	var conn *websocket.Conn

	_ = grpool.AddWithRecover(ctx, func(ctx context.Context) {
		conn, err = util.WebSocketClient(ctx, authorizationUrl, websocket.TextMessage, data, result, c.ProxyURL)
		if err != nil {
			logger.Error(ctx, err)
		}
	}, nil)

	defer func() {
		err = conn.Close()
		if err != nil {
			logger.Error(ctx, err)
		}
	}()

	responseContent := ""
	sparkRes := new(model.SparkRes)
	for {

		message := <-result

		err = gjson.Unmarshal(message, &sparkRes)
		if err != nil {
			logger.Error(ctx, err)
			return
		}

		if sparkRes.Header.Code != 0 {
			err = errors.New(gjson.MustEncodeString(sparkRes))
			return
		}

		responseContent += sparkRes.Payload.Choices.Text[0].Content

		if sparkRes.Header.Status == 2 {
			sparkRes.Content = responseContent
			break
		}
	}

	res = model.ChatCompletionResponse{
		ID:    sparkRes.Header.Sid,
		Model: request.Model,
		Choices: []model.ChatCompletionChoice{{
			Index: sparkRes.Payload.Choices.Seq,
			Message: openai.ChatCompletionMessage{
				Role:    sparkRes.Payload.Choices.Text[0].Role,
				Content: responseContent,
			},
		}},
		Usage: openai.Usage{
			PromptTokens:     sparkRes.Payload.Usage.Text.PromptTokens,
			CompletionTokens: sparkRes.Payload.Usage.Text.CompletionTokens,
			TotalTokens:      sparkRes.Payload.Usage.Text.TotalTokens,
		},
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

	sparkReq := model.SparkReq{
		Header: model.Header{
			AppId: c.AppId,
			Uid:   uuid.NewString(),
		},
		Parameter: model.Parameter{
			Chat: &model.Chat{
				Domain:          c.Domain,
				RandomThreshold: 0,
				MaxTokens:       c.MaxTokens,
			},
		},
		Payload: model.Payload{
			Message: &model.Message{
				Text: request.Messages,
			},
		},
	}

	data, err := gjson.Marshal(sparkReq)
	if err != nil {
		logger.Error(ctx, err)
		return
	}

	authorizationUrl := c.getAuthorizationUrl(ctx)

	logger.Infof(ctx, "ChatCompletionStream Xfyun model: %s, appid: %s, getAuthorizationUrl: %s", request.Model, c.AppId, authorizationUrl)

	result := make(chan []byte)
	var conn *websocket.Conn

	_ = grpool.AddWithRecover(ctx, func(ctx context.Context) {
		conn, err = util.WebSocketClient(ctx, authorizationUrl, websocket.TextMessage, data, result, c.ProxyURL)
		if err != nil {
			logger.Error(ctx, err)
		}
	}, nil)

	duration := gtime.Now().UnixMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.Now().UnixMilli()
			logger.Infof(ctx, "ChatCompletionStream Xfyun model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		for {

			message := <-result

			sparkRes := new(model.SparkRes)
			err := gjson.Unmarshal(message, &sparkRes)
			if err != nil {
				logger.Error(ctx, err)

				err = conn.Close()
				if err != nil {
					logger.Error(ctx, err)
				}

				responseChan <- nil
				time.Sleep(time.Millisecond)
				close(responseChan)

				return
			}

			response := &model.ChatCompletionResponse{
				ID:    sparkRes.Header.Sid,
				Model: request.Model,
				Choices: []model.ChatCompletionChoice{{
					Index: sparkRes.Payload.Choices.Seq,
					Delta: openai.ChatCompletionStreamChoiceDelta{
						Role:    sparkRes.Payload.Choices.Text[0].Role,
						Content: sparkRes.Payload.Choices.Text[0].Content,
					},
				}},
				//Usage: openai.Usage{
				//	PromptTokens:     sparkRes.Payload.Usage.Text.PromptTokens,
				//	CompletionTokens: sparkRes.Payload.Usage.Text.CompletionTokens,
				//	TotalTokens:      sparkRes.Payload.Usage.Text.TotalTokens,
				//},
				ConnTime: duration - now,
			}

			if sparkRes.Header.Status == 2 {

				logger.Infof(ctx, "ChatCompletionStream Xfyun model: %s finished", request.Model)

				response.Choices[0].FinishReason = "stop"

				err = conn.Close()
				if err != nil {
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

func (c *Client) Image(ctx context.Context, request model.ImageRequest) (res model.ImageResponse, err error) {

	return
}

func (c *Client) getAuthorizationUrl(ctx context.Context) string {

	parse, err := url.Parse(c.OriginalURL + c.Path)
	if err != nil {
		logger.Error(ctx, err)
		return ""
	}

	now := gtime.Now()
	loc, _ := time.LoadLocation("GMT")
	zone, _ := now.ToZone(loc.String())
	date := zone.Layout("Mon, 02 Jan 2006 15:04:05 GMT")

	tmp := "host: " + parse.Host + "\n"
	tmp += "date: " + date + "\n"
	tmp += "GET " + parse.Path + " HTTP/1.1"

	hash := hmac.New(sha256.New, []byte(c.Secret))

	_, err = hash.Write([]byte(tmp))
	if err != nil {
		logger.Error(ctx, err)
		return ""
	}

	signature := gbase64.EncodeToString(hash.Sum(nil))

	authorizationOrigin := gbase64.EncodeToString([]byte(fmt.Sprintf("api_key=\"%s\",algorithm=\"%s\",headers=\"%s\",signature=\"%s\"", c.Key, "hmac-sha256", "host date request-line", signature)))

	wsURL := gstr.Replace(gstr.Replace(c.BaseURL+c.Path, "https://", "wss://"), "http://", "ws://")

	return fmt.Sprintf("%s?authorization=%s&date=%s&host=%s", wsURL, authorizationOrigin, gurl.RawEncode(date), parse.Host)
}
