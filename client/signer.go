package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"time"
)

const newLine = "\n"

// Sign returns a signature for the given parameters suitable for connecting to the Pintu API.
func sign(secret string, httpMethod string, dateTime time.Time, host string, path string) string {
	components := strings.Join([]string{
		httpMethod,
		MicrosTimestamp(dateTime).String(),
		host,
		path,
	}, newLine)
	signature := createSignature([]byte(secret), []byte(components))
	return signature
}

func createSignature(secret []byte, message []byte) string {
	hash := hmac.New(sha256.New, secret)
	_, _ = hash.Write(message)
	hex.EncodeToString(hash.Sum(nil))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}
