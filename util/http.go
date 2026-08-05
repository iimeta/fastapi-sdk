package util

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
)

func HttpDo(ctx context.Context, method, rawURL string, header map[string]string, data, result any, timeout time.Duration, proxyURL string, requestErrorHandler RequestErrorHandler) ([]byte, http.Header, error) {

	logger.Debugf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s", method, rawURL, header, mustEncodeString(data), proxyURL)

	client := &http.Client{
		Timeout: timeout,
	}

	if proxyURL != "" {
		if proxyUrl, err := url.Parse(proxyURL); err != nil {
			logger.Errorf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", method, rawURL, header, mustEncodeString(data), proxyURL, err)
			return nil, nil, err
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
		return nil, nil, err
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

	decompressResponse(response)

	if err != nil {

		if response != nil {

			bytes, _ := io.ReadAll(response.Body)

			logger.Errorf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s, statusCode: %d, header: %+v, response: %s, error: %v", method, rawURL, header, mustEncodeString(data), proxyURL, response.StatusCode, response.Header, bytes, err)

			if err := response.Body.Close(); err != nil {
				logger.Error(ctx, err)
			}

			return nil, nil, err
		}

		logger.Errorf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", method, rawURL, header, mustEncodeString(data), proxyURL, err)

		return nil, nil, err
	}

	if isFailureStatusCode(response) {

		defer func() {
			if err := response.Body.Close(); err != nil {
				logger.Error(ctx, err)
			}
		}()

		if requestErrorHandler != nil {
			return nil, nil, requestErrorHandler(ctx, response)
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			logger.Errorf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s, statusCode: %d, header: %+v, error: %v", method, rawURL, header, mustEncodeString(data), proxyURL, response.StatusCode, response.Header, err)
			return nil, nil, err
		}

		return nil, nil, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes))
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			logger.Error(ctx, err)
		}
	}()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Errorf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s, statusCode: %d, header: %+v, error: %v", method, rawURL, header, mustEncodeString(data), proxyURL, response.StatusCode, response.Header, err)
		return nil, nil, err
	}

	logger.Debugf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s, statusCode: %d, header: %+v, response: %s", method, rawURL, header, mustEncodeString(data), proxyURL, response.StatusCode, response.Header, bytes)

	if bytes != nil && len(bytes) > 0 && result != nil {
		if err = json.Unmarshal(bytes, result); err != nil {
			logger.Errorf(ctx, "method: %s, url: %s, header: %+v, data: %s, proxyURL: %s, statusCode: %d, header: %+v, response: %s, error: %v", method, rawURL, header, mustEncodeString(data), proxyURL, response.StatusCode, response.Header, bytes, err)
			return bytes, nil, errors.New(fmt.Sprintf("response: %s, error: %v", bytes, err))
		}
	}

	return bytes, response.Header, nil
}

func HttpGet(ctx context.Context, rawURL string, header map[string]string, data, result any, timeout time.Duration, proxyURL string, requestErrorHandler RequestErrorHandler) ([]byte, http.Header, error) {
	return HttpDo(ctx, http.MethodGet, rawURL, header, data, result, timeout, proxyURL, requestErrorHandler)
}

func HttpPost(ctx context.Context, rawURL string, header map[string]string, data, result any, timeout time.Duration, proxyURL string, requestErrorHandler RequestErrorHandler) ([]byte, http.Header, error) {
	return HttpDo(ctx, http.MethodPost, rawURL, header, data, result, timeout, proxyURL, requestErrorHandler)
}

func HttpDelete(ctx context.Context, rawURL string, header map[string]string, data, result any, timeout time.Duration, proxyURL string, requestErrorHandler RequestErrorHandler) ([]byte, http.Header, error) {
	return HttpDo(ctx, http.MethodDelete, rawURL, header, data, result, timeout, proxyURL, requestErrorHandler)
}

// 当调用方自带 Accept-Encoding: gzip 时, Go 的 Transport 不会自动解压,
// 此时 response.Body 是 gzip 原始字节, 日志会乱码且 json 解析失败, 故在此手动解压.
func decompressResponse(response *http.Response) {

	if response == nil || response.Body == nil {
		return
	}

	if !strings.EqualFold(strings.TrimSpace(response.Header.Get("Content-Encoding")), "gzip") {
		return
	}

	buf := bufio.NewReader(response.Body)

	// 校验 gzip magic number, 避免响应头声明有误时损坏原始内容
	if magic, err := buf.Peek(2); err != nil || magic[0] != 0x1f || magic[1] != 0x8b {
		response.Body = &readCloser{Reader: buf, Closer: response.Body}
		return
	}

	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		response.Body = &readCloser{Reader: buf, Closer: response.Body}
		return
	}

	response.Body = &readCloser{Reader: gzipReader, Closer: response.Body}
	response.Header.Del("Content-Encoding")
	response.Header.Del("Content-Length")
	response.ContentLength = -1
	response.Uncompressed = true
}

// readCloser 用解压后的 Reader 读取, 但仍关闭原始的 response.Body
type readCloser struct {
	io.Reader
	io.Closer
}
