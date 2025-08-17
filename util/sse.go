package util

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
)

var (
	headerData        = []byte("data: ")
	headerDataNoSpace = []byte("data:")
	errorPrefix       = []byte(`data: {"errors":`)
)

var (
	ErrTooManyEmptyStreamMessages = errors.New("stream has sent too many empty messages")
)

type RequestErrorHandler func(ctx context.Context, response *http.Response) (err error)

type StreamReader struct {
	reader             *bufio.Reader
	response           *http.Response
	emptyMessagesLimit uint
	isFinished         bool
	ReqTime            string
}

func SSEClient(ctx context.Context, rawURL string, header map[string]string, data []byte, proxyURL string, requestErrorHandler RequestErrorHandler) (stream *StreamReader, err error) {

	logger.Debugf(ctx, "SSEClient url: %s, header: %+v, data: %s, proxyURL: %s", rawURL, header, gjson.MustEncodeString(data), proxyURL)

	client := &http.Client{
		Timeout: 600 * time.Second,
	}

	request, err := http.NewRequest("POST", rawURL, bytes.NewBuffer(data))
	if err != nil {
		logger.Errorf(ctx, "SSEClient url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
		return nil, err
	}

	reqTime := gtime.TimestampMilliStr()

	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-request-time", reqTime)

	if header != nil {
		for k, v := range header {
			request.Header.Set(k, v)
		}
	}

	if proxyURL != "" {
		if proxyUrl, err := url.Parse(proxyURL); err != nil {
			logger.Errorf(ctx, "SSEClient url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
			return nil, err
		} else {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
		}
	}

	response, err := client.Do(request)
	if err != nil {
		logger.Errorf(ctx, "SSEClient url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
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
			logger.Errorf(ctx, "SSEClient url: %s, header: %+v, data: %s, proxyURL: %s, error: %v", rawURL, header, gjson.MustEncodeString(data), proxyURL, err)
			return nil, err
		}

		return nil, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes))
	}

	stream = &StreamReader{
		reader:             bufio.NewReader(response.Body),
		response:           response,
		emptyMessagesLimit: 300,
		ReqTime:            reqTime,
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

	var (
		emptyMessagesCount uint
		hasErrorPrefix     bool
	)

	for {

		rawLine, readErr := stream.reader.ReadBytes('\n')
		if readErr != nil || hasErrorPrefix {
			return rawLine, readErr
		}

		noSpaceLine := bytes.TrimSpace(rawLine)
		if bytes.HasPrefix(noSpaceLine, errorPrefix) {
			hasErrorPrefix = true
		}

		if (!bytes.HasPrefix(noSpaceLine, headerData) && !bytes.HasPrefix(noSpaceLine, headerDataNoSpace)) || hasErrorPrefix {

			emptyMessagesCount++
			if emptyMessagesCount > stream.emptyMessagesLimit {
				return nil, ErrTooManyEmptyStreamMessages
			}

			continue
		}

		noPrefixLine := bytes.TrimPrefix(bytes.TrimPrefix(noSpaceLine, headerData), headerDataNoSpace)
		if string(noPrefixLine) == "[DONE]" {
			stream.isFinished = true
			return nil, io.EOF
		}

		return noPrefixLine, nil
	}
}

func (stream *StreamReader) Close() error {
	return stream.response.Body.Close()
}

func isFailureStatusCode(resp *http.Response) bool {
	return resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest
}
