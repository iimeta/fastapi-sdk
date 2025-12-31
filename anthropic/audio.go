package anthropic

import (
	"context"

	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (a *Anthropic) AudioSpeech(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) AudioTranscriptions(ctx context.Context, request model.AudioRequest) (response model.AudioResponse, err error) {
	//TODO implement me
	panic("implement me")
}
