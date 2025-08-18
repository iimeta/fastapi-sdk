package sdk

import (
	"context"

	"github.com/iimeta/fastapi-sdk/ai360"
	"github.com/iimeta/fastapi-sdk/aliyun"
	"github.com/iimeta/fastapi-sdk/anthropic"
	"github.com/iimeta/fastapi-sdk/baidu"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/deepseek"
	"github.com/iimeta/fastapi-sdk/google"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/openai"
	"github.com/iimeta/fastapi-sdk/volcengine"
	"github.com/iimeta/fastapi-sdk/xfyun"
	"github.com/iimeta/fastapi-sdk/zhipuai"
)

type Adapter interface {
	ChatCompletions(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error)
	ChatCompletionsStream(ctx context.Context, data []byte) (responseChan chan *model.ChatCompletionResponse, err error)

	ImageGenerations(ctx context.Context, request model.ImageGenerationRequest) (response model.ImageResponse, err error)
	ImageEdits(ctx context.Context, request model.ImageEditRequest) (response model.ImageResponse, err error)

	AudioSpeech(ctx context.Context, request model.SpeechRequest) (response model.SpeechResponse, err error)
	AudioTranscriptions(ctx context.Context, request model.AudioRequest) (response model.AudioResponse, err error)

	TextEmbeddings(ctx context.Context, request model.EmbeddingRequest) (response model.EmbeddingResponse, err error)
	TextModerations(ctx context.Context, request model.ModerationRequest) (response model.ModerationResponse, err error)
}

func NewAdapter(ctx context.Context, corp, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) Adapter {

	logger.Infof(ctx, "NewAdapter corp: %s, model: %s, key: %s", corp, model, key)

	switch corp {
	case consts.CORP_OPENAI:
		return openai.NewAdapter(ctx, model, key, baseURL, path, isSupportSystemRole, isSupportStream, proxyURL...)
	case consts.CORP_AZURE:
		return openai.NewAzureAdapter(ctx, model, key, baseURL, path, isSupportSystemRole, isSupportStream, proxyURL...)
	case consts.CORP_BAIDU:
		return baidu.NewAdapter(ctx, model, key, baseURL, path, isSupportSystemRole, isSupportStream, proxyURL...)
	case consts.CORP_XFYUN:
		return xfyun.NewAdapter(ctx, model, key, baseURL, path, isSupportSystemRole, isSupportStream, proxyURL...)
	case consts.CORP_ALIYUN:
		return aliyun.NewAdapter(ctx, model, key, baseURL, path, isSupportSystemRole, isSupportStream, proxyURL...)
	case consts.CORP_ZHIPUAI:
		return zhipuai.NewAdapter(ctx, model, key, baseURL, path, isSupportSystemRole, isSupportStream, proxyURL...)
	case consts.CORP_GOOGLE:
		return google.NewAdapter(ctx, model, key, baseURL, path, isSupportSystemRole, isSupportStream, proxyURL...)
	case consts.CORP_GCP_GEMINI:
		return google.NewGcpAdapter(ctx, model, key, baseURL, path, isSupportSystemRole, isSupportStream, proxyURL...)
	case consts.CORP_DEEPSEEK:
		return deepseek.NewAdapter(ctx, model, key, baseURL, path, isSupportSystemRole, isSupportStream, proxyURL...)
	case consts.CORP_DEEPSEEK_BAIDU:
		return deepseek.NewAdapterBaidu(ctx, model, key, baseURL, path, isSupportSystemRole, isSupportStream, proxyURL...)
	case consts.CORP_360AI:
		return ai360.NewAdapter(ctx, model, key, baseURL, path, isSupportSystemRole, isSupportStream, proxyURL...)
	case consts.CORP_ANTHROPIC:
		return anthropic.NewAdapter(ctx, model, key, baseURL, path, isSupportSystemRole, isSupportStream, proxyURL...)
	case consts.CORP_GCP_CLAUDE:
		return anthropic.NewGcpAdapter(ctx, model, key, baseURL, path, isSupportSystemRole, isSupportStream, proxyURL...)
	case consts.CORP_AWS_CLAUDE:
		return anthropic.NewAwsAdapter(ctx, model, key, baseURL, path, isSupportSystemRole, isSupportStream, proxyURL...)
	case consts.CORP_VOLC_ENGINE:
		return volcengine.NewAdapter(ctx, model, key, baseURL, path, isSupportSystemRole, isSupportStream, proxyURL...)
	}

	return openai.NewAdapter(ctx, model, key, baseURL, path, isSupportSystemRole, isSupportStream, proxyURL...)
}
