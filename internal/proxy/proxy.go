package proxy

import (
	"errors"
	"net/http"

	"net/http/httputil"
	"net/url"
	"time"

	"github.com/cronenberg64/ergo-api/internal/circuitbreaker"
)

type CircuitBreakerTransport struct {
	RoundTripper http.RoundTripper
	cb           *circuitbreaker.CircuitBreaker
}

func (t *CircuitBreakerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	err := t.cb.Execute(func() error {
		var err error
		resp, err = t.RoundTripper.RoundTrip(req)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 500 {
			return errors.New("server error")
		}
		return nil
	})
	return resp, err
}

func NewProxy(target string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	
	// Initialize Circuit Breaker (5 failures, 10s reset)
	cb := circuitbreaker.NewCircuitBreaker(5, 10*time.Second)
	
	proxy.Transport = &CircuitBreakerTransport{
		RoundTripper: http.DefaultTransport,
		cb:           cb,
	}

	return proxy, nil
}
