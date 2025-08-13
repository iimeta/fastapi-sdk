package google

import (
	"context"
	"errors"
	"fmt"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/sdkerr"
)

type Google struct {
	model               string
	key                 string
	baseURL             string
	path                string
	proxyURL            string
	isSupportSystemRole *bool
	header              map[string]string
	isGcp               bool
}

func NewAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *Google {

	logger.Infof(ctx, "NewAdapter Google model: %s, key: %s", model, key)

	client := &Google{
		model:               model,
		key:                 key,
		baseURL:             "https://generativelanguage.googleapis.com/v1beta",
		path:                "/models/" + model,
		isSupportSystemRole: isSupportSystemRole,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewAdapter Google model: %s, baseURL: %s", model, baseURL)
		client.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewAdapter Google model: %s, path: %s", model, path)
		client.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAdapter Google model: %s, proxyURL: %s", model, proxyURL[0])
		client.proxyURL = proxyURL[0]
	}

	return client
}

func NewGcpAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *Google {

	logger.Infof(ctx, "NewGcpAdapter Google model: %s, key: %s", model, key)

	client := &Google{
		model:               model,
		key:                 key,
		baseURL:             "https://us-east5-aiplatform.googleapis.com/v1",
		path:                "/projects/%s/locations/us-east5/publishers/google/models/%s",
		isSupportSystemRole: isSupportSystemRole,
		isGcp:               true,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewGcpAdapter Google model: %s, baseURL: %s", model, baseURL)
		client.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewGcpAdapter Google model: %s, path: %s", model, path)
		client.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewGcpAdapter Google model: %s, proxyURL: %s", model, proxyURL[0])
		client.proxyURL = proxyURL[0]
	}

	client.header = make(map[string]string)
	client.header["Authorization"] = "Bearer " + key

	return client
}

func (g *Google) requestErrorHandler(ctx context.Context, response *gclient.Response) (err error) {
	return sdkerr.NewRequestError(500, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, response.ReadAllString())))
}

func (g *Google) apiErrorHandler(response *model.GoogleChatCompletionRes) error {
	return sdkerr.NewApiError(500, response.Error.Code, gjson.MustEncodeString(response), "api_error", "")
}
