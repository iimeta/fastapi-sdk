package util

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// 按代理地址缓存 Transport, 以复用连接.
// 每次请求都新建 Transport 会导致无法复用连接, 且旧 Transport 的空闲连接永远不会被回收(连接泄漏).
var proxyTransports sync.Map // map[string]*http.Transport

// 解析并校验代理地址, 支持: http, https, socks5, socks5h
// 示例: http://127.0.0.1:7890, socks5://user:pass@127.0.0.1:1080
func ParseProxyURL(proxyURL string) (*url.URL, error) {

	proxyUrl, err := url.Parse(strings.TrimSpace(proxyURL))
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(proxyUrl.Scheme) {
	case "http", "https", "socks5", "socks5h":
	default:
		return nil, fmt.Errorf("unsupported proxy scheme %q, supported schemes: http, https, socks5, socks5h", proxyUrl.Scheme)
	}

	if proxyUrl.Host == "" {
		return nil, fmt.Errorf("invalid proxyURL %q: missing host", proxyURL)
	}

	return proxyUrl, nil
}

// 返回对应代理地址的 Transport, proxyURL 为空时使用 http.DefaultTransport
func httpTransport(proxyURL string) (http.RoundTripper, error) {

	if proxyURL == "" {
		return http.DefaultTransport, nil
	}

	if transport, ok := proxyTransports.Load(proxyURL); ok {
		return transport.(*http.Transport), nil
	}

	proxyUrl, err := ParseProxyURL(proxyURL)
	if err != nil {
		return nil, err
	}

	// 基于默认 Transport 克隆, 保留默认的连接超时、空闲连接回收、HTTP/2 等配置
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = http.ProxyURL(proxyUrl)

	actual, _ := proxyTransports.LoadOrStore(proxyURL, transport)

	return actual.(*http.Transport), nil
}
