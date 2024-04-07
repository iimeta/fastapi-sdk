package sdk

import (
	"context"
	"github.com/iimeta/fastapi-sdk/aliyun"
	"github.com/iimeta/fastapi-sdk/baidu"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/openai"
	"github.com/iimeta/fastapi-sdk/xfyun"
)

type Chat interface {
	ChatCompletion(ctx context.Context, request model.ChatCompletionRequest) (res model.ChatCompletionResponse, err error)
	ChatCompletionStream(ctx context.Context, request model.ChatCompletionRequest) (responseChan chan *model.ChatCompletionResponse, err error)
	Image(ctx context.Context, request model.ImageRequest) (res model.ImageResponse, err error)
}

func NewClient(ctx context.Context, corp, model, key string, baseURL ...string) Chat {

	logger.Infof(ctx, "NewClient corp: %s, model: %s, key: %s", corp, model, key)

	switch corp {
	case consts.CORP_OPENAI:
		return openai.NewClient(ctx, model, key, baseURL...)
	case consts.CORP_BAIDU:
		return baidu.NewClient(ctx, model, key, baseURL...)
	case consts.CORP_XFYUN:
		return xfyun.NewClient(ctx, model, key, baseURL...)
	case consts.CORP_ALIYUN:
		return aliyun.NewClient(ctx, model, key, baseURL...)
	}

	return xfyun.NewClient(ctx, model, key, baseURL...)
}

func NewProxyClient(ctx context.Context, corp, model, key string, proxyURL ...string) Chat {

	logger.Infof(ctx, "NewProxyClient corp: %s, model: %s, key: %s", corp, model, key)

	switch corp {
	case consts.CORP_OPENAI:
		return openai.NewProxyClient(ctx, model, key, proxyURL...)
	case consts.CORP_BAIDU:
		return baidu.NewProxyClient(ctx, model, key, proxyURL...)
	case consts.CORP_XFYUN:
		return xfyun.NewProxyClient(ctx, model, key, proxyURL...)
	case consts.CORP_ALIYUN:
		return aliyun.NewProxyClient(ctx, model, key, proxyURL...)
	}

	return openai.NewProxyClient(ctx, model, key, proxyURL...)
}
