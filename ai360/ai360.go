package ai360

import (
	"context"
	"errors"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/sdkerr"
	"github.com/iimeta/go-openai"
)

type AI360 struct {
	model               string
	key                 string
	baseURL             string
	path                string
	proxyURL            string
	header              map[string]string
	isSupportSystemRole *bool
	isSupportStream     *bool
}

func NewAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *AI360 {

	logger.Infof(ctx, "NewAdapter 360AI model: %s, key: %s", model, key)

	ai360 := &AI360{
		model:   model,
		key:     key,
		baseURL: "https://api.360.cn/v1",
		path:    "/chat/completions",
		header: g.MapStrStr{
			"Authorization": "Bearer " + key,
		},
		isSupportSystemRole: isSupportSystemRole,
		isSupportStream:     isSupportStream,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewAdapter 360AI model: %s, baseURL: %s", model, baseURL)
		ai360.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewAdapter 360AI model: %s, path: %s", model, path)
		ai360.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAdapter 360AI model: %s, proxyURL: %s", model, proxyURL[0])
		ai360.proxyURL = proxyURL[0]
	}

	return ai360
}

func (a *AI360) apiErrorHandler(err error) error {

	apiError := &openai.APIError{}
	if errors.As(err, &apiError) {

		switch apiError.HTTPStatusCode {
		case 400:
			if apiError.Code == "1001" {
				return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
			}
		case 401:

			if apiError.Code == "1002" {
				return sdkerr.ERR_INVALID_API_KEY
			}

			if apiError.Code == "1004" || apiError.Code == "1006" {
				return sdkerr.ERR_INSUFFICIENT_QUOTA
			}

		case 404:
			return sdkerr.ERR_MODEL_NOT_FOUND
		case 429:
			if apiError.Code == "1005" {
				return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
			}
		}

		return err
	}

	reqError := &openai.RequestError{}
	if errors.As(err, &reqError) {
		return sdkerr.NewRequestError(apiError.HTTPStatusCode, reqError.Err)
	}

	return err
}
