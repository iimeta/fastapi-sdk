package anthropic

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/logger"
	"github.com/iimeta/fastapi-sdk/model"
	"github.com/iimeta/fastapi-sdk/sdkerr"
)

type Client struct {
	key                 string
	baseURL             string
	path                string
	proxyURL            string
	isSupportSystemRole *bool
	header              map[string]string
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

func NewClient(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole *bool, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewClient Anthropic model: %s, key: %s", model, key)

	client := &Client{
		key:                 key,
		baseURL:             "https://api.anthropic.com/v1",
		path:                "/messages",
		isSupportSystemRole: isSupportSystemRole,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewClient Anthropic model: %s, baseURL: %s", model, baseURL)
		client.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewClient Anthropic model: %s, path: %s", model, path)
		client.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewClient Anthropic model: %s, proxyURL: %s", model, proxyURL[0])
		client.proxyURL = proxyURL[0]
	}

	client.header = make(map[string]string)
	client.header["x-api-key"] = key
	client.header["anthropic-version"] = "2023-06-01"
	client.header["anthropic-beta"] = "prompt-caching-2024-07-31"

	return client
}

func NewGcpClient(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole *bool, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewGcpClient Anthropic model: %s, key: %s", model, key)

	client := &Client{
		key:                 key,
		baseURL:             "https://us-east5-aiplatform.googleapis.com/v1",
		path:                "/projects/%s/locations/us-east5/publishers/anthropic/models/%s:streamRawPredict",
		isSupportSystemRole: isSupportSystemRole,
		isGcp:               true,
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewGcpClient Anthropic model: %s, baseURL: %s", model, baseURL)
		client.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewGcpClient Anthropic model: %s, path: %s", model, path)
		client.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewGcpClient Anthropic model: %s, proxyURL: %s", model, proxyURL[0])
		client.proxyURL = proxyURL[0]
	}

	client.header = make(map[string]string)
	client.header["Authorization"] = "Bearer " + key

	return client
}

func NewAwsClient(ctx context.Context, model, key, baseURL, path string, isSupportSystemRole *bool, proxyURL ...string) *Client {

	logger.Infof(ctx, "NewAwsClient Anthropic model: %s, key: %s", model, key)

	result := gstr.Split(key, "|")

	client := &Client{
		isSupportSystemRole: isSupportSystemRole,
		isAws:               true,
		awsClient: bedrockruntime.New(bedrockruntime.Options{
			Region:      result[0],
			Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(result[1], result[2], "")),
		}),
	}

	if baseURL != "" {
		logger.Infof(ctx, "NewAwsClient Anthropic model: %s, baseURL: %s", model, baseURL)
		client.baseURL = baseURL
	}

	if path != "" {
		logger.Infof(ctx, "NewAwsClient Anthropic model: %s, path: %s", model, path)
		client.path = path
	}

	if len(proxyURL) > 0 && proxyURL[0] != "" {
		logger.Infof(ctx, "NewAwsClient Anthropic model: %s, proxyURL: %s", model, proxyURL[0])
		client.proxyURL = proxyURL[0]
	}

	return client
}

func (c *Client) requestErrorHandler(ctx context.Context, response *gclient.Response) error {

	bytes := response.ReadAll()

	errRes := model.AnthropicErrorResponse{}
	if err := gjson.Unmarshal(bytes, &errRes); err != nil || errRes.Error == nil {

		reqErr := &sdkerr.RequestError{
			HttpStatusCode: response.StatusCode,
			Err:            errors.New(fmt.Sprintf("response: %s, err: %v", bytes, err)),
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

func (c *Client) apiErrorHandler(response *model.AnthropicChatCompletionRes) error {

	switch response.Error.Type {
	}

	return sdkerr.NewApiError(500, response.Error.Type, gjson.MustEncodeString(response), "api_error", "")
}
