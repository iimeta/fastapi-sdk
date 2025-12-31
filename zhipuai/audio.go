package zhipuai

import (
	"context"

	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (z *ZhipuAI) AudioSpeech(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZhipuAI) AudioTranscriptions(ctx context.Context, request model.AudioRequest) (response model.AudioResponse, err error) {
	//TODO implement me
	panic("implement me")
}
