package main

import (
	"net/http"

	"github.com/m-mizutani/zlog"
)

type httpAuthFilter struct{}

func (x *httpAuthFilter) ReplaceString(s string) string {
	return s
}

func (x *httpAuthFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	if fieldName != "Authorization" {
		return false
	}
	if _, ok := value.([]string); !ok {
		return false
	}

	return true
}

func main() {
	logger := zlog.New()
	logger.Filters = zlog.Filters{
		&httpAuthFilter{},
	}
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	req.Header.Add("Authorization", "Barer xxxx")

	logger.With("req", req).Info("send request")
}
