package google

import (
	"context"
	"encoding/base64"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (g *Google) VideoCreate(ctx context.Context, request model.VideoCreateRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoCreate Google model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoCreate Google model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	var data any = request
	if !slices.Contains(g.ReqPassthroughParams, "req_data") {
		if data, err = g.ConvVideoCreateRequest(ctx, request); err != nil {
			logger.Errorf(ctx, "VideoCreate Google ConvVideoCreateRequest error: %v", err)
			return response, err
		}
	} else {
		data = forceInteractionBackground(gjson.MustEncode(request))
	}

	bytes, responseHeader, err := util.HttpPost(ctx, g.interactionsURL("", false), g.header, data, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoCreate Google model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoCreate Google ConvVideoJobResponse error: %v", err)
		return response, err
	}

	response.Model = g.Model
	response.Prompt = request.Prompt
	response.ResponseBytes = bytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "VideoCreate Google model: %s finished, id: %s", g.Model, response.Id)

	return response, nil
}

func (g *Google) VideoRemix(ctx context.Context, request model.VideoRemixRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoRemix Google model: %s start", g.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoRemix Google model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	body := map[string]any{
		"model":                   g.Model,
		"input":                   request.Prompt,
		"previous_interaction_id": request.VideoId,
		"background":              true,
		"response_format": map[string]any{
			"type":     "video",
			"delivery": "uri",
		},
	}

	bytes, responseHeader, err := util.HttpPost(ctx, g.interactionsURL("", false), g.header, body, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoRemix Google model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoRemix Google ConvVideoJobResponse error: %v", err)
		return response, err
	}

	response.Model = g.Model
	response.Prompt = request.Prompt
	response.RemixedFromVideoId = &request.VideoId
	response.ResponseBytes = bytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "VideoRemix Google model: %s finished, id: %s", g.Model, response.Id)

	return response, nil
}

func (g *Google) VideoList(ctx context.Context, request model.VideoListRequest) (response model.VideoListResponse, err error) {
	return response, errors.New("Google video does not support list")
}

func (g *Google) VideoRetrieve(ctx context.Context, request model.VideoRetrieveRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoRetrieve Google model: %s, videoId: %s start", g.Model, request.VideoId)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoRetrieve Google model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	bytes, responseHeader, err := util.HttpGet(ctx, g.interactionsURL(request.VideoId, false), g.header, nil, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoRetrieve Google model: %s, error: %v", g.Model, err)
		return response, err
	}

	if response, err = g.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoRetrieve Google ConvVideoJobResponse error: %v", err)
		return response, err
	}

	if response.Model == "" {
		response.Model = g.Model
	}

	response.ResponseBytes = bytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "VideoRetrieve Google model: %s, videoId: %s, status: %s finished", g.Model, request.VideoId, response.Status)

	return response, nil
}

func (g *Google) VideoDelete(ctx context.Context, request model.VideoDeleteRequest) (response model.VideoJobResponse, err error) {
	return response, errors.New("Google video does not support delete")
}

func (g *Google) VideoContent(ctx context.Context, request model.VideoContentRequest) (response model.VideoContentResponse, err error) {

	logger.Infof(ctx, "VideoContent Google model: %s, videoId: %s start", g.Model, request.VideoId)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoContent Google model: %s totalTime: %d ms", g.Model, response.TotalTime)
	}()

	retrieve, err := g.VideoRetrieve(ctx, model.VideoRetrieveRequest{VideoId: request.VideoId})
	if err != nil {
		logger.Errorf(ctx, "VideoContent Google VideoRetrieve error: %v", err)
		return response, err
	}

	interactionRes := model.GoogleInteractionRes{}
	if len(retrieve.ResponseBytes) > 0 {
		_ = gjson.Unmarshal(retrieve.ResponseBytes, &interactionRes)
	}

	uri, videoData, _ := extractInteractionVideo(interactionRes.Steps)
	if uri == "" {
		uri = retrieve.VideoUrl
	}

	if videoData != "" {
		decoded, decodeErr := base64.StdEncoding.DecodeString(videoData)
		if decodeErr != nil {
			logger.Errorf(ctx, "VideoContent Google decode error: %v", decodeErr)
			return response, decodeErr
		}
		response.Data = decoded
		logger.Infof(ctx, "VideoContent Google model: %s, videoId: %s finished, size: %d bytes", g.Model, request.VideoId, len(response.Data))
		return response, nil
	}

	if uri == "" {
		return response, fmt.Errorf("VideoContent Google: video is empty for videoId %s", request.VideoId)
	}

	content, err := g.downloadVideoURI(ctx, uri)
	if err != nil {
		logger.Errorf(ctx, "VideoContent Google download error: %v", err)
		return response, err
	}

	response.Data = content

	logger.Infof(ctx, "VideoContent Google model: %s, videoId: %s finished, size: %d bytes", g.Model, request.VideoId, len(response.Data))

	return response, nil
}

func (g *Google) downloadVideoURI(ctx context.Context, uri string) ([]byte, error) {

	if strings.Contains(uri, "/files/") && !strings.Contains(uri, ":download") && !strings.Contains(uri, "alt=media") {
		fileId := extractGoogleFileId(uri)
		if fileId != "" {
			if err := g.waitFileActive(ctx, fileId); err != nil {
				return nil, err
			}
			uri = g.fileDownloadURL(fileId)
		}
	}

	data, _, err := util.HttpGet(ctx, uri, g.header, nil, nil, g.Timeout, g.ProxyUrl, g.requestErrorHandler)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func extractGoogleFileId(uri string) string {

	idx := strings.Index(uri, "/files/")
	if idx < 0 {
		return ""
	}

	id := uri[idx+len("/files/"):]
	if i := strings.IndexAny(id, "?:"); i >= 0 {
		id = id[:i]
	}

	return strings.Trim(id, "/")
}

func (g *Google) fileDownloadURL(fileId string) string {
	return fmt.Sprintf("%s/files/%s:download?alt=media", g.BaseUrl, fileId)
}

func (g *Google) waitFileActive(ctx context.Context, fileId string) error {

	deadline := time.Now().Add(g.Timeout)
	if g.Timeout <= 0 {
		deadline = time.Now().Add(2 * time.Minute)
	}

	for time.Now().Before(deadline) {

		file, err := g.FileRetrieve(ctx, model.FileRetrieveRequest{FileId: fileId})
		if err != nil {
			return err
		}

		switch strings.ToLower(file.Status) {
		case "processed", "active":
			return nil
		case "failed", "error":
			return errors.New("Google file processing failed")
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}

	return errors.New("Google file processing timeout")
}
