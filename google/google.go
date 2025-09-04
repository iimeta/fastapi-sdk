package google

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/iimeta/fastapi-sdk/errors"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/options"
)

type Google struct {
	*options.AdapterOptions
	header map[string]string
	isGcp  bool
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) *Google {

	google := &Google{
		AdapterOptions: options,
	}

	if google.BaseUrl == "" {
		google.BaseUrl = "https://generativelanguage.googleapis.com/v1beta"
	}

	if google.Path == "" {
		google.Path = "/models/" + google.Model
	}

	logger.Infof(ctx, "NewAdapter Google model: %s, key: %s", google.Model, google.Key)

	return google
}

func NewGcpAdapter(ctx context.Context, options *options.AdapterOptions) *Google {

	gcp := &Google{
		AdapterOptions: options,
		header: g.MapStrStr{
			"Authorization": "Bearer " + options.Key,
		},
		isGcp: true,
	}

	if gcp.BaseUrl == "" {
		gcp.BaseUrl = "https://us-east5-aiplatform.googleapis.com/v1"
	}

	if gcp.Path == "" {
		gcp.Path = "/projects/%s/locations/us-east5/publishers/google/models/%s"
	}

	logger.Infof(ctx, "NewGcpAdapter Google model: %s, key: %s", gcp.Model, gcp.Key)

	return gcp
}

func (g *Google) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return errors.NewRequestError(response.StatusCode, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)))
}

func (g *Google) apiErrorHandler(response *model.GoogleChatCompletionRes) error {
	return errors.NewApiError(response.Error.Code, response.Error.Code, gjson.MustEncodeString(response), "api_error", "")
}
