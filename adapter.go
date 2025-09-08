package sdk

import (
	"context"

	"github.com/iimeta/fastapi-sdk/ai360"
	"github.com/iimeta/fastapi-sdk/aliyun"
	"github.com/iimeta/fastapi-sdk/anthropic"
	"github.com/iimeta/fastapi-sdk/baidu"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/deepseek"
	"github.com/iimeta/fastapi-sdk/general"
	"github.com/iimeta/fastapi-sdk/google"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/openai"
	"github.com/iimeta/fastapi-sdk/options"
	"github.com/iimeta/fastapi-sdk/volcengine"
	"github.com/iimeta/fastapi-sdk/xfyun"
	"github.com/iimeta/fastapi-sdk/zhipuai"
)

type AdapterGroup interface {
	Adapter
	Converter
}

type Adapter interface {
	ChatCompletions(ctx context.Context, data any) (response model.ChatCompletionResponse, err error)
	ChatCompletionsStream(ctx context.Context, data any) (responseChan chan *model.ChatCompletionResponse, err error)

	ImageGenerations(ctx context.Context, data []byte) (response model.ImageResponse, err error)
	ImageEdits(ctx context.Context, request model.ImageEditRequest) (response model.ImageResponse, err error)

	AudioSpeech(ctx context.Context, data []byte) (response model.SpeechResponse, err error)
	AudioTranscriptions(ctx context.Context, request model.AudioRequest) (response model.AudioResponse, err error)

	TextEmbeddings(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error)
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) AdapterGroup {

	logger.Infof(ctx, "NewAdapter provider: %s", options.Provider)

	switch options.Provider {
	case consts.PROVIDER_OPENAI:
		return openai.NewAdapter(ctx, options)
	case consts.PROVIDER_ANTHROPIC:
		return anthropic.NewAdapter(ctx, options)
	case consts.PROVIDER_GOOGLE:
		return google.NewAdapter(ctx, options)
	case consts.PROVIDER_AZURE:
		return openai.NewAzureAdapter(ctx, options)
	case consts.PROVIDER_DEEPSEEK:
		return deepseek.NewAdapter(ctx, options)
	case consts.PROVIDER_DEEPSEEK_BAIDU:
		return deepseek.NewAdapterBaidu(ctx, options)
	case consts.PROVIDER_BAIDU:
		return baidu.NewAdapter(ctx, options)
	case consts.PROVIDER_ALIYUN:
		return aliyun.NewAdapter(ctx, options)
	case consts.PROVIDER_XFYUN:
		return xfyun.NewAdapter(ctx, options)
	case consts.PROVIDER_ZHIPUAI:
		return zhipuai.NewAdapter(ctx, options)
	case consts.PROVIDER_VOLC_ENGINE:
		return volcengine.NewAdapter(ctx, options)
	case consts.PROVIDER_360AI:
		return ai360.NewAdapter(ctx, options)
	case consts.PROVIDER_AWS_CLAUDE:
		return anthropic.NewAwsAdapter(ctx, options)
	case consts.PROVIDER_GCP_CLAUDE:
		return anthropic.NewGcpAdapter(ctx, options)
	case consts.PROVIDER_GCP_GEMINI:
		return google.NewGcpAdapter(ctx, options)
	}

	return general.NewAdapter(ctx, options)
}
