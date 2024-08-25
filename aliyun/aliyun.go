package aliyun

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/sdkerr"
)

type Client struct {
	key                 string
	baseURL             string
	path                string
	proxyURL            string
	isSupportSystemRole *bool
}

func NewClient(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole *bool, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewClient Aliyun model: %s, key: %s", model, key)

	client := &Client{
		key:                 key,
		baseURL:             "https://dashscope.aliyuncs.com/api/v1",
		path:                "/services/aigc/text-generation/generation",
		isSupportSystemRole: isSupportSystemRole,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewClient Aliyun model: %s, baseURL: %s", model, baseURL)
		client.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewClient Aliyun model: %s, path: %s", model, path)
		client.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewClient Aliyun model: %s, proxyURL: %s", model, proxyURL[0])
		client.proxyURL = proxyURL[0]
	}

	return client
}

func (c *Client) requestErrorHandler(ctx context.Context, response *gclient.Response) (err error) {
	return sdkerr.NewRequestError(500, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, response.ReadAllString())))
}

func (c *Client) apiErrorHandler(response *model.AliyunChatCompletionRes) error {

	switch response.Code {
	case "InvalidParameter":
		if gstr.Contains(response.Message, "Range of input length") {
			return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
		}
	case "BadRequest.TooLarge":
		return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
	case "InvalidApiKey":
		return sdkerr.ERR_INVALID_API_KEY
	case "Throttling.AllocationQuota":
		return sdkerr.ERR_INSUFFICIENT_QUOTA
	}

	return sdkerr.NewApiError(500, response.Code, gjson.MustEncodeString(response), "api_error", "")
}
