package failover

import (
	"gopkg.in/resilient-http/resilient.go.v0/strategies"
	"net/http"
	"time"
)

type Options struct {
	Retries    int
	RetryDelay time.Duration
}

func New(o Options) func() strategies.Handler {
	return func() strategies.Handler {
		retries := o.Retries

		return func(req *http.Request, res *http.Response, err error) bool {
			var retry bool

			defer (func() {
				if retry == false {
					return
				}

				retries -= 1
				if retries <= 0 {
					retry = false
					return
				}

				if o.RetryDelay > 0 {
					time.Sleep(o.RetryDelay)
				}
			})()

			if err != nil {
				retry = true
				return retry
			}

			if res.StatusCode >= 500 {
				retry = true
			}

			return retry
		}
	}
}
