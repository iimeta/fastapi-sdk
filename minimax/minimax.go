package minimax

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/options"
)

type MiniMax struct {
	*options.AdapterOptions
	header map[string]string
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) *MiniMax {

	minimax := &MiniMax{
		AdapterOptions: options,
		header: map[string]string{
			"Authorization": "Bearer " + options.Key,
		},
	}

	if minimax.BaseUrl == "" {
		minimax.BaseUrl = "https://api.minimax.cn"
	}

	for k, v := range minimax.PassthroughHeader {
		minimax.header[k] = v
	}

	for k, v := range minimax.Header {
		minimax.header[k] = v
	}

	logger.Infof(ctx, "NewAdapter MiniMax model: %s, key: %s", minimax.Model, minimax.Key)

	return minimax
}

func (m *MiniMax) requestErrorHandler(ctx context.Context, response *http.Response) (err error) {

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	res := model.MiniMaxErrorRes{}
	if err = json.Unmarshal(bytes, &res); err == nil && res.Error != nil {
		return m.apiErrorHandler(response.StatusCode, res.Error)
	}

	return errors.NewRequestError(response.StatusCode, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, bytes)))
}

func (m *MiniMax) apiErrorHandler(statusCode int, detail *model.MiniMaxErrorDetail) error {

	switch detail.Type {
	case "authorized_error":
		return errors.ERR_INVALID_API_KEY
	case "insufficient_balance_error":
		return errors.ERR_INSUFFICIENT_QUOTA
	case "rate_limit_error":
		return errors.ERR_RATE_LIMIT_EXCEEDED
	}

	if statusCode == 0 {
		statusCode = 500
	}

	return errors.NewApiError(statusCode, detail.Type, detail.Message, detail.Type, "")
}
