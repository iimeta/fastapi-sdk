package google

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/iimeta/fastapi-sdk/common"
	"github.com/iimeta/fastapi-sdk/consts"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
	"github.com/iimeta/go-openai"
	"io"
)

func (c *Client) ChatCompletion(ctx context.Context, request model.ChatCompletionRequest) (res model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletion Google model: %s start", request.Model)

	now := gtime.Now().UnixMilli()
	defer func() {
		res.TotalTime = gtime.Now().UnixMilli() - now
		logger.Infof(ctx, "ChatCompletion Google model: %s totalTime: %d ms", request.Model, res.TotalTime)
	}()

	var messages []model.ChatCompletionMessage
	if c.isSupportSystemRole != nil {
		messages = common.HandleMessages(request.Messages, *c.isSupportSystemRole)
	} else {
		messages = common.HandleMessages(request.Messages, false)
	}

	contents := make([]model.Content, 0)
	for _, message := range messages {

		role := message.Role

		if role == consts.ROLE_ASSISTANT {
			role = consts.ROLE_MODEL
		}

		parts := make([]model.Part, 0)

		if contents, ok := message.Content.([]interface{}); ok {

			for _, value := range contents {

				if content, ok := value.(map[string]interface{}); ok {

					if content["type"] == "image_url" {

						if imageUrl, ok := content["image_url"].(map[string]interface{}); ok {

							url := gconv.String(imageUrl["url"])

							if gstr.HasPrefix(url, "data:image/") {
								base64 := gstr.Split(url, "base64,")
								if len(base64) > 1 {
									// data:image/jpeg;base64,
									mimeType := fmt.Sprintf("image/%s", gstr.Split(base64[0][11:], ";")[0])
									parts = append(parts, model.Part{
										InlineData: &model.InlineData{
											MimeType: mimeType,
											Data:     base64[1],
										},
									})
								} else {
									parts = append(parts, model.Part{
										InlineData: &model.InlineData{
											MimeType: "image/jpeg",
											Data:     base64[0],
										},
									})
								}
							} else {
								parts = append(parts, model.Part{
									InlineData: &model.InlineData{
										MimeType: "image/jpeg",
										Data:     url,
									},
								})
							}
						}

					} else if content["type"] == "video_url" {
						if videoUrl, ok := content["video_url"].(map[string]interface{}); ok {

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

	chatCompletionRes := new(model.GoogleChatCompletionRes)
	err = util.HttpPost(ctx, fmt.Sprintf("%s:generateContent?key=%s", c.baseURL+c.path, c.key), nil, chatCompletionReq, &chatCompletionRes, c.proxyURL)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletion Google model: %s, error: %v", request.Model, err)
		return
	}

	if chatCompletionRes.Error.Code != 0 || chatCompletionRes.Candidates[0].FinishReason != "STOP" {
		logger.Errorf(ctx, "ChatCompletion Google model: %s, chatCompletionRes: %s", request.Model, gjson.MustEncodeString(chatCompletionRes))

		err = c.apiErrorHandler(chatCompletionRes)
		logger.Errorf(ctx, "ChatCompletion Google model: %s, error: %v", request.Model, err)

		return
	}

	res = model.ChatCompletionResponse{
		ID:      consts.COMPLETION_ID_PREFIX + grand.S(29),
		Object:  consts.COMPLETION_OBJECT,
		Created: gtime.Now().Unix(),
		Model:   request.Model,
		Choices: []model.ChatCompletionChoice{{
			Message: &openai.ChatCompletionMessage{
				Role:    consts.ROLE_ASSISTANT,
				Content: chatCompletionRes.Candidates[0].Content.Parts[0].Text,
			},
			FinishReason: openai.FinishReasonStop,
		}},
		Usage: &model.Usage{
			PromptTokens:     chatCompletionRes.UsageMetadata.PromptTokenCount,
			CompletionTokens: chatCompletionRes.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      chatCompletionRes.UsageMetadata.TotalTokenCount,
		},
	}

	return res, nil
}

func (c *Client) ChatCompletionStream(ctx context.Context, request model.ChatCompletionRequest) (responseChan chan *model.ChatCompletionResponse, err error) {

	logger.Infof(ctx, "ChatCompletionStream Google model: %s start", request.Model)

	now := gtime.Now().UnixMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ChatCompletionStream Google model: %s totalTime: %d ms", request.Model, gtime.Now().UnixMilli()-now)
		}
	}()

	var messages []model.ChatCompletionMessage
	if c.isSupportSystemRole != nil {
		messages = common.HandleMessages(request.Messages, *c.isSupportSystemRole)
	} else {
		messages = common.HandleMessages(request.Messages, false)
	}

	contents := make([]model.Content, 0)
	for _, message := range messages {

		role := message.Role

		if role == consts.ROLE_ASSISTANT {
			role = consts.ROLE_MODEL
		}

		parts := make([]model.Part, 0)

		if contents, ok := message.Content.([]interface{}); ok {

			for _, value := range contents {

				if content, ok := value.(map[string]interface{}); ok {

					if content["type"] == "image_url" {

						if imageUrl, ok := content["image_url"].(map[string]interface{}); ok {

							url := gconv.String(imageUrl["url"])

							if gstr.HasPrefix(url, "data:image/") {
								base64 := gstr.Split(url, "base64,")
								if len(base64) > 1 {
									// data:image/jpeg;base64,
									mimeType := fmt.Sprintf("image/%s", gstr.Split(base64[0][11:], ";")[0])
									parts = append(parts, model.Part{
										InlineData: &model.InlineData{
											MimeType: mimeType,
											Data:     base64[1],
										},
									})
								} else {
									parts = append(parts, model.Part{
										InlineData: &model.InlineData{
											MimeType: "image/jpeg",
											Data:     base64[0],
										},
									})
								}
							} else {
								parts = append(parts, model.Part{
									InlineData: &model.InlineData{
										MimeType: "image/jpeg",
										Data:     url,
									},
								})
							}
						}

					} else if content["type"] == "video_url" {
						if videoUrl, ok := content["video_url"].(map[string]interface{}); ok {

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

	stream, err := util.SSEClient(ctx, fmt.Sprintf("%s:streamGenerateContent?alt=sse&key=%s", c.baseURL+c.path, c.key), nil, chatCompletionReq, c.proxyURL, c.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ChatCompletionStream Google model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	duration := gtime.Now().UnixMilli()

	responseChan = make(chan *model.ChatCompletionResponse)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.Now().UnixMilli()
			logger.Infof(ctx, "ChatCompletionStream Google model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", request.Model, duration-now, end-duration, end-now)

			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream Google model: %s, stream.Close error: %v", request.Model, err)
			}
		}()

		var (
			usage   *model.Usage
			created = gtime.Now().Unix()
			id      = consts.COMPLETION_ID_PREFIX + grand.S(29)
		)

		for {

			streamResponse, err := stream.Recv()
			if err != nil && !errors.Is(err, io.EOF) {

				if !errors.Is(err, context.Canceled) {
					logger.Errorf(ctx, "ChatCompletionStream Google model: %s, error: %v", request.Model, err)
				}

				end := gtime.Now().UnixMilli()
				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     err,
				}

				return
			}

			if errors.Is(err, io.EOF) {
				logger.Infof(ctx, "ChatCompletionStream Google model: %s finished", request.Model)

				end := gtime.Now().UnixMilli()
				responseChan <- &model.ChatCompletionResponse{
					ID:      id,
					Object:  consts.COMPLETION_STREAM_OBJECT,
					Created: created,
					Model:   request.Model,
					Choices: []model.ChatCompletionChoice{{
						Delta:        &openai.ChatCompletionStreamChoiceDelta{},
						FinishReason: openai.FinishReasonStop,
					}},
					Usage:     usage,
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
				}

				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     io.EOF,
				}

				return
			}

			chatCompletionRes := new(model.GoogleChatCompletionRes)
			if err := gjson.Unmarshal(streamResponse, &chatCompletionRes); err != nil {
				logger.Errorf(ctx, "ChatCompletionStream Google model: %s, streamResponse: %s, error: %v", request.Model, streamResponse, err)

				end := gtime.Now().UnixMilli()
				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     errors.New(fmt.Sprintf("streamResponse: %s, error: %v", streamResponse, err)),
				}

				return
			}

			if chatCompletionRes.Error.Code != 0 {
				logger.Errorf(ctx, "ChatCompletionStream Google model: %s, chatCompletionRes: %s", request.Model, gjson.MustEncodeString(chatCompletionRes))

				err = c.apiErrorHandler(chatCompletionRes)
				logger.Errorf(ctx, "ChatCompletionStream Google model: %s, error: %v", request.Model, err)

				end := gtime.Now().UnixMilli()
				responseChan <- &model.ChatCompletionResponse{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Error:     err,
				}

				return
			}

			if chatCompletionRes.UsageMetadata != nil {
				usage = &model.Usage{
					PromptTokens:     chatCompletionRes.UsageMetadata.PromptTokenCount,
					CompletionTokens: chatCompletionRes.UsageMetadata.CandidatesTokenCount,
					TotalTokens:      chatCompletionRes.UsageMetadata.TotalTokenCount,
				}
			}

			response := &model.ChatCompletionResponse{
				ID:      id,
				Object:  consts.COMPLETION_STREAM_OBJECT,
				Created: created,
				Model:   request.Model,
				Choices: []model.ChatCompletionChoice{{
					Index: chatCompletionRes.Candidates[0].Index,
					Delta: &openai.ChatCompletionStreamChoiceDelta{
						Role:    consts.ROLE_ASSISTANT,
						Content: chatCompletionRes.Candidates[0].Content.Parts[0].Text,
					},
				}},
				Usage:    usage,
				ConnTime: duration - now,
			}

			end := gtime.Now().UnixMilli()
			response.Duration = end - duration
			response.TotalTime = end - now

			responseChan <- response
		}
	}, nil); err != nil {
		logger.Errorf(ctx, "ChatCompletionStream Google model: %s, error: %v", request.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}
