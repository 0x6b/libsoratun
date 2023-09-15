package internal

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// UnifiedEndpointHTTPClient is an HTTP client that can be used to communicate with SORACOM Unified Endpoint.
type UnifiedEndpointHTTPClient struct {
	HttpClient *http.Client
	Endpoint   *url.URL
	Headers    []string
	Verbose    bool
}

type Params struct {
	Path    string
	Body    io.Reader
	Method  string
	Headers []string
}

func (c *UnifiedEndpointHTTPClient) MakeRequest(params *Params) (*http.Request, error) {
	req, err := http.NewRequest(
		params.Method,
		fmt.Sprintf("%s://%s:%s/%s",
			c.Endpoint.Scheme,
			c.Endpoint.Hostname(),
			c.Endpoint.Port(),
			params.Path,
		),
		params.Body,
	)
	if err != nil {
		return nil, err
	}

	for _, h := range params.Headers {
		header := strings.SplitN(h, ":", 2)
		if len(header) == 2 {
			req.Header.Set(strings.TrimSpace(header[0]), strings.TrimSpace(header[1]))
		}
	}

	return req, nil
}

func (c *UnifiedEndpointHTTPClient) DoRequest(req *http.Request) (*http.Response, error) {
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		defer func() {
			err := res.Body.Close()
			if err != nil {
				fmt.Println("failed to close response", err)
			}
		}()
		r, _ := io.ReadAll(res.Body)
		return res, fmt.Errorf("%s: %s %s: %s", res.Status, req.Method, req.URL, r)
	}

	return res, nil
}
