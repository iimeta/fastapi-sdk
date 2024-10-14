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

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "ChatCompletion OpenAI model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	messages := make([]openai.ChatCompletionMessage, 0)
	for _, message := range request.Messages {

		chatCompletionMessage := openai.ChatCompletionMessage{
			Role:         message.Role,
			Name:         message.Name,
			Content:      message.Content,
			FunctionCall: message.FunctionCall,
			ToolCalls:    message.ToolCalls,
			ToolCallID:   message.ToolCallID,
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
			CompletionTokensDetails: response.Usage.CompletionTokensDetails,
		},
		SystemFingerprint: response.SystemFingerprint,
	}

	for _, choice := range response.Choices {
		res.Choices = append(res.Choices, model.ChatCompletionChoice{
			Index:        choice.Index,
			Message:      &choice.Message,
			FinishReason: choice.FinishReason,
			LogProbs:     choice.LogProbs,
		})
	}

	return res, nil
}

func (c *Client) ChatCompletionStream(ctx context.Context, request model.ChatCompletionRequest) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s start", request.Model)

	now := gtime.Now().UnixMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s totalTime: %d ms", request.Model, gtime.Now().UnixMilli()-now)
		}
	}()

	if gstr.HasPrefix(request.Model, "o1-") {
		return c.O1ChatCompletionStream(ctx, request)
	}

	messages := make([]openai.ChatCompletionMessage, 0)
	for _, message := range request.Messages {

		chatCompletionMessage := openai.ChatCompletionMessage{
			Role:         message.Role,
			Name:         message.Name,
			Content:      message.Content,
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

	duration := gtime.Now().UnixMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream OpenAI model: %s, stream.Close error: %v", request.Model, err)
			}

			end := gtime.Now().UnixMilli()
			logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, streamResponse, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {

				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "ChatCompletionStream OpenAI model: %s, error: %v", request.Model, err)
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
					Index:        choice.Index,
					Delta:        &choice.Delta,
					FinishReason: choice.FinishReason,
					//ContentFilterResults: &choice.ContentFilterResults,
				})
			}

			if streamResponse.Usage != nil {

				response.Usage = &model.Usage{
					PromptTokens:            streamResponse.Usage.PromptTokens,
					CompletionTokens:        streamResponse.Usage.CompletionTokens,
					TotalTokens:             streamResponse.Usage.TotalTokens,
					CompletionTokensDetails: streamResponse.Usage.CompletionTokensDetails,
				}

				if len(response.Choices) == 0 {
					response.Choices = append(response.Choices, model.ChatCompletionChoice{
						Delta:        new(openai.ChatCompletionStreamChoiceDelta),
						FinishReason: openai.FinishReasonStop,
					})
				}
			}

			if errors.Is(err, io.EOF) {
				logger.Infof(ctx, "ChatCompletionStream OpenAI model: %s finished", request.Model)

				end := gtime.Now().UnixMilli()
				response.Duration = end - duration
				response.TotalTime = end - now
				response.Error = io.EOF
				responseChan <- response

				return
			}

			end := gtime.Now().UnixMilli()
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

	now := gtime.Now().UnixMilli()
	duration := now

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.Now().UnixMilli()
			logger.Infof(ctx, "O1ChatCompletionStream OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		request.Stream = false

		streamResponse, err := c.ChatCompletion(ctx, request)
		if err != nil {

			if !errors.Is(err, context.Canceled) {
				logger.Errorf(ctx, "O1ChatCompletionStream OpenAI model: %s, error: %v", request.Model, err)
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

		response := &model.ChatCompletionResponse{
			ID:                streamResponse.ID,
			Object:            streamResponse.Object,
			Created:           streamResponse.Created,
			Model:             streamResponse.Model,
			PromptAnnotations: streamResponse.PromptAnnotations,
			ConnTime:          duration - now,
			SystemFingerprint: streamResponse.SystemFingerprint,
		}

		response.Choices = []model.ChatCompletionChoice{{
			Delta: &openai.ChatCompletionStreamChoiceDelta{
				Content:      gconv.String(streamResponse.Choices[0].Message.Content),
				Role:         streamResponse.Choices[0].Message.Role,
				FunctionCall: streamResponse.Choices[0].Message.FunctionCall,
				ToolCalls:    streamResponse.Choices[0].Message.ToolCalls,
			},
			FinishReason: openai.FinishReasonStop,
		}}

		response.Usage = &model.Usage{
			PromptTokens:            streamResponse.Usage.PromptTokens,
			CompletionTokens:        streamResponse.Usage.CompletionTokens,
			TotalTokens:             streamResponse.Usage.TotalTokens,
			CompletionTokensDetails: streamResponse.Usage.CompletionTokensDetails,
		}

		end := gtime.Now().UnixMilli()
		response.Duration = end - duration
		response.TotalTime = end - now

		responseChan <- response

		response = &model.ChatCompletionResponse{}
		end = gtime.Now().UnixMilli()
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
