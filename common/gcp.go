package common

import (
	"context"
	"encoding/json"
	"fmt"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"google.golang.org/api/option"
)

type ApplicationDefaultCredentials struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

func GetGcpToken(ctx context.Context, credential, proxyUrl string) (string, error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "GetGcpToken time: %d", gtime.TimestampMilli()-now)
	}()

	adc := ApplicationDefaultCredentials{}
	if err := json.Unmarshal([]byte(credential), &adc); err != nil {
		logger.Errorf(ctx, "GetGcpToken json.Unmarshal key: %s, error: %v", credential, err)
		return "", err
	}

	client, err := credentials.NewIamCredentialsClient(ctx, option.WithCredentialsJSON([]byte(credential)))
	if err != nil {
		logger.Errorf(ctx, "GetGcpToken NewIamCredentialsClient key: %s, error: %v", credential, err)
		return "", err
	}

	defer func() {
		if err = client.Close(); err != nil {
			logger.Error(ctx, err)
		}
	}()

	request := &credentialspb.GenerateAccessTokenRequest{
		Name:  fmt.Sprintf("projects/-/serviceAccounts/%s", adc.ClientEmail),
		Scope: []string{"https://www.googleapis.com/auth/cloud-platform"},
	}

	response, err := client.GenerateAccessToken(ctx, request)
	if err != nil {
		logger.Errorf(ctx, "GetGcpToken GenerateAccessToken key: %s, error: %v", credential, err)
		return "", err
	}

	return response.AccessToken, nil
}
