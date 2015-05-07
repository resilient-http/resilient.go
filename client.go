package resilient

import (
	"net/http"
)

type Middleware func(params interface{}) func(next)

type Client struct {
	*http.Client
	middlewares []Middleware
}

func (c *Client) Use(m Middleware) {

}

func (c *Client) Get() {

}

func (c *Client) Post() {

}
