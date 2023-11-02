// Package fetch provides a cache client that can fetch data from a URL.
package fetch

import (
	"context"
	"fmt"
	"net/http"

	"git.sr.ht/~jamesponddotco/httpx-go"
	"git.sr.ht/~jamesponddotco/sitred"
	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"golang.org/x/time/rate"
)

// ErrFetchData is returned when the client fails to fetch data from a URL.
const ErrFetchData xerrors.Error = "failed to fetch data"

// Client represents a client that can fetch data from a URL.
type Client struct {
	// httpc is the underlying HTTP client used to fetch data.
	httpc *httpx.Client
}

// New creates a new client that can fetch data from a URL.
func New(serviceName, serviceContact string) *Client {
	return &Client{
		httpc: &httpx.Client{
			RateLimiter: rate.NewLimiter(rate.Limit(2), 1),
			RetryPolicy: httpx.DefaultRetryPolicy(),
			UserAgent: &httpx.UserAgent{
				Token:   serviceName,
				Version: sitred.Version,
				Comment: []string{serviceContact},
			},
			Cache: nil,
		},
	}
}

// Remote fetches data from a URL and returns it as a raw http.Response.
func (c *Client) Remote(ctx context.Context, uri string) (*http.Response, error) {
	resp, err := c.httpc.Get(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFetchData, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s", ErrFetchData, resp.Status)
	}

	return resp, nil
}
