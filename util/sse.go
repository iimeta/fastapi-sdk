package util

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/iimeta/fastapi-sdk/logger"
	"io"
	"net/http"
	"time"
)

var (
	headerData        = []byte("data: ")
	headerDataNoSpace = []byte("data:")
	errorPrefix       = []byte(`data: {"errors":`)
)

var (
	ErrTooManyEmptyStreamMessages = errors.New("stream has sent too many empty messages")
)

type RequestErrorHandler func(ctx context.Context, response *gclient.Response) (err error)

type StreamReader struct {
	reader             *bufio.Reader
	response           *gclient.Response
	emptyMessagesLimit uint
	isFinished         bool
}

func SSEClient(ctx context.Context, url string, header map[string]string, data interface{}, proxyURL string, requestErrorHandler RequestErrorHandler) (stream *StreamReader, err error) {

	logger.Debugf(ctx, "SSEClient url: %s, header: %+v, data: %s, proxyURL: %s", url, header, gjson.MustEncodeString(data), proxyURL)

	client := g.Client().Timeout(600 * time.Second)

	if header != nil {
		client.SetHeaderMap(header)
	}

	if proxyURL != "" {
		client.SetProxy(proxyURL)
	}

	client.SetHeader("Accept", "text/event-stream")
	client.SetHeader("Cache-Control", "no-cache")
	client.SetHeader("Connection", "keep-alive")
	client.SetHeader("Content-Type", "application/json")

	response, err := client.Post(ctx, url, data)
	if err != nil {
		logger.Errorf(ctx, "SSEClient url: %s, header: %+v, data: %s, proxyURL: %s, err: %v", url, header, gjson.MustEncodeString(data), proxyURL, err)
		if response != nil {
			if err := response.Close(); err != nil {
				logger.Error(ctx, err)
			}
		}
		return nil, err
	}

	if isFailureStatusCode(response) {

		defer func() {
			if err := response.Close(); err != nil {
				logger.Error(ctx, err)
			}
		}()

		if requestErrorHandler != nil {
			return nil, requestErrorHandler(ctx, response)
		}

		return nil, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, response.ReadAllString()))
	}

	stream = &StreamReader{
		reader:             bufio.NewReader(response.Body),
		response:           response,
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

func isFailureStatusCode(resp *gclient.Response) bool {
	return resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest
}
