package anthropic

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
	"github.com/iimeta/go-openai"
	"io"
)

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

	for _, message := range messages {

		if contents, ok := message.Content.([]interface{}); ok {

			for _, value := range contents {

				if content, ok := value.(map[string]interface{}); ok {

					if content["type"] == "image_url" {

						if imageUrl, ok := content["image_url"].(map[string]interface{}); ok {

							url := gconv.String(imageUrl["url"])

							if gstr.HasPrefix(url, "data:image/") {
								base64 := gstr.Split(url, "base64,")
								if len(base64) > 1 {
									// data:image/jpeg;base64,
									mimeType := fmt.Sprintf("image/%s", gstr.Split(base64[0][11:], ";")[0])
									content["source"] = model.Source{
										Type:      "base64",
										MediaType: mimeType,
										Data:      base64[1],
									}
								} else {
									content["source"] = model.Source{
										Type:      "base64",
										MediaType: "image/jpeg",
										Data:      base64[0],
									}
								}
							} else {
								content["source"] = model.Source{
									Type:      "base64",
									MediaType: "image/jpeg",
									Data:      url,
								}
							}

							content["type"] = "image"
							delete(content, "image_url")
						}
					}
				}
			}
		}
	}

	if c.isGcp {
		chatCompletionReq.Model = ""
	}

	chatCompletionRes := new(model.AnthropicChatCompletionRes)

	if c.isAws {

		chatCompletionReq.AnthropicVersion = "bedrock-2023-05-31"
		chatCompletionReq.Metadata = nil

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

	for _, message := range messages {

		if contents, ok := message.Content.([]interface{}); ok {

			for _, value := range contents {

				if content, ok := value.(map[string]interface{}); ok {

					if content["type"] == "image_url" {

						if imageUrl, ok := content["image_url"].(map[string]interface{}); ok {

							url := gconv.String(imageUrl["url"])

							if gstr.HasPrefix(url, "data:image/") {
								base64 := gstr.Split(url, "base64,")
								if len(base64) > 1 {
									// data:image/jpeg;base64,
									mimeType := fmt.Sprintf("image/%s", gstr.Split(base64[0][11:], ";")[0])
									content["source"] = model.Source{
										Type:      "base64",
										MediaType: mimeType,
										Data:      base64[1],
									}
								} else {
									content["source"] = model.Source{
										Type:      "base64",
										MediaType: "image/jpeg",
										Data:      base64[0],
									}
								}
							} else {
								content["source"] = model.Source{
									Type:      "base64",
									MediaType: "image/jpeg",
									Data:      url,
								}
							}

							content["type"] = "image"
							delete(content, "image_url")
						}
					}
				}
			}
		}
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

			var id string

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

				if chatCompletionRes.Message.Id != "" {
					id = chatCompletionRes.Message.Id
				}

				response := &model.ChatCompletionResponse{
					ID:       consts.COMPLETION_ID_PREFIX + id,
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

			var id string

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

				if chatCompletionRes.Message.Id != "" {
					id = chatCompletionRes.Message.Id
				}

				response := &model.ChatCompletionResponse{
					ID:       consts.COMPLETION_ID_PREFIX + id,
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
