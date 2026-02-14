package anthropic

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/iimeta/fastapi-sdk/v2/util"
)

func signHeader(path, region, accessKey, secretKey string, data []byte) map[string]string {

	now := time.Now().UTC()
	header := make(map[string]string)

	header["accept"] = "application/json"
	header["host"] = fmt.Sprintf("bedrock-runtime.%s.amazonaws.com", region)
	header["x-amz-date"] = now.Format("20060102T150405Z")

	payloadHash := sha256.Sum256(data)
	payloadHashHex := hex.EncodeToString(payloadHash[:])

	canonicalRequest := createCanonicalRequest(path, header["host"], header["x-amz-date"], payloadHashHex)

	stringToSign := createStringToSign(now, region, canonicalRequest)

	signature := calculateSignature(now, region, secretKey, stringToSign)

	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", now.Format("20060102"), region, "bedrock")

	signedHeaders := "content-type;host;x-amz-date"

	header["Authorization"] = fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s", accessKey, credentialScope, signedHeaders, signature)

	return header
}

func createCanonicalRequest(path, host, xAmzDate string, payloadHashHex string) string {

	uri := strings.ReplaceAll(path, ":", "%3A")
	if uri == "" {
		uri = "/"
	}

	canonicalHeaders := "content-type:application/json\n"
	var signedHeadersList []string

	signedHeadersList = append(signedHeadersList, "content-type")

	canonicalHeaders += fmt.Sprintf("host:%s\n", host)
	signedHeadersList = append(signedHeadersList, "host")

	canonicalHeaders += fmt.Sprintf("x-amz-date:%s\n", xAmzDate)
	signedHeadersList = append(signedHeadersList, "x-amz-date")

	signedHeaders := strings.Join(signedHeadersList, ";")

	canonicalRequest := fmt.Sprintf("%s\n%s\n\n%s\n%s\n%s", http.MethodPost, uri, canonicalHeaders, signedHeaders, payloadHashHex)

	return canonicalRequest
}

func createStringToSign(t time.Time, region, canonicalRequest string) string {

	hash := sha256.Sum256([]byte(canonicalRequest))
	hashHex := hex.EncodeToString(hash[:])

	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", t.Format("20060102"), region, "bedrock")

	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s", t.Format("20060102T150405Z"), credentialScope, hashHex)

	return stringToSign
}

func calculateSignature(t time.Time, region, secretKey, stringToSign string) string {

	dateKey := util.HMACSHA256([]byte("AWS4"+secretKey), []byte(t.Format("20060102")))
	regionKey := util.HMACSHA256(dateKey, []byte(region))
	serviceKey := util.HMACSHA256(regionKey, []byte("bedrock"))
	signingKey := util.HMACSHA256(serviceKey, []byte("aws4_request"))
	signature := util.HMACSHA256(signingKey, []byte(stringToSign))

	return hex.EncodeToString(signature)
}
