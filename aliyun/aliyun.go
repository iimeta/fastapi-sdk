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
	"github.com/iimeta/fastapi-sdk/options"
	"github.com/iimeta/fastapi-sdk/sdkerr"
)

type Aliyun struct {
	*options.AdapterOptions
	header map[string]string
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) *Aliyun {

	aliyun := &Aliyun{
		AdapterOptions: options,
		header: g.MapStrStr{
			"Authorization": "Bearer " + options.Key,
		},
	}

	if aliyun.BaseUrl == "" {
		aliyun.BaseUrl = "https://dashscope.aliyuncs.com/api/v1"
	}

	if aliyun.Path == "" {
		aliyun.Path = "/services/aigc/text-generation/generation"
	}

	logger.Infof(ctx, "NewAdapter Aliyun model: %s, key: %s", aliyun.Model, aliyun.Key)

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
