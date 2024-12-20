package httpx

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/pechorka/gostdlib/pkg/errs"
)

const (
	defaultTimeout               = 10 * time.Second
	defaultKeepAlive             = 30 * time.Second
	defaultMaxIdleConns          = 100
	defaultMaxIdleConnsPerHost   = 100
	defaultIdleConnTimeout       = 90 * time.Second
	defaultTLSHandshakeTimeout   = 10 * time.Second
	defaultExpectContinueTimeout = 1 * time.Second
)

// Client is a wrapper around http.Client with sensible defaults
type Client struct {
	*http.Client
}

// Option is a function that configures a Client
type Option func(*http.Client)

// NewClient creates a new Client with sensible defaults
func NewClient(opts ...Option) *Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   defaultTimeout,
			KeepAlive: defaultKeepAlive,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          defaultMaxIdleConns,
		MaxIdleConnsPerHost:   defaultMaxIdleConnsPerHost,
		IdleConnTimeout:       defaultIdleConnTimeout,
		TLSHandshakeTimeout:   defaultTLSHandshakeTimeout,
		ExpectContinueTimeout: defaultExpectContinueTimeout,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   defaultTimeout,
	}

	// Apply any custom options
	for _, opt := range opts {
		opt(client)
	}

	return &Client{
		Client: client,
	}
}

// Do wraps http.Client.Do
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.Client.Do(req)
}

var defaultClient = NewClient()

// Do makes a request with the default client
func Do(req *http.Request) (*http.Response, error) {
	return defaultClient.Do(req)
}

// GetJSON makes a GET request and decodes the response as JSON
func GetJSON[Resp any](ctx context.Context, url string) (resp Resp, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return resp, errs.Wrap(err, "failed to create request")
	}

	httpResp, err := defaultClient.Do(req)
	if err != nil {
		return resp, errs.Wrap(err, "failed to do request")
	}
	if httpResp.StatusCode >= 400 {
		body, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return resp, errs.Wrapf(err, "request failed with status code %d", httpResp.StatusCode)
		}
		return resp, errs.Newf("request failed with status code %d: %s", httpResp.StatusCode, string(body))
	}

	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		err = errs.Wrap(err, "failed to decode response")
		if cerr := httpResp.Body.Close(); cerr != nil {
			err = errs.Join(err, errs.Wrap(cerr, "failed to close response body"))
		}
		return resp, err
	}

	err = httpResp.Body.Close()
	if err != nil {
		return resp, errs.Wrap(err, "failed to close response body")
	}

	return resp, nil
}