package bailian

import (
	"context"
	"slices"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

// 官方格式请求体直接透传, OpenAI 兼容格式转换为官方格式后请求
func (b *Bailian) ImageGenerations(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerations Bailian model: %s start", b.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageGenerations Bailian model: %s totalTime: %d ms", b.Model, response.TotalTime)
	}()

	var request any = data
	if !slices.Contains(b.ReqPassthroughParams, "req_data") && !isOfficialRequest(data) {

		imageGenerationRequest, err := b.ConvImageGenerationsRequest(ctx, data)
		if err != nil {
			logger.Errorf(ctx, "ImageGenerations Bailian ConvImageGenerationsRequest error: %v", err)
			return response, err
		}

		if request, err = b.ConvImageGenerationsRequestOfficial(ctx, imageGenerationRequest); err != nil {
			logger.Errorf(ctx, "ImageGenerations Bailian ConvImageGenerationsRequestOfficial error: %v", err)
			return response, err
		}
	}

	if b.Path == "" {
		b.Path = "/services/aigc/multimodal-generation/generation"
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, b.BaseUrl+b.Path, b.header, request, nil, b.Timeout, b.ProxyUrl, b.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageGenerations Bailian model: %s, error: %v", b.Model, err)
		return response, err
	}

	if response, err = b.ConvImageGenerationsResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "ImageGenerations Bailian ConvImageGenerationsResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "ImageGenerations Bailian model: %s finished", b.Model)

	return response, nil
}

// 百炼图像编辑与生成为同一接口, 编辑请求转换为官方格式后请求
func (b *Bailian) ImageEdits(ctx context.Context, request model.ImageEditRequest) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageEdits Bailian model: %s start", b.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageEdits Bailian model: %s totalTime: %d ms", b.Model, response.TotalTime)
	}()

	data, err := b.ConvImageEditsRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "ImageEdits Bailian ConvImageEditsRequest error: %v", err)
		return response, err
	}

	if b.Path == "" {
		b.Path = "/services/aigc/multimodal-generation/generation"
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, b.BaseUrl+b.Path, b.header, data, nil, b.Timeout, b.ProxyUrl, b.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ImageEdits Bailian model: %s, error: %v", b.Model, err)
		return response, err
	}

	if response, err = b.ConvImageEditsResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "ImageEdits Bailian ConvImageEditsResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "ImageEdits Bailian model: %s finished", b.Model)

	return response, nil
}

func (b *Bailian) ImageGenerationsStream(ctx context.Context, data []byte) (responseChan chan *model.ImageResponse, err error) {
	return nil, errors.New("Bailian image generations does not support stream")
}

func (b *Bailian) ImageEditsStream(ctx context.Context, request model.ImageEditRequest) (responseChan chan *model.ImageResponse, err error) {
	return nil, errors.New("Bailian image edits does not support stream")
}
