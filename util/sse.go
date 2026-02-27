package util

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
)

var (
	headerEvent        = []byte("event: ")
	headerEventNoSpace = []byte("event:")
	headerData         = []byte("data: ")
	headerDataNoSpace  = []byte("data:")
	errorPrefix        = []byte(`data: {"errors":`)
)

var (
	ErrTooManyEmptyStreamMessages = errors.New("stream has sent too many empty messages")
)

type RequestErrorHandler func(ctx context.Context, response *http.Response) (err error)

type StreamReader struct {
	Response           *http.Response
	reader             *bufio.Reader
	emptyMessagesLimit uint
	event              string
	isFinished         bool
}

func SSEClient(ctx context.Context, rawURL string, header map[string]string, data any, timeout time.Duration, proxyURL string, requestErrorHandler RequestErrorHandler) (stream *StreamReader, err error) {

	logger.Debugf(ctx, "SSEClient url: %s, header: %+v, data: %s, proxyURL: %s", rawURL, header, mustEncodeString(data), proxyURL)

	client := &http.Client{
		Timeout: timeout,
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

	request, err := http.NewRequestWithContext(ctx, "POST", rawURL, bodyReader)
	if err != nil {
		logger.Errorf(ctx, "SSEClient url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, mustEncodeString(data), proxyURL, err)
		return nil, err
	}

	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Trace-Id", gtrace.GetTraceID(ctx))

	if header != nil {
		for k, v := range header {
			request.Header.Set(k, v)
		}
	}

	if proxyURL != "" {
		if proxyUrl, err := url.Parse(proxyURL); err != nil {
			logger.Errorf(ctx, "SSEClient url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, mustEncodeString(data), proxyURL, err)
			return nil, err
		} else {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
		}
	}

	response, err := client.Do(request)
	if err != nil {
		logger.Errorf(ctx, "SSEClient url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, mustEncodeString(data), proxyURL, err)
		if response != nil {
			if err := response.Body.Close(); err != nil {
				logger.Error(ctx, err)
			}
		}
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
			logger.Errorf(ctx, "SSEClient url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, mustEncodeString(data), proxyURL, err)
			return nil, err
		}

		return nil, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes))
	}

	stream = &StreamReader{
		Response:           response,
		reader:             bufio.NewReader(response.Body),
		emptyMessagesLimit: 300,
	}

	return stream, nil
}

func (stream *StreamReader) Recv() (response []byte, err error) {

	if stream.isFinished {
		return nil, io.EOF
	}

	return stream.processLines()
}

func (stream *StreamReader) processLines() ([]byte, error) {

	var emptyMessagesCount uint

	for {

		rawLine, readErr := stream.reader.ReadBytes('\n')
		if readErr != nil {
			return rawLine, readErr
		}

		line := bytes.TrimSpace(rawLine)
		if len(line) == 0 {
			stream.event = ""
			continue
		}

		if bytes.HasPrefix(line, errorPrefix) {
			return rawLine, fmt.Errorf("received error line: %s", rawLine)
		}

		if bytes.HasPrefix(line, headerEvent) || bytes.HasPrefix(line, headerEventNoSpace) {
			eventVal := bytes.TrimPrefix(bytes.TrimPrefix(line, headerEvent), headerEventNoSpace)
			stream.event = string(bytes.TrimSpace(eventVal))
			continue
		}

		if bytes.HasPrefix(line, headerData) || bytes.HasPrefix(line, headerDataNoSpace) {
			dataVal := bytes.TrimPrefix(bytes.TrimPrefix(line, headerData), headerDataNoSpace)
			trimmedData := bytes.TrimSpace(dataVal)
			if string(trimmedData) == "[DONE]" {
				stream.isFinished = true
				return nil, io.EOF
			}
			return trimmedData, nil
		}

		emptyMessagesCount++
		if emptyMessagesCount > stream.emptyMessagesLimit {
			return nil, ErrTooManyEmptyStreamMessages
		}
	}
}

func (stream *StreamReader) Event() string {
	return stream.event
}

func (stream *StreamReader) Close() error {
	return stream.Response.Body.Close()
}

func isFailureStatusCode(resp *http.Response) bool {
	return resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest
}

func mustEncodeString(data any) string {

	if data != nil {
		if v, ok := data.([]byte); ok {
			return string(v)
		} else if v, ok := data.(string); ok {
			return v
		}
		return gjson.MustEncodeString(data)
	}

	return ""
}
