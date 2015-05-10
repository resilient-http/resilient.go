package client

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

// HTTP client request struct storing params and request configuration.
type Options struct {
	URL         string
	Method      string
	Header      map[string]string
	TargetType  string
	ForceType   string
	Data        map[string]interface{}
	FormData    url.Values
	Body        io.Reader
	QueryData   url.Values
	Client      *http.Client
	Cookies     []*http.Cookie
	DisableTLS  bool
	ProxyURL    string
	Timeout     time.Duration
	DialTimeout time.Duration
}
