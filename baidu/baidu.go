package baidu

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/options"
)

type Baidu struct {
	*options.AdapterOptions
	header map[string]string
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) *Baidu {

	baidu := &Baidu{
		AdapterOptions: options,
	}

	if baidu.BaseUrl == "" {
		baidu.BaseUrl = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1"
	}

	baidu.Path += "?access_token=" + options.Key

	logger.Infof(ctx, "NewAdapter Baidu model: %s, key: %s", baidu.Model, baidu.Key)

	return baidu
}

func (b *Baidu) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return errors.NewRequestError(response.StatusCode, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)))
}

func (b *Baidu) apiErrorHandler(response *model.BaiduChatCompletionRes) error {

	switch response.ErrorCode {
	case 336103, 336007:
		return errors.ERR_CONTEXT_LENGTH_EXCEEDED
	case 4, 18, 336501:
		return errors.ERR_RATE_LIMIT_EXCEEDED
	}

	return errors.NewApiError(500, response.ErrorCode, gjson.MustEncodeString(response), "api_error", "")
}
