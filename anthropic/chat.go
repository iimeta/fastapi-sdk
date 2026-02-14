package anthropic

import (
	"context"
	"encoding/json"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (a *Anthropic) ChatCompletions(ctx context.Context, data any) (response model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletions Anthropic model: %s start", a.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions Anthropic model: %s totalTime: %d ms", a.Model, response.TotalTime)
	}()

	if !a.IsOfficialFormatRequest {

		request, err := a.ConvChatCompletionsRequest(ctx, data)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletions Anthropic ConvChatCompletionsRequest error: %v", err)
			return response, err
		}

		if data, err = a.ConvChatCompletionsRequestOfficial(ctx, request); err != nil {
			logger.Errorf(ctx, "ChatCompletions Anthropic ConvChatCompletionsRequestOfficial error: %v", err)
			return response, err
		}
	}

	if a.isAws {

		chatCompletionReq := model.AnthropicChatCompletionReq{}
		if v, ok := data.([]byte); ok {
			if err = json.Unmarshal(v, &chatCompletionReq); err != nil {
				logger.Errorf(ctx, "ChatCompletions Anthropic model: %s, request: %s, json.Unmarshal error: %v", a.Model, data, err)
				return response, err
			}
		}

		chatCompletionReq.AnthropicVersion = "bedrock-2023-05-31"
		chatCompletionReq.Model = ""
		chatCompletionReq.Metadata = nil

		data = gjson.MustEncode(chatCompletionReq)

		a.header = signHeader(a.Path, a.region, a.accessKey, a.secretKey, data.([]byte))
	}

	if a.Path == "" {
		a.Path = "/messages"
	}

	if response.ResponseBytes, err = util.HttpPost(ctx, a.BaseUrl+a.Path, a.header, data, nil, a.Timeout, a.ProxyUrl, a.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "ChatCompletions Anthropic model: %s, error: %v", a.Model, err)
		return response, err
	}

	if response, err = a.ConvChatCompletionsResponse(ctx, response.ResponseBytes); err != nil {
		logger.Errorf(ctx, "ChatCompletions Anthropic ConvChatCompletionsResponse error: %v", err)
		return response, err
	}

	return response, nil
}

func (a *Anthropic) ChatCompletionsStream(ctx context.Context, data any) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionsStream Anthropic model: %s start", a.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream Anthropic model: %s totalTime: %d ms", a.Model, gtime.TimestampMilli()-now)
		}
	}()

	if !a.IsOfficialFormatRequest {

		request, err := a.ConvChatCompletionsRequest(ctx, data)
		if err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream Anthropic ConvChatCompletionsRequest error: %v", err)
			return responseChan, err
		}

		if data, err = a.ConvChatCompletionsRequestOfficial(ctx, request); err != nil {
			logger.Errorf(ctx, "ChatCompletionsStream Anthropic ConvChatCompletionsRequest error: %v", err)
			return responseChan, err
		}
	}

	if a.isAws {

		chatCompletionReq := model.AnthropicChatCompletionReq{}

		if v, ok := data.([]byte); ok {
			if err = json.Unmarshal(v, &chatCompletionReq); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream Anthropic model: %s, request: %s, json.Unmarshal error: %v", a.Model, data, err)
				return responseChan, err
			}
		}

		chatCompletionReq.AnthropicVersion = "bedrock-2023-05-31"
		chatCompletionReq.Stream = false

		invokeModelStreamInput := &bedrockruntime.InvokeModelWithResponseStreamInput{
			ModelId:     aws.String(chatCompletionReq.Model),
			Accept:      aws.String("application/json"),
			ContentType: aws.String("application/json"),
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

			var id string
			var promptTokens int

			for {

				event, ok := <-stream.Events()
				if !ok {

					logger.Errorf(ctx, "ChatCompletionsStream Anthropic model: %s, error: %v", a.Model, context.Canceled)

					end := gtime.TimestampMilli()
					responseChan <- &model.ChatCompletionResponse{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Error:     context.Canceled,
					}

					return
				}

				var bytes []byte

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

				response, err := a.ConvChatCompletionsStreamResponse(ctx, bytes)
				if err != nil {
					logger.Errorf(ctx, "ChatCompletionsStream Anthropic ConvChatCompletionsStreamResponse error: %v", err)

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

					if response.Usage.PromptTokens == 0 {
						response.Usage.PromptTokens = promptTokens
					} else {
						promptTokens = response.Usage.PromptTokens
					}

					if response.Usage.PromptTokens+response.Usage.CompletionTokens > response.Usage.TotalTokens {
						response.Usage.TotalTokens = response.Usage.PromptTokens + response.Usage.CompletionTokens
					}
				}

				if len(response.Choices) > 0 && response.Choices[0].FinishReason != "" {
					logger.Infof(ctx, "ChatCompletionsStream Anthropic model: %s, finishReason: %s finished", a.Model, response.Choices[0].FinishReason)

					end := gtime.TimestampMilli()

					response.ConnTime = duration - now
					response.Duration = end - duration
					response.TotalTime = end - now
					responseChan <- &response

					responseChan <- &model.ChatCompletionResponse{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Error:     io.EOF,
					}

					return
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

	} else {

		if a.Path == "" {
			a.Path = "/messages"
		}

		stream, err := util.SSEClient(ctx, a.BaseUrl+a.Path, a.header, data, a.Timeout, a.ProxyUrl, a.requestErrorHandler)
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

					if errors.Is(err, io.EOF) {
						logger.Infof(ctx, "ChatCompletionsStream Anthropic model: %s finished", a.Model)
					} else {
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

				response, err := a.ConvChatCompletionsStreamResponse(ctx, responseBytes)
				if err != nil {
					logger.Errorf(ctx, "ChatCompletionsStream Anthropic ConvChatCompletionsStreamResponse error: %v", err)

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

					if response.Usage.PromptTokens == 0 {
						response.Usage.PromptTokens = promptTokens
					} else {
						promptTokens = response.Usage.PromptTokens
					}

					if response.Usage.PromptTokens+response.Usage.CompletionTokens > response.Usage.TotalTokens {
						response.Usage.TotalTokens = response.Usage.PromptTokens + response.Usage.CompletionTokens
					}
				}

				if len(response.Choices) > 0 && response.Choices[0].FinishReason != "" {
					logger.Infof(ctx, "ChatCompletionsStream Anthropic model: %s, finishReason: %s finished", a.Model, response.Choices[0].FinishReason)

					end := gtime.TimestampMilli()

					response.ConnTime = duration - now
					response.Duration = end - duration
					response.TotalTime = end - now
					responseChan <- &response

					responseChan <- &model.ChatCompletionResponse{
						ConnTime:  duration - now,
						Duration:  end - duration,
						TotalTime: end - now,
						Error:     io.EOF,
					}

					return
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
