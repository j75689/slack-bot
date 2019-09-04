package valuechain

import (
	"net/url"
	"strings"
)

// URLEncoder process url encode
type URLEncoder struct {
	Basic
	next Handler
}

func (handler *URLEncoder) resolve(tag, value string) string {
	if strings.ToUpper(tag) == "URLENCODE" {
		return url.QueryEscape(value)
	}

	return handler.Basic.next.resolve(tag, value)
}
