package google

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/iimeta/fastapi-sdk/common"
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

	logger.Infof(ctx, "NewClient Google model: %s, key: %s", model, key)

	client := &Client{
		Key:     key,
		BaseURL: "https://generativelanguage.googleapis.com/v1beta",
		Path:    "/models/" + model,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewClient Google model: %s, baseURL: %s", model, baseURL)
		client.BaseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewClient Google model: %s, path: %s", model, path)
		client.Path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewClient Google model: %s, proxyURL: %s", model, proxyURL[0])
		client.ProxyURL = proxyURL[0]
	}

	return client
}

func (c *Client) ChatCompletion(ctx context.Context, request model.ChatCompletionRequest) (res model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletion Google model: %s start", request.Model)

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "ChatCompletion Google model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	contents := make([]model.Content, 0)
	messages := common.HandleMessages(request.Messages, false)

	for _, message := range messages {

		role := message.Role

		if role == consts.ROLE_ASSISTANT {
			role = consts.ROLE_MODEL
		}

		contents = append(contents, model.Content{
			Role: role,
			Parts: []model.Part{{
				Text: message.Content,
			}},
		})
	}

	chatCompletionReq := model.GoogleChatCompletionReq{
		Contents: contents,
	}

	chatCompletionRes := new(model.GoogleChatCompletionRes)
	err = util.HttpPostJson(ctx, fmt.Sprintf("%s:generateContent?key=%s", c.BaseURL+c.Path, c.Key), nil, chatCompletionReq, &chatCompletionRes, c.ProxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletion Google model: %s, error: %v", request.Model, err)
		return
	}

	if chatCompletionRes.Error.Code != 0 {
		err = c.handleErrorResp(chatCompletionRes)
		logger.Errorf(ctx, "ChatCompletion Google model: %s, error: %v", request.Model, err)
		return
	}

	res = model.ChatCompletionResponse{
		ID:      consts.COMPLETION_ID_PREFIX + grand.S(29),
		Object:  consts.COMPLETION_OBJECT,
		Created: gtime.Now().Unix(),
		Model:   request.Model,
		Choices: []model.ChatCompletionChoice{{
			Message: &openai.ChatCompletionMessage{
				Role:    consts.ROLE_ASSISTANT,
				Content: chatCompletionRes.Candidates[0].Content.Parts[0].Text,
			},
			FinishReason: openai.FinishReasonStop,
		}},
		Usage: &model.Usage{
			PromptTokens:     chatCompletionRes.UsageMetadata.PromptTokenCount,
			CompletionTokens: chatCompletionRes.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      chatCompletionRes.UsageMetadata.TotalTokenCount,
		},
	}

	return res, nil
}

func (c *Client) ChatCompletionStream(ctx context.Context, request model.ChatCompletionRequest) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionStream Google model: %s start", request.Model)

	now := gtime.Now().UnixMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionStream Google model: %s totalTime: %d ms", request.Model, gtime.Now().UnixMilli()-now)
		}
	}()

	contents := make([]model.Content, 0)
	messages := common.HandleMessages(request.Messages, false)

	for _, message := range messages {

		role := message.Role

		if role == consts.ROLE_ASSISTANT {
			role = consts.ROLE_MODEL
		}

		contents = append(contents, model.Content{
			Role: role,
			Parts: []model.Part{{
				Text: message.Content,
			}},
		})
	}

	chatCompletionReq := model.GoogleChatCompletionReq{
		Contents: contents,
	}

	stream, err := util.SSEClient(ctx, fmt.Sprintf("%s:streamGenerateContent?alt=sse&key=%s", c.BaseURL+c.Path, c.Key), nil, chatCompletionReq, c.ProxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStream Google model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	duration := gtime.Now().UnixMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.Now().UnixMilli()
			logger.Infof(ctx, "ChatCompletionStream Google model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
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
					logger.Errorf(ctx, "ChatCompletionStream Google model: %s, error: %v", request.Model, err)
				}

				responseChan <- &model.ChatCompletionResponse{Error: err}
				time.Sleep(time.Millisecond)
				close(responseChan)

				return
			}

			if errors.Is(err, io.EOF) {

				logger.Infof(ctx, "ChatCompletionStream Google model: %s finished", request.Model)

				if err = stream.Close(); err != nil {
					logger.Errorf(ctx, "ChatCompletionStream Google model: %s, stream.Close error: %v", request.Model, err)
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

			chatCompletionRes := new(model.GoogleChatCompletionRes)
			if err := gjson.Unmarshal(streamResponse, &chatCompletionRes); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream Google model: %s, error: %v", request.Model, err)

				responseChan <- &model.ChatCompletionResponse{Error: err}
				time.Sleep(time.Millisecond)
				close(responseChan)

				return
			}

			if chatCompletionRes.Error.Code != 0 {

				err = c.handleErrorResp(chatCompletionRes)
				logger.Errorf(ctx, "ChatCompletionStream Google model: %s, error: %v", request.Model, err)

				if err = stream.Close(); err != nil {
					logger.Errorf(ctx, "ChatCompletionStream Google model: %s, stream.Close error: %v", request.Model, err)
				}

				end := gtime.Now().UnixMilli()

				responseChan <- &model.ChatCompletionResponse{
					ID:      id,
					Object:  consts.COMPLETION_STREAM_OBJECT,
					Created: created,
					Model:   request.Model,
					Choices: []model.ChatCompletionChoice{{
						Delta: &openai.ChatCompletionStreamChoiceDelta{
							Role:    consts.ROLE_ASSISTANT,
							Content: chatCompletionRes.Error.Message,
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

			usage = &model.Usage{
				PromptTokens:     chatCompletionRes.UsageMetadata.PromptTokenCount,
				CompletionTokens: chatCompletionRes.UsageMetadata.CandidatesTokenCount,
				TotalTokens:      chatCompletionRes.UsageMetadata.TotalTokenCount,
			}

			response := &model.ChatCompletionResponse{
				ID:      id,
				Object:  consts.COMPLETION_STREAM_OBJECT,
				Created: created,
				Model:   request.Model,
				Choices: []model.ChatCompletionChoice{{
					Index: chatCompletionRes.Candidates[0].Index,
					Delta: &openai.ChatCompletionStreamChoiceDelta{
						Role:    consts.ROLE_ASSISTANT,
						Content: chatCompletionRes.Candidates[0].Content.Parts[0].Text,
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
		logger.Errorf(ctx, "ChatCompletionStream Google model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (c *Client) Image(ctx context.Context, request model.ImageRequest) (res model.ImageResponse, err error) {

	return
}

func (c *Client) handleErrorResp(response *model.GoogleChatCompletionRes) error {

	return errors.New(gjson.MustEncodeString(response))
}
