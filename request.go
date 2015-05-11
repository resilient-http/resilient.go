package resilient

import (
	"errors"
	"gopkg.in/resilient-http/resilient.go.v0/client"
	"net/http"
	"net/url"
)

type Request struct {
	Requests []string
}

func NewRequest(r *Resilient, o *client.Options) (*http.Response, error) {
	var err error

	err = r.middlewares.DispatchOut(o)
	if err != nil {
		return nil, err
	}

	c, err := client.New(o)
	if err != nil {
		return nil, err
	}

	strategies := initStrategies(r)
	return send(c, strategies, *r)
}

func initStrategies(r *Resilient) []StrategyHandler {
	strategies := make([]StrategyHandler, len(r.strategies))
	for i, strategy := range r.strategies {
		strategies[i] = strategy()
	}
	return strategies
}

func send(c *client.Client, strategies []StrategyHandler, r Resilient) (*http.Response, error) {
	req := c.Request
	hasServers := len(r.Servers) > 0

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

	err = r.middlewares.DispatchIn(req, res, err)
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
		return send(c, strategies, r)
	}

	return res, err
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
