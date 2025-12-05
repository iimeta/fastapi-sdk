package anthropic

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/model"
)

func (a *Anthropic) ConvChatCompletionsRequestOfficial(ctx context.Context, request model.ChatCompletionRequest) ([]byte, error) {

	chatCompletionReq := model.AnthropicChatCompletionReq{
		Model:         request.Model,
		Messages:      request.Messages,
		MaxTokens:     request.MaxTokens,
		StopSequences: request.Stop,
		Stream:        request.Stream,
		Temperature:   request.Temperature,
		ToolChoice:    request.ToolChoice,
		TopK:          request.TopK,
		TopP:          request.TopP,
		Tools:         request.Tools,
	}

	if chatCompletionReq.Messages[0].Role == consts.ROLE_SYSTEM {
		chatCompletionReq.System = chatCompletionReq.Messages[0].Content
		chatCompletionReq.Messages = chatCompletionReq.Messages[1:]
	}

	if request.User != "" {
		chatCompletionReq.Metadata = &model.Metadata{
			UserId: request.User,
		}
	}

	if chatCompletionReq.MaxTokens == 0 {
		chatCompletionReq.MaxTokens = 4096
	}

	for _, message := range chatCompletionReq.Messages {

		if contents, ok := message.Content.([]any); ok {

			for _, value := range contents {

				if content, ok := value.(map[string]any); ok {

					if content["type"] == "image_url" {

						if imageUrl, ok := content["image_url"].(map[string]any); ok {

							mimeType, data := common.GetMime(gconv.String(imageUrl["url"]))

							content["source"] = model.Source{
								Type:      "base64",
								MediaType: mimeType,
								Data:      data,
							}

							content["type"] = "image"
							delete(content, "image_url")
						}
					}
				}
			}
		}
	}

	if a.isGcp {
		chatCompletionReq.Model = ""
		chatCompletionReq.AnthropicVersion = "vertex-2023-10-16"
	}

	return gjson.MustEncode(chatCompletionReq), nil
}

func (a *Anthropic) ConvChatCompletionsResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Anthropic) ConvChatCompletionsStreamResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
