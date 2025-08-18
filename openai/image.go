package openai

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/go-openai"
)

func (o *OpenAI) ImageGenerations(ctx context.Context, request model.ImageGenerationRequest) (res model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerations OpenAI model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageGenerations OpenAI model: %s totalTime: %d ms", request.Model, gtime.TimestampMilli()-now)
	}()

	response, err := o.client.CreateImage(ctx, openai.ImageRequest{
		Prompt:            request.Prompt,
		Background:        request.Background,
		Model:             request.Model,
		Moderation:        request.Moderation,
		N:                 request.N,
		OutputCompression: request.OutputCompression,
		OutputFormat:      request.OutputFormat,
		Quality:           request.Quality,
		ResponseFormat:    request.ResponseFormat,
		Size:              request.Size,
		Style:             request.Style,
		User:              request.User,
	})
	if err != nil {
		logger.Errorf(ctx, "ImageGenerations OpenAI model: %s, error: %v", request.Model, err)
		return res, err
	}

	data := make([]model.ImageResponseDataInner, 0)
	for _, d := range response.Data {
		data = append(data, model.ImageResponseDataInner{
			URL:           d.URL,
			B64JSON:       d.B64JSON,
			RevisedPrompt: d.RevisedPrompt,
		})
	}

	res = model.ImageResponse{
		Created: response.Created,
		Data:    data,
		Usage: model.Usage{
			TotalTokens:  response.Usage.TotalTokens,
			InputTokens:  response.Usage.InputTokens,
			OutputTokens: response.Usage.OutputTokens,
			InputTokensDetails: model.InputTokensDetails{
				TextTokens:   response.Usage.InputTokensDetails.TextTokens,
				ImageTokens:  response.Usage.InputTokensDetails.ImageTokens,
				CachedTokens: response.Usage.InputTokensDetails.CachedTokens,
			},
		},
	}

	return res, nil
}

func (o *OpenAI) ImageEdits(ctx context.Context, request model.ImageEditRequest) (res model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageEdits OpenAI model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageEdits OpenAI model: %s totalTime: %d ms", request.Model, gtime.TimestampMilli()-now)
	}()

	response, err := o.client.CreateEditImage(ctx, openai.ImageEditRequest{
		Image:          request.Image,
		Prompt:         request.Prompt,
		Background:     request.Background,
		Mask:           request.Mask,
		Model:          request.Model,
		N:              request.N,
		Quality:        request.Quality,
		ResponseFormat: request.ResponseFormat,
		Size:           request.Size,
		User:           request.User,
	})
	if err != nil {
		logger.Errorf(ctx, "ImageEdits OpenAI model: %s, error: %v", request.Model, err)
		return res, err
	}

	data := make([]model.ImageResponseDataInner, 0)
	for _, d := range response.Data {
		data = append(data, model.ImageResponseDataInner{
			URL:           d.URL,
			B64JSON:       d.B64JSON,
			RevisedPrompt: d.RevisedPrompt,
		})
	}

	res = model.ImageResponse{
		Created: response.Created,
		Data:    data,
		Usage: model.Usage{
			TotalTokens:  response.Usage.TotalTokens,
			InputTokens:  response.Usage.InputTokens,
			OutputTokens: response.Usage.OutputTokens,
			InputTokensDetails: model.InputTokensDetails{
				TextTokens:   response.Usage.InputTokensDetails.TextTokens,
				ImageTokens:  response.Usage.InputTokensDetails.ImageTokens,
				CachedTokens: response.Usage.InputTokensDetails.CachedTokens,
			},
		},
	}

	return res, nil
}
