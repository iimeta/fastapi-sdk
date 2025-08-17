package baidu

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/sdkerr"
)

type Baidu struct {
	model               string
	key                 string
	baseURL             string
	path                string
	proxyURL            string
	header              map[string]string
	isSupportSystemRole *bool
	isSupportStream     *bool
	accessToken         string
}

func NewAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *Baidu {

	logger.Infof(ctx, "NewAdapter Baidu model: %s, key: %s", model, key)

	baidu := &Baidu{
		model:               model,
		key:                 key,
		baseURL:             "https://aip.baidubce.com/rpc/2.0/ai_custom/v1",
		path:                "/wenxinworkshop/chat/completions_pro",
		isSupportSystemRole: isSupportSystemRole,
		isSupportStream:     isSupportStream,
		accessToken:         key,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewAdapter Baidu model: %s, baseURL: %s", model, baseURL)
		baidu.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewAdapter Baidu model: %s, path: %s", model, path)
		baidu.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAdapter Baidu model: %s, proxyURL: %s", model, proxyURL[0])
		baidu.proxyURL = proxyURL[0]
	}

	return baidu
}

func (b *Baidu) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return sdkerr.NewRequestError(500, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)))
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
