package deepseek

import (
	"context"

	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (d *DeepSeek) AudioSpeech(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) AudioTranscriptions(ctx context.Context, request model.AudioRequest) (response model.AudioResponse, err error) {
	//TODO implement me
	panic("implement me")
}
