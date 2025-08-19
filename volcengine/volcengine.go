package volcengine

import (
	"context"
	"errors"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/sdkerr"
)

type VolcEngine struct {
	model               string
	key                 string
	baseURL             string
	path                string
	proxyURL            string
	header              map[string]string
	isSupportSystemRole *bool
	isSupportStream     *bool
}

func NewAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *VolcEngine {

	logger.Infof(ctx, "NewAdapter VolcEngine model: %s, key: %s", model, key)

	split := gstr.Split(key, "|")

	volcengine := &VolcEngine{
		model:   split[0],
		key:     split[1],
		baseURL: "https://ark.cn-beijing.volces.com/api/v3",
		path:    "/chat/completions",
		header: g.MapStrStr{
			"Authorization": "Bearer " + split[1],
		},
		isSupportSystemRole: isSupportSystemRole,
		isSupportStream:     isSupportStream,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewAdapter VolcEngine model: %s, baseURL: %s", model, baseURL)
		volcengine.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewAdapter VolcEngine model: %s, path: %s", model, path)
		volcengine.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAdapter VolcEngine model: %s, proxyURL: %s", model, proxyURL[0])
		volcengine.proxyURL = proxyURL[0]
	}

	return volcengine
}

func (v *VolcEngine) apiErrorHandler(err error) error {

	apiError := &sdkerr.ApiError{}
	if errors.As(err, &apiError) {

		switch apiError.HttpStatusCode {
		case 400:
			if apiError.Code == "context_length_exceeded" {
				return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
			}
		case 401:
			if apiError.Code == "invalid_api_key" {
				return sdkerr.ERR_INVALID_API_KEY
			}
		case 404:
			return sdkerr.ERR_MODEL_NOT_FOUND
		case 429:
			if apiError.Code == "insufficient_quota" {
				return sdkerr.ERR_INSUFFICIENT_QUOTA
			}
		}

		return err
	}

	reqError := &sdkerr.RequestError{}
	if errors.As(err, &reqError) {
		return sdkerr.NewRequestError(apiError.HttpStatusCode, reqError.Err)
	}

	return err
}
