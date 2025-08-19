package openai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/sdkerr"
)

type OpenAI struct {
	model               string
	key                 string
	baseURL             string
	path                string
	proxyURL            string
	header              map[string]string
	isSupportSystemRole *bool
	isSupportStream     *bool
	isAzure             bool
	apiVersion          string
}

func NewAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *OpenAI {

	logger.Infof(ctx, "NewAdapter OpenAI model: %s, key: %s", model, key)

	openai := &OpenAI{
		model:   model,
		key:     key,
		baseURL: "https://api.openai.com/v1",
		path:    "/chat/completions",
		header: g.MapStrStr{
			"Authorization": "Bearer " + key,
		},
		isSupportSystemRole: isSupportSystemRole,
		isSupportStream:     isSupportStream,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewAdapter OpenAI model: %s, baseURL: %s", model, baseURL)
		openai.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewAdapter OpenAI model: %s, path: %s", model, path)
		openai.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAdapter OpenAI model: %s, proxyURL: %s", model, proxyURL[0])
		openai.proxyURL = proxyURL[0]
	}

	return openai
}

func NewAzureAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *OpenAI {

	logger.Infof(ctx, "NewAzureAdapter OpenAI model: %s, baseURL: %s, key: %s", model, baseURL, key)

	azure := &OpenAI{
		model:   model,
		key:     key,
		baseURL: baseURL,
		path:    "/chat/completions",
		header: g.MapStrStr{
			"api-key": key,
		},
		isSupportSystemRole: isSupportSystemRole,
		isSupportStream:     isSupportStream,
		isAzure:             true,
		apiVersion:          "2024-10-01-preview",
	}

	if path != "" {
		logger.Infof(ctx, "NewAzureAdapter OpenAI model: %s, path: %s", model, path)

		split := gstr.Split(path, "?api-version=")

		if len(split) > 1 && split[1] != "" {
			azure.apiVersion = split[1]
		}
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAzureAdapter OpenAI model: %s, proxyURL: %s", model, proxyURL[0])
		azure.proxyURL = proxyURL[0]
	}

	return azure
}

func (o *OpenAI) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return sdkerr.NewRequestError(500, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)))
}

func (o *OpenAI) apiErrorHandler(err error) error {

	//apiError := &sdkerr.ApiError{}
	//if errors.As(err, &apiError) {
	//
	//	switch apiError.HttpStatusCode {
	//	case 400:
	//		if apiError.Code == "context_length_exceeded" {
	//			return sdkerr.ERR_CONTEXT_LENGTH_EXCEEDED
	//		}
	//	case 401:
	//		if apiError.Code == "invalid_api_key" {
	//			return sdkerr.ERR_INVALID_API_KEY
	//		}
	//	case 404:
	//		return sdkerr.ERR_MODEL_NOT_FOUND
	//	case 429:
	//		if apiError.Code == "insufficient_quota" {
	//			return sdkerr.ERR_INSUFFICIENT_QUOTA
	//		}
	//	}
	//
	//	return err
	//}
	//
	//reqError := &sdkerr.RequestError{}
	//if errors.As(err, &reqError) {
	//	return sdkerr.NewRequestError(apiError.HttpStatusCode, reqError.Err)
	//}

	return err
}
