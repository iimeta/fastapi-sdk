package xai

import (
	"context"
	"slices"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (x *XAI) AudioSpeech(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {

	logger.Infof(ctx, "AudioSpeech xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "AudioSpeech xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	var request any = data
	if !slices.Contains(x.ReqPassthroughParams, "req_data") {
		if request, err = x.ConvAudioSpeechRequest(ctx, data); err != nil {
			logger.Errorf(ctx, "AudioSpeech xAI ConvAudioSpeechRequest error: %v", err)
			return response, err
		}
	}

	if x.Path == "" {
		x.Path = "/audio/speech"
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, x.BaseUrl+x.Path, x.header, request, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "AudioSpeech xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvAudioSpeechResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "AudioSpeech xAI ConvAudioSpeechResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "AudioSpeech xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) AudioTranscriptions(ctx context.Context, request model.AudioRequest) (response model.AudioResponse, err error) {

	logger.Infof(ctx, "AudioTranscriptions xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "AudioTranscriptions xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	data, err := x.ConvAudioTranscriptionsRequest(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "AudioTranscriptions xAI ConvAudioTranscriptionsRequest error: %v", err)
		return response, err
	}

	if x.Path == "" {
		x.Path = "/audio/transcriptions"
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, x.BaseUrl+x.Path, x.header, data, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "AudioTranscriptions xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvAudioTranscriptionsResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "AudioTranscriptions xAI ConvAudioTranscriptionsResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "AudioTranscriptions xAI model: %s finished", x.Model)

	return response, nil
}
