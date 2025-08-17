package google

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
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
	header              map[string]string
	isSupportSystemRole *bool
	isSupportStream     *bool
	isGcp               bool
}

func NewAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *Google {

	logger.Infof(ctx, "NewAdapter Google model: %s, key: %s", model, key)

	google := &Google{
		model:               model,
		key:                 key,
		baseURL:             "https://generativelanguage.googleapis.com/v1beta",
		path:                "/models/" + model,
		isSupportSystemRole: isSupportSystemRole,
		isSupportStream:     isSupportStream,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewAdapter Google model: %s, baseURL: %s", model, baseURL)
		google.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewAdapter Google model: %s, path: %s", model, path)
		google.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAdapter Google model: %s, proxyURL: %s", model, proxyURL[0])
		google.proxyURL = proxyURL[0]
	}

	return google
}

func NewGcpAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *Google {

	logger.Infof(ctx, "NewGcpAdapter Google model: %s, key: %s", model, key)

	gcp := &Google{
		model:   model,
		key:     key,
		baseURL: "https://us-east5-aiplatform.googleapis.com/v1",
		path:    "/projects/%s/locations/us-east5/publishers/google/models/%s",
		header: g.MapStrStr{
			"Authorization": "Bearer " + key,
		},
		isSupportSystemRole: isSupportSystemRole,
		isSupportStream:     isSupportStream,
		isGcp:               true,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewGcpAdapter Google model: %s, baseURL: %s", model, baseURL)
		gcp.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewGcpAdapter Google model: %s, path: %s", model, path)
		gcp.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewGcpAdapter Google model: %s, proxyURL: %s", model, proxyURL[0])
		gcp.proxyURL = proxyURL[0]
	}

	return gcp
}

func (g *Google) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return sdkerr.NewRequestError(500, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)))
}

func (g *Google) apiErrorHandler(response *model.GoogleChatCompletionRes) error {
	return sdkerr.NewApiError(500, response.Error.Code, gjson.MustEncodeString(response), "api_error", "")
}
