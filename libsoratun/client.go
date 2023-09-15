package libsoratun

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

const (
	UnifiedEndpointHostname = "100.127.69.42"
	UnifiedEndpointPort     = 80
)

// UnifiedEndpointHTTPClient is an HTTP client that can be used to communicate with SORACOM Unified Endpoint.
type UnifiedEndpointHTTPClient struct {
	HttpClient *http.Client
	endpoint   *url.URL
	headers    []string
}

type Params struct {
	Path   string
	Body   io.Reader
	Method string
}

func NewUnifiedEndpointHTTPClient(context func(ctx context.Context, network, addr string) (net.Conn, error)) *UnifiedEndpointHTTPClient {
	endpoint, _ := url.Parse(fmt.Sprintf("http://%s:%d", UnifiedEndpointHostname, UnifiedEndpointPort))
	c := &UnifiedEndpointHTTPClient{
		HttpClient: &http.Client{
			Transport: &http.Transport{
				DialContext: context,
			},
		},
		endpoint: endpoint,
		headers:  []string{"User-Agent: libsoratun/0.0.1"},
	}

	return c
}

func (c *UnifiedEndpointHTTPClient) MakeRequest(params *Params) (*http.Request, error) {
	req, err := http.NewRequest(
		params.Method,
		fmt.Sprintf("%s://%s:%s/%s",
			c.endpoint.Scheme,
			c.endpoint.Hostname(),
			c.endpoint.Port(),
			params.Path,
		),
		params.Body,
	)
	if err != nil {
		return nil, err
	}

	for _, h := range c.headers {
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
