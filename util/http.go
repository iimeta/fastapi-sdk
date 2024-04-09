package util

import (
	"bufio"
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gorilla/websocket"
	"github.com/iimeta/fastapi-sdk/logger"
	"io"
	"net/http"
	"net/url"
	"time"
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
			if err = response.Close(); err != nil {
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
			if err = response.Close(); err != nil {
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
	if err != nil {
		logger.Error(ctx, err)
		return err
	}

	defer func() {
		err = response.Close()
		if err != nil {
			logger.Error(ctx, err)
		}
	}()

	bytes := response.ReadAll()
	logger.Infof(ctx, "HttpPost url: %s, header: %+v, data: %+v, response: %s", url, header, data, string(bytes))

	if bytes != nil && len(bytes) > 0 {
		err = gjson.Unmarshal(bytes, result)
		if err != nil {
			logger.Error(ctx, err)
			return err
		}
	}

	return nil
}

func WebSocketClient(ctx context.Context, wsURL string, messageType int, message []byte, result chan []byte, proxyURL ...string) (*websocket.Conn, error) {

	logger.Infof(ctx, "WebSocketClient wsURL: %s", wsURL)

	client := gclient.NewWebSocket()

	client.HandshakeTimeout = 60 * time.Second // 设置超时时间
	//client.TLSClientConfig = &tls.Config{}   // 设置 tls 配置

	// 设置代理
	if len(proxyURL) > 0 && proxyURL[0] != "" {
		if proxyUrl, err := url.Parse(proxyURL[0]); err != nil {
			logger.Error(ctx, err)
		} else {
			client.Proxy = http.ProxyURL(proxyUrl)
		}
	}

	conn, _, err := client.Dial(wsURL, nil)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	err = conn.WriteMessage(messageType, message)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	_ = grpool.AddWithRecover(ctx, func(ctx context.Context) {
		for {

			_, message, err := conn.ReadMessage()
			if err != nil && websocket.IsUnexpectedCloseError(err) {
				logger.Error(ctx, err)
				return
			}

			_ = grpool.AddWithRecover(ctx, func(ctx context.Context) {
				result <- message
			}, nil)
		}
	}, nil)

	return conn, nil
}

func SSEClient(ctx context.Context, method, url string, header map[string]string, data interface{}, result chan []byte) error {

	logger.Infof(ctx, "SSEClient method: %s, url: %s, header: %+v, data: %+v", method, url, header, data)

	client := g.Client().Timeout(600 * time.Second)
	if header != nil {
		client.SetHeaderMap(header)
	}

	client.SetHeader("Accept", "text/event-stream")

	response, err := client.DoRequest(ctx, method, url, data)
	if err != nil {
		logger.Error(ctx, err)
		return err
	}

	defer func() {
		err = response.Close()
		if err != nil {
			logger.Error(ctx, err)
		}
	}()

	// 使用bufio.NewReader读取响应正文
	reader := bufio.NewReader(response.Body)

	isClose := false
	for {

		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Infof(ctx, "SSEClient method: %s, url: %s, header: %+v, data: %+v done", method, url, header, data)
				return nil
			}
			logger.Error(ctx, err)
			return err
		}

		logger.Infof(ctx, "SSEClient method: %s, url: %s, header: %+v, data: %+v, message: %s", method, url, header, data, message)

		_ = grpool.AddWithRecover(ctx, func(ctx context.Context) {
			result <- []byte(message)
		}, nil)

		if isClose {
			logger.Infof(ctx, "SSEClient method: %s, url: %s, header: %+v, data: %+v done", method, url, header, data)
			return nil
		}

		isClose = message == "event: close"
	}
}
