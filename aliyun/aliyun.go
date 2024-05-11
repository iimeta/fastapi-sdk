package aliyun

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
	"github.com/sashabaranov/go-openai"
	"io"
	"time"
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
		BaseURL: "https://dashscope.aliyuncs.com/api/v1",
		Path:    "/services/aigc/text-generation/generation",
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewClient Aliyun model: %s, baseURL: %s", model, baseURL)
		client.BaseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewClient Aliyun model: %s, path: %s", model, path)
		client.Path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewClient Aliyun model: %s, proxyURL: %s", model, proxyURL[0])
		client.ProxyURL = proxyURL[0]
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

	chatCompletionReq := model.AliyunChatCompletionReq{
		Model: request.Model,
		Input: model.Input{
			Messages: request.Messages,
		},
		Parameters: model.Parameters{
			MaxTokens:         request.MaxTokens,
			Temperature:       request.Temperature,
			TopP:              request.TopP,
			TopK:              request.N,
			Stop:              request.Stop,
			RepetitionPenalty: request.FrequencyPenalty,
			Seed:              request.Seed,
			Tools:             request.Tools,
		},
	}

	if request.ResponseFormat != nil {
		chatCompletionReq.Parameters.ResultFormat = gconv.String(request.ResponseFormat.Type)
	}

	header := make(map[string]string)
	header["Authorization"] = "Bearer " + c.Key

	chatCompletionRes := new(model.AliyunChatCompletionRes)
	err = util.HttpPostJson(ctx, c.BaseURL+c.Path, header, chatCompletionReq, &chatCompletionRes, c.ProxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletion Aliyun model: %s, error: %v", request.Model, err)
		return
	}

	if chatCompletionRes.Code != "" {
		err = errors.New(gjson.MustEncodeString(chatCompletionRes))
		logger.Errorf(ctx, "ChatCompletion Aliyun model: %s, error: %v", request.Model, err)
		return
	}

	res = model.ChatCompletionResponse{
		ID:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.RequestId,
		Object:  consts.COMPLETION_OBJECT,
		Created: gtime.Now().Unix(),
		Model:   request.Model,
		Choices: []model.ChatCompletionChoice{{
			Message: &openai.ChatCompletionMessage{
				Role:    consts.ROLE_ASSISTANT,
				Content: chatCompletionRes.Output.Text,
			},
		}},
		Usage: &model.Usage{
			PromptTokens:     chatCompletionRes.Usage.InputTokens,
			CompletionTokens: chatCompletionRes.Usage.OutputTokens,
			TotalTokens:      chatCompletionRes.Usage.InputTokens + chatCompletionRes.Usage.OutputTokens,
		},
	}

	return res, nil
}

func (c *Client) ChatCompletionStream(ctx context.Context, request model.ChatCompletionRequest) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionStream Aliyun model: %s start", request.Model)

	now := gtime.Now().UnixMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionStream Aliyun model: %s totalTime: %d ms", request.Model, gtime.Now().UnixMilli()-now)
		}
	}()

	chatCompletionReq := model.AliyunChatCompletionReq{
		Model: request.Model,
		Input: model.Input{
			Messages: request.Messages,
		},
		Parameters: model.Parameters{
			ResultFormat:      "message",
			MaxTokens:         request.MaxTokens,
			Temperature:       request.Temperature,
			TopP:              request.TopP,
			TopK:              request.N,
			Stop:              request.Stop,
			RepetitionPenalty: request.FrequencyPenalty,
			Seed:              request.Seed,
			Tools:             request.Tools,
			IncrementalOutput: true,
		},
	}

	if request.ResponseFormat != nil {
		chatCompletionReq.Parameters.ResultFormat = gconv.String(request.ResponseFormat.Type)
	}

	header := make(map[string]string)
	header["Authorization"] = "Bearer " + c.Key

	stream, err := util.SSEClient(ctx, c.BaseURL+c.Path, header, chatCompletionReq, c.ProxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStream Aliyun model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	duration := gtime.Now().UnixMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.Now().UnixMilli()
			logger.Infof(ctx, "ChatCompletionStream Aliyun model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		var (
			usage   *model.Usage
			created = gtime.Now().Unix()
			id      = consts.COMPLETION_ID_PREFIX + grand.S(29)
		)

		for {

			streamResponse, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {

				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "ChatCompletionStream Aliyun model: %s, error: %v", request.Model, err)
				}

				responseChan <- &model.ChatCompletionResponse{Error: err}
				time.Sleep(time.Millisecond)
				close(responseChan)

				return
			}

			if errors.Is(err, io.EOF) {

				logger.Infof(ctx, "ChatCompletionStream Aliyun model: %s finished", request.Model)

				if err = stream.Close(); err != nil {
					logger.Errorf(ctx, "ChatCompletionStream Aliyun model: %s, stream.Close error: %v", request.Model, err)
				}

				end := gtime.Now().UnixMilli()

				responseChan <- &model.ChatCompletionResponse{
					ID:      id,
					Object:  consts.COMPLETION_STREAM_OBJECT,
					Created: created,
					Model:   request.Model,
					Choices: []model.ChatCompletionChoice{{
						Delta:        &openai.ChatCompletionStreamChoiceDelta{},
						FinishReason: openai.FinishReasonStop,
					}},
					Usage:     usage,
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
				}

				return
			}

			chatCompletionRes := new(model.AliyunChatCompletionRes)
			if err := gjson.Unmarshal(streamResponse, &chatCompletionRes); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream Aliyun model: %s, error: %v", request.Model, err)

				responseChan <- &model.ChatCompletionResponse{Error: err}
				time.Sleep(time.Millisecond)
				close(responseChan)

				return
			}

			if chatCompletionRes.Code != "" {

				err = errors.New(gjson.MustEncodeString(chatCompletionRes))
				logger.Errorf(ctx, "ChatCompletionStream Aliyun model: %s, error: %v", request.Model, err)

				if err = stream.Close(); err != nil {
					logger.Errorf(ctx, "ChatCompletionStream Aliyun model: %s, stream.Close error: %v", request.Model, err)
				}

				end := gtime.Now().UnixMilli()

				responseChan <- &model.ChatCompletionResponse{
					ID:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.RequestId,
					Object:  consts.COMPLETION_STREAM_OBJECT,
					Created: created,
					Model:   request.Model,
					Choices: []model.ChatCompletionChoice{{
						Delta: &openai.ChatCompletionStreamChoiceDelta{
							Role:    consts.ROLE_ASSISTANT,
							Content: chatCompletionRes.Message,
						},
						FinishReason: openai.FinishReasonStop,
					}},
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     err,
				}

				return
			}

			id = consts.COMPLETION_ID_PREFIX + chatCompletionRes.RequestId
			usage = &model.Usage{
				PromptTokens:     chatCompletionRes.Usage.InputTokens,
				CompletionTokens: chatCompletionRes.Usage.OutputTokens,
				TotalTokens:      chatCompletionRes.Usage.InputTokens + chatCompletionRes.Usage.OutputTokens,
			}

			response := &model.ChatCompletionResponse{
				ID:      id,
				Object:  consts.COMPLETION_STREAM_OBJECT,
				Created: created,
				Model:   request.Model,
				Choices: []model.ChatCompletionChoice{{
					Delta: &openai.ChatCompletionStreamChoiceDelta{
						Role:    consts.ROLE_ASSISTANT,
						Content: chatCompletionRes.Output.Text,
					},
				}},
				Usage:    usage,
				ConnTime: duration - now,
			}

			end := gtime.Now().UnixMilli()
			response.Duration = end - duration
			response.TotalTime = end - now

			responseChan <- response
		}
	}, nil); err != nil {
		logger.Errorf(ctx, "ChatCompletionStream Aliyun model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (c *Client) Image(ctx context.Context, request model.ImageRequest) (res model.ImageResponse, err error) {

	return
}
