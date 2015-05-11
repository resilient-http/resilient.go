package lookup

import (
	"encoding/json"
	"gopkg.in/resilient-http/resilient.go.v0/client"
	"gopkg.in/resilient-http/resilient.go.v0/middlewares"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	MISSING_SERVERS = iota
)

var Error = func(msg string, code int) error {
	return middlewares.NewMiddlewareError("lookup", msg, code)
}

type Options struct {
	Servers      []string
	Refresh      time.Duration
	Timeout      time.Duration
	Retries      int
	Background   bool
	SelfLookup   bool
	SelfRefresh  time.Duration
}

type Discovery struct {
	servers       []string
	lookupServers []string
	updated       int64
	options       Options
	resilient     middlewares.Resilient
}

func New(opts Options) middlewares.Handler {
	return func(r middlewares.Resilient) (middlewares.Middleware, error) {

		if len(opts.Servers) == 0 {
			return nil, Error("Missing servers", MISSING_SERVERS)
		}

		if opts.Refresh == 0 {
			opts.Refresh = 5 * time.Minute
		}

		return &Discovery{
			servers:   opts.Servers,
			options:   opts,
			updated:   0,
			resilient: r,
		}, nil
	}
}

func (m *Discovery) Out(o *client.Options) error {
	if len(m.servers) > 0 && (time.Now().UnixNano()-m.updated) < int64(m.options.Refresh) {
		return nil
	}

	err := m.Discover()
	if err != nil {
		return err
	}

	m.resilient.UseServers(m.servers)
	return nil
}

func (m *Discovery) In(*http.Request, *http.Response, error) error {
	return nil
}

type URLs struct {
	urls []string
}

func (m *Discovery) Discover() error {
	req, err := client.New(&client.Options{
		URL: m.servers[0],
	})
	if err != nil {
		return err
	}

	res, err := req.Do()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	t := []string{}
	err = json.Unmarshal(body, &t)
	if err != nil {
		return err
	}

	m.resilient.UseServers(t)

	return nil
}
