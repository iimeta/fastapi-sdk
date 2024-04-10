package util

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/iimeta/fastapi-sdk/logger"
)

func HttpGet(ctx context.Context, url string, header map[string]string, data g.Map, result interface{}, proxyURL ...string) error {

	logger.Infof(ctx, "HttpGet url: %s, header: %+v, data: %s, proxyURL: %v", url, header, gjson.MustEncodeString(data), proxyURL)

	client := g.Client()

	if header != nil {
		client.SetHeaderMap(header)
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		client.SetProxy(proxyURL[0])
	}

	response, err := client.Get(ctx, url, data)

	if response != nil {
		defer func() {
			if err := response.Close(); err != nil {
				logger.Error(ctx, err)
			}
		}()
	}

	if err != nil {
		logger.Error(ctx, err)
		return err
	}

	bytes := response.ReadAll()
	logger.Infof(ctx, "HttpGet url: %s, statusCode: %d, header: %+v, data: %s, response: %s", url, response.StatusCode, header, gjson.MustEncodeString(data), string(bytes))

	if bytes != nil && len(bytes) > 0 {
		if err = gjson.Unmarshal(bytes, result); err != nil {
			logger.Error(ctx, err)
			return err
		}
	}

	return nil
}

func HttpPostJson(ctx context.Context, url string, header map[string]string, data, result interface{}, proxyURL ...string) error {

	logger.Infof(ctx, "HttpPostJson url: %s, header: %+v, data: %s, proxyURL: %v", url, header, gjson.MustEncodeString(data), proxyURL)

	client := g.Client()

	if header != nil {
		client.SetHeaderMap(header)
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		client.SetProxy(proxyURL[0])
	}

	response, err := client.ContentJson().Post(ctx, url, data)

	if response != nil {
		defer func() {
			if err := response.Close(); err != nil {
				logger.Error(ctx, err)
			}
		}()
	}

	if err != nil {
		logger.Error(ctx, err)
		return err
	}

	bytes := response.ReadAll()
	logger.Infof(ctx, "HttpPostJson url: %s, statusCode: %d, header: %+v, data: %s, response: %s", url, response.StatusCode, header, gjson.MustEncodeString(data), string(bytes))

	if bytes != nil && len(bytes) > 0 {
		if err = gjson.Unmarshal(bytes, result); err != nil {
			logger.Error(ctx, err)
			return err
		}
	}

	return nil
}

func HttpPost(ctx context.Context, url string, header map[string]string, data, result interface{}, proxyURL ...string) error {

	logger.Infof(ctx, "HttpPost url: %s, header: %+v, data: %+v, proxyURL: %v", url, header, data, proxyURL)

	client := g.Client()

	if header != nil {
		client.SetHeaderMap(header)
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		client.SetProxy(proxyURL[0])
	}

	response, err := client.Post(ctx, url, data)

	defer func() {
		if err := response.Close(); err != nil {
			logger.Error(ctx, err)
		}
	}()

	if err != nil {
		logger.Error(ctx, err)
		return err
	}

	bytes := response.ReadAll()
	logger.Infof(ctx, "HttpPost url: %s, header: %+v, data: %+v, response: %s", url, header, data, string(bytes))

	if bytes != nil && len(bytes) > 0 {
		if err = gjson.Unmarshal(bytes, result); err != nil {
			logger.Error(ctx, err)
			return err
		}
	}

	return nil
}
