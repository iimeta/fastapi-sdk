package xai

import (
	"encoding/json"
	"strconv"

	"github.com/iimeta/fastapi-sdk/v2/model"
)

func convXAIChatToUnified(official model.XAIChatCompletionRes, raw []byte) model.ChatCompletionResponse {
	response := model.ChatCompletionResponse{ResponseBytes: raw}
	_ = json.Unmarshal(raw, &response)
	response.ResponseBytes = raw
	response.Usage = ConvertUsage(official.Usage)
	if len(official.Citations) > 0 {
		response.Citations = official.Citations
	}
	if official.OutputFiles != nil {
		response.OutputFiles = official.OutputFiles
	}
	if response.Id == "" {
		response.Id = official.RequestId
	}
	for i := range response.Choices {
		if response.Choices[i].Message != nil && response.Choices[i].Message.Annotations == nil {
			response.Choices[i].Message.Annotations = []any{}
		}
	}
	return response
}

func convXAIResponsesToOpenAI(official model.XAIResponsesRes, raw []byte) model.OpenAIResponsesRes {
	res := model.OpenAIResponsesRes{ResponseBytes: raw}
	_ = json.Unmarshal(raw, &res)
	res.ResponseBytes = raw
	res.Usage = ConvertUsage(official.Usage)
	if official.Error != nil {
		res.Error = convXAIResponsesError(official.Error)
	}
	return res
}

func convXAIResponsesStreamToOpenAI(official model.XAIResponsesStreamRes, raw []byte) model.OpenAIResponsesStreamRes {
	res := model.OpenAIResponsesStreamRes{ResponseBytes: raw}
	_ = json.Unmarshal(raw, &res)
	res.ResponseBytes = raw
	if official.Response.Usage != nil || official.Response.Id != "" {
		res.Response.Usage = ConvertUsage(official.Response.Usage)
	}
	if official.Error != nil {
		res.Error = convXAIResponsesError(official.Error)
	} else if official.Response.Error != nil {
		res.Error = convXAIResponsesError(official.Response.Error)
		res.Response.Error = res.Error
	}
	return res
}

func convXAIResponsesCompactToOpenAI(official model.XAIResponsesCompactRes, raw []byte) model.OpenAIResponsesRes {
	res := model.OpenAIResponsesRes{ResponseBytes: raw}
	_ = json.Unmarshal(raw, &res)
	res.ResponseBytes = raw
	res.Id = official.Id
	res.Object = official.Object
	res.Model = official.Model
	res.CreatedAt = official.CreatedAt
	if res.Status == "" {
		res.Status = "completed"
	}
	res.Usage = ConvertUsage(official.Usage)
	if official.Error != nil {
		res.Error = convXAIResponsesError(official.Error)
	}
	return res
}

func convXAIResponsesError(err *model.XAIResponsesError) *model.OpenAIResponsesError {
	if err == nil {
		return nil
	}
	return &model.OpenAIResponsesError{
		Code:           err.Code,
		Message:        err.Message,
		Param:          err.Param,
		SequenceNumber: err.SequenceNumber,
		Type:           err.Type,
	}
}

func convXAIImageToUnified(official model.XAIImageRes, raw []byte) model.ImageResponse {
	response := model.ImageResponse{
		ResponseBytes: raw,
	}
	if u := ConvertUsage(official.Usage); u != nil {
		response.Usage = *u
	}
	for _, item := range official.Data {
		entry := model.ImageResponseData{
			Url:          item.Url,
			B64Json:      item.B64Json,
			MimeType:     item.MimeType,
			StorageError: item.StorageError,
		}
		if item.FileOutput != nil {
			entry.FileOutput = &model.ImageFileOutput{
				FileId:             item.FileOutput.FileId,
				Filename:           item.FileOutput.Filename,
				PublicUrl:          item.FileOutput.PublicUrl,
				PublicUrlError:     item.FileOutput.PublicUrlError,
				ExpiresAt:          item.FileOutput.ExpiresAt,
				PublicUrlExpiresAt: item.FileOutput.PublicUrlExpiresAt,
			}
		}
		response.Data = append(response.Data, entry)
	}
	return response
}

func convXAIVideoToUnified(official model.XAIVideoJobRes, raw []byte) model.VideoJobResponse {
	response := model.VideoJobResponse{
		Id:            firstNonEmpty(official.Id, official.RequestId),
		Object:        official.Object,
		Model:         official.Model,
		Prompt:        official.Prompt,
		Usage:         ConvertUsage(official.Usage),
		ResponseBytes: raw,
	}
	response.Status = convVideoStatus(official.Status, nil)
	if official.Duration > 0 {
		response.Seconds = itoa(official.Duration)
	}
	if official.Video != nil {
		response.VideoUrl = official.Video.Url
		if response.Seconds == "" && official.Video.Duration > 0 {
			response.Seconds = itoa(official.Video.Duration)
		}
	}
	if official.Error != nil {
		response.Error = &model.VideoError{
			Code:    official.Error.Code,
			Message: official.Error.Message,
		}
	}
	return response
}

func convChatToXAIReq(request model.ChatCompletionRequest) model.XAIChatCompletionReq {
	official := model.XAIChatCompletionReq{
		Model:               request.Model,
		Messages:            request.Messages,
		Deferred:            request.Deferred,
		FrequencyPenalty:    request.FrequencyPenalty,
		LogitBias:           request.LogitBias,
		Logprobs:            request.LogProbs,
		MaxCompletionTokens: request.MaxCompletionTokens,
		MaxTokens:           request.MaxTokens,
		N:                   request.N,
		ParallelToolCalls:   request.ParallelToolCalls,
		PresencePenalty:     request.PresencePenalty,
		PromptCacheKey:      request.PromptCacheKey,
		ReasoningEffort:     request.ReasoningEffort,
		ResponseFormat:      request.ResponseFormat,
		SafetyIdentifier:    request.SafetyIdentifier,
		SearchParameters:    request.SearchParameters,
		Seed:                request.Seed,
		ServiceTier:         request.ServiceTier,
		Stop:                request.Stop,
		Stream:              request.Stream,
		Temperature:         request.Temperature,
		ToolChoice:          request.ToolChoice,
		Tools:               request.Tools,
		TopLogprobs:         request.TopLogProbs,
		TopP:                request.TopP,
		User:                request.User,
		WebSearchOptions:    request.WebSearchOptions,
	}
	if request.StreamOptions != nil {
		official.StreamOptions = &model.XAIStreamOptions{
			IncludeUsage: request.StreamOptions.IncludeUsage,
		}
	}
	return official
}

func convImageToXAIReq(request model.ImageGenerationRequest) model.XAIImageReq {
	official := model.XAIImageReq{
		Prompt:         request.Prompt,
		Model:          request.Model,
		N:              request.N,
		AspectRatio:    request.AspectRatio,
		Resolution:     request.Resolution,
		ResponseFormat: request.ResponseFormat,
		Quality:        request.Quality,
		User:           request.User,
		Image:          convImageContent(request.Image),
	}
	if official.AspectRatio == "" && request.Size != "" {
		official.AspectRatio = sizeToAspectRatio(request.Size)
	}
	if request.StorageOptions != nil {
		if b, err := json.Marshal(request.StorageOptions); err == nil {
			var storage model.XAIStorageOptions
			if json.Unmarshal(b, &storage) == nil {
				official.StorageOptions = &storage
			}
		}
	}
	if len(request.Images) > 0 {
		images := make([]any, 0, len(request.Images))
		for _, item := range request.Images {
			if converted := convImageEditImage(item); converted != nil {
				images = append(images, converted)
			}
		}
		if len(images) == 1 && official.Image == nil {
			official.Image = images[0]
		} else if len(images) > 0 {
			official.Images = images
		}
	}
	return official
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func itoa(v int) string {
	return strconv.Itoa(v)
}
