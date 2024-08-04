package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/iimeta/fastapi-sdk/logger"
)

func HttpGet(ctx context.Context, url string, header map[string]string, data g.Map, result interface{}, proxyURL string) error {

	logger.Debugf(ctx, "HttpGet url: %s, header: %+v, data: %s, proxyURL: %v", url, header, gjson.MustEncodeString(data), proxyURL)

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
		logger.Error(ctx, err)
		return err
	}

	bytes := response.ReadAll()
	logger.Debugf(ctx, "HttpGet url: %s, statusCode: %d, header: %+v, data: %s, response: %s", url, response.StatusCode, header, gjson.MustEncodeString(data), string(bytes))

	if bytes != nil && len(bytes) > 0 {
		if err = gjson.Unmarshal(bytes, result); err != nil {
			logger.Error(ctx, err)
			return err
		}
	}

	return nil
}

func HttpPost(ctx context.Context, url string, header map[string]string, data, result interface{}, proxyURL string) error {

	logger.Debugf(ctx, "HttpPost url: %s, header: %+v, data: %s, proxyURL: %v", url, header, gjson.MustEncodeString(data), proxyURL)

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
		logger.Error(ctx, err)
		return err
	}

	bytes := response.ReadAll()
	logger.Debugf(ctx, "HttpPost url: %s, statusCode: %d, header: %+v, data: %s, response: %s", url, response.StatusCode, header, gjson.MustEncodeString(data), string(bytes))

	if bytes != nil && len(bytes) > 0 {
		if err = gjson.Unmarshal(bytes, result); err != nil {
			logger.Error(ctx, err)
			return errors.New(fmt.Sprintf("response: %s, err: %v", bytes, err))
		}
	}

	return nil
}
