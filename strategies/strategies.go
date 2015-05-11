package strategies

import (
	"net/http"
)

type Handler func(*http.Request, *http.Response, error) bool

type Strategy func() Handler
