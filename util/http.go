package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/iimeta/fastapi-sdk/logger"
)

func HttpGet(ctx context.Context, url string, header map[string]string, data g.Map, result interface{}, proxyURL string) ([]byte, error) {

	logger.Debugf(ctx, "HttpGet url: %s, header: %+v, data: %s, proxyURL: %s", url, header, gjson.MustEncodeString(data), proxyURL)

	client := g.Client()

	if header != nil {
		client.SetHeaderMap(header)
	}

	if proxyURL != "" {
		client.SetProxy(proxyURL)
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
		logger.Errorf(ctx, "HttpGet url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", url, header, gjson.MustEncodeString(data), proxyURL, err)
		return nil, err
	}

	bytes := response.ReadAll()
	logger.Debugf(ctx, "HttpGet url: %s, statusCode: %d, header: %+v, data: %s, proxyURL: %s, response: %s", url, response.StatusCode, header, gjson.MustEncodeString(data), proxyURL, string(bytes))

	if bytes != nil && len(bytes) > 0 {
		if err = gjson.Unmarshal(bytes, result); err != nil {
			logger.Errorf(ctx, "HttpGet url: %s, statusCode: %d, header: %+v, data: %s, proxyURL: %s, response: %s, error: %v", url, response.StatusCode, header, gjson.MustEncodeString(data), proxyURL, string(bytes), err)
			return bytes, err
		}
	}

	return bytes, nil
}

func HttpPost(ctx context.Context, url string, header map[string]string, data, result interface{}, proxyURL string) ([]byte, error) {

	logger.Debugf(ctx, "HttpPost url: %s, header: %+v, data: %s, proxyURL: %s", url, header, gjson.MustEncodeString(data), proxyURL)

	client := g.Client()

	if header != nil {
		client.SetHeaderMap(header)
	}

	if proxyURL != "" {
		client.SetProxy(proxyURL)
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
		logger.Errorf(ctx, "HttpPost url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", url, header, gjson.MustEncodeString(data), proxyURL, err)
		return nil, err
	}

	bytes := response.ReadAll()
	logger.Debugf(ctx, "HttpPost url: %s, statusCode: %d, header: %+v, data: %s, proxyURL: %s, response: %s", url, response.StatusCode, header, gjson.MustEncodeString(data), proxyURL, string(bytes))

	if bytes != nil && len(bytes) > 0 {
		if err = gjson.Unmarshal(bytes, result); err != nil {
			logger.Errorf(ctx, "HttpPost url: %s, statusCode: %d, header: %+v, data: %s, proxyURL: %s, response: %s, error: %v", url, response.StatusCode, header, gjson.MustEncodeString(data), proxyURL, string(bytes), err)
			return bytes, errors.New(fmt.Sprintf("response: %s, error: %v", bytes, err))
		}
	}

	return bytes, nil
}
