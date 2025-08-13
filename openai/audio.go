package openai

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/go-openai"
)

func (o *OpenAI) AudioSpeech(ctx context.Context, request model.SpeechRequest) (res model.SpeechResponse, err error) {

	logger.Infof(ctx, "AudioSpeech OpenAI model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "AudioSpeech OpenAI model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	response, err := o.client.CreateSpeech(ctx, openai.CreateSpeechRequest{
		Model:          request.Model,
		Input:          request.Input,
		Voice:          request.Voice,
		ResponseFormat: request.ResponseFormat,
		Speed:          request.Speed,
	})

	if err != nil {
		logger.Errorf(ctx, "AudioSpeech OpenAI model: %s, error: %v", request.Model, err)
		return res, o.apiErrorHandler(err)
	}

	logger.Infof(ctx, "AudioSpeech OpenAI model: %s finished", request.Model)

	res = model.SpeechResponse{
		ReadCloser: response.ReadCloser,
	}

	return res, nil
}

func (o *OpenAI) AudioTranscriptions(ctx context.Context, request model.AudioRequest) (res model.AudioResponse, err error) {

	logger.Infof(ctx, "AudioTranscriptions OpenAI model: %s start", request.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "AudioTranscriptions OpenAI model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	response, err := o.client.CreateTranscription(ctx, openai.AudioRequest{
		Model:                  request.Model,
		FilePath:               request.FilePath,
		Reader:                 request.Reader,
		Prompt:                 request.Prompt,
		Temperature:            request.Temperature,
		Language:               request.Language,
		Format:                 request.Format,
		TimestampGranularities: request.TimestampGranularities,
	})

	if err != nil {
		logger.Errorf(ctx, "AudioTranscriptions OpenAI model: %s, error: %v", request.Model, err)
		return res, o.apiErrorHandler(err)
	}

	logger.Infof(ctx, "AudioTranscriptions OpenAI model: %s finished", request.Model)

	res = model.AudioResponse{
		Task:     response.Task,
		Language: response.Language,
		Duration: response.Duration,
		Segments: response.Segments,
		Words:    response.Words,
		Text:     response.Text,
	}

	return res, nil
}
