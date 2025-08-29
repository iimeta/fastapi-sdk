package general

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/iimeta/fastapi-sdk/errors"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/options"
)

type General struct {
	*options.AdapterOptions
	header map[string]string
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) *General {

	general := &General{
		AdapterOptions: options,
		header: g.MapStrStr{
			"Authorization": "Bearer " + options.Key,
		},
	}

	logger.Infof(ctx, "NewAdapter General model: %s, key: %s", general.Model, general.Key)

	return general
}

func (g *General) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return errors.NewRequestError(response.StatusCode, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)))
}

func (g *General) apiErrorHandler(err error) error {

	apiError := &errors.ApiError{}
	if errors.As(err, &apiError) {
		return errors.NewRequestError(apiError.HttpStatusCode, apiError)
	}

	reqError := &errors.RequestError{}
	if errors.As(err, &reqError) {
		return errors.NewRequestError(apiError.HttpStatusCode, reqError.Err)
	}

	return err
}
