package main

import (
	"C"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/0x6b/libsoratun/libsoratun"
)

func main() {
}

// Send sends a request to the unified endpoint with given Config and HTTP headers.
//
//export Send
func Send(configJson *C.char, method, path, body *C.char) *C.char {
	config, err := libsoratun.NewConfig([]byte(C.GoString(configJson)))
	if err != nil {
		return nil
	}

	t, err := libsoratun.NewTunnel(config)
	if err != nil {
		return nil
	}

	c := libsoratun.NewUnifiedEndpointHTTPClient(t.DialContext)

	m := C.GoString(method)
	if !(m == http.MethodGet || m == http.MethodPost) {
		return nil
	}

	req, err := c.MakeRequest(&libsoratun.Params{
		Path:   strings.TrimPrefix(C.GoString(path), "/"),
		Body:   strings.NewReader(C.GoString(body)),
		Method: m,
	})
	if err != nil {
		return nil
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil
	}

	response, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("failed to read response from Unified Endpoint: %v\n", err)
		panic(err)
	}

	return C.CString(string(response))
}
