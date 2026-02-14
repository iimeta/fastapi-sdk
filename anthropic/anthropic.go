package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/iimeta/fastapi-sdk/v2/consts"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/model"
	"github.com/iimeta/fastapi-sdk/v2/options"
)

type Anthropic struct {
	*options.AdapterOptions
	header    map[string]string
	isGcp     bool
	isAws     bool
	awsClient *bedrockruntime.Client
	region    string
	accessKey string
	secretKey string
}

func NewAdapter(ctx context.Context, options *options.AdapterOptions) *Anthropic {

	anthropic := &Anthropic{
		AdapterOptions: options,
		header: map[string]string{
			"x-api-key":         options.Key,
			"anthropic-version": "2023-06-01",
			"anthropic-beta":    "prompt-caching-2024-07-31",
		},
	}

	if anthropic.BaseUrl == "" {
		anthropic.BaseUrl = "https://api.anthropic.com/v1"
	}

	for k, v := range anthropic.Header {
		anthropic.header[k] = v
	}

	logger.Infof(ctx, "NewAdapter Anthropic model: %s, key: %s", anthropic.Model, anthropic.Key)

	return anthropic
}

func NewGcpAdapter(ctx context.Context, options *options.AdapterOptions) *Anthropic {

	gcp := &Anthropic{
		AdapterOptions: options,
		header: map[string]string{
			"Authorization": "Bearer " + options.Key,
		},
		isGcp: true,
	}

	if gcp.BaseUrl == "" {
		gcp.BaseUrl = "https://us-east5-aiplatform.googleapis.com/v1"
	}

	//if gcp.Path == "" {
	//	gcp.Path = "/projects/%s/locations/us-east5/publishers/anthropic/models/%s:streamRawPredict"
	//}

	for k, v := range gcp.Header {
		gcp.header[k] = v
	}

	logger.Infof(ctx, "NewGcpAdapter Anthropic model: %s, key: %s", gcp.Model, gcp.Key)

	return gcp
}

func NewAwsAdapter(ctx context.Context, options *options.AdapterOptions) *Anthropic {

	result := gstr.Split(options.Key, "|")

	aws := &Anthropic{
		AdapterOptions: options,
		isAws:          true,
		region:         result[0],
		accessKey:      result[1],
		secretKey:      result[2],
	}

	if aws.BaseUrl == "" || aws.BaseUrl == consts.PROVIDER_AWS_CLAUDE {
		aws.BaseUrl = fmt.Sprintf("https://bedrock-runtime.%s.amazonaws.com", result[0])
	} else if strings.HasSuffix(aws.BaseUrl, "/") {
		aws.BaseUrl = aws.BaseUrl[:len(aws.BaseUrl)-1]
	}

	if options.Path == "" {
		options.Path = fmt.Sprintf("/model/%s/invoke", options.Model)
	}

	logger.Infof(ctx, "NewAwsAdapter Anthropic model: %s, key: %s", aws.Model, aws.Key)

	return aws
}

func (a *Anthropic) requestErrorHandler(ctx context.Context, response *http.Response) error {

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	errRes := model.AnthropicErrorResponse{}
	if err := json.Unmarshal(bytes, &errRes); err != nil || errRes.Error == nil {

		reqErr := &errors.RequestError{
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

	return errors.NewRequestError(response.StatusCode, errors.New(fmt.Sprintf("error, status code: %d, response: %s", response.StatusCode, gjson.MustEncodeString(errRes.Error))))
}

func (a *Anthropic) apiErrorHandler(response *model.AnthropicChatCompletionRes) error {

	switch response.Error.Type {
	}

	return errors.NewApiError(500, response.Error.Type, gjson.MustEncodeString(response), "api_error", "")
}
