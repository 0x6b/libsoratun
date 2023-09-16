package main

import (
	"C"
	"fmt"
	"io"

	"github.com/0x6b/libsoratun/libsoratun"
)

func main() {}

// Send sends a request to the unified endpoint with given Config and HTTP headers.
// Args:
//
// - `configJson`: JSON string of soratun config.
// - `method`: HTTP method of the request. Only GET or POST is supported.
// - `path`: Path of the request.
// - `body`: Body of the request.
//
//export Send
func Send(configJson, method, path, body *C.char) *C.char {
	config, err := libsoratun.NewConfig([]byte(C.GoString(configJson)))
	if err != nil {
		return nil
	}

	c, err := libsoratun.NewUnifiedEndpointHTTPClient(*config)
	if err != nil {
		return nil
	}

	req, err := c.MakeRequest(&libsoratun.Params{
		Method: C.GoString(method),
		Path:   C.GoString(path),
		Body:   C.GoString(body),
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
