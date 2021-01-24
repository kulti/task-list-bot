package main

import (
	"context"
	"net/http"
)

type contextRoundTripper struct {
	rt  http.RoundTripper
	ctx context.Context
}

func (r contextRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return r.rt.RoundTrip(req.WithContext(r.ctx))
}
