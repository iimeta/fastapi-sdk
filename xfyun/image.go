package xfyun

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (x *Xfyun) ImageGenerations(ctx context.Context, data []byte) (response model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGenerations Xfyun model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageGenerations Xfyun model: %s totalTime: %d ms", x.Model, gtime.TimestampMilli()-now)
	}()

	request, err := x.ConvImageGenerationsRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "ImageGenerations Xfyun ConvImageGenerationsRequest error: %v", err)
		return response, err
	}

	width := 512
	height := 512

	if request.Size != "" {

		size := gstr.Split(request.Size, `×`)

		if len(size) != 2 {
			size = gstr.Split(request.Size, `x`)
		}

		if len(size) != 2 {
			size = gstr.Split(request.Size, `X`)
		}

		if len(size) != 2 {
			size = gstr.Split(request.Size, `*`)
		}

		if len(size) != 2 {
			size = gstr.Split(request.Size, `:`)
		}

		if len(size) == 2 {
			width = gconv.Int(size[0])
			height = gconv.Int(size[1])
		}
	}

	imageReq := model.XfyunChatCompletionReq{
		Header: model.Header{
			AppId: x.appId,
			Uid:   grand.Digits(10),
		},
		Parameter: model.Parameter{
			Chat: &model.Chat{
				Domain: "general",
				Width:  width,
				Height: height,
			},
		},
		Payload: model.Payload{
			Message: &model.Message{
				Text: []model.ChatCompletionMessage{{
					Role:    consts.ROLE_USER,
					Content: request.Prompt,
				}},
			},
		},
	}

	imageRes := model.XfyunChatCompletionRes{}
	if _, err = util.HttpPost(ctx, x.getHttpUrl(ctx), x.header, gjson.MustEncode(imageReq), &imageRes, x.Timeout, x.ProxyUrl, x.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "ImageGenerations Xfyun model: %s, error: %v", x.Model, err)
		return response, err
	}

	response = model.ImageResponse{
		Created: gtime.Timestamp(),
		Data: []model.ImageResponseData{{
			B64Json: imageRes.Payload.Choices.Text[0].Content,
		}},
	}

	return response, nil
}

func (x *Xfyun) ImageEdits(ctx context.Context, request model.ImageEditRequest) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}
