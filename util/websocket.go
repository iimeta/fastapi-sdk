package util

import (
	"context"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gorilla/websocket"
	"github.com/iimeta/fastapi-sdk/logger"
	"net/http"
	"net/url"
	"time"
)

type WebSocketConn struct {
	conn     *websocket.Conn
	response *http.Response
}

func WebSocketClient(ctx context.Context, wsURL string, messageType int, message []byte, proxyURL ...string) (*WebSocketConn, error) {

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

	conn, response, err := client.Dial(wsURL, nil)
	if err != nil {

		if response != nil {
			if err := response.Body.Close(); err != nil {
				logger.Error(ctx, err)
			}
		}

		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.Error(ctx, err)
			}
		}

		logger.Error(ctx, err)
		return nil, err
	}

	if err = conn.WriteMessage(messageType, message); err != nil {

		logger.Error(ctx, err)

		if err := response.Body.Close(); err != nil {
			logger.Error(ctx, err)
		}

		if err := conn.Close(); err != nil {
			logger.Error(ctx, err)
		}

		return nil, err
	}

	return &WebSocketConn{
		conn:     conn,
		response: response,
	}, nil
}

func (c *WebSocketConn) ReadMessage(ctx context.Context) ([]byte, error) {

	for {

		_, message, err := c.conn.ReadMessage()
		if err != nil && websocket.IsUnexpectedCloseError(err) {

			logger.Error(ctx, err)

			if err := c.Close(); err != nil {
				logger.Error(ctx, err)
			}

			return nil, err
		}

		return message, nil
	}
}

func (c *WebSocketConn) Close() (err error) {

	if e := c.response.Body.Close(); e != nil {
		err = e
	}

	if e := c.conn.Close(); e != nil {
		err = e
	}

	return err
}
