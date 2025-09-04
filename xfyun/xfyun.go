package xfyun

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/encoding/gurl"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/errors"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/options"
)

type Xfyun struct {
	*options.AdapterOptions
	header      map[string]string
	appId       string
	secret      string
	originalUrl string
	domain      string
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) *Xfyun {

	result := gstr.Split(options.Key, "|")

	xfyun := &Xfyun{
		AdapterOptions: options,
		appId:          result[0],
		secret:         result[1],
		originalUrl:    "https://spark-api.xf-yun.com",
		domain:         "4.0Ultra",
	}

	xfyun.Key = result[2]

	if xfyun.BaseUrl == "" {
		xfyun.BaseUrl = "https://spark-api.xf-yun.com/v4.0"
	}

	if xfyun.Path == "" {
		xfyun.Path = "/chat"
	}

	version := xfyun.BaseUrl[strings.LastIndex(xfyun.BaseUrl, "/")+1:]

	switch version {
	case "v4.0":
		xfyun.domain = "4.0Ultra"
	case "v3.5":
		xfyun.domain = "generalv3.5"
	case "v3.1":
		xfyun.domain = "generalv3"
	case "v2.1":
		xfyun.domain = "generalv2"
	case "v1.1":
		xfyun.domain = "general"
	default:
		v := gconv.Float64(version[1:])
		if math.Round(v) > v {
			xfyun.domain = fmt.Sprintf("general%s", version)
		} else {
			xfyun.domain = fmt.Sprintf("generalv%0.f", math.Round(v))
		}
	}

	logger.Infof(ctx, "NewAdapter Xfyun model: %s, key: %s", xfyun.Model, xfyun.Key)

	return xfyun
}

func (x *Xfyun) getWebSocketUrl(ctx context.Context) string {

	date, host, signature, err := x.getSignature(ctx, http.MethodGet)
	if err != nil {
		logger.Errorf(ctx, "getWebSocketUrl Xfyun client: %+v, error: %s", x, err)
		return ""
	}

	authorizationOrigin := gbase64.EncodeToString([]byte(fmt.Sprintf("api_key=\"%s\",algorithm=\"%s\",headers=\"%s\",signature=\"%s\"", x.Key, "hmac-sha256", "host date request-line", signature)))

	wsURL := gstr.Replace(gstr.Replace(x.BaseUrl+x.Path, "https://", "wss://"), "http://", "ws://")

	return fmt.Sprintf("%s?authorization=%s&date=%s&host=%s", wsURL, authorizationOrigin, date, host)
}

func (x *Xfyun) getHttpUrl(ctx context.Context) string {

	x.originalUrl = "https://spark-api.cn-huabei-1.xf-yun.com"

	date, host, signature, err := x.getSignature(ctx, http.MethodPost)
	if err != nil {
		logger.Errorf(ctx, "getHttpUrl Xfyun client: %+v, error: %s", x, err)
		return ""
	}

	authorizationOrigin := gbase64.EncodeToString([]byte(fmt.Sprintf("api_key=\"%s\",algorithm=\"%s\",headers=\"%s\",signature=\"%s\"", x.Key, "hmac-sha256", "host date request-line", signature)))

	return fmt.Sprintf("%s?authorization=%s&date=%s&host=%s", x.BaseUrl+x.Path, authorizationOrigin, date, host)
}

func (x *Xfyun) getSignature(ctx context.Context, method string) (date, host, signature string, err error) {

	parse, err := url.Parse(x.originalUrl + x.BaseUrl[strings.LastIndex(x.BaseUrl, "/"):] + x.Path)
	if err != nil {
		logger.Errorf(ctx, "getSignature Xfyun client: %+v, error: %s", x, err)
		return "", "", "", err
	}

	now := gtime.Now()
	loc, _ := time.LoadLocation("GMT")
	zone, _ := now.ToZone(loc.String())
	date = zone.Layout("Mon, 02 Jan 2006 15:04:05 GMT")

	tmp := "host: " + parse.Host + "\n"
	tmp += "date: " + date + "\n"
	tmp += method + " " + parse.Path + " HTTP/1.1"

	hash := hmac.New(sha256.New, []byte(x.secret))

	if _, err = hash.Write([]byte(tmp)); err != nil {
		logger.Errorf(ctx, "getSignature Xfyun client: %+v, error: %s", x, err)
		return "", "", "", err
	}

	return gurl.RawEncode(date), parse.Host, gbase64.EncodeToString(hash.Sum(nil)), nil
}

func (x *Xfyun) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return errors.NewRequestError(response.StatusCode, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)))
}

func (x *Xfyun) apiErrorHandler(response *model.XfyunChatCompletionRes) error {

	switch response.Header.Code {
	case 10163, 10907:
		return errors.ERR_CONTEXT_LENGTH_EXCEEDED
	}

	return errors.NewApiError(500, response.Header.Code, gjson.MustEncodeString(response), "api_error", "")
}
