package general

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (g *General) AudioSpeech(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {

	logger.Infof(ctx, "AudioSpeech General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "AudioSpeech General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	request, err := g.ConvAudioSpeechRequest(ctx, data)
	if err != nil {
		logger.Errorf(ctx, "AudioSpeech General ConvAudioSpeechRequest error: %v", err)
		return response, err
	}

	bytes, err := util.HttpPost(ctx, g.BaseUrl+g.Path, g.header, request, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "AudioSpeech General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvAudioSpeechResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "AudioSpeech General ConvAudioSpeechResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "AudioSpeech General model: %s finished", g.Model)

	return response, nil
}

func (g *General) AudioTranscriptions(ctx context.Context, request model.AudioRequest) (response model.AudioResponse, err error) {

	logger.Infof(ctx, "AudioTranscriptions General model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "AudioTranscriptions General model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	data, err := g.ConvAudioTranscriptionsRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "AudioTranscriptions General ConvAudioTranscriptionsRequest error: %v", err)
		return response, err
	}

	bytes, err := util.HttpPost(ctx, g.BaseUrl+g.Path, g.header, data, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "AudioTranscriptions General model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvAudioTranscriptionsResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "AudioTranscriptions General ConvAudioTranscriptionsResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "AudioTranscriptions General model: %s finished", g.Model)

	return response, nil
}
