package sdk

import (
	"context"

	"github.com/iimeta/fastapi-sdk/anthropic"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/general"
	"github.com/iimeta/fastapi-sdk/google"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/options"
)

type AdapterOfficialGroup interface {
	AdapterOfficial
	Converter
}

type AdapterOfficial interface {
	ChatCompletionsOfficial(ctx context.Context, data []byte) (response any, err error)
	ChatCompletionsStreamOfficial(ctx context.Context, data []byte) (responseChan chan any, err error)
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
	default:
		return general.NewAdapter(ctx, options)
	}
}
