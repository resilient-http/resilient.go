package failover

import (
	"net/http"
	"time"
)

type Options struct {
	Retries    int
	RetryDelay time.Duration
}

func Failover(o Options) func(*http.Request, *http.Response, error) bool {
	return func(req *http.Request, res *http.Response, err error) bool {
		var retry bool

		defer (func() {
			if retry == false {
				return
			}

			o.Retries -= 1
			if o.Retries == 0 {
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
