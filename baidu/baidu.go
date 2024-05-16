package baidu

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/sdkerr"
	"github.com/iimeta/fastapi-sdk/util"
	"github.com/sashabaranov/go-openai"
	"io"
	"time"
)

type Client struct {
	AccessToken string
	BaseURL     string
	Path        string
	ProxyURL    string
}

func NewClient(ctx context.Context, model, key, baseURL, path string, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewClient Baidu model: %s, key: %s", model, key)

	client := &Client{
		AccessToken: key,
		BaseURL:     "https://aip.baidubce.com/rpc/2.0/ai_custom/v1",
		Path:        "/wenxinworkshop/chat/completions_pro",
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewClient Baidu model: %s, baseURL: %s", model, baseURL)
		client.BaseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewClient Baidu model: %s, path: %s", model, path)
		client.Path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewClient Baidu model: %s, proxyURL: %s", model, proxyURL[0])
		client.ProxyURL = proxyURL[0]
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

	messages := common.HandleMessages(request.Messages, true)

	chatCompletionReq := model.BaiduChatCompletionReq{
		Messages:        messages,
		MaxOutputTokens: request.MaxTokens,
		Temperature:     request.Temperature,
		TopP:            request.TopP,
		Stream:          request.Stream,
		Stop:            request.Stop,
		PenaltyScore:    request.FrequencyPenalty,
		UserId:          request.User,
	}

	if chatCompletionReq.Messages[0].Role == consts.ROLE_SYSTEM {
		chatCompletionReq.System = chatCompletionReq.Messages[0].Content
		chatCompletionReq.Messages = chatCompletionReq.Messages[1:]
	}

	if request.ResponseFormat != nil {
		chatCompletionReq.ResponseFormat = gconv.String(request.ResponseFormat.Type)
	}

	chatCompletionRes := new(model.BaiduChatCompletionRes)
	err = util.HttpPost(ctx, fmt.Sprintf("%s?access_token=%s", c.BaseURL+c.Path, c.AccessToken), nil, chatCompletionReq, &chatCompletionRes, c.ProxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletion Baidu model: %s, error: %v", request.Model, err)
		return
	}

	if chatCompletionRes.ErrorCode != 0 {
		logger.Errorf(ctx, "ChatCompletion Baidu model: %s, chatCompletionRes: %s", request.Model, gjson.MustEncodeString(chatCompletionRes))

		err = c.apiErrorHandler(chatCompletionRes)
		logger.Errorf(ctx, "ChatCompletion Baidu model: %s, error: %v", request.Model, err)

		return
	}

	res = model.ChatCompletionResponse{
		ID:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Id,
		Object:  consts.COMPLETION_OBJECT,
		Created: chatCompletionRes.Created,
		Model:   request.Model,
		Choices: []model.ChatCompletionChoice{{
			Message: &openai.ChatCompletionMessage{
				Role:    consts.ROLE_ASSISTANT,
				Content: chatCompletionRes.Result,
			},
		}},
		Usage: chatCompletionRes.Usage,
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

	messages := common.HandleMessages(request.Messages, true)

	chatCompletionReq := model.BaiduChatCompletionReq{
		Messages:        messages,
		MaxOutputTokens: request.MaxTokens,
		Temperature:     request.Temperature,
		TopP:            request.TopP,
		Stream:          request.Stream,
		Stop:            request.Stop,
		PenaltyScore:    request.FrequencyPenalty,
		UserId:          request.User,
	}

	if chatCompletionReq.Messages[0].Role == consts.ROLE_SYSTEM {
		chatCompletionReq.System = chatCompletionReq.Messages[0].Content
		chatCompletionReq.Messages = chatCompletionReq.Messages[1:]
	}

	if request.ResponseFormat != nil {
		chatCompletionReq.ResponseFormat = gconv.String(request.ResponseFormat.Type)
	}

	stream, err := util.SSEClient(ctx, fmt.Sprintf("%s?access_token=%s", c.BaseURL+c.Path, c.AccessToken), nil, chatCompletionReq, c.ProxyURL, c.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStream Baidu model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	duration := gtime.Now().UnixMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream Baidu model: %s, stream.Close error: %v", request.Model, err)
			}

			end := gtime.Now().UnixMilli()
			logger.Infof(ctx, "ChatCompletionStream Baidu model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		for {

			streamResponse, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {

				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "ChatCompletionStream Baidu model: %s, error: %v", request.Model, err)
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

			chatCompletionRes := new(model.BaiduChatCompletionRes)
			if err = gjson.Unmarshal(streamResponse, &chatCompletionRes); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream Baidu model: %s, error: %v", request.Model, err)

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

			if chatCompletionRes.ErrorCode != 0 {
				logger.Errorf(ctx, "ChatCompletionStream Baidu model: %s, chatCompletionRes: %s", request.Model, gjson.MustEncodeString(chatCompletionRes))

				err = c.apiErrorHandler(chatCompletionRes)
				logger.Errorf(ctx, "ChatCompletionStream Baidu model: %s, error: %v", request.Model, err)

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
				ID:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Id,
				Object:  consts.COMPLETION_STREAM_OBJECT,
				Created: chatCompletionRes.Created,
				Model:   request.Model,
				Choices: []model.ChatCompletionChoice{{
					Index: chatCompletionRes.SentenceId,
					Delta: &openai.ChatCompletionStreamChoiceDelta{
						Role:    consts.ROLE_ASSISTANT,
						Content: chatCompletionRes.Result,
					},
				}},
				Usage:    chatCompletionRes.Usage,
				ConnTime: duration - now,
			}

			if errors.Is(err, io.EOF) || chatCompletionRes.IsEnd {
				logger.Infof(ctx, "ChatCompletionStream Baidu model: %s finished", request.Model)

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
		logger.Errorf(ctx, "ChatCompletionStream Baidu model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (c *Client) Image(ctx context.Context, request model.ImageRequest) (res model.ImageResponse, err error) {

	return
}

func (c *Client) requestErrorHandler(ctx context.Context, response *gclient.Response) (err error) {
	return errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, response.ReadAllString()))
}

func (c *Client) apiErrorHandler(response *model.BaiduChatCompletionRes) error {

	switch response.ErrorCode {
	case 336103, 336007:
		return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
	}

	return errors.New(gjson.MustEncodeString(response))
}
