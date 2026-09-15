package minimax

import (
	"context"

	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (m *MiniMax) ImageGenerations(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ImageGenerationsStream(ctx context.Context, data []byte) (responseChan chan *model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ImageEdits(ctx context.Context, request model.ImageEditRequest) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ImageEditsStream(ctx context.Context, request model.ImageEditRequest) (responseChan chan *model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}
