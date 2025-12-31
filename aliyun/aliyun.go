package aliyun

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/options"
)

type Aliyun struct {
	*options.AdapterOptions
	header map[string]string
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) *Aliyun {

	aliyun := &Aliyun{
		AdapterOptions: options,
		header: map[string]string{
			"Authorization": "Bearer " + options.Key,
		},
	}

	if aliyun.BaseUrl == "" {
		aliyun.BaseUrl = "https://dashscope.aliyuncs.com/api/v1"
	}

	logger.Infof(ctx, "NewAdapter Aliyun model: %s, key: %s", aliyun.Model, aliyun.Key)

	return aliyun
}

func (a *Aliyun) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return errors.NewRequestError(response.StatusCode, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)))
}

func (a *Aliyun) apiErrorHandler(response *model.AliyunChatCompletionRes) error {

	switch response.Code {
	case "InvalidParameter":
		if gstr.Contains(response.Message, "Range of input length") {
			return errors.ERR_CONTEXT_LENGTH_EXCEEDED
		}
	case "BadRequest.TooLarge":
		return errors.ERR_CONTEXT_LENGTH_EXCEEDED
	case "InvalidApiKey":
		return errors.ERR_INVALID_API_KEY
	case "Throttling.AllocationQuota":
		return errors.ERR_INSUFFICIENT_QUOTA
	}

	return errors.NewApiError(500, response.Code, gjson.MustEncodeString(response), "api_error", "")
}
