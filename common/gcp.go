package common

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iimeta/fastapi-sdk/v2/errors"
	"github.com/iimeta/fastapi-sdk/v2/logger"
	"github.com/iimeta/fastapi-sdk/v2/util"
)

const (
	gcpTokenUrl   = "https://oauth2.googleapis.com/token"
	gcpTokenScope = "https://www.googleapis.com/auth/cloud-platform"
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

type GcpTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// 使用服务账号密钥签发JWT, 通过OAuth2 JWT Bearer方式直接换取access_token, 有效期1小时
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

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(adc.PrivateKey))
	if err != nil {
		logger.Errorf(ctx, "GetGcpToken ParseRSAPrivateKeyFromPEM key: %s, error: %v", credential, err)
		return "", err
	}

	iat := time.Now().Add(-10 * time.Second)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":   adc.ClientEmail,
		"scope": gcpTokenScope,
		"aud":   gcpTokenUrl,
		"iat":   iat.Unix(),
		"exp":   iat.Add(time.Hour).Unix(),
	})

	token.Header["kid"] = adc.PrivateKeyId

	assertion, err := token.SignedString(privateKey)
	if err != nil {
		logger.Errorf(ctx, "GetGcpToken SignedString key: %s, error: %v", credential, err)
		return "", err
	}

	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	data := url.Values{
		"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
		"assertion":  {assertion},
	}

	response := GcpTokenResponse{}
	if _, _, err = util.HttpPost(ctx, gcpTokenUrl, header, []byte(data.Encode()), &response, 30*time.Second, proxyUrl, nil); err != nil {
		logger.Errorf(ctx, "GetGcpToken HttpPost key: %s, error: %v", credential, err)
		return "", err
	}

	if response.AccessToken == "" {
		err = errors.New("GetGcpToken access_token is empty")
		logger.Errorf(ctx, "GetGcpToken key: %s, error: %v", credential, err)
		return "", err
	}

	return response.AccessToken, nil
}
