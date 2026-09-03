package bailian

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/options"
)

type Bailian struct {
	*options.AdapterOptions
	header map[string]string
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) *Bailian {

	bailian := &Bailian{
		AdapterOptions: options,
		header: map[string]string{
			"Authorization": "Bearer " + options.Key,
		},
	}

	if bailian.BaseUrl == "" {
		bailian.BaseUrl = "https://dashscope.aliyuncs.com/api/v1"
	}

	for k, v := range bailian.PassthroughHeader {
		bailian.header[k] = v
	}

	for k, v := range bailian.Header {
		bailian.header[k] = v
	}

	logger.Infof(ctx, "NewAdapter Bailian model: %s, key: %s", bailian.Model, bailian.Key)

	return bailian
}

func (b *Bailian) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	res := model.BailianImageGenerationRes{}
	if err = json.Unmarshal(bytes, &res); err == nil && res.Code != "" {
		return b.apiErrorHandler(response.StatusCode, res.Code, res.Message)
	}

	return errors.NewRequestError(response.StatusCode, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)))
}

func (b *Bailian) apiErrorHandler(statusCode int, code, message string) error {

	switch code {
	case "InvalidApiKey", "InvalidAccessKeyId":
		return errors.ERR_INVALID_API_KEY
	case "Arrearage", "Throttling.AllocationQuota", "AllocationQuota.FreeTierOnly":
		return errors.ERR_INSUFFICIENT_QUOTA
	case "Throttling", "Throttling.RateQuota":
		return errors.ERR_RATE_LIMIT_EXCEEDED
	case "InvalidParameter":
		if gstr.Contains(message, "Model not exist") {
			return errors.ERR_MODEL_NOT_FOUND
		}
	}

	if statusCode == 0 {
		statusCode = 500
	}

	return errors.NewApiError(statusCode, code, message, "api_error", "")
}

// 视频生成等异步接口须带 X-DashScope-Async: enable, 在基础请求头上复制一份追加, 避免影响同步接口
func (b *Bailian) asyncHeader() map[string]string {

	header := make(map[string]string, len(b.header)+1)
	for k, v := range b.header {
		header[k] = v
	}
	header["X-DashScope-Async"] = "enable"

	return header
}
