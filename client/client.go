package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Client struct {
	HttpClient *http.Client
	Request    *http.Request
	Options    *Options
}

func New(o *Options) (*Client, error) {
	if o.Client == nil {
		jar, _ := cookiejar.New(&cookiejar.Options{})
		o.Client = &http.Client{Jar: jar}
	}

	req, _ := http.NewRequest(o.Method, o.URL, o.Body)

	for k, v := range o.Header {
		req.Header.Set(k, v)
	}

	// Add all querystring from Query func
	q := req.URL.Query()
	for k, v := range o.QueryData {
		for _, vv := range v {
			q.Add(k, vv)
		}
	}
	req.URL.RawQuery = q.Encode()

	// Add cookies
	for _, cookie := range o.Cookies {
		req.AddCookie(cookie)
	}

	if o.Timeout != 0 {
		o.Client.Timeout = o.Timeout
	}

	// Set Transport
	transport := &http.Transport{}

	if o.DialTimeout > 0 {
		transport.Dial = func(network, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(network, addr, o.Timeout)
			if err != nil {
				return nil, err
			}
			conn.SetDeadline(time.Now().Add(o.Timeout))
			return conn, nil
		}
	}

	if o.DisableTLS {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if o.ProxyURL != "" {
		parsedProxyUrl, err := url.Parse(o.ProxyURL)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(parsedProxyUrl)
	}

	o.Client.Transport = transport

	return &Client{
		Options:    o,
		Request:    req,
		HttpClient: o.Client,
	}, nil
}

func (c *Client) Do() (*http.Response, error) {
	return c.HttpClient.Do(c.Request)
}
