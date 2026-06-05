package volcengine

import (
	"context"

	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (v *VolcEngine) ImageGenerations(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ImageEdits(ctx context.Context, request model.ImageEditRequest) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ImageGenerationsStream(ctx context.Context, data []byte) (responseChan chan *model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ImageEditsStream(ctx context.Context, request model.ImageEditRequest) (responseChan chan *model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}
