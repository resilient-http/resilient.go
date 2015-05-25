package compression

import (
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"gopkg.in/resilient-http/resilient.go.v0/client"
	"gopkg.in/resilient-http/resilient.go.v0/middlewares"
	"net/http"
	"strings"
	"time"
)

const (
	MISSING_SERVERS = iota
)

var Error = func(msg string, code int) error {
	return middlewares.NewMiddlewareError("lookup", msg, code)
}

type Options struct {
	ContentEncoding string
}

type Compression struct {
	options   Options
	resilient middlewares.Resilient
}

type compression struct {
	writer          func(buffer io.Writer) (io.WriteCloser, error)
	reader          func(buffer io.Reader) (io.ReadCloser, error)
	ContentEncoding string
}

func New(opts Options) middlewares.Handler {
	return func(r middlewares.Resilient) (middlewares.Middleware, error) {

		if opts.ContentEncoding == "" {
			opts.ContentEncoding = "gzip"
		}

		return &Compression{
			options:   opts,
			resilient: r,
		}, nil
	}
}

func (m *Compression) Out(o *client.Options) error {
	switch m.options.ContentEncoding {
	case "gzip":
		o.Compression = Gzip()
		break
	case "deflate":
		o.Compression = Zlip()
		break
	}

	return nil
}

func (m *Compression) In(*http.Request, *http.Response, error) error {
	return nil
}

func Gzip() *compression {
	reader := func(buffer io.Reader) (io.ReadCloser, error) {
		return gzip.NewReader(buffer)
	}
	writer := func(buffer io.Writer) (io.WriteCloser, error) {
		return gzip.NewWriter(buffer), nil
	}
	return &compression{writer: writer, reader: reader, ContentEncoding: "gzip"}
}

func Zlib() *compression {
	reader := func(buffer io.Reader) (io.ReadCloser, error) {
		return zlib.NewReader(buffer)
	}
	writer := func(buffer io.Writer) (io.WriteCloser, error) {
		return zlib.NewWriter(buffer), nil
	}
	return &compression{writer: writer, reader: reader, ContentEncoding: "deflate"}
}
