package xai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (x *XAI) Responses(ctx context.Context, data []byte) (res model.OpenAIResponsesRes, err error) {

	logger.Infof(ctx, "Responses xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "Responses xAI model: %s totalTime: %d ms", x.Model, res.TotalTime)
	}()

	if x.Path == "" {
		x.Path = "/responses"
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, x.BaseUrl+x.Path, x.header, data, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "Responses xAI model: %s, error: %v", x.Model, err)
		return res, err
	}

	var official model.XAIResponsesRes
	if err = json.Unmarshal(responseBytes, &official); err != nil {
		logger.Errorf(ctx, "Responses xAI model: %s, response: %s, error: %v", x.Model, responseBytes, err)
		return res, err
	}
	res = convXAIResponsesToOpenAI(official, responseBytes)
	res.ResponseHeaders = responseHeader

	if res.Error != nil {
		logger.Errorf(ctx, "Responses xAI model: %s, responsesRes: %s", x.Model, gjson.MustEncodeString(res))

		err = x.responsesErrorHandler(res.Error)
		logger.Errorf(ctx, "Responses xAI model: %s, error: %v", x.Model, err)

		return res, err
	}

	return res, nil
}

func (x *XAI) ResponsesStream(ctx context.Context, data []byte) (responseChan chan *model.OpenAIResponsesStreamRes, err error) {

	logger.Infof(ctx, "ResponsesStream xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ResponsesStream xAI model: %s totalTime: %d ms", x.Model, gtime.TimestampMilli()-now)
		}
	}()

	if x.IsSupportStream != nil && !*x.IsSupportStream {
		return x.ResponsesStreamToNonStream(ctx, data)
	}

	if x.Path == "" {
		x.Path = "/responses"
	}

	stream, err := util.SSEClient(ctx, x.BaseUrl+x.Path, x.header, data, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ResponsesStream xAI model: %s, error: %v", x.Model, err)
		return responseChan, err
	}

	streamResponseHeaders := stream.Response.Header

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.OpenAIResponsesStreamRes)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ResponsesStream xAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", x.Model, duration-now, end-duration, end-now)

			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ResponsesStream xAI model: %s, stream.Close error: %v", x.Model, err)
			}
		}()

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ResponsesStream xAI model: %s finished", x.Model)
				} else {
					logger.Errorf(ctx, "ResponsesStream xAI model: %s, error: %v", x.Model, err)
				}

				end := gtime.TimestampMilli()
				responseChan <- &model.OpenAIResponsesStreamRes{
					SSEEvent:  stream.Event(),
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Err:       err,
				}

				return
			}

			var official model.XAIResponsesStreamRes
			if err := json.Unmarshal(responseBytes, &official); err != nil {
				logger.Errorf(ctx, "ResponsesStream xAI model: %s, response: %s, error: %v", x.Model, responseBytes, err)

				end := gtime.TimestampMilli()
				responseChan <- &model.OpenAIResponsesStreamRes{
					SSEEvent:  stream.Event(),
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Err:       errors.New(fmt.Sprintf("response: %s, error: %v", responseBytes, err)),
				}

				return
			}

			converted := convXAIResponsesStreamToOpenAI(official, responseBytes)
			if apiErr := streamResponsesAPIError(converted, responseBytes); apiErr != nil {
				logger.Errorf(ctx, "ResponsesStream xAI model: %s, responsesRes: %s", x.Model, gjson.MustEncodeString(converted))

				err = x.responsesErrorHandler(apiErr)
				logger.Errorf(ctx, "ResponsesStream xAI model: %s, error: %v", x.Model, err)

				end := gtime.TimestampMilli()
				responseChan <- &model.OpenAIResponsesStreamRes{
					SSEEvent:      stream.Event(),
					ResponseBytes: responseBytes,
					ConnTime:      duration - now,
					Duration:      end - duration,
					TotalTime:     end - now,
					Err:           err,
				}

				return
			}

			end := gtime.TimestampMilli()
			converted.SSEEvent = stream.Event()
			converted.ResponseHeaders = streamResponseHeaders
			converted.ConnTime = duration - now
			converted.Duration = end - duration
			converted.TotalTime = end - now

			responseChan <- &converted
		}
	}, nil); err != nil {
		logger.Errorf(ctx, "ResponsesStream xAI model: %s, error: %v", x.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (x *XAI) ResponsesCompact(ctx context.Context, data []byte) (res model.OpenAIResponsesRes, err error) {

	logger.Infof(ctx, "ResponsesCompact xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "ResponsesCompact xAI model: %s totalTime: %d ms", x.Model, res.TotalTime)
	}()

	if x.Path == "" {
		x.Path = "/responses/compact"
	}

	compactBytes, compactResponseHeader, err := util.HttpPost(ctx, x.BaseUrl+x.Path, x.header, data, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ResponsesCompact xAI model: %s, error: %v", x.Model, err)
		return res, err
	}

	var official model.XAIResponsesCompactRes
	if err = json.Unmarshal(compactBytes, &official); err != nil {
		logger.Errorf(ctx, "ResponsesCompact xAI model: %s, response: %s, error: %v", x.Model, compactBytes, err)
		return res, err
	}
	res = convXAIResponsesCompactToOpenAI(official, compactBytes)
	res.ResponseHeaders = compactResponseHeader

	if res.Error != nil {
		logger.Errorf(ctx, "ResponsesCompact xAI model: %s, responsesRes: %s", x.Model, gjson.MustEncodeString(res))

		err = x.responsesErrorHandler(res.Error)
		logger.Errorf(ctx, "ResponsesCompact xAI model: %s, error: %v", x.Model, err)

		return res, err
	}

	return res, nil
}

func (x *XAI) ResponsesStreamToNonStream(ctx context.Context, data []byte) (responseChan chan *model.OpenAIResponsesStreamRes, err error) {

	logger.Infof(ctx, "ResponsesStreamToNonStream xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	duration := now

	responseChan = make(chan *model.OpenAIResponsesStreamRes)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ResponsesStreamToNonStream xAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", x.Model, duration-now, end-duration, end-now)
		}()

		request := make(map[string]any)
		if err = json.Unmarshal(data, &request); err != nil {
			logger.Errorf(ctx, "ResponsesStreamToNonStream xAI model: %s, data: %s, error: %v", x.Model, data, err)

			end := gtime.TimestampMilli()
			responseChan <- &model.OpenAIResponsesStreamRes{
				ConnTime:  gtime.TimestampMilli() - now,
				Duration:  end - gtime.TimestampMilli(),
				TotalTime: end - now,
				Err:       err,
			}
		}

		request["stream"] = false

		responses, err := x.Responses(ctx, gjson.MustEncode(request))
		if err != nil {

			if errors.Is(err, io.EOF) {
				logger.Infof(ctx, "ResponsesStreamToNonStream xAI model: %s finished", x.Model)
			} else {
				logger.Errorf(ctx, "ResponsesStreamToNonStream xAI model: %s, error: %v", x.Model, err)
			}

			end := gtime.TimestampMilli()
			responseChan <- &model.OpenAIResponsesStreamRes{
				ConnTime:  gtime.TimestampMilli() - now,
				Duration:  end - gtime.TimestampMilli(),
				TotalTime: end - now,
				Err:       err,
			}

			return
		}

		duration = gtime.TimestampMilli()

		responsesRes := model.OpenAIResponsesRes{}
		if err := json.Unmarshal(responses.ResponseBytes, &responsesRes); err != nil {
			logger.Errorf(ctx, "ResponsesStreamToNonStream xAI model: %s, responses: %s, error: %v", x.Model, responses.ResponseBytes, err)

			end := gtime.TimestampMilli()
			responseChan <- &model.OpenAIResponsesStreamRes{
				ConnTime:  duration - now,
				Duration:  end - duration,
				TotalTime: end - now,
				Err:       errors.New(fmt.Sprintf("response: %s, error: %v", responses.ResponseBytes, err)),
			}

			return
		}

		if responsesRes.Error != nil {
			logger.Errorf(ctx, "ResponsesStreamToNonStream xAI model: %s, responsesRes: %s", x.Model, gjson.MustEncodeString(responsesRes))

			err = x.responsesErrorHandler(responsesRes.Error)
			logger.Errorf(ctx, "ResponsesStreamToNonStream xAI model: %s, error: %v", x.Model, err)

			end := gtime.TimestampMilli()
			responseChan <- &model.OpenAIResponsesStreamRes{
				ConnTime:  duration - now,
				Duration:  end - duration,
				TotalTime: end - now,
				Err:       err,
			}

			return
		}

		response := &model.OpenAIResponsesStreamRes{
			Type: "response.created",
			Response: model.OpenAIResponsesResponse{
				Id:                 responses.Id,
				Object:             responses.Object,
				CreatedAt:          responses.CreatedAt,
				Status:             "in_progress",
				Background:         responses.Background,
				IncompleteDetails:  responses.IncompleteDetails,
				Instructions:       responses.Instructions,
				MaxOutputTokens:    responses.MaxOutputTokens,
				Model:              responses.Model,
				ParallelToolCalls:  responses.ParallelToolCalls,
				PreviousResponseId: responses.PreviousResponseId,
				Reasoning:          responses.Reasoning,
				ServiceTier:        responses.ServiceTier,
				Store:              responses.Store,
				Temperature:        responses.Temperature,
				Text:               responses.Text,
				ToolChoice:         responses.ToolChoice,
				Tools:              responses.Tools,
				TopP:               responses.TopP,
				Truncation:         responses.Truncation,
				User:               responses.User,
				Metadata:           responses.Metadata,
			},
		}

		response.ResponseBytes = gjson.MustEncode(response)

		end := gtime.TimestampMilli()
		response.ConnTime = duration - now
		response.Duration = end - duration
		response.TotalTime = end - now

		responseChan <- response

		delta := ""
		for _, output := range responsesRes.Output {
			if output.Status == "completed" {
				if output.Type == "function_call" {
					delta += gconv.String(output.Arguments)
				} else if len(output.Content) > 0 {
					delta += output.Content[0].Text
				}
			}
		}

		response = &model.OpenAIResponsesStreamRes{
			Type:           "response.output_text.delta",
			SequenceNumber: 1,
			ItemId:         responses.Id,
			OutputIndex:    1,
			Delta:          delta,
		}

		response.ResponseBytes = gjson.MustEncode(response)

		end = gtime.TimestampMilli()
		response.ConnTime = duration - now
		response.Duration = end - duration
		response.TotalTime = end - now

		responseChan <- response

		response = &model.OpenAIResponsesStreamRes{
			Type:           "response.completed",
			SequenceNumber: 2,
			Response: model.OpenAIResponsesResponse{
				Id:                 responses.Id,
				Object:             responses.Object,
				CreatedAt:          responses.CreatedAt,
				Status:             responses.Status,
				Background:         responses.Background,
				IncompleteDetails:  responses.IncompleteDetails,
				Instructions:       responses.Instructions,
				MaxOutputTokens:    responses.MaxOutputTokens,
				Model:              responses.Model,
				Output:             responses.Output,
				ParallelToolCalls:  responses.ParallelToolCalls,
				PreviousResponseId: responses.PreviousResponseId,
				Reasoning:          responses.Reasoning,
				ServiceTier:        responses.ServiceTier,
				Store:              responses.Store,
				Temperature:        responses.Temperature,
				Text:               responses.Text,
				ToolChoice:         responses.ToolChoice,
				Tools:              responses.Tools,
				TopP:               responses.TopP,
				Truncation:         responses.Truncation,
				User:               responses.User,
				Metadata:           responses.Metadata,
				Usage:              responses.Usage,
				Error:              responses.Error,
			},
		}

		response.ResponseBytes = gjson.MustEncode(response)

		end = gtime.TimestampMilli()
		response.ConnTime = duration - now
		response.Duration = end - duration
		response.TotalTime = end - now

		responseChan <- response

		end = gtime.TimestampMilli()
		responseChan <- &model.OpenAIResponsesStreamRes{
			ConnTime:  duration - now,
			Duration:  end - duration,
			TotalTime: end - now,
			Err:       io.EOF,
		}

	}, nil); err != nil {
		logger.Errorf(ctx, "ResponsesStreamToNonStream xAI model: %s, error: %v", x.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (x *XAI) responsesErrorHandler(err *model.OpenAIResponsesError) error {
	return errors.NewRequestError(502, errors.New(fmt.Sprintf("error, status code: %s, error: %s", err.Code, gjson.MustEncodeString(err))))
}

// 流式错误可能在顶层 error / type=error, 也可能在 response.error / response.failed
func streamResponsesAPIError(res model.OpenAIResponsesStreamRes, responseBytes []byte) *model.OpenAIResponsesError {

	if res.Error != nil {
		return res.Error
	}

	if res.Response.Error != nil {
		return res.Response.Error
	}

	if res.Type == "error" || res.Type == "response.failed" || res.Response.Status == "failed" {
		return &model.OpenAIResponsesError{
			Type:           res.Type,
			SequenceNumber: res.SequenceNumber,
			Message:        string(responseBytes),
		}
	}

	return nil
}
