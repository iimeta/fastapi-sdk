package google

import (
	"context"
	"encoding/base64"
	"io"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/v2/common"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
)

func (g *Google) ConvChatCompletionsRequestOfficial(ctx context.Context, request model.ChatCompletionRequest) ([]byte, error) {

	contents := make([]model.Content, 0)
	for _, message := range request.Messages {

		role := message.Role

		if role == consts.ROLE_ASSISTANT {
			role = consts.ROLE_MODEL
		}

		parts := make([]model.Part, 0)

		if contents, ok := message.Content.([]any); ok {

			for _, value := range contents {

				if content, ok := value.(map[string]any); ok {

					if content["type"] == "image_url" {

						if imageUrl, ok := content["image_url"].(map[string]any); ok {

							mimeType, data := common.GetMime(gconv.String(imageUrl["url"]))

							parts = append(parts, model.Part{
								InlineData: &model.InlineData{
									MimeType: mimeType,
									Data:     data,
								},
							})
						}

					} else if content["type"] == "video_url" {
						if videoUrl, ok := content["video_url"].(map[string]any); ok {

							url := gconv.String(videoUrl["url"])
							format := gconv.String(videoUrl["format"])

							parts = append(parts, model.Part{
								FileData: &model.FileData{
									MimeType: "video/" + format,
									FileUri:  url,
								},
							})
						}
					} else {
						parts = append(parts, model.Part{
							Text: gconv.String(content["text"]),
						})
					}
				}
			}

		} else {
			parts = append(parts, model.Part{
				Text: gconv.String(message.Content),
			})
		}

		contents = append(contents, model.Content{
			Role:  role,
			Parts: parts,
		})
	}

	chatCompletionReq := model.GoogleChatCompletionReq{
		Contents: contents,
		GenerationConfig: model.GenerationConfig{
			MaxOutputTokens: request.MaxTokens,
			Temperature:     request.Temperature,
			TopP:            request.TopP,
		},
	}

	if request.Tools != nil {
		if tools, ok := request.Tools.([]any); ok {

			var functionDeclarations []any

			for _, value := range tools {
				if tool, ok := value.(map[string]any); ok {
					functionDeclarations = append(functionDeclarations, tool["function"])
				}
			}

			chatCompletionReq.Tools = map[string]any{
				"functionDeclarations": functionDeclarations,
			}

		} else {
			chatCompletionReq.Tools = request.Tools
		}
	}

	return gjson.MustEncode(chatCompletionReq), nil
}

func (g *Google) ConvChatCompletionsResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvChatCompletionsStreamResponseOfficial(ctx context.Context, response model.ChatCompletionResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvImageGenerationsRequestOfficial(ctx context.Context, request model.ImageGenerationRequest) ([]byte, error) {

	parts := make([]model.Part, 0)

	parts = append(parts, model.Part{
		Text: gconv.String(request.Prompt),
	})

	contents := make([]model.Content, 0)

	contents = append(contents, model.Content{
		Role:  consts.ROLE_USER,
		Parts: parts,
	})

	chatCompletionReq := model.GoogleChatCompletionReq{
		Contents: contents,
		GenerationConfig: model.GenerationConfig{
			ImageConfig: &model.ImageConfig{
				AspectRatio: request.AspectRatio,
				ImageSize:   request.Quality,
			},
		},
	}

	return gjson.MustEncode(chatCompletionReq), nil
}

func (g *Google) ConvImageGenerationsResponseOfficial(ctx context.Context, response model.ImageResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Google) ConvImageEditsRequestOfficial(ctx context.Context, request model.ImageEditRequest) ([]byte, error) {

	parts := make([]model.Part, 0)

	parts = append(parts, model.Part{
		Text: gconv.String(request.Prompt),
	})

	for _, image := range request.Image {

		file, err := image.Open()
		if err != nil {
			logger.Error(ctx, err)
			return nil, err
		}

		fileBytes, err := io.ReadAll(file)

		if err := file.Close(); err != nil {
			logger.Error(ctx, err)
		}

		if err != nil {
			logger.Error(ctx, err)
			return nil, err
		}

		parts = append(parts, model.Part{
			InlineData: &model.InlineData{
				MimeType: image.Header.Get("Content-Type"),
				Data:     base64.StdEncoding.EncodeToString(fileBytes),
			},
		})
	}

	contents := make([]model.Content, 0)

	contents = append(contents, model.Content{
		Role:  consts.ROLE_USER,
		Parts: parts,
	})

	chatCompletionReq := model.GoogleChatCompletionReq{
		Contents: contents,
		GenerationConfig: model.GenerationConfig{
			ImageConfig: &model.ImageConfig{
				AspectRatio: request.AspectRatio,
				ImageSize:   request.Quality,
			},
		},
	}

	return gjson.MustEncode(chatCompletionReq), nil
}

func (g *Google) ConvImageEditsResponseOfficial(ctx context.Context, response model.ImageResponse) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
