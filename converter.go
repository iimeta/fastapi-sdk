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
	"github.com/iimeta/fastapi-sdk/google"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/openai"
	"github.com/iimeta/fastapi-sdk/volcengine"
	"github.com/iimeta/fastapi-sdk/xfyun"
	"github.com/iimeta/fastapi-sdk/zhipuai"
)

type Converter interface {
	ConvChatCompletionsRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error)
	ConvChatCompletionsResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error)
	ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error)

	ConvChatCompletionsRequestOfficial(ctx context.Context, data []byte) ([]byte, error)
	ConvChatCompletionsResponseOfficial(ctx context.Context, data []byte) (model.ChatCompletionResponse, error)
	ConvChatCompletionsStreamResponseOfficial(ctx context.Context, data []byte) (model.ChatCompletionResponse, error)

	ConvChatResponsesRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error)
	ConvChatResponsesResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error)
	ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error)

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

func NewConverter(ctx context.Context, corp string) Converter {

	logger.Infof(ctx, "NewConverter corp: %s", corp)

	switch corp {
	case consts.CORP_OPENAI:
		return &openai.OpenAI{}
	case consts.CORP_AZURE:
		return &openai.OpenAI{}
	case consts.CORP_BAIDU:
		return &baidu.Baidu{}
	case consts.CORP_XFYUN:
		return &xfyun.Xfyun{}
	case consts.CORP_ALIYUN:
		return &aliyun.Aliyun{}
	case consts.CORP_ZHIPUAI:
		return &zhipuai.ZhipuAI{}
	case consts.CORP_GOOGLE:
		return &google.Google{}
	case consts.CORP_GCP_GEMINI:
		return &google.Google{}
	case consts.CORP_DEEPSEEK:
		return &deepseek.DeepSeek{}
	case consts.CORP_DEEPSEEK_BAIDU:
		return &deepseek.DeepSeek{}
	case consts.CORP_360AI:
		return &ai360.AI360{}
	case consts.CORP_ANTHROPIC:
		return &anthropic.Anthropic{}
	case consts.CORP_GCP_CLAUDE:
		return &anthropic.Anthropic{}
	case consts.CORP_AWS_CLAUDE:
		return &anthropic.Anthropic{}
	case consts.CORP_VOLC_ENGINE:
		return &volcengine.VolcEngine{}
	}

	return &openai.OpenAI{}
}
