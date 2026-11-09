package util

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gorilla/websocket"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
)

type WebSocketConn struct {
	conn      *websocket.Conn
	response  *http.Response
	closeOnce sync.Once
	closeErr  error
}

func WebSocketClient(ctx context.Context, wsURL string, requestHeader http.Header, messageType int, message []byte, proxyURL string) (*WebSocketConn, error) {

	logger.Infof(ctx, "WebSocketClient wsURL: %s, proxyURL: %s", wsURL, proxyURL)

	client := gclient.NewWebSocket()

	client.HandshakeTimeout = 60 * time.Second // 设置超时时间

	// 设置代理
	if proxyURL != "" {

		proxyUrl, err := webSocketProxyURL(proxyURL)
		if err != nil {
			logger.Errorf(ctx, "WebSocketClient wsURL: %s, proxyURL: %s, error: %v", wsURL, proxyURL, err)
			return nil, err
		}

		client.Proxy = http.ProxyURL(proxyUrl)
	}

	if requestHeader == nil {
		requestHeader = http.Header{}
	}

	if requestHeader.Get(consts.TRACE_ID) == "" {
		requestHeader.Set(consts.TRACE_ID, gtrace.GetTraceID(ctx))
	}

	if requestHeader.Get(consts.FASTAPI_REQUEST_ID) == "" {
		requestHeader.Set(consts.FASTAPI_REQUEST_ID, gtrace.GetTraceID(ctx))
	}

	conn, response, err := client.DialContext(ctx, wsURL, requestHeader)
	noteUpstreamRequestId(ctx, response)
	if err != nil {

		if response != nil && response.Body != nil {

			// 握手失败时 ErrBadHandshake 不包含任何上游信息, 补充状态码和响应内容便于排查
			if errors.Is(err, websocket.ErrBadHandshake) {
				body, _ := io.ReadAll(response.Body)
				err = fmt.Errorf("%w, status code: %d, response: %s", err, response.StatusCode, body)
			}

			if e := response.Body.Close(); e != nil {
				logger.Error(ctx, e)
			}
		}

		if conn != nil {
			if e := conn.Close(); e != nil {
				logger.Error(ctx, e)
			}
		}

		logger.Errorf(ctx, "WebSocketClient wsURL: %s, proxyURL: %s, error: %v", wsURL, proxyURL, err)

		return nil, err
	}

	webSocketConn := &WebSocketConn{
		conn:     conn,
		response: response,
	}

	if message != nil {
		if err = conn.WriteMessage(messageType, message); err != nil {
			logger.Errorf(ctx, "WebSocketClient wsURL: %s, WriteMessage error: %v", wsURL, err)

			if e := webSocketConn.Close(); e != nil {
				logger.Error(ctx, e)
			}

			return nil, err
		}
	}

	return webSocketConn, nil
}

// 读取消息
//   - 本端已调用 Close 关闭连接: 返回 io.EOF
//   - 对端发送关闭帧: 返回 *websocket.CloseError
//   - 其他读取错误(网络断开、超时等): 原样返回
//
// 读取出错后连接已不可用, 会自动关闭. 出错时 messageType 为 -1.
func (c *WebSocketConn) ReadMessage(ctx context.Context) (int, []byte, error) {

	messageType, message, err := c.conn.ReadMessage()
	if err == nil {
		return messageType, message, nil
	}

	// 本端主动关闭后的读取错误, 属于正常结束
	if errors.Is(err, net.ErrClosed) {
		return messageType, nil, io.EOF
	}

	logger.Errorf(ctx, "WebSocketConn ReadMessage error: %v", err)

	if e := c.Close(); e != nil && !errors.Is(e, net.ErrClosed) {
		logger.Error(ctx, e)
	}

	return messageType, nil, err
}

func (c *WebSocketConn) WriteMessage(ctx context.Context, messageType int, message []byte) error {

	if messageType != 0 && message != nil {
		if err := c.conn.WriteMessage(messageType, message); err != nil {
			logger.Error(ctx, err)
			return err
		}
	}

	return nil
}

func (c *WebSocketConn) WriteJSON(ctx context.Context, message any) error {

	if message != nil {
		if err := c.conn.WriteJSON(message); err != nil {
			logger.Error(ctx, err)
			return err
		}
	}

	return nil
}

// 关闭连接, 可重复调用, 只有第一次调用会真正关闭
func (c *WebSocketConn) Close() error {

	c.closeOnce.Do(func() {

		if c.response != nil && c.response.Body != nil {
			if err := c.response.Body.Close(); err != nil {
				c.closeErr = err
			}
		}

		if err := c.conn.Close(); err != nil {
			c.closeErr = err
		}
	})

	return c.closeErr
}

// 解析 WebSocket 代理地址
// gorilla/websocket 仅支持 http 和 socks5 代理:
//   - socks5h 转为 socks5 (其 socks5 实现本身就是由代理端解析域名, 两者等价)
//   - https 代理不支持, 直接返回错误, 避免静默失败
func webSocketProxyURL(proxyURL string) (*url.URL, error) {

	proxyUrl, err := ParseProxyURL(proxyURL)
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(proxyUrl.Scheme) {
	case "http", "socks5":
		return proxyUrl, nil
	case "socks5h":
		u := *proxyUrl
		u.Scheme = "socks5"
		return &u, nil
	default:
		return nil, fmt.Errorf("unsupported websocket proxy scheme %q, supported schemes: http, socks5, socks5h", proxyUrl.Scheme)
	}
}
