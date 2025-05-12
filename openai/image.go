package openai

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/go-openai"
)

func (c *Client) Image(ctx context.Context, request model.ImageRequest) (res model.ImageResponse, err error) {

	logger.Infof(ctx, "Image OpenAI model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "Image OpenAI model: %s totalTime: %d ms", request.Model, gtime.TimestampMilli()-now)
	}()

	response, err := c.client.CreateImage(ctx, openai.ImageRequest{
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
		logger.Errorf(ctx, "Image OpenAI model: %s, error: %v", request.Model, err)
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
		Usage: &model.Usage{
			InputTokens:        response.Usage.InputTokens,
			OutputTokens:       response.Usage.OutputTokens,
			InputTokensDetails: response.Usage.InputTokensDetails,
		},
	}

	return res, nil
}
