package logging

import (
	base642 "encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var API_KEY = "8215_8323f00494625bce:dfd821b90482bbb7a883135fb355e2a859923d188a5e2a2a8f7f3403cf0f1139"

func Timber(errorText string) {
	req, err := http.NewRequest("POST", "https://logs.timber.io/frames", strings.NewReader(errorText))
	if err != nil {
		gin.ErrorLogger()
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("authorization", "Basic "+base642.StdEncoding.EncodeToString([]byte(API_KEY)))
	_, err = http.DefaultClient.Do(req)
}

func base64(key string) string {
	return base642.StdEncoding.EncodeToString([]byte(key))
}
