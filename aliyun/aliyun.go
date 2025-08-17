package aliyun

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/sdkerr"
)

type Aliyun struct {
	model               string
	key                 string
	baseURL             string
	path                string
	proxyURL            string
	header              map[string]string
	isSupportSystemRole *bool
	isSupportStream     *bool
}

func NewAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *Aliyun {

	logger.Infof(ctx, "NewAdapter Aliyun model: %s, key: %s", model, key)

	aliyun := &Aliyun{
		model:   model,
		key:     key,
		baseURL: "https://dashscope.aliyuncs.com/api/v1",
		path:    "/services/aigc/text-generation/generation",
		header: g.MapStrStr{
			"Authorization": "Bearer " + key,
		},
		isSupportSystemRole: isSupportSystemRole,
		isSupportStream:     isSupportStream,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewAdapter Aliyun model: %s, baseURL: %s", model, baseURL)
		aliyun.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewAdapter Aliyun model: %s, path: %s", model, path)
		aliyun.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAdapter Aliyun model: %s, proxyURL: %s", model, proxyURL[0])
		aliyun.proxyURL = proxyURL[0]
	}

	return aliyun
}

func (a *Aliyun) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return sdkerr.NewRequestError(500, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)))
}

func (a *Aliyun) apiErrorHandler(response *model.AliyunChatCompletionRes) error {

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
