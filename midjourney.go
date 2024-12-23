package sdk

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"net/http"
)

type MidjourneyClient struct {
	baseURL         string
	path            string
	apiSecret       string
	apiSecretHeader string
	proxyURL        string
	method          string
}

func NewMidjourneyClient(ctx context.Context, baseURL, path, apiSecret, apiSecretHeader, method string, proxyURL ...string) *MidjourneyClient {

	client := &MidjourneyClient{
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

func (c *MidjourneyClient) Request(ctx context.Context, data interface{}) (res model.MidjourneyResponse, err error) {

	logger.Infof(ctx, "Midjourney Request data: %s start", gjson.MustEncodeString(data))

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "Midjourney Request data: %s totalTime: %d ms", gjson.MustEncodeString(data), gtime.TimestampMilli()-now)
	}()

	if res.Response, err = request(ctx, c.method, c.baseURL+c.path, c.apiSecretHeader, c.apiSecret, data, c.proxyURL); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	return res, nil
}

func request(ctx context.Context, method, url, apiSecretHeader, apiSecret string, data interface{}, proxyURL string) ([]byte, error) {

	logger.Debugf(ctx, "Midjourney Request url: %s, apiSecretHeader: %s, apiSecret: %s, data: %s, proxyURL: %v", url, apiSecretHeader, apiSecret, gjson.MustEncodeString(data), proxyURL)

	var (
		client   = g.Client()
		response *gclient.Response
		err      error
	)

	client.SetHeaderMap(g.MapStrStr{apiSecretHeader: apiSecret})

	if proxyURL != "" {
		client.SetProxy(proxyURL)
	}

	if method == http.MethodGet {
		response, err = client.Get(ctx, url, data)
	} else {
		response, err = client.ContentJson().Post(ctx, url, data)
	}

	if response != nil {
		defer func() {
			if err := response.Close(); err != nil {
				logger.Error(ctx, err)
			}
		}()
	}

	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	bytes := response.ReadAll()
	logger.Debugf(ctx, "Midjourney Request url: %s, statusCode: %d, apiSecretHeader: %s, apiSecret: %s, data: %s, response: %s", url, response.StatusCode, apiSecretHeader, apiSecret, gjson.MustEncodeString(data), string(bytes))

	return bytes, nil
}
