package openai

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/go-openai"
	"io"
)

func (c *Client) ChatCompletion(ctx context.Context, request model.ChatCompletionRequest) (res model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletion OpenAI model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletion OpenAI model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	messages := make([]openai.ChatCompletionMessage, 0)
	for _, message := range request.Messages {

		chatCompletionMessage := openai.ChatCompletionMessage{
			Role:         message.Role,
			Content:      message.Content,
			Refusal:      message.Refusal,
			Name:         message.Name,
			FunctionCall: message.FunctionCall,
			ToolCalls:    message.ToolCalls,
			ToolCallID:   message.ToolCallID,
			Audio:        message.Audio,
		}

		messages = append(messages, chatCompletionMessage)
	}

	chatCompletionRequest := openai.ChatCompletionRequest{
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
		ParallelToolCalls:   request.ParallelToolCalls,
		Store:               request.Store,
		Metadata:            request.Metadata,
		ReasoningEffort:     request.ReasoningEffort,
		Modalities:          request.Modalities,
		Audio:               request.Audio,
	}

	if gstr.HasPrefix(chatCompletionRequest.Model, "o1-") {
		if chatCompletionRequest.MaxCompletionTokens == 0 && chatCompletionRequest.MaxTokens != 0 {
			chatCompletionRequest.MaxCompletionTokens = chatCompletionRequest.MaxTokens
		}
		chatCompletionRequest.MaxTokens = 0
	}

	response, err := c.client.CreateChatCompletion(ctx, chatCompletionRequest)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletion OpenAI model: %s, error: %v", request.Model, err)
		return res, c.apiErrorHandler(err)
	}

	logger.Infof(ctx, "ChatCompletion OpenAI model: %s finished", request.Model)

	res = model.ChatCompletionResponse{
		ID:      response.ID,
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

	logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s totalTime: %d ms", request.Model, gtime.TimestampMilli()-now)
		}
	}()

	if gstr.HasPrefix(request.Model, "o1-") && c.isAzure {
		return c.O1ChatCompletionStream(ctx, request)
	}

	messages := make([]openai.ChatCompletionMessage, 0)
	for _, message := range request.Messages {

		chatCompletionMessage := openai.ChatCompletionMessage{
			Role:         message.Role,
			Content:      message.Content,
			Refusal:      message.Refusal,
			Name:         message.Name,
			FunctionCall: message.FunctionCall,
			ToolCalls:    message.ToolCalls,
			ToolCallID:   message.ToolCallID,
			Audio:        message.Audio,
		}

		messages = append(messages, chatCompletionMessage)
	}

	// 默认让流式返回usage
	if request.StreamOptions == nil { // request.Tools == nil &&
		request.StreamOptions = &openai.StreamOptions{
			IncludeUsage: true,
		}
	}

	chatCompletionRequest := openai.ChatCompletionRequest{
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
	}

	if gstr.HasPrefix(chatCompletionRequest.Model, "o1-") {
		if chatCompletionRequest.MaxCompletionTokens == 0 && chatCompletionRequest.MaxTokens != 0 {
			chatCompletionRequest.MaxCompletionTokens = chatCompletionRequest.MaxTokens
		}
		chatCompletionRequest.MaxTokens = 0
	}

	stream, err := c.client.CreateChatCompletionStream(ctx, chatCompletionRequest)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStream OpenAI model: %s, error: %v", request.Model, err)
		return responseChan, c.apiErrorHandler(err)
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream OpenAI model: %s, stream.Close error: %v", request.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, streamResponse, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {

				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "ChatCompletionStream OpenAI model: %s, error: %v", request.Model, err)
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
				ID:                streamResponse.ID,
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
				logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s finished", request.Model)

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
		logger.Errorf(ctx, "ChatCompletionStream OpenAI model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (c *Client) O1ChatCompletionStream(ctx context.Context, request model.ChatCompletionRequest) (responseChan chan *model.ChatCompletionResponse, err error) {

	responseChan = make(chan *model.ChatCompletionResponse)

	now := gtime.TimestampMilli()
	duration := now

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.TimestampMilli()
			logger.Infof(ctx, "O1ChatCompletionStream OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		request.Stream = false

		streamResponse, err := c.ChatCompletion(ctx, request)
		if err != nil {

			if !errors.Is(err, context.Canceled) {
				logger.Errorf(ctx, "O1ChatCompletionStream OpenAI model: %s, error: %v", request.Model, err)
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
			ID:                streamResponse.ID,
			Object:            streamResponse.Object,
			Created:           streamResponse.Created,
			Model:             streamResponse.Model,
			PromptAnnotations: streamResponse.PromptAnnotations,
			ConnTime:          duration - now,
			SystemFingerprint: streamResponse.SystemFingerprint,
		}

		choices := make([]model.ChatCompletionChoice, 0)
		for i, choice := range streamResponse.Choices {
			choices = append(choices, model.ChatCompletionChoice{
				Index: i,
				Delta: &model.ChatCompletionStreamChoiceDelta{
					Content:          gconv.String(choice.Message.Content),
					ReasoningContent: choice.Message.ReasoningContent,
					Role:             choice.Message.Role,
					FunctionCall:     choice.Message.FunctionCall,
					ToolCalls:        choice.Message.ToolCalls,
				},
				FinishReason: openai.FinishReasonStop,
			})
		}
		response.Choices = choices

		response.Usage = &model.Usage{
			PromptTokens:            streamResponse.Usage.PromptTokens,
			CompletionTokens:        streamResponse.Usage.CompletionTokens,
			TotalTokens:             streamResponse.Usage.TotalTokens,
			CompletionTokensDetails: streamResponse.Usage.CompletionTokensDetails,
		}

		end := gtime.TimestampMilli()
		response.Duration = end - duration
		response.TotalTime = end - now

		responseChan <- response

		response = &model.ChatCompletionResponse{}
		end = gtime.TimestampMilli()
		response.Duration = end - duration
		response.TotalTime = end - now
		response.Error = io.EOF
		responseChan <- response

	}, nil); err != nil {
		logger.Errorf(ctx, "O1ChatCompletionStream OpenAI model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
