package sdk

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

type Client struct {
	ApiSecret       string
	ApiSecretHeader string
	FetchUrl        string
	baseURL         string
	path            string
	proxyURL        string
}

func NewMidjourneyClient(ctx context.Context, baseURL, path, apiSecret, apiSecretHeader string, proxyURL ...string) *Client {

	client := &Client{
		ApiSecret:       apiSecret,
		ApiSecretHeader: apiSecretHeader,
		FetchUrl:        baseURL + "/mj/task/${taskId}/fetch",
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewMidjourneyClient baseURL: %s", baseURL)
		client.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewMidjourneyClient path: %s", path)
		client.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewMidjourneyClient proxyURL: %s", proxyURL[0])
		client.proxyURL = proxyURL[0]
	}

	return client
}

func (c *Client) Main(ctx context.Context, request interface{}) (res model.MidjourneyResponse, err error) {

	logger.Infof(ctx, "Midjourney Main request: %s start", gjson.MustEncodeString(request))

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney Main request: %s totalTime: %d ms", gjson.MustEncodeString(request), gtime.Now().UnixMilli()-now)
	}()

	if res.Response, err = util.HttpPostByMidjourney(ctx, c.baseURL+c.path, g.MapStrStr{c.ApiSecretHeader: c.ApiSecret}, request, c.proxyURL); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	return res, nil
}

func (c *Client) Fetch(ctx context.Context, request model.MidjourneyProxyRequest) (res model.MidjourneyProxyFetchResponse, err error) {

	logger.Infof(ctx, "Midjourney Fetch taskId: %s start", request.TaskId)

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney Fetch taskId: %s totalTime: %d ms", request.TaskId, gtime.Now().UnixMilli()-now)
	}()

	fetchUrl := gstr.Replace(c.FetchUrl, "${taskId}", request.TaskId, -1)

	if err = util.HttpGet(ctx, fetchUrl, g.MapStrStr{c.ApiSecretHeader: c.ApiSecret}, nil, &res, ""); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	logger.Infof(ctx, "Midjourney Fetch Response: %s", gjson.MustEncodeString(res))

	return res, nil
}
