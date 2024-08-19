package anthropic

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
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
	header              map[string]string
	isGcp               bool
	isAws               bool
	awsClient           *bedrockruntime.Client
}

// https://docs.aws.amazon.com/bedrock/latest/userguide/model-ids.html
var AwsModelIDMap = map[string]string{
	"claude-2.0":                 "anthropic.claude-v2",
	"claude-2.1":                 "anthropic.claude-v2:1",
	"claude-3-sonnet-20240229":   "anthropic.claude-3-sonnet-20240229-v1:0",
	"claude-3-5-sonnet-20240620": "anthropic.claude-3-5-sonnet-20240620-v1:0",
	"claude-3-haiku-20240307":    "anthropic.claude-3-haiku-20240307-v1:0",
	"claude-3-opus-20240229":     "anthropic.claude-3-opus-20240229-v1:0",
	"claude-instant-1.2":         "anthropic.claude-instant-v1",
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

	client.header = make(map[string]string)
	client.header["x-api-key"] = key
	client.header["anthropic-version"] = "2023-06-01"
	client.header["anthropic-beta"] = "prompt-caching-2024-07-31"

	return client
}

func NewGcpClient(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole *bool, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewGcpClient Anthropic model: %s, key: %s", model, key)

	client := &Client{
		key:                 key,
		baseURL:             "https://us-east5-aiplatform.googleapis.com/v1",
		path:                "/projects/%s/locations/us-east5/publishers/anthropic/models/%s:streamRawPredict",
		isSupportSystemRole: isSupportSystemRole,
		isGcp:               true,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewGcpClient Anthropic model: %s, baseURL: %s", model, baseURL)
		client.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewGcpClient Anthropic model: %s, path: %s", model, path)
		client.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewGcpClient Anthropic model: %s, proxyURL: %s", model, proxyURL[0])
		client.proxyURL = proxyURL[0]
	}

	client.header = make(map[string]string)
	client.header["Authorization"] = "Bearer " + key

	return client
}

func NewAwsClient(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole *bool, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewAwsClient Anthropic model: %s, key: %s", model, key)

	result := gstr.Split(key, "|")

	client := &Client{
		isSupportSystemRole: isSupportSystemRole,
		isAws:               true,
		awsClient: bedrockruntime.New(bedrockruntime.Options{
			Region:      result[0],
			Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(result[1], result[2], "")),
		}),
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewAwsClient Anthropic model: %s, baseURL: %s", model, baseURL)
		client.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewAwsClient Anthropic model: %s, path: %s", model, path)
		client.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAwsClient Anthropic model: %s, proxyURL: %s", model, proxyURL[0])
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
		Model:            request.Model,
		Messages:         messages,
		MaxTokens:        request.MaxTokens,
		StopSequences:    request.Stop,
		Stream:           request.Stream,
		Temperature:      request.Temperature,
		ToolChoice:       request.ToolChoice,
		TopK:             request.TopK,
		TopP:             request.TopP,
		AnthropicVersion: "vertex-2023-10-16",
	}

	if chatCompletionReq.Messages[0].Role == consts.ROLE_SYSTEM {
		chatCompletionReq.System = chatCompletionReq.Messages[0].Content
		chatCompletionReq.Messages = chatCompletionReq.Messages[1:]
	}

	if request.User != "" {
		chatCompletionReq.Metadata = &model.Metadata{
			UserId: request.User,
		}
	}

	for _, tool := range request.Tools {
		chatCompletionReq.Tools = append(chatCompletionReq.Tools, model.AnthropicTool{
			Name:        tool.Function.Name,
			Description: tool.Function.Description,
			InputSchema: tool.Function.Parameters,
		})
	}

	if chatCompletionReq.MaxTokens == 0 {
		chatCompletionReq.MaxTokens = 4096
	}

	if c.isGcp {
		chatCompletionReq.Model = ""
	}

	chatCompletionRes := new(model.AnthropicChatCompletionRes)

	if c.isAws {

		chatCompletionReq.AnthropicVersion = "bedrock-2023-05-31"

		invokeModelInput := &bedrockruntime.InvokeModelInput{
			ModelId:     aws.String(AwsModelIDMap[chatCompletionReq.Model]),
			Accept:      aws.String("application/json"),
			ContentType: aws.String("application/json"),
		}
		chatCompletionReq.Model = ""
		invokeModelInput.Body, err = gjson.Marshal(chatCompletionReq)
		if err != nil {
			logger.Error(ctx, err)
			return
		}

		var invokeModelOutput *bedrockruntime.InvokeModelOutput
		invokeModelOutput, err = c.awsClient.InvokeModel(ctx, invokeModelInput)
		if err != nil {
			logger.Error(ctx, err)
			return
		}

		if err = gjson.Unmarshal(invokeModelOutput.Body, &chatCompletionRes); err != nil {
			logger.Error(ctx, err)
			return
		}

	} else {
		if err = util.HttpPost(ctx, c.baseURL+c.path, c.header, chatCompletionReq, &chatCompletionRes, c.proxyURL); err != nil {
			logger.Errorf(ctx, "ChatCompletion Anthropic model: %s, error: %v", request.Model, err)
			return
		}
	}

	if chatCompletionRes.Error != nil && chatCompletionRes.Error.Type != "" {
		logger.Errorf(ctx, "ChatCompletion Anthropic model: %s, chatCompletionRes: %s", request.Model, gjson.MustEncodeString(chatCompletionRes))

		err = c.apiErrorHandler(chatCompletionRes)
		logger.Errorf(ctx, "ChatCompletion Anthropic model: %s, error: %v", request.Model, err)

		return
	}

	res = model.ChatCompletionResponse{
		ID:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Id,
		Object:  consts.COMPLETION_OBJECT,
		Created: gtime.Now().Unix(),
		Model:   request.Model,
		Usage: &model.Usage{
			PromptTokens:     chatCompletionRes.Usage.InputTokens,
			CompletionTokens: chatCompletionRes.Usage.OutputTokens,
			TotalTokens:      chatCompletionRes.Usage.InputTokens + chatCompletionRes.Usage.OutputTokens,
		},
	}

	for _, content := range chatCompletionRes.Content {
		if content.Type == consts.DELTA_TYPE_INPUT_JSON {
			res.Choices = append(res.Choices, model.ChatCompletionChoice{
				Delta: &openai.ChatCompletionStreamChoiceDelta{
					Role: consts.ROLE_ASSISTANT,
					ToolCalls: []openai.ToolCall{{
						Function: openai.FunctionCall{
							Arguments: content.PartialJson,
						},
					}},
				},
			})
		} else {
			res.Choices = append(res.Choices, model.ChatCompletionChoice{
				Message: &openai.ChatCompletionMessage{
					Role:    chatCompletionRes.Role,
					Content: content.Text,
				},
				FinishReason: "stop",
			})
		}
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
		Model:            request.Model,
		Messages:         messages,
		MaxTokens:        request.MaxTokens,
		StopSequences:    request.Stop,
		Stream:           request.Stream,
		Temperature:      request.Temperature,
		ToolChoice:       request.ToolChoice,
		TopK:             request.TopK,
		TopP:             request.TopP,
		AnthropicVersion: "vertex-2023-10-16",
	}

	if chatCompletionReq.Messages[0].Role == consts.ROLE_SYSTEM {
		chatCompletionReq.System = chatCompletionReq.Messages[0].Content
		chatCompletionReq.Messages = chatCompletionReq.Messages[1:]
	}

	if request.User != "" {
		chatCompletionReq.Metadata = &model.Metadata{
			UserId: request.User,
		}
	}

	for _, tool := range request.Tools {
		chatCompletionReq.Tools = append(chatCompletionReq.Tools, model.AnthropicTool{
			Name:        tool.Function.Name,
			Description: tool.Function.Description,
			InputSchema: tool.Function.Parameters,
		})
	}

	if chatCompletionReq.MaxTokens == 0 {
		chatCompletionReq.MaxTokens = 4096
	}

	if c.isGcp {
		chatCompletionReq.Model = ""
	}

	if c.isAws {

		chatCompletionReq.AnthropicVersion = "bedrock-2023-05-31"
		chatCompletionReq.Stream = false

		invokeModelStreamInput := &bedrockruntime.InvokeModelWithResponseStreamInput{
			ModelId:     aws.String(AwsModelIDMap[chatCompletionReq.Model]),
			Accept:      aws.String("application/json"),
			ContentType: aws.String("application/json"),
		}
		chatCompletionReq.Model = ""
		invokeModelStreamInput.Body, err = gjson.Marshal(chatCompletionReq)
		if err != nil {
			logger.Error(ctx, err)
			return
		}

		var invokeModelStreamOutput *bedrockruntime.InvokeModelWithResponseStreamOutput
		invokeModelStreamOutput, err = c.awsClient.InvokeModelWithResponseStream(ctx, invokeModelStreamInput)
		if err != nil {
			logger.Error(ctx, err)
			return
		}

		stream := invokeModelStreamOutput.GetStream()

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

				event, ok := <-stream.Events()
				if !ok {

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
				switch v := event.(type) {
				case *types.ResponseStreamMemberChunk:
					if err := gjson.Unmarshal(v.Value.Bytes, &chatCompletionRes); err != nil {
						logger.Errorf(ctx, "ChatCompletionStream Anthropic model: %s, v.Value.Bytes: %s, error: %v", request.Model, v.Value.Bytes, err)

						end := gtime.Now().UnixMilli()
						responseChan <- &model.ChatCompletionResponse{
							ConnTime:  duration - now,
							Duration:  end - duration,
							TotalTime: end - now,
							Error:     errors.New(fmt.Sprintf("v.Value.Bytes: %s, error: %v", v.Value.Bytes, err)),
						}

						return
					}
				case *types.UnknownUnionMember:

					end := gtime.Now().UnixMilli()
					responseChan <- &model.ChatCompletionResponse{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Error:     errors.New("unknown tag:" + v.Tag),
					}

					return
				default:

					end := gtime.Now().UnixMilli()
					responseChan <- &model.ChatCompletionResponse{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Error:     errors.New("unknown type"),
					}

					return
				}

				if chatCompletionRes.Error != nil && chatCompletionRes.Error.Type != "" {
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
					Created:  gtime.Now().Unix(),
					Model:    request.Model,
					ConnTime: duration - now,
				}

				if chatCompletionRes.Usage != nil {
					response.Usage = &model.Usage{
						PromptTokens:     chatCompletionRes.Usage.InputTokens,
						CompletionTokens: chatCompletionRes.Usage.OutputTokens,
						TotalTokens:      chatCompletionRes.Usage.InputTokens + chatCompletionRes.Usage.OutputTokens,
					}
				}

				if chatCompletionRes.Message.Usage != nil {
					response.Usage = &model.Usage{
						PromptTokens: chatCompletionRes.Message.Usage.InputTokens,
					}
				}

				if chatCompletionRes.Delta.StopReason != "" {
					response.Choices = append(response.Choices, model.ChatCompletionChoice{
						FinishReason: openai.FinishReasonStop,
					})
				} else {
					if chatCompletionRes.Delta.Type == consts.DELTA_TYPE_INPUT_JSON {
						response.Choices = append(response.Choices, model.ChatCompletionChoice{
							Delta: &openai.ChatCompletionStreamChoiceDelta{
								Role: consts.ROLE_ASSISTANT,
								ToolCalls: []openai.ToolCall{{
									Function: openai.FunctionCall{
										Arguments: chatCompletionRes.Delta.PartialJson,
									},
								}},
							},
						})
					} else {
						response.Choices = append(response.Choices, model.ChatCompletionChoice{
							Delta: &openai.ChatCompletionStreamChoiceDelta{
								Role:    consts.ROLE_ASSISTANT,
								Content: chatCompletionRes.Delta.Text,
							},
						})
					}
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

					responseChan <- &model.ChatCompletionResponse{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Error:     io.EOF,
					}

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

	} else {

		stream, err := util.SSEClient(ctx, c.baseURL+c.path, c.header, chatCompletionReq, c.proxyURL, c.requestErrorHandler)
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

				if chatCompletionRes.Error != nil && chatCompletionRes.Error.Type != "" {
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
					Created:  gtime.Now().Unix(),
					Model:    request.Model,
					ConnTime: duration - now,
				}

				if chatCompletionRes.Usage != nil {
					response.Usage = &model.Usage{
						PromptTokens:     chatCompletionRes.Usage.InputTokens,
						CompletionTokens: chatCompletionRes.Usage.OutputTokens,
						TotalTokens:      chatCompletionRes.Usage.InputTokens + chatCompletionRes.Usage.OutputTokens,
					}
				}

				if chatCompletionRes.Message.Usage != nil {
					response.Usage = &model.Usage{
						PromptTokens: chatCompletionRes.Message.Usage.InputTokens,
					}
				}

				if chatCompletionRes.Delta.StopReason != "" {
					response.Choices = append(response.Choices, model.ChatCompletionChoice{
						FinishReason: openai.FinishReasonStop,
					})
				} else {
					if chatCompletionRes.Delta.Type == consts.DELTA_TYPE_INPUT_JSON {
						response.Choices = append(response.Choices, model.ChatCompletionChoice{
							Delta: &openai.ChatCompletionStreamChoiceDelta{
								Role: consts.ROLE_ASSISTANT,
								ToolCalls: []openai.ToolCall{{
									Function: openai.FunctionCall{
										Arguments: chatCompletionRes.Delta.PartialJson,
									},
								}},
							},
						})
					} else {
						response.Choices = append(response.Choices, model.ChatCompletionChoice{
							Delta: &openai.ChatCompletionStreamChoiceDelta{
								Role:    consts.ROLE_ASSISTANT,
								Content: chatCompletionRes.Delta.Text,
							},
						})
					}
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

					responseChan <- &model.ChatCompletionResponse{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Error:     io.EOF,
					}

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
	}

	return responseChan, nil
}

func (c *Client) Image(ctx context.Context, request model.ImageRequest) (res model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) requestErrorHandler(ctx context.Context, response *gclient.Response) error {

	bytes := response.ReadAll()

	errRes := model.AnthropicErrorResponse{}
	if err := gjson.Unmarshal(bytes, &errRes); err != nil || errRes.Error == nil {

		reqErr := &sdkerr.RequestError{
			HttpStatusCode: response.StatusCode,
			Err:            errors.New(fmt.Sprintf("response: %s, err: %v", bytes, err)),
		}

		if errRes.Error != nil {
			reqErr.Err = errors.New(gjson.MustEncodeString(errRes.Error))
		}

		return reqErr
	}

	switch errRes.Error.Type {
	}

	return sdkerr.NewRequestError(500, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, gjson.MustEncodeString(errRes.Error))))
}

func (c *Client) apiErrorHandler(response *model.AnthropicChatCompletionRes) error {

	switch response.Error.Type {
	}

	return sdkerr.NewApiError(500, response.Error.Type, gjson.MustEncodeString(response), "api_error", "")
}
