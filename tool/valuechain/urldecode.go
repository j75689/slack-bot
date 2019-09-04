package valuechain

import (
	"net/url"
	"strings"
)

// URLDecoder process url decode
type URLDecoder struct {
	Basic
	next Handler
}

func (handler *URLDecoder) resolve(tag, value string) string {
	if strings.ToUpper(tag) == "URLDECODE" {
		if data, err := url.QueryUnescape(value); err == nil {
			return data
		}
	}

	return handler.Basic.next.resolve(tag, value)
}
