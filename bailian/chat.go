package bailian

import (
	"context"

	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (b *Bailian) ChatCompletions(ctx context.Context, data any) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bailian) ChatCompletionsStream(ctx context.Context, data any) (responseChan chan *model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}
