package openai

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/go-openai"
)

func (c *Client) Moderations(ctx context.Context, request model.ModerationRequest) (res model.ModerationResponse, err error) {

	logger.Infof(ctx, "Moderations OpenAI model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "Moderations OpenAI model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	response, err := c.client.Moderations(ctx, openai.ModerationRequest{
		Input: request.Input,
		Model: request.Model,
	})
	if err != nil {
		logger.Errorf(ctx, "Moderations OpenAI model: %s, error: %v", request.Model, err)
		return res, err
	}

	logger.Infof(ctx, "Moderations OpenAI model: %s finished", request.Model)

	res = model.ModerationResponse{
		Id:      response.Id,
		Model:   response.Model,
		Results: response.Results,
		Usage:   &model.Usage{},
	}

	return res, nil
}
