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
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
	"io"
)

func (c *Client) ChatCompletionOfficial(ctx context.Context, data []byte) (res model.AnthropicChatCompletionRes, err error) {

	logger.Infof(ctx, "ChatCompletionOfficial Anthropic model: %s start", c.model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletionOfficial Anthropic model: %s totalTime: %d ms", c.model, res.TotalTime)
	}()

	request := make(map[string]interface{})
	if err = gjson.Unmarshal(data, &request); err != nil {
		logger.Errorf(ctx, "ChatCompletionOfficial Anthropic model: %s, data: %s, gjson.Unmarshal error: %v", c.model, data, err)
		return res, err
	}

	if c.isGcp {
		delete(request, "model")
	}

	if c.isAws {

		request["anthropic_version"] = "bedrock-2023-05-31"
		delete(request, "metadata")

		invokeModelInput := &bedrockruntime.InvokeModelInput{
			ModelId:     aws.String(AwsModelIDMap[gconv.String(request["model"])]),
			Accept:      aws.String("application/json"),
			ContentType: aws.String("application/json"),
		}

		delete(request, "model")

		if invokeModelInput.Body, err = gjson.Marshal(request); err != nil {
			logger.Errorf(ctx, "ChatCompletionOfficial Anthropic model: %s, request: %s, gjson.Marshal error: %v", c.model, gjson.MustEncodeString(request), err)
			return res, err
		}

		invokeModelOutput, err := c.awsClient.InvokeModel(ctx, invokeModelInput)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionOfficial Anthropic model: %s, invokeModelInput: %s, awsClient.InvokeModel error: %v", c.model, gjson.MustEncodeString(invokeModelInput), err)
			return res, err
		}

		if err = gjson.Unmarshal(invokeModelOutput.Body, &res); err != nil {
			logger.Errorf(ctx, "ChatCompletionOfficial Anthropic model: %s, invokeModelOutput.Body: %s, gjson.Unmarshal error: %v", c.model, invokeModelOutput.Body, err)
			return res, err
		}

	} else {
		if res.ResponseBytes, err = util.HttpPost(ctx, c.baseURL+c.path, c.header, request, &res, c.proxyURL); err != nil {
			logger.Errorf(ctx, "ChatCompletionOfficial Anthropic model: %s, error: %v", c.model, err)
			return res, err
		}
	}

	if res.Error != nil && res.Error.Type != "" {
		logger.Errorf(ctx, "ChatCompletionOfficial Anthropic model: %s, chatCompletionRes: %s", c.model, gjson.MustEncodeString(res))

		err = c.apiErrorHandler(&res)
		logger.Errorf(ctx, "ChatCompletionOfficial Anthropic model: %s, error: %v", c.model, err)

		return res, err
	}

	return res, nil
}

func (c *Client) ChatCompletionStreamOfficial(ctx context.Context, data []byte) (responseChan chan *model.AnthropicChatCompletionRes, err error) {

	logger.Infof(ctx, "ChatCompletionStreamOfficial Anthropic model: %s start", c.model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionStreamOfficial Anthropic model: %s totalTime: %d ms", c.model, gtime.TimestampMilli()-now)
		}
	}()

	request := make(map[string]interface{})
	if err = gjson.Unmarshal(data, &request); err != nil {
		logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, data: %s, gjson.Unmarshal error: %v", c.model, data, err)
		return responseChan, err
	}

	if c.isGcp {
		delete(request, "model")
	}

	if c.isAws {

		request["anthropic_version"] = "bedrock-2023-05-31"
		delete(request, "stream")

		invokeModelStreamInput := &bedrockruntime.InvokeModelWithResponseStreamInput{
			ModelId:     aws.String(AwsModelIDMap[gconv.String(request["model"])]),
			Accept:      aws.String("application/json"),
			ContentType: aws.String("application/json"),
		}

		delete(request, "model")

		invokeModelStreamInput.Body, err = gjson.Marshal(request)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, request: %s, gjson.Marshal error: %v", c.model, gjson.MustEncodeString(request), err)
			return
		}

		invokeModelStreamOutput, err := c.awsClient.InvokeModelWithResponseStream(ctx, invokeModelStreamInput)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, invokeModelStreamInput: %s, awsClient.InvokeModelWithResponseStream error: %v", c.model, gjson.MustEncodeString(invokeModelStreamInput), err)
			return responseChan, err
		}

		stream := invokeModelStreamOutput.GetStream()

		duration := gtime.TimestampMilli()

		responseChan = make(chan *model.AnthropicChatCompletionRes)

		if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

			defer func() {
				if err := stream.Close(); err != nil {
					logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, stream.Close error: %v", c.model, err)
				}

				end := gtime.TimestampMilli()
				logger.Infof(ctx, "ChatCompletionStreamOfficial Anthropic model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", c.model, duration-now, end-duration, end-now)
			}()

			for {

				event, ok := <-stream.Events()
				if !ok {

					if !errors.Is(err, context.Canceled) {
						logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, error: %v", c.model, err)
					}

					end := gtime.TimestampMilli()
					responseChan <- &model.AnthropicChatCompletionRes{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Err:       err,
					}

					return
				}

				var responseBytes []byte
				chatCompletionRes := new(model.AnthropicChatCompletionRes)

				switch v := event.(type) {
				case *types.ResponseStreamMemberChunk:
					responseBytes = v.Value.Bytes
					if err := gjson.Unmarshal(v.Value.Bytes, &chatCompletionRes); err != nil {
						logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, v.Value.Bytes: %s, error: %v", c.model, v.Value.Bytes, err)

						end := gtime.TimestampMilli()
						responseChan <- &model.AnthropicChatCompletionRes{
							ConnTime:  duration - now,
							Duration:  end - duration,
							TotalTime: end - now,
							Err:       errors.New(fmt.Sprintf("v.Value.Bytes: %s, error: %v", v.Value.Bytes, err)),
						}

						return
					}
				case *types.UnknownUnionMember:

					end := gtime.TimestampMilli()
					responseChan <- &model.AnthropicChatCompletionRes{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Err:       errors.New("unknown tag:" + v.Tag),
					}

					return
				default:

					end := gtime.TimestampMilli()
					responseChan <- &model.AnthropicChatCompletionRes{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Err:       errors.New("unknown type"),
					}

					return
				}

				if chatCompletionRes.Error != nil && chatCompletionRes.Error.Type != "" {
					logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, chatCompletionRes: %s", c.model, gjson.MustEncodeString(chatCompletionRes))

					err = c.apiErrorHandler(chatCompletionRes)
					logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, error: %v", c.model, err)

					end := gtime.TimestampMilli()
					responseChan <- &model.AnthropicChatCompletionRes{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Err:       err,
					}

					return
				}

				response := &model.AnthropicChatCompletionRes{
					Id:            chatCompletionRes.Id,
					Type:          chatCompletionRes.Type,
					Role:          chatCompletionRes.Role,
					Content:       chatCompletionRes.Content,
					Model:         chatCompletionRes.Model,
					StopReason:    chatCompletionRes.StopReason,
					StopSequence:  chatCompletionRes.StopSequence,
					Message:       chatCompletionRes.Message,
					Index:         chatCompletionRes.Index,
					Delta:         chatCompletionRes.Delta,
					Usage:         chatCompletionRes.Usage,
					Error:         chatCompletionRes.Error,
					ResponseBytes: responseBytes,
					ConnTime:      duration - now,
				}

				if errors.Is(err, io.EOF) || chatCompletionRes.Delta.StopReason != "" {
					logger.Infof(ctx, "ChatCompletionStreamOfficial Anthropic model: %s finished", c.model)

					end := gtime.TimestampMilli()
					response.Duration = end - duration
					response.TotalTime = end - now
					responseChan <- response

					responseChan <- &model.AnthropicChatCompletionRes{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Err:       io.EOF,
					}

					return
				}

				end := gtime.TimestampMilli()
				response.Duration = end - duration
				response.TotalTime = end - now

				responseChan <- response
			}
		}, nil); err != nil {
			logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, error: %v", c.model, err)
			return responseChan, err
		}

	} else {

		stream, err := util.SSEClient(ctx, c.baseURL+c.path, c.header, request, c.proxyURL, c.requestErrorHandler)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, error: %v", c.model, err)
			return responseChan, err
		}

		duration := gtime.TimestampMilli()

		responseChan = make(chan *model.AnthropicChatCompletionRes)

		if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

			defer func() {
				if err := stream.Close(); err != nil {
					logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, stream.Close error: %v", c.model, err)
				}

				end := gtime.TimestampMilli()
				logger.Infof(ctx, "ChatCompletionStreamOfficial Anthropic model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", c.model, duration-now, end-duration, end-now)
			}()

			for {

				streamResponse, err := stream.Recv()
				if err != nil && !errors.Is(err, io.EOF) {

					if !errors.Is(err, context.Canceled) {
						logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, error: %v", c.model, err)
					}

					end := gtime.TimestampMilli()
					responseChan <- &model.AnthropicChatCompletionRes{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Err:       err,
					}

					return
				}

				chatCompletionRes := new(model.AnthropicChatCompletionRes)
				if err := gjson.Unmarshal(streamResponse, &chatCompletionRes); err != nil {
					logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, streamResponse: %s, error: %v", c.model, streamResponse, err)

					end := gtime.TimestampMilli()
					responseChan <- &model.AnthropicChatCompletionRes{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Err:       errors.New(fmt.Sprintf("streamResponse: %s, error: %v", streamResponse, err)),
					}

					return
				}

				if chatCompletionRes.Error != nil && chatCompletionRes.Error.Type != "" {
					logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, chatCompletionRes: %s", c.model, gjson.MustEncodeString(chatCompletionRes))

					err = c.apiErrorHandler(chatCompletionRes)
					logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, error: %v", c.model, err)

					end := gtime.TimestampMilli()
					responseChan <- &model.AnthropicChatCompletionRes{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Err:       err,
					}

					return
				}

				response := &model.AnthropicChatCompletionRes{
					Id:            chatCompletionRes.Id,
					Type:          chatCompletionRes.Type,
					Role:          chatCompletionRes.Role,
					Content:       chatCompletionRes.Content,
					Model:         chatCompletionRes.Model,
					StopReason:    chatCompletionRes.StopReason,
					StopSequence:  chatCompletionRes.StopSequence,
					Message:       chatCompletionRes.Message,
					Index:         chatCompletionRes.Index,
					Delta:         chatCompletionRes.Delta,
					Usage:         chatCompletionRes.Usage,
					Error:         chatCompletionRes.Error,
					ResponseBytes: streamResponse,
					ConnTime:      duration - now,
				}

				if errors.Is(err, io.EOF) || chatCompletionRes.Delta.StopReason != "" {
					logger.Infof(ctx, "ChatCompletionStreamOfficial Anthropic model: %s finished", c.model)

					end := gtime.TimestampMilli()
					response.Duration = end - duration
					response.TotalTime = end - now
					responseChan <- response

					responseChan <- &model.AnthropicChatCompletionRes{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Err:       io.EOF,
					}

					return
				}

				end := gtime.TimestampMilli()
				response.Duration = end - duration
				response.TotalTime = end - now

				responseChan <- response
			}
		}, nil); err != nil {
			logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, error: %v", c.model, err)
			return responseChan, err
		}
	}

	return responseChan, nil
}
