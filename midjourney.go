package sdk

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func NewMidjourneyProxy(ctx context.Context, baseURL, apiSecret, apiSecretHeader string) *model.MidjourneyProxy {
	return &model.MidjourneyProxy{
		ApiSecret:       apiSecret,
		ApiSecretHeader: apiSecretHeader,
		ImagineUrl:      baseURL + "/submit/imagine",
		ChangeUrl:       baseURL + "/submit/change",
		DescribeUrl:     baseURL + "/submit/describe",
		BlendUrl:        baseURL + "/submit/blend",
		FetchUrl:        baseURL + "/task/${task_id}/fetch",
	}
}

func Imagine(ctx context.Context, midjourneyProxy *model.MidjourneyProxy, request model.MidjourneyProxyImagineReq) (res model.MidjourneyProxyImagineRes, err error) {

	logger.Infof(ctx, "Midjourney Imagine prompt: %s start", request.Prompt)

	now := gtime.Now().UnixMilli()

	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney Imagine prompt: %s totalTime: %d ms", request.Prompt, gtime.Now().UnixMilli()-now)
	}()

	header := make(map[string]string)
	header[midjourneyProxy.ApiSecretHeader] = midjourneyProxy.ApiSecret

	if err = util.HttpPostJson(ctx, midjourneyProxy.ImagineUrl, header, request, &res); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	return res, nil
}

func Change(ctx context.Context, midjourneyProxy *model.MidjourneyProxy, request model.MidjourneyProxyChangeReq) (res model.MidjourneyProxyChangeRes, err error) {

	logger.Infof(ctx, "Midjourney Change request: %s start", gjson.MustEncodeString(request))

	now := gtime.Now().UnixMilli()

	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney Change request: %s totalTime: %d ms", gjson.MustEncodeString(request), gtime.Now().UnixMilli()-now)
	}()

	header := make(map[string]string)
	header[midjourneyProxy.ApiSecretHeader] = midjourneyProxy.ApiSecret

	if err = util.HttpPostJson(ctx, midjourneyProxy.ChangeUrl, header, request, &res); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	return res, nil
}

func Describe(ctx context.Context, midjourneyProxy *model.MidjourneyProxy, request model.MidjourneyProxyDescribeReq) (res model.MidjourneyProxyDescribeRes, err error) {

	logger.Info(ctx, "Midjourney Describe start")

	now := gtime.Now().UnixMilli()

	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney Describe totalTime: %d ms", gtime.Now().UnixMilli()-now)
	}()

	header := make(map[string]string)
	header[midjourneyProxy.ApiSecretHeader] = midjourneyProxy.ApiSecret

	if err = util.HttpPostJson(ctx, midjourneyProxy.DescribeUrl, header, request, &res); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	return res, nil
}

func Blend(ctx context.Context, midjourneyProxy *model.MidjourneyProxy, request model.MidjourneyProxyBlendReq) (res model.MidjourneyProxyBlendRes, err error) {

	logger.Info(ctx, "Midjourney Blend start")

	now := gtime.Now().UnixMilli()

	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney Blend totalTime: %d ms", gtime.Now().UnixMilli()-now)
	}()

	header := make(map[string]string)
	header[midjourneyProxy.ApiSecretHeader] = midjourneyProxy.ApiSecret

	if err = util.HttpPostJson(ctx, midjourneyProxy.BlendUrl, header, request, &res); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	return res, nil
}

func Fetch(ctx context.Context, midjourneyProxy *model.MidjourneyProxy, taskId string) (res model.MidjourneyProxyFetchRes, err error) {

	logger.Infof(ctx, "Midjourney Fetch taskId: %s start", taskId)

	now := gtime.Now().UnixMilli()

	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Midjourney Fetch taskId: %s totalTime: %d ms", taskId, gtime.Now().UnixMilli()-now)
	}()

	header := make(map[string]string)
	header[midjourneyProxy.ApiSecretHeader] = midjourneyProxy.ApiSecret

	fetchUrl := gstr.Replace(midjourneyProxy.FetchUrl, "${task_id}", taskId, -1)

	if err = util.HttpGet(ctx, fetchUrl, header, nil, &res); err != nil {
		logger.Error(ctx, err)
		return res, err
	}

	logger.Infof(ctx, "midjourneyProxyFetchRes: %s", gjson.MustEncodeString(res))

	return res, nil
}
