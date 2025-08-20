package unstructured

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// Client represents an HTTP client for interacting with the Unstructured.io API.
// It handles authentication, request formatting, and response parsing.
type Client struct {
	hc       *http.Client
	endpoint *url.URL
}

// Option is a function that configures a Client instance.
// Options are used to set the endpoint URL and API key during client creation.
type Option func(*Client) error

// WithEndpoint returns an Option that sets the API endpoint URL.
// The endpoint should be the base URL for the Unstructured.io API, including the base path like "/api/v1".
// Without this option, the client will default to `https://platform.unstructuredapp.io/api/v1`.
func WithEndpoint(endpoint string) Option {
	return func(c *Client) error {
		u, err := url.Parse(endpoint)
		if err != nil {
			return fmt.Errorf("failed to parse endpoint URL: %w", err)
		}

		c.endpoint = u

		return nil
	}
}

// WithKey returns an Option that sets the API key for authentication.
// The API key is used to authenticate all requests to the Unstructured.io API.
// This is accomplished using a [http.RoundTripper] that sets the key as the value of the `Unstructured-API-Key` header on all requests.
func WithKey(key string) Option {
	return func(c *Client) error {
		if b, ok := c.hc.Transport.(*bearer); ok {
			b.key = key
			return nil
		}

		c.hc.Transport = &bearer{
			key: key,
			rt:  cmp.Or(c.hc.Transport, http.DefaultTransport),
		}

		return nil
	}
}

// WithClient returns an Option that sets the HTTP client to use for requests.
// If no client is provided, the client will default to [http.DefaultClient].
func WithClient(hc *http.Client) Option {
	return func(c *Client) error {
		c.hc = hc
		return nil
	}
}

// New creates a new Client instance with the provided options.
// If the `UNSTRUCTURED_API_KEY` environment variable is set, it will be used as the API key for authentication.
// If the `UNSTRUCTURED_API_URL` environment variable is set to a valid URL, it will be used as the base URL for the Unstructured.io API.
// If no endpoint option is given via options or environment variables, the endpoint will default to the Unstructured.io platform at `https://platform.unstructuredapp.io/api/v1`.
// In order to configure the client properly, an API key must be provided via [WithKey] or the `UNSTRUCTURED_API_KEY` environment variable.
func New(opts ...Option) (*Client, error) {
	c := Client{
		hc: http.DefaultClient,
		endpoint: &url.URL{
			Scheme: "https",
			Host:   "platform.unstructuredapp.io",
			Path:   "/api/v1",
		},
	}

	// attempt to set endpoint from environment variable
	if v := os.Getenv("UNSTRUCTURED_API_URL"); v != "" {
		if u, err := url.Parse(v); err == nil {
			c.endpoint = u
		}
	}

	// attempt to set API key from environment variable
	if v := os.Getenv("UNSTRUCTURED_API_KEY"); v != "" {
		c.hc.Transport = &bearer{
			key: v,
			rt:  cmp.Or(c.hc.Transport, http.DefaultTransport),
		}
	}

	// apply options
	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		// Handle 422 validation errors specifically
		if resp.StatusCode == http.StatusUnprocessableEntity {
			var validationErr HTTPValidationError
			if err := json.Unmarshal(body, &validationErr); err == nil {
				return &APIError{
					Code: resp.StatusCode,
					Err:  &validationErr,
				}
			}
		}

		return &APIError{
			Code: resp.StatusCode,
			Err:  errors.New(string(body)),
		}
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
