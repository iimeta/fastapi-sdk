package bailian

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (b *Bailian) ConvChatCompletionsRequestOfficial(ctx context.Context, request model.ChatCompletionRequest) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvChatCompletionsResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ConvChatCompletionsStreamResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

// OpenAI 兼容格式请求转为百炼官方格式请求体
func (b *Bailian) ConvImageGenerationsRequestOfficial(ctx context.Context, request model.ImageGenerationRequest) ([]byte, error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvImageGenerationsRequestOfficial time: %d", gtime.TimestampMilli()-now)
	}()

	model_ := request.Model
	if model_ == "" {
		model_ = b.Model
	}

	req := buildOfficialRequest(model_, request.Prompt, collectImages(request.Image, request.Images), request.N, request.Size, request.Quality, request.Watermark, request.SequentialImageGeneration)

	return json.Marshal(req)
}

func (b *Bailian) ConvImageGenerationsResponseOfficial(ctx context.Context, response model.ImageResponse) ([]byte, error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvImageGenerationsResponseOfficial time: %d", gtime.TimestampMilli()-now)
	}()

	if len(response.ResponseBytes) > 0 {
		return response.ResponseBytes, nil
	}

	return json.Marshal(response)
}

// OpenAI 兼容格式编辑请求转为百炼官方格式请求体(官方生成与编辑为同一接口)
func (b *Bailian) ConvImageEditsRequestOfficial(ctx context.Context, request model.ImageEditRequest) ([]byte, error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "ConvImageEditsRequestOfficial time: %d", gtime.TimestampMilli()-now)
	}()

	model_ := request.Model
	if model_ == "" {
		model_ = b.Model
	}

	req := buildOfficialRequest(model_, request.Prompt, collectImages(request.Image, request.Images), request.N, request.Size, request.Quality, nil, "")

	return json.Marshal(req)
}

func (b *Bailian) ConvImageEditsResponseOfficial(ctx context.Context, response model.ImageResponse) ([]byte, error) {
	return b.ConvImageGenerationsResponseOfficial(ctx, response)
}
