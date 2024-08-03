package anthropic

import (
	"context"
	"encoding/json"
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
	"github.com/iimeta/go-openai"
	"io"
)

type Client struct {
	key                 string
	baseURL             string
	path                string
	proxyURL            string
	isSupportSystemRole *bool
}

func NewClient(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole *bool, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewClient Anthropic model: %s, key: %s", model, key)

	client := &Client{
		key:                 key,
		baseURL:             "https://api.anthropic.com/v1",
		path:                "/messages",
		isSupportSystemRole: isSupportSystemRole,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewClient Anthropic model: %s, baseURL: %s", model, baseURL)
		client.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewClient Anthropic model: %s, path: %s", model, path)
		client.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewClient Anthropic model: %s, proxyURL: %s", model, proxyURL[0])
		client.proxyURL = proxyURL[0]
	}

	return client
}

func (c *Client) ChatCompletion(ctx context.Context, request model.ChatCompletionRequest) (res model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletion Anthropic model: %s start", request.Model)

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "ChatCompletion Anthropic model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	var messages []model.ChatCompletionMessage
	if c.isSupportSystemRole != nil {
		messages = common.HandleMessages(request.Messages, *c.isSupportSystemRole)
	} else {
		messages = common.HandleMessages(request.Messages, true)
	}

	chatCompletionReq := model.AnthropicChatCompletionReq{
		Model:         request.Model,
		Messages:      messages,
		MaxTokens:     request.MaxTokens,
		StopSequences: request.Stop,
		Stream:        request.Stream,
		Temperature:   request.Temperature,
		ToolChoice:    request.ToolChoice,
		TopK:          request.TopK,
		TopP:          request.TopP,
	}

	if messages[0].Role == consts.ROLE_SYSTEM {
		chatCompletionReq.System = gconv.String(messages[0].Content)
		messages = messages[1:]
	}

	chatCompletionReq.Metadata.UserId = request.User

	for _, tool := range request.Tools {
		chatCompletionReq.Tools = append(chatCompletionReq.Tools, model.AnthropicTool{
			Name:        tool.Function.Name,
			Description: tool.Function.Description,
			InputSchema: tool.Function.Parameters,
		})
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

	header := make(map[string]string)
	header["x-api-key"] = c.key
	header["anthropic-version"] = "2023-06-01"

	chatCompletionRes := new(model.AnthropicChatCompletionRes)
	if err = util.HttpPost(ctx, c.baseURL+c.path, header, chatCompletionReq, &chatCompletionRes, c.proxyURL); err != nil {
		logger.Errorf(ctx, "ChatCompletion Anthropic model: %s, error: %v", request.Model, err)
		return
	}

	if chatCompletionRes.Error.Code != "" && chatCompletionRes.Error.Code != "200" {
		logger.Errorf(ctx, "ChatCompletion Anthropic model: %s, chatCompletionRes: %s", request.Model, gjson.MustEncodeString(chatCompletionRes))

		err = c.apiErrorHandler(chatCompletionRes)
		logger.Errorf(ctx, "ChatCompletion Anthropic model: %s, error: %v", request.Model, err)

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

	logger.Infof(ctx, "ChatCompletionStream Anthropic model: %s start", request.Model)

	now := gtime.Now().UnixMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionStream Anthropic model: %s totalTime: %d ms", request.Model, gtime.Now().UnixMilli()-now)
		}
	}()

	var messages []model.ChatCompletionMessage
	if c.isSupportSystemRole != nil {
		messages = common.HandleMessages(request.Messages, *c.isSupportSystemRole)
	} else {
		messages = common.HandleMessages(request.Messages, true)
	}

	chatCompletionReq := model.AnthropicChatCompletionReq{
		Model:       request.Model,
		Messages:    messages,
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

	header := make(map[string]string)
	header["x-api-key"] = c.key
	header["anthropic-version"] = "2023-06-01"

	stream, err := util.SSEClient(ctx, c.baseURL+c.path, header, chatCompletionReq, c.proxyURL, c.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStream Anthropic model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	duration := gtime.Now().UnixMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream Anthropic model: %s, stream.Close error: %v", request.Model, err)
			}

			end := gtime.Now().UnixMilli()
			logger.Infof(ctx, "ChatCompletionStream Anthropic model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		for {

			streamResponse, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {

				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "ChatCompletionStream Anthropic model: %s, error: %v", request.Model, err)
				}

				end := gtime.Now().UnixMilli()
				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     err,
				}

				return
			}

			chatCompletionRes := new(model.AnthropicChatCompletionRes)
			if err := gjson.Unmarshal(streamResponse, &chatCompletionRes); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream Anthropic model: %s, streamResponse: %s, error: %v", request.Model, streamResponse, err)

				end := gtime.Now().UnixMilli()
				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     errors.New(fmt.Sprintf("streamResponse: %s, error: %v", streamResponse, err)),
				}

				return
			}

			if chatCompletionRes.Error.Code != "" && chatCompletionRes.Error.Code != "200" {
				logger.Errorf(ctx, "ChatCompletionStream Anthropic model: %s, chatCompletionRes: %s", request.Model, gjson.MustEncodeString(chatCompletionRes))

				err = c.apiErrorHandler(chatCompletionRes)
				logger.Errorf(ctx, "ChatCompletionStream Anthropic model: %s, error: %v", request.Model, err)

				end := gtime.Now().UnixMilli()
				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     err,
				}

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
				logger.Infof(ctx, "ChatCompletionStream Anthropic model: %s finished", request.Model)

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
		logger.Errorf(ctx, "ChatCompletionStream Anthropic model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (c *Client) Image(ctx context.Context, request model.ImageRequest) (res model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) requestErrorHandler(ctx context.Context, response *gclient.Response) error {

	errRes := model.AnthropicErrorResponse{}
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

	return sdkerr.NewRequestError(500, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, gjson.MustEncodeString(errRes.Error))))
}

func (c *Client) apiErrorHandler(response *model.AnthropicChatCompletionRes) error {

	switch response.Error.Code {
	case "1261":
		return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
	case "1113":
		return sdkerr.ERR_INSUFFICIENT_QUOTA
	}

	return sdkerr.NewApiError(500, response.Error.Code, gjson.MustEncodeString(response), "api_error", "")
}
