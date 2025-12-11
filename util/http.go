package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/iimeta/fastapi-sdk/errors"
	"github.com/iimeta/fastapi-sdk/logger"
)

// d
func HttpDo(ctx context.Context, method, rawURL string, header map[string]string, data, result any, timeout time.Duration, proxyURL string, requestErrorHandler RequestErrorHandler) ([]byte, error) {

	logger.Debugf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s", method, rawURL, header, mustEncodeString(data), proxyURL)

	client := &http.Client{
		Timeout: timeout,
	}

	if proxyURL != "" {
		if proxyUrl, err := url.Parse(proxyURL); err != nil {
			logger.Errorf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", method, rawURL, header, mustEncodeString(data), proxyURL, err)
			return nil, err
		} else {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
		}
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

	request, err := http.NewRequestWithContext(ctx, method, rawURL, bodyReader)
	if err != nil {
		logger.Errorf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", method, rawURL, header, mustEncodeString(data), proxyURL, err)
		return nil, err
	}

	contentType := request.Header.Get("Content-Type")
	if contentType == "" && method == http.MethodPost {
		request.Header.Set("Content-Type", "application/json")
	}

	request.Header.Set("Trace-Id", gtrace.GetTraceID(ctx))

	if header != nil {
		for k, v := range header {
			request.Header.Set(k, v)
		}
	}

	response, err := client.Do(request)
	if err != nil {

		if response != nil {

			bytes, _ := io.ReadAll(response.Body)

			logger.Errorf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s, statusCode: %d, header: %+v, response: %s, error: %v", method, rawURL, header, mustEncodeString(data), proxyURL, response.StatusCode, response.Header, bytes, err)

			if err := response.Body.Close(); err != nil {
				logger.Error(ctx, err)
			}

			return nil, err
		}

		logger.Errorf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", method, rawURL, header, mustEncodeString(data), proxyURL, err)

		return nil, err
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
			logger.Errorf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s, statusCode: %d, header: %+v, error: %v", method, rawURL, header, mustEncodeString(data), proxyURL, response.StatusCode, response.Header, err)
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
		logger.Errorf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s, statusCode: %d, header: %+v, error: %v", method, rawURL, header, mustEncodeString(data), proxyURL, response.StatusCode, response.Header, err)
		return nil, err
	}

	logger.Debugf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s, statusCode: %d, header: %+v, response: %s", method, rawURL, header, mustEncodeString(data), proxyURL, response.StatusCode, response.Header, bytes)

	if bytes != nil && len(bytes) > 0 && result != nil {
		if err = json.Unmarshal(bytes, result); err != nil {
			logger.Errorf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s, statusCode: %d, header: %+v, response: %s, error: %v", method, rawURL, header, mustEncodeString(data), proxyURL, response.StatusCode, response.Header, bytes, err)
			return bytes, errors.New(fmt.Sprintf("response: %s, error: %v", bytes, err))
		}
	}

	return bytes, nil
}

func HttpGet(ctx context.Context, rawURL string, header map[string]string, data, result any, timeout time.Duration, proxyURL string, requestErrorHandler RequestErrorHandler) ([]byte, error) {
	return HttpDo(ctx, http.MethodGet, rawURL, header, data, result, timeout, proxyURL, requestErrorHandler)
}

func HttpPost(ctx context.Context, rawURL string, header map[string]string, data, result any, timeout time.Duration, proxyURL string, requestErrorHandler RequestErrorHandler) ([]byte, error) {
	return HttpDo(ctx, http.MethodPost, rawURL, header, data, result, timeout, proxyURL, requestErrorHandler)
}

func HttpDelete(ctx context.Context, rawURL string, header map[string]string, data, result any, timeout time.Duration, proxyURL string, requestErrorHandler RequestErrorHandler) ([]byte, error) {
	return HttpDo(ctx, http.MethodDelete, rawURL, header, data, result, timeout, proxyURL, requestErrorHandler)
}
