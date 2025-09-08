package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/errors"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/util"
)

func (o *OpenAI) Responses(ctx context.Context, data []byte) (res model.OpenAIResponsesRes, err error) {

	logger.Infof(ctx, "Responses OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "Responses OpenAI model: %s totalTime: %d ms", o.Model, res.TotalTime)
	}()

	if o.Path == "" {
		if o.isAzure {
			o.Path = "/openai/responses?api-version=" + o.apiVersion
		} else {
			o.Path = "/responses"
		}
	}

	if res.ResponseBytes, err = util.HttpPost(ctx, o.BaseUrl+o.Path, o.header, data, &res, o.Timeout, o.ProxyUrl, o.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "Responses OpenAI model: %s, error: %v", o.Model, err)
		return res, err
	}

	if res.Error != nil {
		logger.Errorf(ctx, "Responses OpenAI model: %s, responsesRes: %s", o.Model, gjson.MustEncodeString(res))

		err = o.responsesErrorHandler(res.Error)
		logger.Errorf(ctx, "Responses OpenAI model: %s, error: %v", o.Model, err)

		return res, err
	}

	return res, nil
}

func (o *OpenAI) ResponsesStream(ctx context.Context, data []byte) (responseChan chan *model.OpenAIResponsesStreamRes, err error) {

	logger.Infof(ctx, "ResponsesStream OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "ResponsesStream OpenAI model: %s totalTime: %d ms", o.Model, gtime.TimestampMilli()-now)
		}
	}()

	if (o.IsSupportStream != nil && !*o.IsSupportStream) || (gstr.HasPrefix(o.Model, "o") && o.isAzure) {
		return o.ResponsesStreamToNonStream(ctx, data)
	}

	if o.Path == "" {
		if o.isAzure {
			o.Path = "/openai/responses?api-version=" + o.apiVersion
		} else {
			o.Path = "/responses"
		}
	}

	stream, err := util.SSEClient(ctx, o.BaseUrl+o.Path, o.header, data, o.Timeout, o.ProxyUrl, o.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "ResponsesStream OpenAI model: %s, error: %v", o.Model, err)
		return responseChan, err
	}

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.OpenAIResponsesStreamRes)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ResponsesStream OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", o.Model, duration-now, end-duration, end-now)

			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "ResponsesStream OpenAI model: %s, stream.Close error: %v", o.Model, err)
			}
		}()

		for {

			streamResponse, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "ResponsesStream OpenAI model: %s finished", o.Model)
				} else {
					logger.Errorf(ctx, "ResponsesStream OpenAI model: %s, error: %v", o.Model, err)
				}

				end := gtime.TimestampMilli()
				responseChan <- &model.OpenAIResponsesStreamRes{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Err:       err,
				}

				return
			}

			responsesRes := model.OpenAIResponsesStreamRes{}
			if err := json.Unmarshal(streamResponse, &responsesRes); err != nil {
				logger.Errorf(ctx, "ResponsesStream OpenAI model: %s, streamResponse: %s, error: %v", o.Model, streamResponse, err)

				end := gtime.TimestampMilli()
				responseChan <- &model.OpenAIResponsesStreamRes{
					ConnTime:  duration - now,
					Duration:  end - duration,
					TotalTime: end - now,
					Err:       errors.New(fmt.Sprintf("streamResponse: %s, error: %v", streamResponse, err)),
				}

				return
			}

			if responsesRes.Response.Error != nil {
				logger.Errorf(ctx, "ResponsesStream OpenAI model: %s, responsesRes: %s", o.Model, gjson.MustEncodeString(responsesRes))

				err = o.responsesErrorHandler(responsesRes.Response.Error)
				logger.Errorf(ctx, "ResponsesStream OpenAI model: %s, error: %v", o.Model, err)

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
				ResponseBytes: streamResponse,
				ConnTime:      duration - now,
			}

			end := gtime.TimestampMilli()
			response.Duration = end - duration
			response.TotalTime = end - now

			responseChan <- response
		}
	}, nil); err != nil {
		logger.Errorf(ctx, "ResponsesStream OpenAI model: %s, error: %v", o.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (o *OpenAI) ResponsesStreamToNonStream(ctx context.Context, data []byte) (responseChan chan *model.OpenAIResponsesStreamRes, err error) {

	logger.Infof(ctx, "ResponsesStreamToNonStream OpenAI model: %s start", o.Model)

	now := gtime.TimestampMilli()
	duration := now

	responseChan = make(chan *model.OpenAIResponsesStreamRes)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.TimestampMilli()
			logger.Infof(ctx, "ResponsesStreamToNonStream OpenAI model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", o.Model, duration-now, end-duration, end-now)
		}()

		request := make(map[string]interface{})
		if err = json.Unmarshal(data, &request); err != nil {
			logger.Errorf(ctx, "ResponsesStreamToNonStream OpenAI model: %s, data: %s, error: %v", o.Model, data, err)

			end := gtime.TimestampMilli()
			responseChan <- &model.OpenAIResponsesStreamRes{
				ConnTime:  gtime.TimestampMilli() - now,
				Duration:  end - gtime.TimestampMilli(),
				TotalTime: end - now,
				Err:       err,
			}
		}

		request["stream"] = false

		responses, err := o.Responses(ctx, gjson.MustEncode(request))
		if err != nil {

			if errors.Is(err, io.EOF) {
				logger.Infof(ctx, "ResponsesStreamToNonStream OpenAI model: %s finished", o.Model)
			} else {
				logger.Errorf(ctx, "ResponsesStreamToNonStream OpenAI model: %s, error: %v", o.Model, err)
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
			logger.Errorf(ctx, "ResponsesStreamToNonStream OpenAI model: %s, responses: %s, error: %v", o.Model, responses.ResponseBytes, err)

			end := gtime.TimestampMilli()
			responseChan <- &model.OpenAIResponsesStreamRes{
				ConnTime:  duration - now,
				Duration:  end - duration,
				TotalTime: end - now,
				Err:       errors.New(fmt.Sprintf("streamResponse: %s, error: %v", responses.ResponseBytes, err)),
			}

			return
		}

		if responsesRes.Error != nil {
			logger.Errorf(ctx, "ResponsesStreamToNonStream OpenAI model: %s, responsesRes: %s", o.Model, gjson.MustEncodeString(responsesRes))

			err = o.responsesErrorHandler(responsesRes.Error)
			logger.Errorf(ctx, "ResponsesStreamToNonStream OpenAI model: %s, error: %v", o.Model, err)

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
					delta += output.Arguments
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
		logger.Errorf(ctx, "ResponsesStreamToNonStream OpenAI model: %s, error: %v", o.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (o *OpenAI) responsesErrorHandler(err *model.OpenAIResponsesError) error {
	return errors.NewRequestError(500, errors.New(fmt.Sprintf("error, status code: %s, error: %s", err.Code, gjson.MustEncodeString(err))))
}
