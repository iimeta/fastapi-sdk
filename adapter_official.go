package sdk

import (
	"context"

	"github.com/iimeta/fastapi-sdk/v2/anthropic"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/general"
	"github.com/iimeta/fastapi-sdk/v2/google"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/options"
	"github.com/iimeta/fastapi-sdk/v2/volcengine"
)

type AdapterOfficialGroup interface {
	AdapterOfficial
	Converter
}

type AdapterOfficial interface {
	ChatCompletionsOfficial(ctx context.Context, data []byte) (response any, err error)
	ChatCompletionsStreamOfficial(ctx context.Context, data []byte) (responseChan chan any, err error)
	VideoCreateOfficial(ctx context.Context, data []byte) (responseBytes []byte, err error)
	VideoListOfficial(ctx context.Context, params model.VolcVideoListReq) (responseBytes []byte, err error)
	VideoRetrieveOfficial(ctx context.Context, taskId string) (responseBytes []byte, err error)
	VideoDeleteOfficial(ctx context.Context, taskId string) (err error)
}

func NewAdapterOfficial(ctx context.Context, options *options.AdapterOptions) AdapterOfficialGroup {

	logger.Infof(ctx, "NewAdapterOfficial provider: %s", options.Provider)

	switch options.Provider {
	case consts.PROVIDER_ANTHROPIC:
		return anthropic.NewAdapter(ctx, options)
	case consts.PROVIDER_GOOGLE:
		return google.NewAdapter(ctx, options)
	case consts.PROVIDER_AWS_CLAUDE:
		return anthropic.NewAwsAdapter(ctx, options)
	case consts.PROVIDER_GCP_CLAUDE:
		return anthropic.NewGcpAdapter(ctx, options)
	case consts.PROVIDER_GCP_GEMINI:
		return google.NewGcpAdapter(ctx, options)
	case consts.PROVIDER_VOLC_ENGINE:
		return volcengine.NewAdapter(ctx, options)
	default:
		return general.NewAdapter(ctx, options)
	}
}
