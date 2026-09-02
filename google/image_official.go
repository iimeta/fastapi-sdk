package google

import (
	"context"
	"net/http"

	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (g *Google) ImageGenerationsOfficial(ctx context.Context, data []byte) (responseBytes []byte, responseHeader http.Header, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ImageGenerationsStreamOfficial(ctx context.Context, data []byte) (responseChan chan *model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}
