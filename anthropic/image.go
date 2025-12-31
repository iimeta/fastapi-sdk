package anthropic

import (
	"context"

	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (a *Anthropic) ImageGenerations(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ImageEdits(ctx context.Context, request model.ImageEditRequest) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}
