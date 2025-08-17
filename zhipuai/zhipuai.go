package zhipuai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/sdkerr"
)

type ZhipuAI struct {
	model               string
	key                 string
	baseURL             string
	path                string
	proxyURL            string
	header              map[string]string
	isSupportSystemRole *bool
	isSupportStream     *bool
}

func NewAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *ZhipuAI {

	logger.Infof(ctx, "NewAdapter ZhipuAI model: %s, key: %s", model, key)

	zhipuai := &ZhipuAI{
		model:               model,
		key:                 key,
		baseURL:             "https://open.bigmodel.cn/api/paas/v4",
		path:                "/chat/completions",
		isSupportSystemRole: isSupportSystemRole,
		isSupportStream:     isSupportStream,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewAdapter ZhipuAI model: %s, baseURL: %s", model, baseURL)
		zhipuai.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewAdapter ZhipuAI model: %s, path: %s", model, path)
		zhipuai.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAdapter ZhipuAI model: %s, proxyURL: %s", model, proxyURL[0])
		zhipuai.proxyURL = proxyURL[0]
	}

	return zhipuai
}

func (z *ZhipuAI) generateToken(ctx context.Context) string {

	split := strings.Split(z.key, ".")
	if len(split) != 2 {
		return z.key
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

		reqErr := &sdkerr.RequestError{
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
		return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
	case "1113":
		return sdkerr.ERR_INSUFFICIENT_QUOTA
	}

	return sdkerr.NewRequestError(500, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, gjson.MustEncodeString(errRes.Error))))
}

func (z *ZhipuAI) apiErrorHandler(response *model.ZhipuAIChatCompletionRes) error {

	switch response.Error.Code {
	case "1261":
		return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
	case "1113":
		return sdkerr.ERR_INSUFFICIENT_QUOTA
	}

	return sdkerr.NewApiError(500, response.Error.Code, gjson.MustEncodeString(response), "api_error", "")
}
