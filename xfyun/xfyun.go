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
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/gorilla/websocket"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
	"github.com/sashabaranov/go-openai"
	"io"
	"math"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	AppId       string
	Key         string
	Secret      string
	OriginalURL string
	BaseURL     string
	Path        string
	ProxyURL    string
	Domain      string
}

func NewClient(ctx context.Context, model, key, baseURL, path string, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewClient Xfyun model: %s, key: %s", model, key)

	result := gstr.Split(key, "|")

	client := &Client{
		AppId:       result[0],
		Key:         result[1],
		Secret:      result[2],
		OriginalURL: "https://spark-api.xf-yun.com",
		BaseURL:     "https://spark-api.xf-yun.com/v3.5",
		Path:        "/chat",
		Domain:      "generalv3.5",
	}

	if baseURL != "" {

		logger.Infof(ctx, "NewClient Xfyun model: %s, baseURL: %s", model, baseURL)
		client.BaseURL = baseURL

		version := baseURL[strings.LastIndex(baseURL, "/")+1:]

		switch version {
		case "v3.5":
			client.Domain = "generalv3.5"
		case "v3.1":
			client.Domain = "generalv3"
		case "v2.1":
			client.Domain = "generalv2"
		case "v1.1":
			client.Domain = "general"
		default:
			v := gconv.Float64(version[1:])
			if math.Round(v) > v {
				client.Domain = fmt.Sprintf("general%s", version)
			} else {
				client.Domain = fmt.Sprintf("generalv%0.f", math.Round(v))
			}
		}
	}

	if path != "" {
		logger.Infof(ctx, "NewClient Xfyun model: %s, path: %s", model, path)
		client.Path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewClient Xfyun model: %s, proxyURL: %s", model, proxyURL[0])
		client.ProxyURL = proxyURL[0]
	}

	return client
}

func (c *Client) ChatCompletion(ctx context.Context, request model.ChatCompletionRequest) (res model.ChatCompletionResponse, err error) {

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "ChatCompletion Xfyun model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, res.ConnTime, res.Duration, res.TotalTime)
	}()

	maxTokens := request.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	chatCompletionReq := model.XfyunChatCompletionReq{
		Header: model.Header{
			AppId: c.AppId,
			Uid:   grand.Digits(10),
		},
		Parameter: model.Parameter{
			Chat: &model.Chat{
				Domain:      c.Domain,
				MaxTokens:   maxTokens,
				Temperature: request.Temperature,
				TopK:        request.N,
				ChatId:      request.User,
			},
		},
		Payload: model.Payload{
			Message: &model.Message{
				Text: request.Messages,
			},
		},
	}

	if request.Functions != nil && len(request.Functions) > 0 {
		chatCompletionReq.Payload.Functions = new(model.Functions)
		chatCompletionReq.Payload.Functions.Text = append(chatCompletionReq.Payload.Functions.Text, request.Functions...)
	}

	data, err := gjson.Marshal(chatCompletionReq)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletion Xfyun model: %s, error: %v", request.Model, err)
		return res, err
	}

	conn, err := util.WebSocketClient(ctx, c.getAuthorizationUrl(ctx), websocket.TextMessage, data, c.ProxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletion Xfyun model: %s, error: %v", request.Model, err)
		return res, err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Errorf(ctx, "ChatCompletion Xfyun model: %s, conn.Close error: %v", request.Model, err)
		}
	}()

	duration := gtime.Now().UnixMilli()

	responseContent := ""
	chatCompletionRes := new(model.XfyunChatCompletionRes)

	for {

		message, err := conn.ReadMessage(ctx)
		if err != nil && !errors.Is(err, io.EOF) {
			logger.Errorf(ctx, "ChatCompletion Xfyun model: %s, error: %v", request.Model, err)
			return res, err
		}

		if err = gjson.Unmarshal(message, &chatCompletionRes); err != nil {
			logger.Errorf(ctx, "ChatCompletion Xfyun model: %s, error: %v", request.Model, err)
			return res, err
		}

		if chatCompletionRes.Header.Code != 0 {
			err = errors.New(gjson.MustEncodeString(chatCompletionRes))
			logger.Errorf(ctx, "ChatCompletion Xfyun model: %s, error: %v", request.Model, err)
			return res, err
		}

		responseContent += chatCompletionRes.Payload.Choices.Text[0].Content

		if chatCompletionRes.Header.Status == 2 {
			break
		}
	}

	res = model.ChatCompletionResponse{
		ID:    chatCompletionRes.Header.Sid,
		Model: request.Model,
		Choices: []model.ChatCompletionChoice{{
			Index: chatCompletionRes.Payload.Choices.Seq,
			Message: &openai.ChatCompletionMessage{
				Role:         chatCompletionRes.Payload.Choices.Text[0].Role,
				Content:      responseContent,
				FunctionCall: chatCompletionRes.Payload.Choices.Text[0].FunctionCall,
			},
		}},
		Usage: &openai.Usage{
			PromptTokens:     chatCompletionRes.Payload.Usage.Text.PromptTokens,
			CompletionTokens: chatCompletionRes.Payload.Usage.Text.CompletionTokens,
			TotalTokens:      chatCompletionRes.Payload.Usage.Text.TotalTokens,
		},
		ConnTime: duration - now,
		Duration: gtime.Now().UnixMilli() - duration,
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

	maxTokens := request.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	chatCompletionReq := model.XfyunChatCompletionReq{
		Header: model.Header{
			AppId: c.AppId,
			Uid:   grand.Digits(10),
		},
		Parameter: model.Parameter{
			Chat: &model.Chat{
				Domain:      c.Domain,
				MaxTokens:   maxTokens,
				Temperature: request.Temperature,
				TopK:        request.N,
				ChatId:      request.User,
			},
		},
		Payload: model.Payload{
			Message: &model.Message{
				Text: request.Messages,
			},
		},
	}

	if request.Functions != nil && len(request.Functions) > 0 {
		chatCompletionReq.Payload.Functions = new(model.Functions)
		chatCompletionReq.Payload.Functions.Text = append(chatCompletionReq.Payload.Functions.Text, request.Functions...)
	}

	data, err := gjson.Marshal(chatCompletionReq)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStream Xfyun model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	conn, err := util.WebSocketClient(ctx, c.getAuthorizationUrl(ctx), websocket.TextMessage, data, c.ProxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStream Xfyun model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	duration := gtime.Now().UnixMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {

			if err := conn.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream Xfyun model: %s, conn.Close error: %v", request.Model, err)
			}

			end := gtime.Now().UnixMilli()
			logger.Infof(ctx, "ChatCompletionStream Xfyun model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		for {

			message, err := conn.ReadMessage(ctx)
			if err != nil && !errors.Is(err, io.EOF) {

				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "ChatCompletionStream Xfyun model: %s, error: %v", request.Model, err)
				}

				responseChan <- nil
				time.Sleep(time.Millisecond)
				close(responseChan)

				return
			}

			chatCompletionRes := new(model.XfyunChatCompletionRes)
			if err := gjson.Unmarshal(message, &chatCompletionRes); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream Xfyun model: %s, error: %v", request.Model, err)

				responseChan <- nil
				time.Sleep(time.Millisecond)
				close(responseChan)

				return
			}

			if chatCompletionRes.Header.Code != 0 {

				err = errors.New(gjson.MustEncodeString(chatCompletionRes))
				logger.Errorf(ctx, "ChatCompletionStream Xfyun model: %s, error: %v", request.Model, err)

				end := gtime.Now().UnixMilli()

				responseChan <- &model.ChatCompletionResponse{
					ID:    chatCompletionRes.Header.Sid,
					Model: request.Model,
					Choices: []model.ChatCompletionChoice{{
						Index: 0,
						Delta: &openai.ChatCompletionStreamChoiceDelta{
							Role:    consts.ROLE_ASSISTANT,
							Content: chatCompletionRes.Header.Message,
						},
						FinishReason: openai.FinishReasonStop,
					}},
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
				}

				return
			}

			response := &model.ChatCompletionResponse{
				ID:    chatCompletionRes.Header.Sid,
				Model: request.Model,
				Choices: []model.ChatCompletionChoice{{
					Index: chatCompletionRes.Payload.Choices.Seq,
					Delta: &openai.ChatCompletionStreamChoiceDelta{
						Role:         chatCompletionRes.Payload.Choices.Text[0].Role,
						Content:      chatCompletionRes.Payload.Choices.Text[0].Content,
						FunctionCall: chatCompletionRes.Payload.Choices.Text[0].FunctionCall,
					},
				}},
				ConnTime: duration - now,
			}

			if chatCompletionRes.Payload.Usage != nil {
				response.Usage = &openai.Usage{
					PromptTokens:     chatCompletionRes.Payload.Usage.Text.PromptTokens,
					CompletionTokens: chatCompletionRes.Payload.Usage.Text.CompletionTokens,
					TotalTokens:      chatCompletionRes.Payload.Usage.Text.TotalTokens,
				}
			}

			if chatCompletionRes.Header.Status == 2 {

				logger.Infof(ctx, "ChatCompletionStream Xfyun model: %s finished", request.Model)

				response.Choices[0].FinishReason = openai.FinishReasonStop

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
		logger.Errorf(ctx, "ChatCompletionStream Xfyun model: %s, error: %v", request.Model, err)
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
		logger.Errorf(ctx, "getAuthorizationUrl Xfyun client: %+v, error: %s", c, err)
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
		logger.Errorf(ctx, "getAuthorizationUrl Xfyun client: %+v, error: %s", c, err)
		return ""
	}

	signature := gbase64.EncodeToString(hash.Sum(nil))

	authorizationOrigin := gbase64.EncodeToString([]byte(fmt.Sprintf("api_key=\"%s\",algorithm=\"%s\",headers=\"%s\",signature=\"%s\"", c.Key, "hmac-sha256", "host date request-line", signature)))

	wsURL := gstr.Replace(gstr.Replace(c.BaseURL+c.Path, "https://", "wss://"), "http://", "ws://")

	return fmt.Sprintf("%s?authorization=%s&date=%s&host=%s", wsURL, authorizationOrigin, gurl.RawEncode(date), parse.Host)
}
