package failover

import (
	"net/http"
	"time"
)

type Options struct {
	Retries    int
	RetryDelay time.Duration
}

func New(o Options) {
	return func() func(*http.Request, *http.Response, error) bool {
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
				return retry
			}

			return retry
		}
	}
}
