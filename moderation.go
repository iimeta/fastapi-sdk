package sdk

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

type ModerationClient struct {
	model    string
	key      string
	baseURL  string
	path     string
	proxyURL string
}

func NewModerationClient(ctx context.Context, model, key, baseURL, path string, proxyURL ...string) *ModerationClient {

	logger.Infof(ctx, "NewModerationClient OpenAI model: %s, key: %s", model, key)

	moderationClient := &ModerationClient{
		model:   model,
		key:     key,
		baseURL: "https://api.openai.com/v1",
		path:    "/moderations",
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewModerationClient OpenAI model: %s, baseURL: %s", model, baseURL)
		moderationClient.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewModerationClient OpenAI model: %s, path: %s", model, path)
		moderationClient.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewModerationClient OpenAI model: %s, proxyURL: %s", model, proxyURL[0])
		moderationClient.proxyURL = proxyURL[0]
	}

	return moderationClient
}

func (c *ModerationClient) Moderations(ctx context.Context, request model.ModerationRequest) (res model.ModerationResponse, err error) {

	logger.Infof(ctx, "Moderations OpenAI model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "Moderations OpenAI model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	header := make(map[string]string)
	header["Authorization"] = "Bearer " + c.key

	response := new(model.ModerationResponse)

	if err = util.HttpPost(ctx, c.baseURL+c.path, header, request, &response, c.proxyURL); err != nil {
		logger.Errorf(ctx, "Moderations OpenAI model: %s, error: %v", request.Model, err)
		return res, err
	}

	logger.Infof(ctx, "Moderations OpenAI model: %s finished", request.Model)

	if response.Error != nil {
		return res, errors.New(gjson.MustEncodeString(response.Error))
	}

	res = model.ModerationResponse{
		Id:      response.Id,
		Model:   response.Model,
		Results: response.Results,
		Usage:   &model.Usage{},
	}

	return res, nil
}
