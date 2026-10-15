package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (g *Google) Interactions(ctx context.Context, data []byte) (res model.GoogleInteractionRes, err error) {

	logger.Infof(ctx, "Interactions Google model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "Interactions Google model: %s totalTime: %d ms", g.Model, res.TotalTime)
	}()

	var responseHeader map[string][]string
	if res.ResponseBytes, responseHeader, err = util.HttpPost(ctx, g.interactionsURL("", false), g.header, data, &res, g.Timeout, g.ProxyUrl, g.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "Interactions Google model: %s, error: %v", g.Model, err)
		return res, err
	}

	res.ResponseHeaders = responseHeader

	if res.Error.Code != 0 {
		logger.Errorf(ctx, "Interactions Google model: %s, interactionRes: %s", g.Model, gjson.MustEncodeString(res))

		err = g.interactionErrorHandler(&res)
		logger.Errorf(ctx, "Interactions Google model: %s, error: %v", g.Model, err)

		return res, err
	}

	return res, nil
}

func (g *Google) InteractionsStream(ctx context.Context, data []byte) (responseChan chan *model.GoogleInteractionRes, err error) {

	logger.Infof(ctx, "InteractionsStream Google model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		if err != nil {
			logger.Infof(ctx, "InteractionsStream Google model: %s totalTime: %d ms", g.Model, gtime.TimestampMilli()-now)
		}
	}()

	stream, err := util.SSEClient(ctx, g.interactionsURL("", true), g.header, data, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "InteractionsStream Google model: %s, error: %v", g.Model, err)
		return responseChan, err
	}

	streamResponseHeaders := stream.Response.Header

	duration := gtime.TimestampMilli()

	responseChan = make(chan *model.GoogleInteractionRes)

	if err = grpool.AddWithRecover(ctx, func(ctx context.Context) {

		defer func() {
			end := gtime.TimestampMilli()
			logger.Infof(ctx, "InteractionsStream Google model: %s connTime: %d ms, duration: %d ms, totalTime: %d ms", g.Model, duration-now, end-duration, end-now)

			if err := stream.Close(); err != nil {
				logger.Errorf(ctx, "InteractionsStream Google model: %s, stream.Close error: %v", g.Model, err)
			}
		}()

		var usage *model.GoogleInteractionUsage

		for {

			responseBytes, err := stream.Recv()
			if err != nil {

				if errors.Is(err, io.EOF) {
					logger.Infof(ctx, "InteractionsStream Google model: %s finished", g.Model)
				} else {
					logger.Errorf(ctx, "InteractionsStream Google model: %s, error: %v", g.Model, err)
				}

				end := gtime.TimestampMilli()
				responseChan <- &model.GoogleInteractionRes{
					Usage:           usage,
					SSEEvent:        stream.Event(),
					ResponseHeaders: streamResponseHeaders,
					ConnTime:        duration - now,
					Duration:        end - duration,
					TotalTime:       end - now,
					Err:             err,
				}

				return
			}

			interactionRes := model.GoogleInteractionRes{}
			if err := json.Unmarshal(responseBytes, &interactionRes); err != nil {
				logger.Errorf(ctx, "InteractionsStream Google model: %s, response: %s, error: %v", g.Model, responseBytes, err)

				end := gtime.TimestampMilli()
				responseChan <- &model.GoogleInteractionRes{
					SSEEvent:        stream.Event(),
					ResponseHeaders: streamResponseHeaders,
					ConnTime:        duration - now,
					Duration:        end - duration,
					TotalTime:       end - now,
					Err:             errors.New(fmt.Sprintf("response: %s, error: %v", responseBytes, err)),
				}

				return
			}

			if interactionRes.Error.Code != 0 {
				logger.Errorf(ctx, "InteractionsStream Google model: %s, interactionRes: %s", g.Model, gjson.MustEncodeString(interactionRes))

				err = g.interactionErrorHandler(&interactionRes)
				logger.Errorf(ctx, "InteractionsStream Google model: %s, error: %v", g.Model, err)

				end := gtime.TimestampMilli()
				responseChan <- &model.GoogleInteractionRes{
					SSEEvent:        stream.Event(),
					ResponseBytes:   responseBytes,
					ResponseHeaders: streamResponseHeaders,
					ConnTime:        duration - now,
					Duration:        end - duration,
					TotalTime:       end - now,
					Err:             err,
				}

				return
			}

			if interactionRes.Usage != nil {
				usage = interactionRes.Usage
			}

			response := &model.GoogleInteractionRes{
				Object:                interactionRes.Object,
				Id:                    interactionRes.Id,
				Status:                interactionRes.Status,
				Model:                 interactionRes.Model,
				PreviousInteractionId: interactionRes.PreviousInteractionId,
				Steps:                 interactionRes.Steps,
				Usage:                 interactionRes.Usage,
				EventType:             interactionRes.EventType,
				Delta:                 interactionRes.Delta,
				SSEEvent:              stream.Event(),
				ResponseBytes:         responseBytes,
				ResponseHeaders:       streamResponseHeaders,
				ConnTime:              duration - now,
			}

			end := gtime.TimestampMilli()
			response.Duration = end - duration
			response.TotalTime = end - now

			responseChan <- response
		}
	}, nil); err != nil {
		logger.Errorf(ctx, "InteractionsStream Google model: %s, error: %v", g.Model, err)
		return responseChan, err
	}

	return responseChan, nil
}

func (g *Google) InteractionsRetrieve(ctx context.Context, id string) (res model.GoogleInteractionRes, err error) {

	logger.Infof(ctx, "InteractionsRetrieve Google model: %s, id: %s start", g.Model, id)

	now := gtime.TimestampMilli()
	defer func() {
		res.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "InteractionsRetrieve Google model: %s totalTime: %d ms", g.Model, res.TotalTime)
	}()

	var responseHeader map[string][]string
	if res.ResponseBytes, responseHeader, err = util.HttpGet(ctx, g.interactionsURL(id, false), g.header, nil, &res, g.Timeout, g.ProxyUrl, g.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "InteractionsRetrieve Google model: %s, error: %v", g.Model, err)
		return res, err
	}

	res.ResponseHeaders = responseHeader

	if res.Error.Code != 0 {
		logger.Errorf(ctx, "InteractionsRetrieve Google model: %s, interactionRes: %s", g.Model, gjson.MustEncodeString(res))

		err = g.interactionErrorHandler(&res)
		logger.Errorf(ctx, "InteractionsRetrieve Google model: %s, error: %v", g.Model, err)

		return res, err
	}

	return res, nil
}
