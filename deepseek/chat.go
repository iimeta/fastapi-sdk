package deepseek

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/go-openai"
	"io"
)

func (c *Client) ChatCompletion(ctx context.Context, request model.ChatCompletionRequest) (res model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletion DeepSeek model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletion DeepSeek model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	var newMessages []model.ChatCompletionMessage
	if c.isSupportSystemRole != nil {
		newMessages = common.HandleMessages(request.Messages, *c.isSupportSystemRole)
	} else {
		newMessages = common.HandleMessages(request.Messages, true)
	}

	messages := make([]openai.ChatCompletionMessage, 0)
	for _, message := range newMessages {

		chatCompletionMessage := openai.ChatCompletionMessage{
			Role:         message.Role,
			Name:         message.Name,
			Content:      gconv.String(message.Content),
			FunctionCall: message.FunctionCall,
			ToolCalls:    message.ToolCalls,
			ToolCallID:   message.ToolCallID,
		}

		messages = append(messages, chatCompletionMessage)
	}

	response, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:               request.Model,
		Messages:            messages,
		MaxTokens:           request.MaxTokens,
		MaxCompletionTokens: request.MaxCompletionTokens,
		Temperature:         request.Temperature,
		TopP:                request.TopP,
		N:                   request.N,
		Stream:              request.Stream,
		Stop:                request.Stop,
		PresencePenalty:     request.PresencePenalty,
		ResponseFormat:      request.ResponseFormat,
		Seed:                request.Seed,
		FrequencyPenalty:    request.FrequencyPenalty,
		LogitBias:           request.LogitBias,
		LogProbs:            request.LogProbs,
		TopLogProbs:         request.TopLogProbs,
		User:                request.User,
		Functions:           request.Functions,
		FunctionCall:        request.FunctionCall,
		Tools:               request.Tools,
		ToolChoice:          request.ToolChoice,
		StreamOptions:       request.StreamOptions,
		ParallelToolCalls:   request.ParallelToolCalls,
		Store:               request.Store,
		Metadata:            request.Metadata,
		ReasoningEffort:     request.ReasoningEffort,
		Modalities:          request.Modalities,
		Audio:               request.Audio,
	})
	if err != nil {
		logger.Errorf(ctx, "ChatCompletion DeepSeek model: %s, error: %v", request.Model, err)
		return res, c.apiErrorHandler(err)
	}

	logger.Infof(ctx, "ChatCompletion DeepSeek model: %s finished", request.Model)

	res = model.ChatCompletionResponse{
		ID:      consts.COMPLETION_ID_PREFIX + response.ID,
		Object:  response.Object,
		Created: response.Created,
		Model:   response.Model,
		Usage: &model.Usage{
			PromptTokens:            response.Usage.PromptTokens,
			CompletionTokens:        response.Usage.CompletionTokens,
			TotalTokens:             response.Usage.TotalTokens,
			PromptTokensDetails:     response.Usage.PromptTokensDetails,
			CompletionTokensDetails: response.Usage.CompletionTokensDetails,
		},
		SystemFingerprint: response.SystemFingerprint,
	}

	for _, choice := range response.Choices {
		res.Choices = append(res.Choices, model.ChatCompletionChoice{
			Index: choice.Index,
			Message: &model.ChatCompletionMessage{
				Role:             choice.Message.Role,
				Content:          choice.Message.Content,
				ReasoningContent: choice.Message.ReasoningContent,
				Refusal:          choice.Message.Refusal,
				MultiContent:     choice.Message.MultiContent,
				Name:             choice.Message.Name,
				FunctionCall:     choice.Message.FunctionCall,
				ToolCalls:        choice.Message.ToolCalls,
				ToolCallID:       choice.Message.ToolCallID,
				Audio:            choice.Message.Audio,
			},
			FinishReason: choice.FinishReason,
			LogProbs:     choice.LogProbs,
		})
	}

	return res, nil
}

func (c *Client) ChatCompletionStream(ctx context.Context, request model.ChatCompletionRequest) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionStream DeepSeek model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionStream DeepSeek model: %s totalTime: %d ms", request.Model, gtime.TimestampMilli()-now)
		}
	}()

	var newMessages []model.ChatCompletionMessage
	if c.isSupportSystemRole != nil {
		newMessages = common.HandleMessages(request.Messages, *c.isSupportSystemRole)
	} else {
		newMessages = common.HandleMessages(request.Messages, true)
	}

	messages := make([]openai.ChatCompletionMessage, 0)
	for _, message := range newMessages {

		chatCompletionMessage := openai.ChatCompletionMessage{
			Role:         message.Role,
			Name:         message.Name,
			Content:      gconv.String(message.Content),
			FunctionCall: message.FunctionCall,
			ToolCalls:    message.ToolCalls,
			ToolCallID:   message.ToolCallID,
		}

		messages = append(messages, chatCompletionMessage)
	}

	// 默认让流式返回usage
	if request.StreamOptions == nil { // request.Tools == nil &&
		request.StreamOptions = &openai.StreamOptions{
			IncludeUsage: true,
		}
	}

	stream, err := c.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:               request.Model,
		Messages:            messages,
		MaxTokens:           request.MaxTokens,
		MaxCompletionTokens: request.MaxCompletionTokens,
		Temperature:         request.Temperature,
		TopP:                request.TopP,
		N:                   request.N,
		Stream:              request.Stream,
		Stop:                request.Stop,
		PresencePenalty:     request.PresencePenalty,
		ResponseFormat:      request.ResponseFormat,
		Seed:                request.Seed,
		FrequencyPenalty:    request.FrequencyPenalty,
		LogitBias:           request.LogitBias,
		LogProbs:            request.LogProbs,
		TopLogProbs:         request.TopLogProbs,
		User:                request.User,
		Functions:           request.Functions,
		FunctionCall:        request.FunctionCall,
		Tools:               request.Tools,
		ToolChoice:          request.ToolChoice,
		StreamOptions:       request.StreamOptions,
		ParallelToolCalls:   request.ParallelToolCalls,
		Store:               request.Store,
		Metadata:            request.Metadata,
		ReasoningEffort:     request.ReasoningEffort,
		Modalities:          request.Modalities,
		Audio:               request.Audio,
	})
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStream DeepSeek model: %s, error: %v", request.Model, err)
		return responseChan, c.apiErrorHandler(err)
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream DeepSeek model: %s, stream.Close error: %v", request.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionStream DeepSeek model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, streamResponse, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {

				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "ChatCompletionStream DeepSeek model: %s, error: %v", request.Model, err)
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

			response := &model.ChatCompletionResponse{
				ID:                consts.COMPLETION_ID_PREFIX + streamResponse.ID,
				Object:            streamResponse.Object,
				Created:           streamResponse.Created,
				Model:             streamResponse.Model,
				PromptAnnotations: streamResponse.PromptAnnotations,
				ResponseBytes:     responseBytes,
				ConnTime:          duration - now,
			}

			for _, choice := range streamResponse.Choices {
				response.Choices = append(response.Choices, model.ChatCompletionChoice{
					Index: choice.Index,
					Delta: &model.ChatCompletionStreamChoiceDelta{
						Content:          choice.Delta.Content,
						ReasoningContent: choice.Delta.ReasoningContent,
						Role:             choice.Delta.Role,
						FunctionCall:     choice.Delta.FunctionCall,
						ToolCalls:        choice.Delta.ToolCalls,
						Refusal:          choice.Delta.Refusal,
						Audio:            choice.Delta.Audio,
					},
					FinishReason: choice.FinishReason,
					//ContentFilterResults: &choice.ContentFilterResults,
				})
			}

			if streamResponse.Usage != nil {

				response.Usage = &model.Usage{
					PromptTokens:            streamResponse.Usage.PromptTokens,
					CompletionTokens:        streamResponse.Usage.CompletionTokens,
					TotalTokens:             streamResponse.Usage.TotalTokens,
					PromptTokensDetails:     streamResponse.Usage.PromptTokensDetails,
					CompletionTokensDetails: streamResponse.Usage.CompletionTokensDetails,
				}

				if len(response.Choices) == 0 {
					response.Choices = append(response.Choices, model.ChatCompletionChoice{
						Delta:        new(model.ChatCompletionStreamChoiceDelta),
						FinishReason: openai.FinishReasonStop,
					})
				}
			}

			if errors.Is(err, io.EOF) {
				logger.Infof(ctx, "ChatCompletionStream DeepSeek model: %s finished", request.Model)

				end := gtime.TimestampMilli()
				response.Duration = end - duration
				response.TotalTime = end - now
				response.Error = io.EOF
				responseChan <- response

				return
			}

			end := gtime.TimestampMilli()
			response.Duration = end - duration
			response.TotalTime = end - now

			responseChan <- response
		}
	}, nil); err != nil {
		logger.Errorf(ctx, "ChatCompletionStream DeepSeek model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
