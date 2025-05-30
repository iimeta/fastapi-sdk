package xfyun

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (c *Client) ImageGeneration(ctx context.Context, request model.ImageGenerationRequest) (res model.ImageResponse, err error) {

	logger.Infof(ctx, "ImageGeneration Xfyun model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ImageGeneration Xfyun model: %s totalTime: %d ms", request.Model, gtime.TimestampMilli()-now)
	}()

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
			AppId: c.appId,
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

	imageRes := new(model.XfyunChatCompletionRes)
	if _, err = util.HttpPost(ctx, c.getHttpUrl(ctx), nil, imageReq, &imageRes, c.proxyURL); err != nil {
		logger.Errorf(ctx, "ImageGeneration Xfyun model: %s, error: %v", request.Model, err)
		return res, err
	}

	res = model.ImageResponse{
		Created: gtime.Timestamp(),
		Data: []model.ImageResponseDataInner{{
			B64JSON: imageRes.Payload.Choices.Text[0].Content,
		}},
	}

	return res, nil
}

func (c *Client) ImageEdit(ctx context.Context, request model.ImageEditRequest) (res model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}
