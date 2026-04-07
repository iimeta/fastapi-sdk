package volcengine

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (v *VolcEngine) ConvChatCompletionsRequestOfficial(ctx context.Context, request model.ChatCompletionRequest) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvChatCompletionsResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvChatCompletionsStreamResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvImageGenerationsRequestOfficial(ctx context.Context, request model.ImageGenerationRequest) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvImageGenerationsResponseOfficial(ctx context.Context, response model.ImageResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvImageEditsRequestOfficial(ctx context.Context, request model.ImageEditRequest) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvImageEditsResponseOfficial(ctx context.Context, response model.ImageResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VolcEngine) ConvVideoJobResponseOfficial(ctx context.Context, response model.VideoJobResponse) (*model.VolcVideoTaskRes, error) {

	res := &model.VolcVideoTaskRes{
		Id:        response.Id,
		Model:     response.Model,
		Status:    ConvToVolcStatus(response.Status),
		CreatedAt: response.CreatedAt,
	}

	if response.Seconds != "" {
		dur := gconv.Int(response.Seconds)
		res.Duration = &dur
	}

	if response.Size != "" {
		res.Ratio = response.Size
	}

	if response.CompletedAt != nil {
		res.UpdatedAt = *response.CompletedAt
	}

	if response.VideoUrl != "" {
		res.Content = &model.VolcVideoContentResult{
			VideoUrl: response.VideoUrl,
		}
	}

	if response.Usage != nil {
		res.Usage = &model.VolcVideoUsage{
			CompletionTokens: response.Usage.CompletionTokens,
			TotalTokens:      response.Usage.TotalTokens,
		}
	}

	if response.Error != nil {
		res.Error = &model.VolcVideoError{
			Code:    response.Error.Code,
			Message: response.Error.Message,
		}
	}

	return res, nil
}

func (v *VolcEngine) ConvVideoListResponseOfficial(ctx context.Context, response model.VideoListResponse) ([]byte, error) {

	listRes := model.VolcVideoListRes{
		Total: len(response.Data),
	}

	for _, item := range response.Data {
		volcItem, err := v.ConvVideoJobResponseOfficial(ctx, item)
		if err != nil {
			logger.Errorf(ctx, "ConvVideoListResponseOfficial VolcEngine error: %v", err)
			continue
		}
		listRes.Items = append(listRes.Items, volcItem)
	}

	if listRes.Items == nil {
		listRes.Items = make([]*model.VolcVideoTaskRes, 0)
	}

	data, err := json.Marshal(listRes)
	if err != nil {
		return nil, fmt.Errorf("ConvVideoListResponseOfficial json.Marshal error: %v", err)
	}

	return data, nil
}

func ConvToVolcStatus(status string) string {
	switch status {
	case "queued":
		return "queued"
	case "in_progress":
		return "running"
	case "completed":
		return "succeeded"
	case "failed":
		return "failed"
	case "deleted":
		return "cancelled"
	case "expired":
		return "expired"
	default:
		return status
	}
}
