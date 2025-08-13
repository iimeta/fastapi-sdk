package baidu

import (
	"context"
	"errors"
	"fmt"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/sdkerr"
)

type Baidu struct {
	accessToken         string
	baseURL             string
	path                string
	proxyURL            string
	isSupportSystemRole *bool
}

func NewAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *Baidu {

	logger.Infof(ctx, "NewAdapter Baidu model: %s, key: %s", model, key)

	client := &Baidu{
		accessToken:         key,
		baseURL:             "https://aip.baidubce.com/rpc/2.0/ai_custom/v1",
		path:                "/wenxinworkshop/chat/completions_pro",
		isSupportSystemRole: isSupportSystemRole,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewAdapter Baidu model: %s, baseURL: %s", model, baseURL)
		client.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewAdapter Baidu model: %s, path: %s", model, path)
		client.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAdapter Baidu model: %s, proxyURL: %s", model, proxyURL[0])
		client.proxyURL = proxyURL[0]
	}

	return client
}

func (b *Baidu) requestErrorHandler(ctx context.Context, response *gclient.Response) (err error) {
	return sdkerr.NewRequestError(500, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, response.ReadAllString())))
}

func (b *Baidu) apiErrorHandler(response *model.BaiduChatCompletionRes) error {

	switch response.ErrorCode {
	case 336103, 336007:
		return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
	case 4, 18, 336501:
		return sdkerr.ERR_RATE_LIMIT_EXCEEDED
	}

	return sdkerr.NewApiError(500, response.ErrorCode, gjson.MustEncodeString(response), "api_error", "")
}
