package minimax

import (
	"context"

	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (m *MiniMax) ChatCompletions(ctx context.Context, data any) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MiniMax) ChatCompletionsStream(ctx context.Context, data any) (responseChan chan *model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}
