package util

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/iimeta/fastapi-sdk/logger"
)

func HttpGet(ctx context.Context, rawURL string, header map[string]string, data, result any, proxyURL string, requestErrorHandler RequestErrorHandler) ([]byte, error) {

	logger.Debugf(ctx, "HttpGet url: %s, header: %+v, data: %s, proxyURL: %s", rawURL, header, gjson.MustEncodeString(data), proxyURL)

	client := &http.Client{
		Timeout: 600 * time.Second,
	}

	var bodyReader io.Reader

	if data != nil {
		if v, ok := data.([]byte); ok {
			bodyReader = bytes.NewBuffer(v)
		} else if v, ok := data.(io.Reader); ok {
			bodyReader = v
		} else {
			bodyReader = bytes.NewBuffer(gjson.MustEncode(data))
		}
	}

	request, err := http.NewRequest("GET", rawURL, bodyReader)
	if err != nil {
		logger.Errorf(ctx, "HttpGet url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		logger.Errorf(ctx, "HttpGet url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
		if response != nil {
			if err := response.Body.Close(); err != nil {
				logger.Error(ctx, err)
			}
		}
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	if isFailureStatusCode(response) {

		defer func() {
			if err := response.Body.Close(); err != nil {
				logger.Error(ctx, err)
			}
		}()

		if requestErrorHandler != nil {
			return nil, requestErrorHandler(ctx, response)
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			logger.Errorf(ctx, "HttpGet url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
			return nil, err
		}

		return nil, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes))
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			logger.Error(ctx, err)
		}
	}()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Errorf(ctx, "HttpGet url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
		return nil, err
	}

	logger.Debugf(ctx, "HttpGet url: %s, statusCode: %d, header: %+v, data: %s, proxyURL: %s, response: %s", rawURL, response.StatusCode, header, gjson.MustEncodeString(data), proxyURL, string(bytes))

	if bytes != nil && len(bytes) > 0 {
		if err = json.Unmarshal(bytes, result); err != nil {
			logger.Errorf(ctx, "HttpGet url: %s, statusCode: %d, header: %+v, data: %s, proxyURL: %s, response: %s, error: %v", rawURL, response.StatusCode, header, gjson.MustEncodeString(data), proxyURL, string(bytes), err)
			return bytes, err
		}
	}

	return bytes, nil
}

func HttpPost(ctx context.Context, rawURL string, header map[string]string, data, result any, proxyURL string, requestErrorHandler RequestErrorHandler) ([]byte, error) {

	logger.Debugf(ctx, "HttpPost url: %s, header: %+v, data: %s, proxyURL: %s", rawURL, header, gjson.MustEncodeString(data), proxyURL)

	client := &http.Client{
		Timeout: 600 * time.Second,
	}

	var bodyReader io.Reader

	if data != nil {
		if v, ok := data.([]byte); ok {
			bodyReader = bytes.NewBuffer(v)
		} else if v, ok := data.(io.Reader); ok {
			bodyReader = v
		} else {
			bodyReader = bytes.NewBuffer(gjson.MustEncode(data))
		}
	}

	request, err := http.NewRequest("POST", rawURL, bodyReader)
	if err != nil {
		logger.Errorf(ctx, "HttpPost url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
		return nil, err
	}

	contentType := request.Header.Get("Content-Type")
	if contentType == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	if header != nil {
		for k, v := range header {
			request.Header.Set(k, v)
		}
	}

	if proxyURL != "" {
		if proxyUrl, err := url.Parse(proxyURL); err != nil {
			logger.Errorf(ctx, "HttpPost url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
			return nil, err
		} else {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
		}
	}

	response, err := client.Do(request)
	if err != nil {
		logger.Errorf(ctx, "HttpPost url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
		if response != nil {
			if err := response.Body.Close(); err != nil {
				logger.Error(ctx, err)
			}
		}
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	if isFailureStatusCode(response) {

		defer func() {
			if err := response.Body.Close(); err != nil {
				logger.Error(ctx, err)
			}
		}()

		if requestErrorHandler != nil {
			return nil, requestErrorHandler(ctx, response)
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			logger.Errorf(ctx, "HttpPost url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
			return nil, err
		}

		return nil, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes))
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			logger.Error(ctx, err)
		}
	}()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Errorf(ctx, "HttpPost url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
		return nil, err
	}

	logger.Debugf(ctx, "HttpPost url: %s, statusCode: %d, header: %+v, data: %s, proxyURL: %s, response: %s", rawURL, response.StatusCode, header, gjson.MustEncodeString(data), proxyURL, string(bytes))

	if bytes != nil && len(bytes) > 0 && result != nil {
		if err = json.Unmarshal(bytes, result); err != nil {
			logger.Errorf(ctx, "HttpPost url: %s, statusCode: %d, header: %+v, data: %s, proxyURL: %s, response: %s, error: %v", rawURL, response.StatusCode, header, gjson.MustEncodeString(data), proxyURL, string(bytes), err)
			return bytes, errors.New(fmt.Sprintf("response: %s, error: %v", bytes, err))
		}
	}

	return bytes, nil
}
