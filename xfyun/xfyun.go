package xfyun

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/encoding/gurl"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/sdkerr"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	appId               string
	secret              string
	key                 string
	originalURL         string
	baseURL             string
	path                string
	proxyURL            string
	domain              string
	isSupportSystemRole *bool
}

func NewClient(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole *bool, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewClient Xfyun model: %s, key: %s", model, key)

	result := gstr.Split(key, "|")

	client := &Client{
		appId:               result[0],
		secret:              result[1],
		key:                 result[2],
		originalURL:         "https://spark-api.xf-yun.com",
		baseURL:             "https://spark-api.xf-yun.com/v4.0",
		path:                "/chat",
		domain:              "4.0Ultra",
		isSupportSystemRole: isSupportSystemRole,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewClient Xfyun model: %s, baseURL: %s", model, baseURL)

		client.baseURL = baseURL

		version := baseURL[strings.LastIndex(baseURL, "/")+1:]

		switch version {
		case "v4.0":
			client.domain = "4.0Ultra"
		case "v3.5":
			client.domain = "generalv3.5"
		case "v3.1":
			client.domain = "generalv3"
		case "v2.1":
			client.domain = "generalv2"
		case "v1.1":
			client.domain = "general"
		default:
			v := gconv.Float64(version[1:])
			if math.Round(v) > v {
				client.domain = fmt.Sprintf("general%s", version)
			} else {
				client.domain = fmt.Sprintf("generalv%0.f", math.Round(v))
			}
		}
	}

	if path != "" {
		logger.Infof(ctx, "NewClient Xfyun model: %s, path: %s", model, path)
		client.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewClient Xfyun model: %s, proxyURL: %s", model, proxyURL[0])
		client.proxyURL = proxyURL[0]
	}

	return client
}

func (c *Client) getWebSocketUrl(ctx context.Context) string {

	date, host, signature, err := c.getSignature(ctx, http.MethodGet)
	if err != nil {
		logger.Errorf(ctx, "getWebSocketUrl Xfyun client: %+v, error: %s", c, err)
		return ""
	}

	authorizationOrigin := gbase64.EncodeToString([]byte(fmt.Sprintf("api_key=\"%s\",algorithm=\"%s\",headers=\"%s\",signature=\"%s\"", c.key, "hmac-sha256", "host date request-line", signature)))

	wsURL := gstr.Replace(gstr.Replace(c.baseURL+c.path, "https://", "wss://"), "http://", "ws://")

	return fmt.Sprintf("%s?authorization=%s&date=%s&host=%s", wsURL, authorizationOrigin, date, host)
}

func (c *Client) getHttpUrl(ctx context.Context) string {

	c.originalURL = "https://spark-api.cn-huabei-1.xf-yun.com"

	date, host, signature, err := c.getSignature(ctx, http.MethodPost)
	if err != nil {
		logger.Errorf(ctx, "getHttpUrl Xfyun client: %+v, error: %s", c, err)
		return ""
	}

	authorizationOrigin := gbase64.EncodeToString([]byte(fmt.Sprintf("api_key=\"%s\",algorithm=\"%s\",headers=\"%s\",signature=\"%s\"", c.key, "hmac-sha256", "host date request-line", signature)))

	return fmt.Sprintf("%s?authorization=%s&date=%s&host=%s", c.baseURL+c.path, authorizationOrigin, date, host)
}

func (c *Client) getSignature(ctx context.Context, method string) (date, host, signature string, err error) {

	parse, err := url.Parse(c.originalURL + c.baseURL[strings.LastIndex(c.baseURL, "/"):] + c.path)
	if err != nil {
		logger.Errorf(ctx, "getSignature Xfyun client: %+v, error: %s", c, err)
		return "", "", "", err
	}

	now := gtime.Now()
	loc, _ := time.LoadLocation("GMT")
	zone, _ := now.ToZone(loc.String())
	date = zone.Layout("Mon, 02 Jan 2006 15:04:05 GMT")

	tmp := "host: " + parse.Host + "\n"
	tmp += "date: " + date + "\n"
	tmp += method + " " + parse.Path + " HTTP/1.1"

	hash := hmac.New(sha256.New, []byte(c.secret))

	if _, err = hash.Write([]byte(tmp)); err != nil {
		logger.Errorf(ctx, "getSignature Xfyun client: %+v, error: %s", c, err)
		return "", "", "", err
	}

	return gurl.RawEncode(date), parse.Host, gbase64.EncodeToString(hash.Sum(nil)), nil
}

func (c *Client) apiErrorHandler(response *model.XfyunChatCompletionRes) error {

	switch response.Header.Code {
	case 10163, 10907:
		return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
	}

	return sdkerr.NewApiError(500, response.Header.Code, gjson.MustEncodeString(response), "api_error", "")
}
