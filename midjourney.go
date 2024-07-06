package sdk

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
	"net/http"
)

type Client struct {
	baseURL         string
	path            string
	apiSecret       string
	apiSecretHeader string
	proxyURL        string
	method          string
}

func NewMidjourneyClient(ctx context.Context, baseURL, path, apiSecret, apiSecretHeader, method string, proxyURL ...string) *Client {

	client := &Client{
		baseURL:         baseURL,
		path:            path,
		apiSecret:       apiSecret,
		apiSecretHeader: apiSecretHeader,
		method:          method,
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

	if c.method == http.MethodGet {
		if res.Response, err = util.HttpGetByMidjourney(ctx, c.baseURL+c.path, g.MapStrStr{c.apiSecretHeader: c.apiSecret}, nil, c.proxyURL); err != nil {
			logger.Error(ctx, err)
			return res, err
		}
	} else {
		if res.Response, err = util.HttpPostByMidjourney(ctx, c.baseURL+c.path, g.MapStrStr{c.apiSecretHeader: c.apiSecret}, request, c.proxyURL); err != nil {
			logger.Error(ctx, err)
			return res, err
		}
	}

	return res, nil
}
