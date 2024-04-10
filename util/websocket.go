package util

import (
	"context"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gorilla/websocket"
	"github.com/iimeta/fastapi-sdk/logger"
	"net/http"
	"net/url"
	"time"
)

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
