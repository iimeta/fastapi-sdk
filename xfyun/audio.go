package xfyun

import (
	"context"

	"github.com/iimeta/fastapi-sdk/model"
)

func (x *Xfyun) AudioSpeech(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) AudioTranscriptions(ctx context.Context, request model.AudioRequest) (response model.AudioResponse, err error) {
	//TODO implement me
	panic("implement me")
}
