package sdk

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/sdkerr"
	"github.com/iimeta/fastapi-sdk/util"
)

type ModerationClient struct {
	model    string
	key      string
	baseUrl  string
	path     string
	timeout  time.Duration
	proxyURL string
}

func NewModerationClient(ctx context.Context, model, key, baseURL, path string, timeout time.Duration, proxyUrl string) *ModerationClient {

	logger.Infof(ctx, "NewModerationClient OpenAI model: %s, key: %s", model, key)

	moderationClient := &ModerationClient{
		model:    model,
		key:      key,
		baseUrl:  "https://api.openai.com/v1",
		path:     "/moderations",
		timeout:  timeout,
		proxyURL: proxyUrl,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewModerationClient OpenAI model: %s, baseUrl: %s", model, baseURL)
		moderationClient.baseUrl = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewModerationClient OpenAI model: %s, path: %s", model, path)
		moderationClient.path = path
	}

	return moderationClient
}

func (c *ModerationClient) TextModerations(ctx context.Context, request model.ModerationRequest) (res model.ModerationResponse, err error) {

	logger.Infof(ctx, "TextModerations OpenAI model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "TextModerations OpenAI model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	header := make(map[string]string)
	header["Authorization"] = "Bearer " + c.key

	response := model.ModerationResponse{}
	if _, err = util.HttpPost(ctx, c.baseUrl+c.path, header, request, &response, c.timeout, c.proxyURL, c.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "TextModerations OpenAI model: %s, error: %v", request.Model, err)
		return res, err
	}

	logger.Infof(ctx, "TextModerations OpenAI model: %s finished", request.Model)

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

func (c *ModerationClient) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return sdkerr.NewRequestError(500, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)))
}

func (c *ModerationClient) apiErrorHandler(response *model.XfyunChatCompletionRes) error {
	return sdkerr.NewApiError(500, response.Header.Code, gjson.MustEncodeString(response), "api_error", "")
}
