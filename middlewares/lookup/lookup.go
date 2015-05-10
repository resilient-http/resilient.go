package lookup

import (
	"fmt"
	"gopkg.in/resilient-http/resilient.go.v0"
	"gopkg.in/resilient-http/resilient.go.v0/client"
	. "gopkg.in/resilient-http/resilient.go.v0/middlewares"
	"net/http"
	"time"
)

var Error = func(msg string, code int) error {
	return NewMiddlewareError("lookup", msg, code)
}

type Options struct {
	Servers []string
	Refresh time.Duration
	Retries int
}

type Discovery struct {
	servers       []string
	lookupServers []string
	resilient     *resilient.Resilient
	updated       *time.Time
	options       Options
}

func New(opts Options) func(resilient *resilient.Resilient) {
	return func(resilient *resilient.Resilient) (*Discovery, error) {

		if len(opts.Servers) == 0 {
			return nil, Error("Missing servers")
		}

		if opts.Refresh == 0 {
			opts.Refresh = 5 * time.Minute
		}

		return &Discovery{
			servers:   opts.Servers,
			resilient: resilient,
			options:   opts,
		}, nil
	}
}

func (m *Discovery) Out(o *client.Options) error {
	if len(servers) == 0 || m.updated-time.Now() > m.options.Refresh {

	}

	m.resilient.UseServers(m.servers)
	return nil
}

func (m *Discovery) In(*http.Request, *http.Response, error) error {
	return nil
}

func (m *Discovery) Discover() error {
	res, err := http.NewRequest("GET", m.servers[0])
	if err != nil {
		return err
	}

	fmt.Printf("Response: %d", res.StatusCode)
}
