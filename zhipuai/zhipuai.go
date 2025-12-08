package zhipuai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iimeta/fastapi-sdk/errors"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/options"
)

type ZhipuAI struct {
	*options.AdapterOptions
	header map[string]string
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) *ZhipuAI {

	zhipuai := &ZhipuAI{
		AdapterOptions: options,
	}

	if zhipuai.BaseUrl == "" {
		zhipuai.BaseUrl = "https://open.bigmodel.cn/api/paas/v4"
	}

	zhipuai.header = map[string]string{
		"Authorization": "Bearer " + zhipuai.generateToken(ctx),
	}

	logger.Infof(ctx, "NewAdapter ZhipuAI model: %s, key: %s", zhipuai.Model, zhipuai.Key)

	return zhipuai
}

func (z *ZhipuAI) generateToken(ctx context.Context) string {

	split := strings.Split(z.Key, ".")
	if len(split) != 2 {
		return z.Key
	}

	now := gtime.Now()

	claims := jwt.MapClaims{
		"api_key":   split[0],
		"exp":       now.Add(time.Minute * 10).UnixMilli(),
		"timestamp": now.UnixMilli(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token.Header["alg"] = "HS256"
	token.Header["sign_type"] = "SIGN"

	sign, err := token.SignedString([]byte(split[1]))
	if err != nil {
		logger.Error(ctx, err)
	}

	return sign
}

func (z *ZhipuAI) requestErrorHandler(ctx context.Context, response *http.Response) error {

	errRes := model.ZhipuAIErrorResponse{}
	if err := json.NewDecoder(response.Body).Decode(&errRes); err != nil || errRes.Error == nil {

		reqErr := &errors.RequestError{
			HttpStatusCode: response.StatusCode,
			Err:            err,
		}

		if errRes.Error != nil {
			reqErr.Err = errors.New(gjson.MustEncodeString(errRes.Error))
		}

		return reqErr
	}

	switch errRes.Error.Code {
	case "1261":
		return errors.ERR_CONTEXT_LENGTH_EXCEEDED
	case "1113":
		return errors.ERR_INSUFFICIENT_QUOTA
	}

	return errors.NewRequestError(response.StatusCode, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, gjson.MustEncodeString(errRes.Error))))
}

func (z *ZhipuAI) apiErrorHandler(response *model.ZhipuAIChatCompletionRes) error {

	switch response.Error.Code {
	case "1261":
		return errors.ERR_CONTEXT_LENGTH_EXCEEDED
	case "1113":
		return errors.ERR_INSUFFICIENT_QUOTA
	}

	return errors.NewApiError(500, response.Error.Code, gjson.MustEncodeString(response), "api_error", "")
}
