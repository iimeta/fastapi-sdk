package openai

import (
	"context"
	"errors"
	"io"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/go-openai"
)

func (o *OpenAI) ChatCompletions(ctx context.Context, data []byte) (res model.ChatCompletionResponse, err error) {

	request, err := o.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions OpenAI ConvChatCompletionsRequest error: %v", err)
		return res, err
	}

	logger.Infof(ctx, "ChatCompletions OpenAI model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ChatCompletions OpenAI model: %s totalTime: %d ms", request.Model, res.TotalTime)
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

		if chatCompletionMessage.Role == consts.ROLE_SYSTEM && (gstr.HasPrefix(request.Model, "o1") || gstr.HasPrefix(request.Model, "o3")) {
			chatCompletionMessage.Role = consts.ROLE_USER
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
		WebSearchOptions:    request.WebSearchOptions,
	}

	if gstr.HasPrefix(chatCompletionRequest.Model, "o") || gstr.HasPrefix(chatCompletionRequest.Model, "gpt-5") {
		if chatCompletionRequest.MaxCompletionTokens == 0 && chatCompletionRequest.MaxTokens != 0 {
			chatCompletionRequest.MaxCompletionTokens = chatCompletionRequest.MaxTokens
		}
		chatCompletionRequest.MaxTokens = 0
	}

	response, err := o.client.CreateChatCompletion(ctx, chatCompletionRequest)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletions OpenAI model: %s, error: %v", request.Model, err)
		return res, o.apiErrorHandler(err)
	}

	logger.Infof(ctx, "ChatCompletions OpenAI model: %s finished", request.Model)

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
			InputTokens:             response.Usage.InputTokens,
			OutputTokens:            response.Usage.OutputTokens,
			InputTokensDetails:      response.Usage.InputTokensDetails,
			OutputTokensDetails:     response.Usage.OutputTokensDetails,
		},
		SystemFingerprint: response.SystemFingerprint,
	}

	for _, choice := range response.Choices {
		if choice.Message.Annotations == nil {
			choice.Message.Annotations = []any{}
		}
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
				Annotations:      choice.Message.Annotations,
			},
			FinishReason: choice.FinishReason,
			LogProbs:     choice.LogProbs,
		})
	}

	return res, nil
}

func (o *OpenAI) ChatCompletionsStream(ctx context.Context, data []byte) (responseChan chan *model.ChatCompletionResponse, err error) {

	request, err := o.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream OpenAI ConvChatCompletionsRequest error: %v", err)
		return nil, err
	}

	logger.Infof(ctx, "ChatCompletionsStream OpenAI model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionsStream OpenAI model: %s totalTime: %d ms", request.Model, gtime.TimestampMilli()-now)
		}
	}()

	if (o.isSupportStream != nil && !*o.isSupportStream) || (gstr.HasPrefix(request.Model, "o") && o.isAzure) {
		return o.ChatCompletionStreamToNonStream(ctx, data)
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

		if chatCompletionMessage.Role == consts.ROLE_SYSTEM && (gstr.HasPrefix(request.Model, "o1") || gstr.HasPrefix(request.Model, "o3")) {
			chatCompletionMessage.Role = consts.ROLE_DEVELOPER
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
		WebSearchOptions:    request.WebSearchOptions,
	}

	if gstr.HasPrefix(chatCompletionRequest.Model, "o") || gstr.HasPrefix(chatCompletionRequest.Model, "gpt-5") {
		if chatCompletionRequest.MaxCompletionTokens == 0 && chatCompletionRequest.MaxTokens != 0 {
			chatCompletionRequest.MaxCompletionTokens = chatCompletionRequest.MaxTokens
		}
		chatCompletionRequest.MaxTokens = 0
	}

	stream, err := o.client.CreateChatCompletionStream(ctx, chatCompletionRequest)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionsStream OpenAI model: %s, error: %v", request.Model, err)
		return responseChan, o.apiErrorHandler(err)
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionsStream OpenAI model: %s, stream.Close error: %v", request.Model, err)
			}

			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionsStream OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		for {

			responseBytes, streamResponse, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {

				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "ChatCompletionsStream OpenAI model: %s, error: %v", request.Model, err)
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
						Annotations:      choice.Delta.Annotations,
					},
					FinishReason: choice.FinishReason,
				})
			}

			if streamResponse.Usage != nil {

				response.Usage = &model.Usage{
					PromptTokens:            streamResponse.Usage.PromptTokens,
					CompletionTokens:        streamResponse.Usage.CompletionTokens,
					TotalTokens:             streamResponse.Usage.TotalTokens,
					PromptTokensDetails:     streamResponse.Usage.PromptTokensDetails,
					CompletionTokensDetails: streamResponse.Usage.CompletionTokensDetails,
					InputTokens:             streamResponse.Usage.InputTokens,
					OutputTokens:            streamResponse.Usage.OutputTokens,
					InputTokensDetails:      streamResponse.Usage.InputTokensDetails,
					OutputTokensDetails:     streamResponse.Usage.OutputTokensDetails,
				}

				if len(response.Choices) == 0 {
					response.Choices = append(response.Choices, model.ChatCompletionChoice{
						Delta:        new(model.ChatCompletionStreamChoiceDelta),
						FinishReason: openai.FinishReasonStop,
					})
				}
			}

			if errors.Is(err, io.EOF) {
				logger.Infof(ctx, "ChatCompletionsStream OpenAI model: %s finished", request.Model)

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
		logger.Errorf(ctx, "ChatCompletionsStream OpenAI model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (o *OpenAI) ChatCompletionStreamToNonStream(ctx context.Context, data []byte) (responseChan chan *model.ChatCompletionResponse, err error) {

	request, err := o.ConvChatCompletionsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStreamToNonStream OpenAI ConvChatCompletionsRequest error: %v", err)
		return nil, err
	}

	responseChan = make(chan *model.ChatCompletionResponse)

	now := gtime.TimestampMilli()
	duration := now

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ChatCompletionStreamToNonStream OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)
		}()

		request.Stream = false

		streamResponse, err := o.ChatCompletions(ctx, gjson.MustEncode(request))
		if err != nil {

			if !errors.Is(err, context.Canceled) {
				logger.Errorf(ctx, "ChatCompletionStreamToNonStream OpenAI model: %s, error: %v", request.Model, err)
			}

			end := gtime.TimestampMilli()
			responseChan <- &model.ChatCompletionResponse{
				ConnTime:  gtime.TimestampMilli() - now,
				Duration:  end - gtime.TimestampMilli(),
				TotalTime: end - now,
				Error:     err,
			}

			return
		}

		duration = gtime.TimestampMilli()

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
					Audio:            choice.Message.Audio,
					Annotations:      choice.Message.Annotations,
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
			InputTokens:             streamResponse.Usage.InputTokens,
			OutputTokens:            streamResponse.Usage.OutputTokens,
			InputTokensDetails:      streamResponse.Usage.InputTokensDetails,
			OutputTokensDetails:     streamResponse.Usage.OutputTokensDetails,
		}

		end := gtime.TimestampMilli()
		response.Duration = end - duration
		response.TotalTime = end - now

		responseChan <- response

		end = gtime.TimestampMilli()
		responseChan <- &model.ChatCompletionResponse{
			Duration:  end - duration,
			TotalTime: end - now,
			Error:     io.EOF,
		}

	}, nil); err != nil {
		logger.Errorf(ctx, "ChatCompletionStreamToNonStream OpenAI model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
