# Failover strategy

Build-in Resilient.go strategy to handle common HTTP failure/error scenarios, retrying the request if necessary.

## Installation

```
go get gopkg.in/resilient-http/resilient.go.v0/strategies/failover
```

## Options

- Retries `int` - Number of retries per request
- RetryDelay `time.Duration` - Retry attempt delay. Default based on `Retry-After` response header

## Retry policy

- 500-600 status codes
- Socket error
- DNS error
- Timeout exceeded

## Usage

```go
package main

import (
  "time"
  "fmt"
  "gopkg.in/resilient-http/resilient.go.v0"
  "gopkg.in/resilient-http/resilient.go.v0/strategies/failover"
)

func main() {
  r := resilient.New()
  
  r.UseServers([]string{
    "http://foo", "http://bar"
  })

  r.UseStrategy(failover.New(failover.Options{
    Retries: 3,
    RetryDelay: 100 * time.Millisecond,
  }))

  res, err := r.Get("/foo")
  if err != nil {
    fmt.Printf("Error: %#v", err)
  }
  fmt.Printf("Status: %d", res.StatusCode)
}
```
