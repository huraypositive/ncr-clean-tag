package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}
	}
}

func makeSignature(method *string, path *string, timestamp *string) string {
	accessKey := os.Getenv("NCR_ACCESS_KEY")
	secretKey := os.Getenv("NCR_SECRET_KEY")
	digest := hmac.New(sha256.New, []byte(secretKey))
	sig := *method + " " + *path + "\n" + *timestamp + "\n" + accessKey
	digest.Write([]byte(sig))
	return base64.StdEncoding.EncodeToString(digest.Sum(nil))
}

func GetHeader(method *string, path *string) *map[string]string {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	accessKey := os.Getenv("NCR_ACCESS_KEY")
	sig := makeSignature(method, path, &timestamp)
	headers := make(map[string]string)
	headers["x-ncp-apigw-timestamp"] = timestamp
	headers["x-ncp-iam-access-key"] = accessKey
	headers["x-ncp-apigw-signature-v2"] = sig
	return &headers
}
