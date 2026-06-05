package deepseek

import (
	"context"

	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (d *DeepSeek) ImageGenerations(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ImageEdits(ctx context.Context, request model.ImageEditRequest) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ImageGenerationsStream(ctx context.Context, data []byte) (responseChan chan *model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) ImageEditsStream(ctx context.Context, request model.ImageEditRequest) (responseChan chan *model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}
