package volcengine

import (
	"context"

	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (v *VolcEngine) AudioSpeech(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) AudioTranscriptions(ctx context.Context, request model.AudioRequest) (response model.AudioResponse, err error) {
	//TODO implement me
	panic("implement me")
}
