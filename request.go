package resilient

import (
	"errors"
	"gopkg.in/resilient-http/resilient.go.v0/client"
	"gopkg.in/resilient-http/resilient.go.v0/strategies"
	"net/http"
	"net/url"
	"time"
)

type Request struct {
	cancelled bool
	attempts  []*Request
	Time      time.Time
	Error     error
	Response  *http.Response
	Request   *http.Request
	Client    *client.Client
	resilient *Resilient
}

func NewRequest(r *Resilient, o *client.Options) (*http.Response, error) {
	var err error

	err = r.middlewares.DispatchOut(o)
	if err != nil {
		return nil, err
	}

	client, err := client.New(o)
	if err != nil {
		return nil, err
	}

	request := &Request{
		Client:     client.HttpClient,
		Request:    client.Request,
		strategies: initStrategies(r),
		resilient:  r,
	}

	return r.Send(c, strategies, *r)
}

func (r *Request) Attemps() []*Request {
	return r.attempts
}

func (r *Request) Cancel() {
	r.Client.Transport.CancelRequest(r.Request)
	r.cancelled = true
}

func (r *Request) Send() (*http.Response, error) {
	req := c.Request
	hasServers := len(re.Servers) > 0

	if hasServers {
		err := buildServerUrl(req, r)
		if err != nil {
			return nil, err
		}
	}

	if req.URL.Host == "" {
		return nil, errors.New("Missing server URL")
	}

	res, err := c.Do()

	err = re.middlewares.DispatchIn(req, res, err)
	if err != nil {
		return res, err
	}

	var retry bool
	for _, strategy := range strategies {
		if strategy(req, res, err) {
			retry = true
		}
	}

	if retry {
		if hasServers && len(r.Servers) > 1 {
			r.Servers = r.Servers[1:]
		}
		return r.Send()
	}

	return res, err
}

func initStrategies(r *Resilient) []strategies.Handler {
	strategies := make([]strategies.Handler, len(r.strategies))
	for i, strategy := range r.strategies {
		strategies[i] = strategy()
	}
	return strategies
}

func buildServerUrl(req *http.Request, r Resilient) error {
	serverURL, err := url.Parse(r.Servers[0])
	if err != nil {
		return err
	}

	req.URL.Host = serverURL.Host
	req.URL.User = serverURL.User
	req.URL.Scheme = serverURL.Scheme
	req.URL.Opaque = serverURL.Opaque
	req.URL.Path = serverURL.Path + req.URL.Path

	return nil
}
