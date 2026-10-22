package xai

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (x *XAI) VideoCreate(ctx context.Context, request model.VideoCreateRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoCreate xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoCreate xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	var data any = request
	if !slices.Contains(x.ReqPassthroughParams, "req_data") {
		if data, err = x.ConvVideoCreateRequest(ctx, request); err != nil {
			logger.Errorf(ctx, "VideoCreate xAI ConvVideoCreateRequest error: %v", err)
			return response, err
		}
	}

	responseBytes, responseHeader, err := util.HttpPost(ctx, x.BaseUrl+x.videoCreatePath(), x.header, data, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if job, ok := xaiVideoJobFromError(err); ok {
		job.ResponseHeaders = responseHeader
		return job, nil
	}

	if err != nil {
		logger.Errorf(ctx, "VideoCreate xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvVideoJobResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "VideoCreate xAI ConvVideoJobResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "VideoCreate xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) VideoRemix(ctx context.Context, request model.VideoRemixRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoRemix xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoRemix xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	official := model.XAIVideoCreateReq{
		Model:  x.Model,
		Prompt: request.Prompt,
		Video:  convVideoContent(request.VideoId),
	}

	x.header["Content-Type"] = "application/json"

	responseBytes, responseHeader, err := util.HttpPost(ctx, x.BaseUrl+x.videoCreatePath(), x.header, official, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if job, ok := xaiVideoJobFromError(err); ok {
		job.ResponseHeaders = responseHeader
		return job, nil
	}

	if err != nil {
		logger.Errorf(ctx, "VideoRemix xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvVideoJobResponse(ctx, responseBytes); err != nil {
		logger.Errorf(ctx, "VideoRemix xAI ConvVideoJobResponse error: %v", err)
		return response, err
	}

	response.ResponseBytes = responseBytes
	response.ResponseHeaders = responseHeader

	logger.Infof(ctx, "VideoRemix xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) VideoList(ctx context.Context, request model.VideoListRequest) (response model.VideoListResponse, err error) {

	logger.Infof(ctx, "VideoList xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoList xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	bytes, _, err := util.HttpGet(ctx, x.BaseUrl+x.videoCreatePath(), x.header, request, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoList xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvVideoListResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoList xAI ConvVideoListResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoList xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) VideoRetrieve(ctx context.Context, request model.VideoRetrieveRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoRetrieve xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoRetrieve xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	bytes, _, err := util.HttpGet(ctx, x.BaseUrl+x.videoRetrievePath(request.VideoId), x.header, nil, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if job, ok := xaiVideoJobFromError(err); ok {
		return job, nil
	}

	if err != nil {
		logger.Errorf(ctx, "VideoRetrieve xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoRetrieve xAI ConvVideoJobResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoRetrieve xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) VideoDelete(ctx context.Context, request model.VideoDeleteRequest) (response model.VideoJobResponse, err error) {

	logger.Infof(ctx, "VideoDelete xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoDelete xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	bytes, _, err := util.HttpDelete(ctx, x.BaseUrl+x.videoRetrievePath(request.VideoId), x.header, nil, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoDelete xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvVideoJobResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoDelete xAI ConvVideoJobResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoDelete xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) VideoContent(ctx context.Context, request model.VideoContentRequest) (response model.VideoContentResponse, err error) {

	logger.Infof(ctx, "VideoContent xAI model: %s start", x.Model)

	now := gtime.TimestampMilli()
	defer func() {
		response.TotalTime = gtime.TimestampMilli() - now
		logger.Infof(ctx, "VideoContent xAI model: %s totalTime: %d ms", x.Model, response.TotalTime)
	}()

	job, err := x.VideoRetrieve(ctx, model.VideoRetrieveRequest{VideoId: request.VideoId})
	if err != nil {
		logger.Errorf(ctx, "VideoContent xAI model: %s, error: %v", x.Model, err)
		return response, err
	}
	if job.VideoUrl == "" {
		return response, errors.New("video url is empty")
	}

	bytes, _, err := util.HttpGet(ctx, job.VideoUrl, nil, nil, nil, x.Timeout, x.ProxyUrl, x.requestErrorHandler)
	if err != nil {
		logger.Errorf(ctx, "VideoContent xAI model: %s, error: %v", x.Model, err)
		return response, err
	}

	if response, err = x.ConvVideoContentResponse(ctx, bytes); err != nil {
		logger.Errorf(ctx, "VideoContent xAI ConvVideoContentResponse error: %v", err)
		return response, err
	}

	logger.Infof(ctx, "VideoContent xAI model: %s finished", x.Model)

	return response, nil
}

func (x *XAI) videoCreatePath() string {
	if x.Path != "" && !isXAIVideoRetrievePath(x.Path) {
		return x.Path
	}
	return "/videos/generations"
}

func (x *XAI) videoRetrievePath(requestId string) string {
	return fmt.Sprintf("/videos/%s", requestId)
}

func isXAIVideoRetrievePath(path string) bool {
	path = strings.TrimRight(path, "/")
	if path == "" || strings.HasSuffix(path, "/generations") {
		return false
	}
	return strings.Contains(path, "/videos/")
}

type xaiVideoJobError struct {
	HttpStatusCode int
	Job            model.VideoJobResponse
}

func (e *xaiVideoJobError) Error() string {

	if e == nil {
		return ""
	}

	if e.Job.Error != nil && e.Job.Error.Message != "" {
		return e.Job.Error.Message
	}

	if e.Job.Error != nil && e.Job.Error.Code != "" {
		return e.Job.Error.Code
	}

	return fmt.Sprintf("error, status code: %d", e.HttpStatusCode)
}

func isXAIVideoJobError(job model.VideoJobResponse) bool {

	if job.Error != nil && (job.Error.Code != "" || job.Error.Message != "") {
		return true
	}

	switch job.Status {
	case "failed", "expired", "deleted", "cancelled":
		return true
	default:
		return false
	}
}

func isXAIVideoHTTPError(raw []byte, job model.VideoJobResponse) bool {

	if !isXAIVideoJobError(job) {
		return false
	}

	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil {
		return false
	}

	if _, ok := body["video"]; ok {
		return true
	}

	if _, ok := body["progress"]; ok {
		return true
	}

	if usage, ok := body["usage"].(map[string]any); ok && len(usage) > 0 {
		return true
	}

	if _, ok := body["error"].(string); ok {
		return true
	}

	if code, _ := body["code"].(string); strings.Contains(code, "imagine:") {
		return true
	}

	return false
}

func xaiVideoJobFromError(err error) (model.VideoJobResponse, bool) {

	var jobErr *xaiVideoJobError
	if errors.As(err, &jobErr) && jobErr != nil {
		return jobErr.Job, true
	}

	return model.VideoJobResponse{}, false
}
