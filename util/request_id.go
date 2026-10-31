package util

import (
	"context"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/iimeta/fastapi-sdk/v2/consts"
)

func noteUpstreamRequestId(ctx context.Context, response *http.Response) {

	if ctx == nil {
		return
	}

	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return
	}

	if response == nil || len(response.Header) == 0 {
		r.SetCtxVar(consts.UPSTREAM_REQUEST_IDS_KEY, nil)
		return
	}

	ids := requestIdsFromHeader(response.Header)

	if len(ids) == 0 {
		r.SetCtxVar(consts.UPSTREAM_REQUEST_IDS_KEY, nil)
		return
	}

	r.SetCtxVar(consts.UPSTREAM_REQUEST_IDS_KEY, ids)
}

func ResetUpstreamRequestIds(ctx context.Context) {

	if ctx == nil {
		return
	}

	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return
	}

	r.SetCtxVar(consts.UPSTREAM_REQUEST_IDS_KEY, nil)
}

func UpstreamRequestIds(ctx context.Context) map[string]string {

	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return nil
	}

	ids, _ := r.GetCtxVar(consts.UPSTREAM_REQUEST_IDS_KEY).Val().(map[string]string)
	if len(ids) == 0 {
		return nil
	}

	out := make(map[string]string, len(ids))
	for key, value := range ids {
		out[strings.Clone(key)] = strings.Clone(value)
	}

	return out
}

// 保存主操作的 RequestIds, 用于辅助 HTTP 请求结束后恢复.
func PreserveUpstreamRequestIds(ctx context.Context) func() {

	if ctx == nil {
		return func() {}
	}

	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return func() {}
	}

	ids := UpstreamRequestIds(ctx)

	return func() {
		r.SetCtxVar(consts.UPSTREAM_REQUEST_IDS_KEY, ids)
	}
}

func requestIdsFromHeader(header http.Header) map[string]string {

	ids := make(map[string]string)

	for key, values := range header {

		if !isRequestIdHeader(key) {
			continue
		}

		if value := cloneHeaderValue(values); value != "" {
			ids[strings.Clone(key)] = value
		}
	}

	return ids
}

// HTTP/2 HPACK Huffman 解码会把 header 值放进 bytes.Buffer 对象池.
// strings.Join 单值时直接返回原字符串, TrimSpace 无空白时也是原串,
// 存进 map 后若底层 buffer 被下一次解码复用, 多个 *-Request-Id 就会变成同一个值.
func cloneHeaderValue(values []string) string {
	switch len(values) {
	case 0:
		return ""
	case 1:
		return strings.Clone(strings.TrimSpace(values[0]))
	default:
		return strings.Clone(strings.TrimSpace(strings.Join(values, ", ")))
	}
}

func isRequestIdHeader(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), "id")
}
