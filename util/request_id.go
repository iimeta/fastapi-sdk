package util

import (
	"context"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/iimeta/fastapi-sdk/v2/consts"
)

func noteUpstreamRequestId(ctx context.Context, response *http.Response) {

	if ctx == nil || response == nil || len(response.Header) == 0 {
		return
	}

	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return
	}

	ids, _ := r.GetCtxVar(consts.UPSTREAM_REQUEST_IDS_KEY).Val().(map[string]string)
	if ids == nil {
		ids = make(map[string]string)
		r.SetCtxVar(consts.UPSTREAM_REQUEST_IDS_KEY, ids)
	}

	for key, values := range response.Header {

		if ids[key] != "" || !isRequestIdHeader(key) {
			continue
		}

		if value := strings.TrimSpace(strings.Join(values, ", ")); value != "" {
			ids[key] = value
		}
	}
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
		out[key] = value
	}

	return out
}

func isRequestIdHeader(name string) bool {
	n := strings.ToLower(name)
	return n == "trace-id" || n == "request-id" ||
		strings.HasSuffix(n, "-request-id") ||
		strings.HasSuffix(n, "-requestid") ||
		strings.HasSuffix(n, "-log-id") ||
		strings.HasSuffix(n, "-logid")
}
