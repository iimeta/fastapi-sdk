package anthropic

import (
	"context"
	"github.com/iimeta/fastapi-sdk/model"
)

func (c *Client) Speech(ctx context.Context, request model.SpeechRequest) (res model.SpeechResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) Transcription(ctx context.Context, request model.AudioRequest) (res model.AudioResponse, err error) {
	//TODO implement me
	panic("implement me")
}
