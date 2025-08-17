package util

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
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

	reqTime := gtime.TimestampMilliStr()

	client.SetHeader("x-request-time", reqTime)

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

	if bytes != nil && len(bytes) > 0 && result != nil {
		if err = gjson.Unmarshal(bytes, result); err != nil {
			logger.Errorf(ctx, "HttpPost url: %s, statusCode: %d, header: %+v, data: %s, proxyURL: %s, response: %s, error: %v", url, response.StatusCode, header, gjson.MustEncodeString(data), proxyURL, string(bytes), err)
			return bytes, errors.New(fmt.Sprintf("response: %s, error: %v", bytes, err))
		}
	}

	end := gtime.TimestampMilli()
	resTime := response.Header.Get("x-response-time")
	resTotalTime := response.Header.Get("x-response-total-time")
	fmt.Println(reqTime, resTime, end, end-gconv.Int64(resTime), end-gconv.Int64(reqTime)-gconv.Int64(resTotalTime), "end")

	return bytes, nil
}

func HttpPostNew(ctx context.Context, rawURL string, header map[string]string, data []byte, result interface{}, proxyURL string) ([]byte, error) {

	logger.Debugf(ctx, "HttpPostNew url: %s, header: %+v, data: %s, proxyURL: %s", rawURL, header, gjson.MustEncodeString(data), proxyURL)

	client := &http.Client{
		Timeout: 600 * time.Second,
	}

	request, err := http.NewRequest("POST", rawURL, bytes.NewBuffer(data))
	if err != nil {
		logger.Errorf(ctx, "HttpPostNew url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	if header != nil {
		for k, v := range header {
			request.Header.Set(k, v)
		}
	}

	if proxyURL != "" {
		if proxyUrl, err := url.Parse(proxyURL); err != nil {
			logger.Errorf(ctx, "HttpPostNew url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
			return nil, err
		} else {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
		}
	}

	reqTime := gtime.TimestampMilliStr()

	request.Header.Set("x-request-time", reqTime)

	response, err := client.Do(request)
	if response != nil {
		defer func() {
			if err := response.Body.Close(); err != nil {
				logger.Error(ctx, err)
			}
		}()
	}

	if err != nil {
		logger.Errorf(ctx, "HttpPostNew url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
		return nil, err
	}

	if response == nil {
		return []byte{}, nil
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Errorf(ctx, "HttpPostNew url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
		return nil, err
	}

	logger.Debugf(ctx, "HttpPostNew url: %s, statusCode: %d, header: %+v, data: %s, proxyURL: %s, response: %s", rawURL, response.StatusCode, header, gjson.MustEncodeString(data), proxyURL, string(bytes))

	if bytes != nil && len(bytes) > 0 && result != nil {
		if err = gjson.Unmarshal(bytes, result); err != nil {
			logger.Errorf(ctx, "HttpPostNew url: %s, statusCode: %d, header: %+v, data: %s, proxyURL: %s, response: %s, error: %v", rawURL, response.StatusCode, header, gjson.MustEncodeString(data), proxyURL, string(bytes), err)
			return bytes, errors.New(fmt.Sprintf("response: %s, error: %v", bytes, err))
		}
	}

	end := gtime.TimestampMilli()
	resTime := response.Header.Get("x-response-time")
	resTotalTime := response.Header.Get("x-response-total-time")
	fmt.Println(reqTime, resTime, end, end-gconv.Int64(resTime), end-gconv.Int64(reqTime)-gconv.Int64(resTotalTime), "end")

	return bytes, nil
}
