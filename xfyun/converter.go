package xfyun

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/common"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (x *Xfyun) ConvChatCompletionsRequest(ctx context.Context, data any) (request model.ChatCompletionRequest, err error) {

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

	if x.IsSupportSystemRole != nil {
		request.Messages = common.HandleMessages(request.Messages, *x.IsSupportSystemRole)
	} else {
		request.Messages = common.HandleMessages(request.Messages, true)
	}

	return request, nil
}

func (x *Xfyun) ConvChatCompletionsResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	chatCompletionRes := model.XfyunChatCompletionRes{}
	if err = json.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if chatCompletionRes.Header.Code != 0 {
		logger.Errorf(ctx, "ConvChatCompletionsResponse Xfyun model: %s, chatCompletionRes: %s", x.Model, gjson.MustEncodeString(chatCompletionRes))

		err = x.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsResponse Xfyun model: %s, error: %v", x.Model, err)

		return response, err
	}

	response = model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Header.Sid,
		Object:  consts.COMPLETION_OBJECT,
		Created: gtime.Timestamp(),
		Model:   x.Model,
		Choices: []model.ChatCompletionChoice{{
			Index: chatCompletionRes.Payload.Choices.Seq,
			Message: &model.ChatCompletionMessage{
				Role:         chatCompletionRes.Payload.Choices.Text[0].Role,
				FunctionCall: chatCompletionRes.Payload.Choices.Text[0].FunctionCall,
			},
		}},
		Usage: &model.Usage{
			PromptTokens:     chatCompletionRes.Payload.Usage.Text.PromptTokens,
			CompletionTokens: chatCompletionRes.Payload.Usage.Text.CompletionTokens,
			TotalTokens:      chatCompletionRes.Payload.Usage.Text.TotalTokens,
		},
		ResponseBytes: data,
	}

	return response, nil
}

func (x *Xfyun) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {

	chatCompletionRes := model.XfyunChatCompletionRes{}
	if err = json.Unmarshal(data, &chatCompletionRes); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if chatCompletionRes.Header.Code != 0 {
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponse Xfyun model: %s, chatCompletionRes: %s", x.Model, gjson.MustEncodeString(chatCompletionRes))

		err = x.apiErrorHandler(&chatCompletionRes)
		logger.Errorf(ctx, "ConvChatCompletionsStreamResponse Xfyun model: %s, error: %v", x.Model, err)

		return response, err
	}

	response = model.ChatCompletionResponse{
		Id:      consts.COMPLETION_ID_PREFIX + chatCompletionRes.Header.Sid,
		Object:  consts.COMPLETION_STREAM_OBJECT,
		Created: gtime.Timestamp(),
		Model:   x.Model,
		Choices: []model.ChatCompletionChoice{{
			Index: chatCompletionRes.Payload.Choices.Seq,
			Delta: &model.ChatCompletionStreamChoiceDelta{
				Role:         chatCompletionRes.Payload.Choices.Text[0].Role,
				Content:      chatCompletionRes.Payload.Choices.Text[0].Content,
				FunctionCall: chatCompletionRes.Payload.Choices.Text[0].FunctionCall,
			},
		}},
		ResponseBytes: data,
	}

	if chatCompletionRes.Payload.Usage != nil {
		response.Usage = &model.Usage{
			PromptTokens:     chatCompletionRes.Payload.Usage.Text.PromptTokens,
			CompletionTokens: chatCompletionRes.Payload.Usage.Text.CompletionTokens,
			TotalTokens:      chatCompletionRes.Payload.Usage.Text.TotalTokens,
		}
	}

	return response, nil
}

func (x *Xfyun) ConvChatResponsesRequest(ctx context.Context, data []byte) (request model.ChatCompletionRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvChatResponsesResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (response model.ChatCompletionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvImageGenerationsRequest(ctx context.Context, data []byte) (request model.ImageGenerationRequest, err error) {

	if err = json.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	return request, nil
}
func (x *Xfyun) ConvImageGenerationsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvImageEditsResponse(ctx context.Context, data []byte) (response model.ImageResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvAudioSpeechRequest(ctx context.Context, data []byte) (request model.SpeechRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvAudioSpeechResponse(ctx context.Context, data []byte) (response model.SpeechResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (response model.AudioResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (request model.EmbeddingRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (response model.EmbeddingResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvVideoCreateRequest(ctx context.Context, request model.VideoCreateRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvVideoListResponse(ctx context.Context, data []byte) (response model.VideoListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvVideoContentResponse(ctx context.Context, data []byte) (response model.VideoContentResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvVideoJobResponse(ctx context.Context, data []byte) (response model.VideoJobResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvFileUploadRequest(ctx context.Context, request model.FileUploadRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvFileListResponse(ctx context.Context, data []byte) (response model.FileListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvFileContentResponse(ctx context.Context, data []byte) (response model.FileContentResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvFileResponse(ctx context.Context, data []byte) (response model.FileResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvBatchCreateRequest(ctx context.Context, request model.BatchCreateRequest) (data *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvBatchListResponse(ctx context.Context, data []byte) (response model.BatchListResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (x *Xfyun) ConvBatchResponse(ctx context.Context, data []byte) (response model.BatchResponse, err error) {
	//TODO implement me
	panic("implement me")
}
