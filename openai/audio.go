package openai

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (o *OpenAI) AudioSpeech(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {

	logger.Infof(ctx, "AudioSpeech OpenAI model: %s start", o.model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "AudioSpeech OpenAI model: %s totalTime: %d ms", o.model, response.TotalTime)
	}()

	request, err := o.ConvAudioSpeechRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "AudioSpeech OpenAI ConvAudioSpeechRequest error: %v", err)
		return response, err
	}

	bytes, err := util.HttpPost(ctx, o.baseURL+"/audio/speech", o.header, request, nil, o.proxyURL)
	if err != nil {
		logger.Errorf(ctx, "AudioSpeech OpenAI model: %s, error: %v", o.model, err)
		return response, err
	}

	if response, err = o.ConvAudioSpeechResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "AudioSpeech OpenAI ConvAudioSpeechResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "AudioSpeech OpenAI model: %s finished", o.model)

	return response, nil
}

func (o *OpenAI) AudioTranscriptions(ctx context.Context, request model.AudioRequest) (response model.AudioResponse, err error) {

	logger.Infof(ctx, "AudioTranscriptions OpenAI model: %s start", o.model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "AudioTranscriptions OpenAI model: %s totalTime: %d ms", o.model, response.TotalTime)
	}()

	data, err := o.ConvAudioTranscriptionsRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "AudioTranscriptions OpenAI ConvAudioTranscriptionsRequest error: %v", err)
		return response, err
	}

	bytes, err := util.HttpPost(ctx, o.baseURL+"/audio/transcriptions", o.header, data, nil, o.proxyURL)
	if err != nil {
		logger.Errorf(ctx, "AudioTranscriptions OpenAI model: %s, error: %v", o.model, err)
		return response, err
	}

	if response, err = o.ConvAudioTranscriptionsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "AudioTranscriptions OpenAI ConvAudioTranscriptionsResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "AudioTranscriptions OpenAI model: %s finished", o.model)

	return response, nil
}
