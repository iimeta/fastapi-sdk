package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/errors"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/options"
)

type OpenAI struct {
	*options.AdapterOptions
	header  map[string]string
	isAzure bool
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) *OpenAI {

	openai := &OpenAI{
		AdapterOptions: options,
		header: g.MapStrStr{
			"Authorization": "Bearer " + options.Key,
		},
	}

	if openai.BaseUrl == "" {
		openai.BaseUrl = "https://api.openai.com/v1"
	}

	if openai.Path == "" {
		openai.Path = "/chat/completions"
	}

	logger.Infof(ctx, "NewAdapter OpenAI model: %s, key: %s", openai.Model, openai.Key)

	return openai
}

func NewAzureAdapter(ctx context.Context, options *options.AdapterOptions) *OpenAI {

	azure := &OpenAI{
		AdapterOptions: options,
		header: g.MapStrStr{
			"api-key": options.Key,
		},
		isAzure: true,
	}

	if gstr.HasSuffix(azure.BaseUrl, "/openai/deployments") {
		azure.BaseUrl = azure.BaseUrl + "/" + options.Model
	} else if !gstr.HasSuffix(azure.BaseUrl, "/models") {

		azure.BaseUrl = strings.TrimRight(azure.BaseUrl, "/")

		if parse, _ := url.Parse(azure.BaseUrl); parse == nil || parse.Path == "" {
			azure.BaseUrl = azure.BaseUrl + "/openai/deployments/" + options.Model
		}
	}

	if azure.Path == "" {
		azure.Path = "/chat/completions?api-version=2024-05-01-preview"
	}

	logger.Infof(ctx, "NewAzureAdapter OpenAI model: %s, baseUrl: %s, key: %s", azure.Model, azure.BaseUrl, azure.Key)

	return azure
}

func (o *OpenAI) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	errorResponse := errors.ErrorResponse{}
	if err := json.Unmarshal(bytes, &errorResponse); err != nil || errorResponse.Error == nil {
		return &errors.RequestError{
			HttpStatusCode: response.StatusCode,
			Err:            errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)),
		}
	}

	return errors.NewApiError(response.StatusCode, errorResponse.Error.Code, errorResponse.Error.Message, errorResponse.Error.Type, *errorResponse.Error.Param)
}

func (o *OpenAI) apiErrorHandler(err error) error {

	//apiError := &errors.ApiError{}
	//if errors.As(err, &apiError) {
	//
	//	switch apiError.HttpStatusCode {
	//	case 400:
	//		if apiError.Code == "context_length_exceeded" {
	//			return errors.ERR_CONTEXT_LENGTH_EXCEEDED
	//		}
	//	case 401:
	//		if apiError.Code == "invalid_api_key" {
	//			return errors.ERR_INVALID_API_KEY
	//		}
	//	case 404:
	//		return errors.ERR_MODEL_NOT_FOUND
	//	case 429:
	//		if apiError.Code == "insufficient_quota" {
	//			return errors.ERR_INSUFFICIENT_QUOTA
	//		}
	//	}
	//
	//	return err
	//}
	//
	//reqError := &errors.RequestError{}
	//if errors.As(err, &reqError) {
	//	return errors.NewRequestError(apiError.HttpStatusCode, reqError.Err)
	//}

	return err
}
