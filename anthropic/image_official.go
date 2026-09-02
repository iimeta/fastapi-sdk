package anthropic

import (
	"context"
	"net/http"

	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (a *Anthropic) ImageGenerationsOfficial(ctx context.Context, data []byte) (responseBytes []byte, responseHeader http.Header, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ImageGenerationsStreamOfficial(ctx context.Context, data []byte) (responseChan chan *model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}
