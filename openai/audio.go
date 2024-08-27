package openai

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/go-openai"
)

func (c *Client) Speech(ctx context.Context, request model.SpeechRequest) (res model.SpeechResponse, err error) {

	logger.Infof(ctx, "Speech OpenAI model: %s start", request.Model)

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "Speech OpenAI model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	response, err := c.client.CreateSpeech(ctx, openai.CreateSpeechRequest{
		Model:          request.Model,
		Input:          request.Input,
		Voice:          request.Voice,
		ResponseFormat: request.ResponseFormat,
		Speed:          request.Speed,
	})

	if err != nil {
		logger.Errorf(ctx, "Speech OpenAI model: %s, error: %v", request.Model, err)
		return res, c.apiErrorHandler(err)
	}

	logger.Infof(ctx, "Speech OpenAI model: %s finished", request.Model)

	res = model.SpeechResponse{
		ReadCloser: response.ReadCloser,
	}

	return res, nil
}

func (c *Client) Transcription(ctx context.Context, request model.AudioRequest) (res model.AudioResponse, err error) {
	//TODO implement me
	panic("implement me")
}
