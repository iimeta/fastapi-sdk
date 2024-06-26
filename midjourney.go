package sdk

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func NewMidjourneyProxy(ctx context.Context, baseURL, apiSecret, apiSecretHeader string) *model.MidjourneyProxy {
	return &model.MidjourneyProxy{
		ApiSecret:              apiSecret,
		ApiSecretHeader:        apiSecretHeader,
		ImagineUrl:             baseURL + "/submit/imagine",
		ChangeUrl:              baseURL + "/submit/change",
		DescribeUrl:            baseURL + "/submit/describe",
		BlendUrl:               baseURL + "/submit/blend",
		SwapFaceUrl:            baseURL + "/insight-face/swap",
		ActionUrl:              baseURL + "/submit/action",
		ModalUrl:               baseURL + "/submit/modal",
		ShortenUrl:             baseURL + "/submit/shorten",
		UploadDiscordImagesUrl: baseURL + "/submit/upload-discord-images",
		FetchUrl:               baseURL + "/task/${taskId}/fetch",
	}
}

func Imagine(ctx context.Context, midjourneyProxy *model.MidjourneyProxy, request model.MidjourneyProxyRequest) (res model.MidjourneyProxyResponse, err error) {

	logger.Infof(ctx, "Midjourney Imagine prompt: %s start", request.Prompt)

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney Imagine prompt: %s totalTime: %d ms", request.Prompt, gtime.Now().UnixMilli()-now)
	}()

	if err = util.HttpPost(ctx, midjourneyProxy.ImagineUrl, g.MapStrStr{midjourneyProxy.ApiSecretHeader: midjourneyProxy.ApiSecret}, request, &res, ""); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	return res, nil
}

func Change(ctx context.Context, midjourneyProxy *model.MidjourneyProxy, request model.MidjourneyProxyRequest) (res model.MidjourneyProxyResponse, err error) {

	logger.Infof(ctx, "Midjourney Change request: %s start", gjson.MustEncodeString(request))

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney Change request: %s totalTime: %d ms", gjson.MustEncodeString(request), gtime.Now().UnixMilli()-now)
	}()

	if err = util.HttpPost(ctx, midjourneyProxy.ChangeUrl, g.MapStrStr{midjourneyProxy.ApiSecretHeader: midjourneyProxy.ApiSecret}, request, &res, ""); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	return res, nil
}

func Describe(ctx context.Context, midjourneyProxy *model.MidjourneyProxy, request model.MidjourneyProxyRequest) (res model.MidjourneyProxyResponse, err error) {

	logger.Info(ctx, "Midjourney Describe start")

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney Describe totalTime: %d ms", gtime.Now().UnixMilli()-now)
	}()

	if err = util.HttpPost(ctx, midjourneyProxy.DescribeUrl, g.MapStrStr{midjourneyProxy.ApiSecretHeader: midjourneyProxy.ApiSecret}, request, &res, ""); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	return res, nil
}

func Blend(ctx context.Context, midjourneyProxy *model.MidjourneyProxy, request model.MidjourneyProxyRequest) (res model.MidjourneyProxyResponse, err error) {

	logger.Info(ctx, "Midjourney Blend start")

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney Blend totalTime: %d ms", gtime.Now().UnixMilli()-now)
	}()

	if err = util.HttpPost(ctx, midjourneyProxy.BlendUrl, g.MapStrStr{midjourneyProxy.ApiSecretHeader: midjourneyProxy.ApiSecret}, request, &res, ""); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	return res, nil
}

func SwapFace(ctx context.Context, midjourneyProxy *model.MidjourneyProxy, request model.MidjourneyProxyRequest) (res model.MidjourneyProxyResponse, err error) {

	logger.Info(ctx, "Midjourney SwapFace start")

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney SwapFace totalTime: %d ms", gtime.Now().UnixMilli()-now)
	}()

	if err = util.HttpPost(ctx, midjourneyProxy.SwapFaceUrl, g.MapStrStr{midjourneyProxy.ApiSecretHeader: midjourneyProxy.ApiSecret}, request, &res, ""); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	return res, nil
}

func Action(ctx context.Context, midjourneyProxy *model.MidjourneyProxy, request model.MidjourneyProxyRequest) (res model.MidjourneyProxyResponse, err error) {

	logger.Info(ctx, "Midjourney Action start")

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney Action totalTime: %d ms", gtime.Now().UnixMilli()-now)
	}()

	if err = util.HttpPost(ctx, midjourneyProxy.ActionUrl, g.MapStrStr{midjourneyProxy.ApiSecretHeader: midjourneyProxy.ApiSecret}, request, &res, ""); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	return res, nil
}

func Modal(ctx context.Context, midjourneyProxy *model.MidjourneyProxy, request model.MidjourneyProxyRequest) (res model.MidjourneyProxyResponse, err error) {

	logger.Info(ctx, "Midjourney Modal start")

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney Modal totalTime: %d ms", gtime.Now().UnixMilli()-now)
	}()

	if err = util.HttpPost(ctx, midjourneyProxy.ModalUrl, g.MapStrStr{midjourneyProxy.ApiSecretHeader: midjourneyProxy.ApiSecret}, request, &res, ""); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	return res, nil
}

func Shorten(ctx context.Context, midjourneyProxy *model.MidjourneyProxy, request model.MidjourneyProxyRequest) (res model.MidjourneyProxyResponse, err error) {

	logger.Info(ctx, "Midjourney Shorten start")

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney Shorten totalTime: %d ms", gtime.Now().UnixMilli()-now)
	}()

	if err = util.HttpPost(ctx, midjourneyProxy.ShortenUrl, g.MapStrStr{midjourneyProxy.ApiSecretHeader: midjourneyProxy.ApiSecret}, request, &res, ""); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	return res, nil
}

func UploadDiscordImages(ctx context.Context, midjourneyProxy *model.MidjourneyProxy, request model.MidjourneyProxyRequest) (res model.MidjourneyProxyResponse, err error) {

	logger.Info(ctx, "Midjourney UploadDiscordImages start")

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney UploadDiscordImages totalTime: %d ms", gtime.Now().UnixMilli()-now)
	}()

	if err = util.HttpPost(ctx, midjourneyProxy.UploadDiscordImagesUrl, g.MapStrStr{midjourneyProxy.ApiSecretHeader: midjourneyProxy.ApiSecret}, request, &res, ""); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	return res, nil
}

func Fetch(ctx context.Context, midjourneyProxy *model.MidjourneyProxy, request model.MidjourneyProxyRequest) (res model.MidjourneyProxyFetchResponse, err error) {

	logger.Infof(ctx, "Midjourney Fetch taskId: %s start", request.TaskId)

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney Fetch taskId: %s totalTime: %d ms", request.TaskId, gtime.Now().UnixMilli()-now)
	}()

	fetchUrl := gstr.Replace(midjourneyProxy.FetchUrl, "${taskId}", request.TaskId, -1)

	if err = util.HttpGet(ctx, fetchUrl, g.MapStrStr{midjourneyProxy.ApiSecretHeader: midjourneyProxy.ApiSecret}, nil, &res, ""); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	logger.Infof(ctx, "midjourneyProxyFetchResponse: %s", gjson.MustEncodeString(res))

	return res, nil
}
