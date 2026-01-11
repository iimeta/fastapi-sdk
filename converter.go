package sdk

import (
	"bytes"
	"context"

	"github.com/iimeta/fastapi-sdk/v2/aliyun"
	"github.com/iimeta/fastapi-sdk/v2/anthropic"
	"github.com/iimeta/fastapi-sdk/v2/baidu"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/deepseek"
	"github.com/iimeta/fastapi-sdk/v2/general"
	"github.com/iimeta/fastapi-sdk/v2/google"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/openai"
	"github.com/iimeta/fastapi-sdk/v2/options"
	"github.com/iimeta/fastapi-sdk/v2/volcengine"
	"github.com/iimeta/fastapi-sdk/v2/xfyun"
	"github.com/iimeta/fastapi-sdk/v2/zhipuai"
)

type Converter interface {
	ConvChatCompletionsRequest(ctx context.Context, data any) (request model.ChatCompletionRequest, err error)
	ConvChatCompletionsResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error)
	ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error)

	ConvChatCompletionsRequestOfficial(ctx context.Context, request model.ChatCompletionRequest) ([]byte, error)
	ConvChatCompletionsResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error)
	ConvChatCompletionsStreamResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error)

	ConvChatResponsesRequest(ctx context.Context, data []byte) (request model.ChatCompletionRequest, err error)
	ConvChatResponsesResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error)
	ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error)

	ConvImageGenerationsRequest(ctx context.Context, data []byte) (request model.ImageGenerationRequest, err error)
	ConvImageGenerationsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error)
	ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (data *bytes.Buffer, err error)
	ConvImageEditsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error)

	ConvAudioSpeechRequest(ctx context.Context, data []byte) (request model.SpeechRequest, err error)
	ConvAudioSpeechResponse(ctx context.Context, data []byte) (response model.SpeechResponse, err error)
	ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (data *bytes.Buffer, err error)
	ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (response model.AudioResponse, err error)

	ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (request model.EmbeddingRequest, err error)
	ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error)

	ConvVideoCreateRequest(ctx context.Context, request model.VideoCreateRequest) (data *bytes.Buffer, err error)
	ConvVideoListResponse(ctx context.Context, data []byte) (response model.VideoListResponse, err error)
	ConvVideoContentResponse(ctx context.Context, data []byte) (response model.VideoContentResponse, err error)
	ConvVideoJobResponse(ctx context.Context, data []byte) (response model.VideoJobResponse, err error)

	ConvFileUploadRequest(ctx context.Context, request model.FileUploadRequest) (data *bytes.Buffer, err error)
	ConvFileListResponse(ctx context.Context, data []byte) (response model.FileListResponse, err error)
	ConvFileContentResponse(ctx context.Context, data []byte) (response model.FileContentResponse, err error)
	ConvFileResponse(ctx context.Context, data []byte) (response model.FileResponse, err error)

	ConvBatchCreateRequest(ctx context.Context, request model.BatchCreateRequest) (data *bytes.Buffer, err error)
	ConvBatchListResponse(ctx context.Context, data []byte) (response model.BatchListResponse, err error)
	ConvBatchResponse(ctx context.Context, data []byte) (response model.BatchResponse, err error)
}

func NewConverter(ctx context.Context, options *options.AdapterOptions) Converter {

	logger.Infof(ctx, "NewConverter provider: %s", options.Provider)

	switch options.Provider {
	case consts.PROVIDER_OPENAI:
		return &openai.OpenAI{AdapterOptions: options}
	case consts.PROVIDER_ANTHROPIC:
		return &anthropic.Anthropic{AdapterOptions: options}
	case consts.PROVIDER_GOOGLE:
		return &google.Google{AdapterOptions: options}
	case consts.PROVIDER_AZURE:
		return &openai.OpenAI{AdapterOptions: options}
	case consts.PROVIDER_DEEPSEEK:
		return &deepseek.DeepSeek{AdapterOptions: options}
	case consts.PROVIDER_DEEPSEEK_BAIDU:
		return &deepseek.DeepSeek{AdapterOptions: options}
	case consts.PROVIDER_BAIDU:
		return &baidu.Baidu{AdapterOptions: options}
	case consts.PROVIDER_ALIYUN:
		return &aliyun.Aliyun{AdapterOptions: options}
	case consts.PROVIDER_XFYUN:
		return &xfyun.Xfyun{AdapterOptions: options}
	case consts.PROVIDER_ZHIPUAI:
		return &zhipuai.ZhipuAI{AdapterOptions: options}
	case consts.PROVIDER_VOLC_ENGINE:
		return &volcengine.VolcEngine{AdapterOptions: options}
	case consts.PROVIDER_AWS_CLAUDE:
		return &anthropic.Anthropic{AdapterOptions: options}
	case consts.PROVIDER_GCP_CLAUDE:
		return &anthropic.Anthropic{AdapterOptions: options}
	case consts.PROVIDER_GCP_GEMINI:
		return &google.Google{AdapterOptions: options}
	default:
		return &general.General{AdapterOptions: options}
	}
}
