package sdk

import (
	"bytes"
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

type Converter interface {
	ConvChatCompletionsRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error)
	ConvChatCompletionsResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error)
	ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error)

	ConvChatCompletionsRequestOfficial(ctx context.Context, data []byte) ([]byte, error)
	ConvChatCompletionsResponseOfficial(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error)
	ConvChatCompletionsStreamResponseOfficial(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error)

	ConvChatResponsesRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error)
	ConvChatResponsesResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error)
	ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error)

	ConvImageGenerationsRequest(ctx context.Context, data []byte) (model.ImageGenerationRequest, error)
	ConvImageGenerationsResponse(ctx context.Context, data []byte) (model.ImageResponse, error)
	ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (*bytes.Buffer, error)
	ConvImageEditsResponse(ctx context.Context, data []byte) (model.ImageResponse, error)

	ConvAudioSpeechRequest(ctx context.Context, data []byte) (model.SpeechRequest, error)
	ConvAudioSpeechResponse(ctx context.Context, data []byte) (model.SpeechResponse, error)
	ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (*bytes.Buffer, error)
	ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (model.AudioResponse, error)

	ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (model.EmbeddingRequest, error)
	ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (model.EmbeddingResponse, error)
}

func NewConverter(ctx context.Context, corp string, opts ...*options.AdapterOptions) Converter {

	logger.Infof(ctx, "NewConverter corp: %s", corp)

	if len(opts) == 0 {
		opts = append(opts, &options.AdapterOptions{})
	}

	switch corp {
	case consts.CORP_OPENAI:
		return &openai.OpenAI{AdapterOptions: opts[0]}
	case consts.CORP_AZURE:
		return &openai.OpenAI{AdapterOptions: opts[0]}
	case consts.CORP_BAIDU:
		return &baidu.Baidu{AdapterOptions: opts[0]}
	case consts.CORP_XFYUN:
		return &xfyun.Xfyun{AdapterOptions: opts[0]}
	case consts.CORP_ALIYUN:
		return &aliyun.Aliyun{AdapterOptions: opts[0]}
	case consts.CORP_ZHIPUAI:
		return &zhipuai.ZhipuAI{AdapterOptions: opts[0]}
	case consts.CORP_GOOGLE:
		return &google.Google{AdapterOptions: opts[0]}
	case consts.CORP_GCP_GEMINI:
		return &google.Google{AdapterOptions: opts[0]}
	case consts.CORP_DEEPSEEK:
		return &deepseek.DeepSeek{AdapterOptions: opts[0]}
	case consts.CORP_DEEPSEEK_BAIDU:
		return &deepseek.DeepSeek{AdapterOptions: opts[0]}
	case consts.CORP_360AI:
		return &ai360.AI360{AdapterOptions: opts[0]}
	case consts.CORP_ANTHROPIC:
		return &anthropic.Anthropic{AdapterOptions: opts[0]}
	case consts.CORP_GCP_CLAUDE:
		return &anthropic.Anthropic{AdapterOptions: opts[0]}
	case consts.CORP_AWS_CLAUDE:
		return &anthropic.Anthropic{AdapterOptions: opts[0]}
	case consts.CORP_VOLC_ENGINE:
		return &volcengine.VolcEngine{AdapterOptions: opts[0]}
	}

	return &general.General{AdapterOptions: opts[0]}
}
