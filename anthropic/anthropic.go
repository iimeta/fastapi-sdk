package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/sdkerr"
)

type Anthropic struct {
	model               string
	key                 string
	baseURL             string
	path                string
	proxyURL            string
	header              map[string]string
	isSupportSystemRole *bool
	isSupportStream     *bool
	isGcp               bool
	isAws               bool
	awsClient           *bedrockruntime.Client
}

// https://docs.aws.amazon.com/bedrock/latest/userguide/model-ids.html
var AwsModelIDMap = map[string]string{
	"claude-2.0":                 "anthropic.claude-v2",
	"claude-2.1":                 "anthropic.claude-v2:1",
	"claude-3-sonnet-20240229":   "anthropic.claude-3-sonnet-20240229-v1:0",
	"claude-3-5-sonnet-20240620": "anthropic.claude-3-5-sonnet-20240620-v1:0",
	"claude-3-5-sonnet-20241022": "anthropic.claude-3-5-sonnet-20241022-v2:0",
	"claude-3-haiku-20240307":    "anthropic.claude-3-haiku-20240307-v1:0",
	"claude-3-5-haiku-20241022":  "anthropic.claude-3-5-haiku-20241022-v1:0",
	"claude-3-opus-20240229":     "anthropic.claude-3-opus-20240229-v1:0",
	"claude-instant-1.2":         "anthropic.claude-instant-v1",
}

func NewAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *Anthropic {

	logger.Infof(ctx, "NewAdapter Anthropic model: %s, key: %s", model, key)

	anthropic := &Anthropic{
		model:   model,
		key:     key,
		baseURL: "https://api.anthropic.com/v1",
		path:    "/messages",
		header: g.MapStrStr{
			"x-api-key":         key,
			"anthropic-version": "2023-06-01",
			"anthropic-beta":    "prompt-caching-2024-07-31",
		},
		isSupportSystemRole: isSupportSystemRole,
		isSupportStream:     isSupportStream,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewAdapter Anthropic model: %s, baseURL: %s", model, baseURL)
		anthropic.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewAdapter Anthropic model: %s, path: %s", model, path)
		anthropic.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAdapter Anthropic model: %s, proxyURL: %s", model, proxyURL[0])
		anthropic.proxyURL = proxyURL[0]
	}

	return anthropic
}

func NewGcpAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *Anthropic {

	logger.Infof(ctx, "NewGcpAdapter Anthropic model: %s, key: %s", model, key)

	gcp := &Anthropic{
		model:   model,
		key:     key,
		baseURL: "https://us-east5-aiplatform.googleapis.com/v1",
		path:    "/projects/%s/locations/us-east5/publishers/anthropic/models/%s:streamRawPredict",
		header: g.MapStrStr{
			"Authorization": "Bearer " + key,
		},
		isSupportSystemRole: isSupportSystemRole,
		isSupportStream:     isSupportStream,
		isGcp:               true,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewGcpAdapter Anthropic model: %s, baseURL: %s", model, baseURL)
		gcp.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewGcpAdapter Anthropic model: %s, path: %s", model, path)
		gcp.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewGcpAdapter Anthropic model: %s, proxyURL: %s", model, proxyURL[0])
		gcp.proxyURL = proxyURL[0]
	}

	return gcp
}

func NewAwsAdapter(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole, isSupportStream *bool, proxyURL ...string) *Anthropic {

	logger.Infof(ctx, "NewAwsAdapter Anthropic model: %s, key: %s", model, key)

	result := gstr.Split(key, "|")

	aws := &Anthropic{
		model:               model,
		key:                 key,
		isSupportSystemRole: isSupportSystemRole,
		isSupportStream:     isSupportStream,
		isAws:               true,
		awsClient: bedrockruntime.New(bedrockruntime.Options{
			Region:      result[0],
			Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(result[1], result[2], "")),
		}),
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewAwsAdapter Anthropic model: %s, baseURL: %s", model, baseURL)
		aws.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewAwsAdapter Anthropic model: %s, path: %s", model, path)
		aws.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAwsAdapter Anthropic model: %s, proxyURL: %s", model, proxyURL[0])
		aws.proxyURL = proxyURL[0]
	}

	return aws
}

func (a *Anthropic) requestErrorHandler(ctx context.Context, response *http.Response) error {

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	errRes := model.AnthropicErrorResponse{}
	if err := json.Unmarshal(bytes, &errRes); err != nil || errRes.Error == nil {

		reqErr := &sdkerr.RequestError{
			HttpStatusCode: response.StatusCode,
			Err:            errors.New(fmt.Sprintf("response: %s, error: %v", bytes, err)),
		}

		if errRes.Error != nil {
			reqErr.Err = errors.New(gjson.MustEncodeString(errRes.Error))
		}

		return reqErr
	}

	switch errRes.Error.Type {
	}

	return sdkerr.NewRequestError(500, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, gjson.MustEncodeString(errRes.Error))))
}

func (a *Anthropic) apiErrorHandler(response *model.AnthropicChatCompletionRes) error {

	switch response.Error.Type {
	}

	return sdkerr.NewApiError(500, response.Error.Type, gjson.MustEncodeString(response), "api_error", "")
}
