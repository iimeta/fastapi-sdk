package google

import (
	"context"
	"net/http"

	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (g *Google) VideoCreateOfficial(ctx context.Context, data []byte) (responseBytes []byte, responseHeader http.Header, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) VideoListOfficial(ctx context.Context, params model.VolcVideoListReq) (responseBytes []byte, responseHeader http.Header, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) VideoRetrieveOfficial(ctx context.Context, taskId string) (responseBytes []byte, responseHeader http.Header, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) VideoDeleteOfficial(ctx context.Context, taskId string) (err error) {
	//TODO implement me
	panic("implement me")
}
