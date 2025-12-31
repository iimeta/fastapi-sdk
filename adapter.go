package sdk

import (
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

	VideoCreate(ctx context.Context, request model.VideoCreateRequest) (response model.VideoJobResponse, err error)
	VideoRemix(ctx context.Context, request model.VideoRemixRequest) (response model.VideoJobResponse, err error)
	VideoList(ctx context.Context, request model.VideoListRequest) (response model.VideoListResponse, err error)
	VideoRetrieve(ctx context.Context, request model.VideoRetrieveRequest) (response model.VideoJobResponse, err error)
	VideoDelete(ctx context.Context, request model.VideoDeleteRequest) (response model.VideoJobResponse, err error)
	VideoContent(ctx context.Context, request model.VideoContentRequest) (response model.VideoContentResponse, err error)

	FileUpload(ctx context.Context, request model.FileUploadRequest) (response model.FileResponse, err error)
	FileList(ctx context.Context, request model.FileListRequest) (response model.FileListResponse, err error)
	FileRetrieve(ctx context.Context, request model.FileRetrieveRequest) (response model.FileResponse, err error)
	FileDelete(ctx context.Context, request model.FileDeleteRequest) (response model.FileResponse, err error)
	FileContent(ctx context.Context, request model.FileContentRequest) (response model.FileContentResponse, err error)

	BatchCreate(ctx context.Context, request model.BatchCreateRequest) (response model.BatchResponse, err error)
	BatchList(ctx context.Context, request model.BatchListRequest) (response model.BatchListResponse, err error)
	BatchRetrieve(ctx context.Context, request model.BatchRetrieveRequest) (response model.BatchResponse, err error)
	BatchCancel(ctx context.Context, request model.BatchCancelRequest) (response model.BatchResponse, err error)
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
