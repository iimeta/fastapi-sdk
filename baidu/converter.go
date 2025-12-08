package baidu

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func (b *Baidu) ConvChatCompletionsRequest(ctx context.Context, data any) (request model.ChatCompletionRequest, err error) {

	if v, ok := data.(model.ChatCompletionRequest); ok {
		request = v
	} else if v, ok := data.([]byte); ok {
		if err = json.Unmarshal(v, &request); err != nil {
			logger.Error(ctx, err)
			return request, err
		}
	} else {
		if err = json.Unmarshal(gjson.MustEncode(data), &request); err != nil {
			logger.Error(ctx, err)
			return request, err
		}
	}

	if b.IsSupportSystemRole != nil {
		request.Messages = common.HandleMessages(request.Messages, *b.IsSupportSystemRole)
	} else {
		request.Messages = common.HandleMessages(request.Messages, true)
	}

	if len(request.Messages) == 1 && request.Messages[0].Role == consts.ROLE_SYSTEM {
		request.Messages[0].Role = consts.ROLE_USER
	}

	return request, nil
}

func (b *Baidu) ConvChatCompletionsResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	chatCompletionRes := model.BaiduChatCompletionRes{}
	if err = json.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if chatCompletionRes.ErrorCode != 0 {
		logger.Errorf(ctx, "ConvChatCompletionsResponse Baidu model: %s, chatCompletionRes: %s", b.Model, gjson.MustEncodeString(chatCompletionRes))

		err = b.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsResponse Baidu model: %s, error: %v", b.Model, err)

		return response, err
	}

	response = model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Id,
		Object:  consts.COMPLETION_OBJECT,
		Created: chatCompletionRes.Created,
		Model:   b.Model,
		Choices: []model.ChatCompletionChoice{{
			Message: &model.ChatCompletionMessage{
				Role:    consts.ROLE_ASSISTANT,
				Content: chatCompletionRes.Result,
			},
		}},
		Usage:         chatCompletionRes.Usage,
		ResponseBytes: data,
	}

	return response, nil
}

func (b *Baidu) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	chatCompletionRes := model.BaiduChatCompletionRes{}
	if err = json.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if chatCompletionRes.ErrorCode != 0 {
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponse Baidu model: %s, chatCompletionRes: %s", b.Model, gjson.MustEncodeString(chatCompletionRes))

		err = b.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponse Baidu model: %s, error: %v", b.Model, err)

		return response, err
	}

	response = model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Id,
		Object:  consts.COMPLETION_STREAM_OBJECT,
		Created: chatCompletionRes.Created,
		Model:   b.Model,
		Choices: []model.ChatCompletionChoice{{
			Index: chatCompletionRes.SentenceId,
			Delta: &model.ChatCompletionStreamChoiceDelta{
				Role:    consts.ROLE_ASSISTANT,
				Content: chatCompletionRes.Result,
			},
		}},
		Usage:         chatCompletionRes.Usage,
		ResponseBytes: data,
	}

	return response, nil
}

func (b *Baidu) ConvChatResponsesRequest(ctx context.Context, data []byte) (request model.ChatCompletionRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvChatResponsesResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvImageGenerationsRequest(ctx context.Context, data []byte) (model.ImageGenerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvImageGenerationsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (*bytes.Buffer, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvImageEditsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvAudioSpeechRequest(ctx context.Context, data []byte) (model.SpeechRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvAudioSpeechResponse(ctx context.Context, data []byte) (model.SpeechResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (*bytes.Buffer, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (model.AudioResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (model.EmbeddingRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (model.EmbeddingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvVideoCreateRequest(ctx context.Context, request model.VideoCreateRequest) (*bytes.Buffer, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvVideoListResponse(ctx context.Context, data []byte) (model.VideoListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvVideoContentResponse(ctx context.Context, data []byte) (model.VideoContentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Baidu) ConvVideoJobResponse(ctx context.Context, data []byte) (model.VideoJobResponse, error) {
	//TODO implement me
	panic("implement me")
}
