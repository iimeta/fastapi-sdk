package general

import (
	"context"

	"github.com/iimeta/fastapi-sdk/model"
)

func (g *General) AudioSpeech(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *General) AudioTranscriptions(ctx context.Context, request model.AudioRequest) (response model.AudioResponse, err error) {
	//TODO implement me
	panic("implement me")
}
