package volcengine

import (
	"bytes"
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
)

func (v *VolcEngine) ConvChatCompletionsRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {

	request := model.ChatCompletionRequest{}
	if err := gjson.Unmarshal(data, &request); err != nil {
		logger.Error(ctx, err)
		return request, err
	}

	if v.isSupportSystemRole != nil {
		request.Messages = common.HandleMessages(request.Messages, *v.isSupportSystemRole)
	} else {
		request.Messages = common.HandleMessages(request.Messages, true)
	}

	return request, nil
}

func (v *VolcEngine) ConvChatCompletionsResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	response := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if !gstr.HasPrefix(response.Id, consts.COMPLETION_ID_PREFIX) {
		response.Id = consts.COMPLETION_ID_PREFIX + response.Id
	}

	return response, nil
}

func (v *VolcEngine) ConvChatCompletionsStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {

	response := model.ChatCompletionResponse{}
	if err := gjson.Unmarshal(data, &response); err != nil {
		logger.Error(ctx, err)
		return response, err
	}

	if !gstr.HasPrefix(response.Id, consts.COMPLETION_ID_PREFIX) {
		response.Id = consts.COMPLETION_ID_PREFIX + response.Id
	}

	return response, nil
}

func (v *VolcEngine) ConvChatResponsesRequest(ctx context.Context, data []byte) (model.ChatCompletionRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvChatResponsesResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvChatResponsesStreamResponse(ctx context.Context, data []byte) (model.ChatCompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvImageGenerationsRequest(ctx context.Context, data []byte) (model.ImageGenerationRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvImageGenerationsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvImageEditsRequest(ctx context.Context, request model.ImageEditRequest) (*bytes.Buffer, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvImageEditsResponse(ctx context.Context, data []byte) (model.ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvAudioSpeechRequest(ctx context.Context, data []byte) (model.SpeechRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvAudioSpeechResponse(ctx context.Context, data []byte) (model.SpeechResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvAudioTranscriptionsRequest(ctx context.Context, request model.AudioRequest) (*bytes.Buffer, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvAudioTranscriptionsResponse(ctx context.Context, data []byte) (model.AudioResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvTextEmbeddingsRequest(ctx context.Context, data []byte) (model.EmbeddingRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvTextEmbeddingsResponse(ctx context.Context, data []byte) (model.EmbeddingResponse, error) {
	//TODO implement me
	panic("implement me")
}
