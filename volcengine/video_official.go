package volcengine

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

func (v *VolcEngine) VideoCreateOfficial(ctx context.Context, data []byte) (responseBytes []byte, err error) {

	logger.Infof(ctx, "VideoCreateOfficial VolcEngine model: %s start", v.Model)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "VideoCreateOfficial VolcEngine model: %s totalTime: %d ms", v.Model, gtime.TimestampMilli()-now)
	}()

	if v.Path == "" {
		v.Path = "/contents/generations/tasks"
	}

	if responseBytes, err = util.HttpPost(ctx, v.BaseUrl+v.Path, v.header, data, nil, v.Timeout, v.ProxyUrl, v.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "VideoCreateOfficial VolcEngine model: %s, error: %v", v.Model, err)
		return nil, err
	}

	logger.Infof(ctx, "VideoCreateOfficial VolcEngine model: %s finished", v.Model)

	return responseBytes, nil
}

func (v *VolcEngine) VideoListOfficial(ctx context.Context, params model.VolcVideoListReq) (responseBytes []byte, err error) {

	logger.Infof(ctx, "VideoListOfficial VolcEngine model: %s start", v.Model)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "VideoListOfficial VolcEngine model: %s totalTime: %d ms", v.Model, gtime.TimestampMilli()-now)
	}()

	if v.Path == "" {
		v.Path = "/contents/generations/tasks"
	}

	// 构造 query string
	query := url.Values{}
	if params.PageNum != nil {
		query.Set("page_num", gconv.String(*params.PageNum))
	}
	if params.PageSize != nil {
		query.Set("page_size", gconv.String(*params.PageSize))
	}
	if params.FilterStatus != "" {
		query.Set("filter.status", params.FilterStatus)
	}
	if params.FilterTaskIds != "" {
		query.Set("filter.task_ids", params.FilterTaskIds)
	}
	if params.FilterModel != "" {
		query.Set("filter.model", params.FilterModel)
	}
	if params.FilterServiceTier != "" {
		query.Set("filter.service_tier", params.FilterServiceTier)
	}

	reqUrl := v.BaseUrl + v.Path
	if len(query) > 0 {
		reqUrl += "?" + query.Encode()
	}

	if responseBytes, err = util.HttpGet(ctx, reqUrl, v.header, nil, nil, v.Timeout, v.ProxyUrl, v.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "VideoListOfficial VolcEngine model: %s, error: %v", v.Model, err)
		return nil, err
	}

	logger.Infof(ctx, "VideoListOfficial VolcEngine model: %s finished", v.Model)

	return responseBytes, nil
}

func (v *VolcEngine) VideoRetrieveOfficial(ctx context.Context, taskId string) (responseBytes []byte, err error) {

	logger.Infof(ctx, "VideoRetrieveOfficial VolcEngine model: %s, taskId: %s start", v.Model, taskId)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "VideoRetrieveOfficial VolcEngine model: %s totalTime: %d ms", v.Model, gtime.TimestampMilli()-now)
	}()

	if v.Path == "" {
		v.Path = fmt.Sprintf("/contents/generations/tasks/%s", taskId)
	}

	if responseBytes, err = util.HttpGet(ctx, v.BaseUrl+v.Path, v.header, nil, nil, v.Timeout, v.ProxyUrl, v.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "VideoRetrieveOfficial VolcEngine model: %s, error: %v", v.Model, err)
		return nil, err
	}

	logger.Infof(ctx, "VideoRetrieveOfficial VolcEngine model: %s, taskId: %s finished", v.Model, taskId)

	return responseBytes, nil
}

func (v *VolcEngine) VideoDeleteOfficial(ctx context.Context, taskId string) (err error) {

	logger.Infof(ctx, "VideoDeleteOfficial VolcEngine model: %s, taskId: %s start", v.Model, taskId)

	now := gtime.TimestampMilli()
	defer func() {
		logger.Infof(ctx, "VideoDeleteOfficial VolcEngine model: %s totalTime: %d ms", v.Model, gtime.TimestampMilli()-now)
	}()

	if v.Path == "" {
		v.Path = fmt.Sprintf("/contents/generations/tasks/%s", taskId)
	}

	if _, err = util.HttpDelete(ctx, v.BaseUrl+v.Path, v.header, nil, nil, v.Timeout, v.ProxyUrl, v.requestErrorHandler); err != nil {
		logger.Errorf(ctx, "VideoDeleteOfficial VolcEngine model: %s, error: %v", v.Model, err)
		return err
	}

	logger.Infof(ctx, "VideoDeleteOfficial VolcEngine model: %s, taskId: %s finished", v.Model, taskId)

	return nil
}
