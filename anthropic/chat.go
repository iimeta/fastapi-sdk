package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (a *Anthropic) ChatCompletions(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions Anthropic model: %s start", a.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions Anthropic model: %s totalTime: %d ms", a.Model, response.TotalTime)
	}()

	request, err := a.ConvChatCompletionsRequestOfficial(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions Anthropic ConvChatCompletionsRequestOfficial error: %v", err)
		return response, err
	}

	var bytes []byte

	if a.isAws {

		chatCompletionReq := model.AnthropicChatCompletionReq{}
		if err = json.Unmarshal(request, &chatCompletionReq); err != nil {
			logger.Errorf(ctx, "ChatCompletions Anthropic model: %s, request: %s, json.Unmarshal error: %v", a.Model, request, err)
			return response, err
		}

		chatCompletionReq.AnthropicVersion = "bedrock-2023-05-31"
		chatCompletionReq.Metadata = nil

		invokeModelInput := &bedrockruntime.InvokeModelInput{
			ModelId:     aws.String(chatCompletionReq.Model),
			Accept:      aws.String("application/json"),
			ContentType: aws.String("application/json"),
		}

		if modelId, ok := AwsModelIDMap[chatCompletionReq.Model]; ok {
			invokeModelInput.ModelId = aws.String(modelId)
		}

		chatCompletionReq.Model = ""

		if invokeModelInput.Body, err = gjson.Marshal(chatCompletionReq); err != nil {
			logger.Errorf(ctx, "ChatCompletions Anthropic model: %s, chatCompletionReq: %s, gjson.Marshal error: %v", a.Model, gjson.MustEncodeString(chatCompletionReq), err)
			return response, err
		}

		invokeModelOutput, err := a.awsClient.InvokeModel(ctx, invokeModelInput)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletions Anthropic model: %s, invokeModelInput: %s, awsClient.InvokeModel error: %v", a.Model, gjson.MustEncodeString(invokeModelInput), err)
			return response, err
		}

		bytes = invokeModelOutput.Body

	} else {
		if bytes, err = util.HttpPost(ctx, a.BaseUrl+a.Path, a.header, request, nil, a.Timeout, a.ProxyUrl, a.requestErrorHandler); err != nil {
			logger.Errorf(ctx, "ChatCompletions Anthropic model: %s, error: %v", a.Model, err)
			return response, err
		}
	}

	if response, err = a.ConvChatCompletionsResponseOfficial(ctx, bytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions Anthropic ConvChatCompletionsResponseOfficial error: %v", err)
		return response, err
	}

	return response, nil
}

func (a *Anthropic) ChatCompletionsStream(ctx context.Context, data []byte) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream Anthropic model: %s start", a.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream Anthropic model: %s totalTime: %d ms", a.Model, gtime.TimestampMilli()-now)
		}
	}()

	request, err := a.ConvChatCompletionsRequestOfficial(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream Anthropic ConvChatCompletionsRequest error: %v", err)
		return nil, err
	}

	if a.isAws {

		chatCompletionReq := model.AnthropicChatCompletionReq{}
		if err = json.Unmarshal(request, &chatCompletionReq); err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream Anthropic model: %s, request: %s, json.Unmarshal error: %v", a.Model, request, err)
			return nil, err
		}

		chatCompletionReq.AnthropicVersion = "bedrock-2023-05-31"
		chatCompletionReq.Stream = false

		invokeModelStreamInput := &bedrockruntime.InvokeModelWithResponseStreamInput{
			ModelId:     aws.String(chatCompletionReq.Model),
			Accept:      aws.String("application/json"),
			ContentType: aws.String("application/json"),
		}

		if modelId, ok := AwsModelIDMap[chatCompletionReq.Model]; ok {
			invokeModelStreamInput.ModelId = aws.String(modelId)
		}

		chatCompletionReq.Model = ""

		if invokeModelStreamInput.Body, err = gjson.Marshal(chatCompletionReq); err != nil {
			logger.Error(ctx, err)
			return responseChan, err
		}

		invokeModelStreamOutput, err := a.awsClient.InvokeModelWithResponseStream(ctx, invokeModelStreamInput)
		if err != nil {
			logger.Error(ctx, err)
			return responseChan, err
		}

		stream := invokeModelStreamOutput.GetStream()

		duration := gtime.TimestampMilli()

		responseChan = make(chan *model.ChatCompletionResponse)

		if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

			defer func() {
				if err := stream.Close(); err != nil {
					logger.Errorf(ctx, "ChatCompletionsStream Anthropic model: %s, stream.Close error: %v", a.Model, err)
				}

				end := gtime.TimestampMilli()
				logger.Infof(ctx, "ChatCompletionsStream Anthropic model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", a.Model, duration-now, end-duration, end-now)
			}()

			for {

				event, ok := <-stream.Events()
				if !ok {

					if !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
						logger.Errorf(ctx, "ChatCompletionsStream Anthropic model: %s, error: %v", a.Model, err)
					}

					end := gtime.TimestampMilli()
					responseChan <- &model.ChatCompletionResponse{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Error:     err,
					}

					return
				}

				var bytes []byte
				var id string

				switch v := event.(type) {
				case *types.ResponseStreamMemberChunk:

					bytes = v.Value.Bytes

				case *types.UnknownUnionMember:

					end := gtime.TimestampMilli()
					responseChan <- &model.ChatCompletionResponse{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Error:     errors.New("unknown tag:" + v.Tag),
					}

					return
				default:

					end := gtime.TimestampMilli()
					responseChan <- &model.ChatCompletionResponse{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Error:     errors.New("unknown type"),
					}

					return
				}

				response, err := a.ConvChatCompletionsStreamResponseOfficial(ctx, bytes)
				if err != nil {
					logger.Errorf(ctx, "ChatCompletionsStream Anthropic ConvChatCompletionsStreamResponseOfficial error: %v", err)

					end := gtime.TimestampMilli()
					responseChan <- &model.ChatCompletionResponse{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Error:     err,
					}

					return
				}

				if response.Id != "" {
					id = response.Id
				}

				response.Id = consts.COMPLETION_ID_PREFIX + id

				end := gtime.TimestampMilli()

				response.ConnTime = duration - now
				response.Duration = end - duration
				response.TotalTime = end - now

				responseChan <- &response
			}
		}, nil); err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream Anthropic model: %s, error: %v", a.Model, err)
			return responseChan, err
		}

	} else {

		stream, err := util.SSEClient(ctx, a.BaseUrl+a.Path, a.header, request, a.Timeout, a.ProxyUrl, a.requestErrorHandler)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream Anthropic model: %s, error: %v", a.Model, err)
			return responseChan, err
		}

		duration := gtime.TimestampMilli()

		responseChan = make(chan *model.ChatCompletionResponse)

		if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

			defer func() {
				if err := stream.Close(); err != nil {
					logger.Errorf(ctx, "ChatCompletionsStream Anthropic model: %s, stream.Close error: %v", a.Model, err)
				}

				end := gtime.TimestampMilli()
				logger.Infof(ctx, "ChatCompletionsStream Anthropic model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", a.Model, duration-now, end-duration, end-now)
			}()

			var id string
			var promptTokens int

			for {

				responseBytes, err := stream.Recv()
				if err != nil {

					if !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
						logger.Errorf(ctx, "ChatCompletionsStream Anthropic model: %s, error: %v", a.Model, err)
					}

					end := gtime.TimestampMilli()
					responseChan <- &model.ChatCompletionResponse{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Error:     err,
					}

					return
				}

				response, err := a.ConvChatCompletionsStreamResponseOfficial(ctx, responseBytes)
				if err != nil {
					logger.Errorf(ctx, "ChatCompletionsStream Anthropic ConvChatCompletionsStreamResponseOfficial error: %v", err)

					end := gtime.TimestampMilli()
					responseChan <- &model.ChatCompletionResponse{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Error:     err,
					}

					return
				}

				if response.Id != "" {
					id = response.Id
				}

				response.Id = consts.COMPLETION_ID_PREFIX + id

				if response.Usage != nil {
					if response.Usage.InputTokens != 0 {
						promptTokens = response.Usage.InputTokens
					}
					response.Usage = &model.Usage{
						PromptTokens:             promptTokens,
						CompletionTokens:         response.Usage.OutputTokens,
						TotalTokens:              promptTokens + response.Usage.OutputTokens,
						CacheCreationInputTokens: response.Usage.CacheCreationInputTokens,
						CacheReadInputTokens:     response.Usage.CacheReadInputTokens,
					}
				}

				end := gtime.TimestampMilli()

				response.ConnTime = duration - now
				response.Duration = end - duration
				response.TotalTime = end - now

				responseChan <- &response
			}

		}, nil); err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream Anthropic model: %s, error: %v", a.Model, err)
			return responseChan, err
		}
	}

	return responseChan, nil
}
