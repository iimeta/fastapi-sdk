package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

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
)

func (a *Anthropic) ChatCompletionOfficial(ctx context.Context, data []byte) (res model.AnthropicChatCompletionRes, err error) {

	logger.Infof(ctx, "ChatCompletionOfficial Anthropic model: %s start", a.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletionOfficial Anthropic model: %s totalTime: %d ms", a.Model, res.TotalTime)
	}()

	request := make(map[string]interface{})
	if err = json.Unmarshal(data, &request); err != nil {
		logger.Errorf(ctx, "ChatCompletionOfficial Anthropic model: %s, data: %s, json.Unmarshal error: %v", a.Model, data, err)
		return res, err
	}

	if a.isGcp {
		delete(request, "model")
	}

	if a.isAws {

		request["anthropic_version"] = "bedrock-2023-05-31"
		delete(request, "metadata")

		invokeModelInput := &bedrockruntime.InvokeModelInput{
			ModelId:     aws.String(AwsModelIDMap[gconv.String(request["model"])]),
			Accept:      aws.String("application/json"),
			ContentType: aws.String("application/json"),
		}

		delete(request, "model")

		if invokeModelInput.Body, err = gjson.Marshal(request); err != nil {
			logger.Errorf(ctx, "ChatCompletionOfficial Anthropic model: %s, request: %s, gjson.Marshal error: %v", a.Model, gjson.MustEncodeString(request), err)
			return res, err
		}

		invokeModelOutput, err := a.awsClient.InvokeModel(ctx, invokeModelInput)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionOfficial Anthropic model: %s, invokeModelInput: %s, awsClient.InvokeModel error: %v", a.Model, gjson.MustEncodeString(invokeModelInput), err)
			return res, err
		}

		if err = json.Unmarshal(invokeModelOutput.Body, &res); err != nil {
			logger.Errorf(ctx, "ChatCompletionOfficial Anthropic model: %s, invokeModelOutput.Body: %s, json.Unmarshal error: %v", a.Model, invokeModelOutput.Body, err)
			return res, err
		}

		res.ResponseBytes = invokeModelOutput.Body

	} else {
		if res.ResponseBytes, err = util.HttpPost(ctx, a.BaseUrl+a.Path, a.header, request, &res, a.Timeout, a.ProxyUrl, a.requestErrorHandler); err != nil {
			logger.Errorf(ctx, "ChatCompletionOfficial Anthropic model: %s, error: %v", a.Model, err)
			return res, err
		}
	}

	if res.Error != nil && res.Error.Type != "" {
		logger.Errorf(ctx, "ChatCompletionOfficial Anthropic model: %s, chatCompletionRes: %s", a.Model, gjson.MustEncodeString(res))

		err = a.apiErrorHandler(&res)
		logger.Errorf(ctx, "ChatCompletionOfficial Anthropic model: %s, error: %v", a.Model, err)

		return res, err
	}

	return res, nil
}

func (a *Anthropic) ChatCompletionStreamOfficial(ctx context.Context, data []byte) (responseChan chan *model.AnthropicChatCompletionRes, err error) {

	logger.Infof(ctx, "ChatCompletionStreamOfficial Anthropic model: %s start", a.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionStreamOfficial Anthropic model: %s totalTime: %d ms", a.Model, gtime.TimestampMilli()-now)
		}
	}()

	request := make(map[string]interface{})
	if err = json.Unmarshal(data, &request); err != nil {
		logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, data: %s, json.Unmarshal error: %v", a.Model, data, err)
		return responseChan, err
	}

	if a.isGcp {
		delete(request, "model")
	}

	if a.isAws {

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
			logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, request: %s, gjson.Marshal error: %v", a.Model, gjson.MustEncodeString(request), err)
			return
		}

		invokeModelStreamOutput, err := a.awsClient.InvokeModelWithResponseStream(ctx, invokeModelStreamInput)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, invokeModelStreamInput: %s, awsClient.InvokeModelWithResponseStream error: %v", a.Model, gjson.MustEncodeString(invokeModelStreamInput), err)
			return responseChan, err
		}

		stream := invokeModelStreamOutput.GetStream()

		duration := gtime.TimestampMilli()

		responseChan = make(chan *model.AnthropicChatCompletionRes)

		if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

			defer func() {
				if err := stream.Close(); err != nil {
					logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, stream.Close error: %v", a.Model, err)
				}

				end := gtime.TimestampMilli()
				logger.Infof(ctx, "ChatCompletionStreamOfficial Anthropic model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", a.Model, duration-now, end-duration, end-now)
			}()

			for {

				event, ok := <-stream.Events()
				if !ok { // todo

					if errors.Is(err, io.EOF) {
						logger.Infof(ctx, "ChatCompletionStreamOfficial Anthropic model: %s finished", a.Model)
					} else {
						logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, error: %v", a.Model, err)
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

				var (
					responseBytes     []byte
					chatCompletionRes model.AnthropicChatCompletionRes
				)

				switch v := event.(type) {
				case *types.ResponseStreamMemberChunk:
					responseBytes = v.Value.Bytes
					if err := json.Unmarshal(v.Value.Bytes, &chatCompletionRes); err != nil {
						logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, v.Value.Bytes: %s, error: %v", a.Model, v.Value.Bytes, err)

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
					logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, chatCompletionRes: %s", a.Model, gjson.MustEncodeString(chatCompletionRes))

					err = a.apiErrorHandler(&chatCompletionRes)
					logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, error: %v", a.Model, err)

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

				if chatCompletionRes.Delta.StopReason != "" {
					logger.Infof(ctx, "ChatCompletionStreamOfficial Anthropic model: %s finished", a.Model)

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
			logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, error: %v", a.Model, err)
			return responseChan, err
		}

	} else {

		stream, err := util.SSEClient(ctx, a.BaseUrl+a.Path, a.header, gjson.MustEncode(request), a.Timeout, a.ProxyUrl, a.requestErrorHandler)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, error: %v", a.Model, err)
			return responseChan, err
		}

		duration := gtime.TimestampMilli()

		responseChan = make(chan *model.AnthropicChatCompletionRes)

		if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

			defer func() {
				if err := stream.Close(); err != nil {
					logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, stream.Close error: %v", a.Model, err)
				}

				end := gtime.TimestampMilli()
				logger.Infof(ctx, "ChatCompletionStreamOfficial Anthropic model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", a.Model, duration-now, end-duration, end-now)
			}()

			for {

				streamResponse, err := stream.Recv()
				if err != nil {

					if errors.Is(err, io.EOF) {
						logger.Infof(ctx, "ChatCompletionStreamOfficial Anthropic model: %s finished", a.Model)
					} else {
						logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, error: %v", a.Model, err)
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

				chatCompletionRes := model.AnthropicChatCompletionRes{}
				if err := json.Unmarshal(streamResponse, &chatCompletionRes); err != nil {
					logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, streamResponse: %s, error: %v", a.Model, streamResponse, err)

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
					logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, chatCompletionRes: %s", a.Model, gjson.MustEncodeString(chatCompletionRes))

					err = a.apiErrorHandler(&chatCompletionRes)
					logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, error: %v", a.Model, err)

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

				if chatCompletionRes.Delta.StopReason != "" {
					logger.Infof(ctx, "ChatCompletionStreamOfficial Anthropic model: %s finished", a.Model)

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
			logger.Errorf(ctx, "ChatCompletionStreamOfficial Anthropic model: %s, error: %v", a.Model, err)
			return responseChan, err
		}
	}

	return responseChan, nil
}
