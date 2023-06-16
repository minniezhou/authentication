package roundtrip

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func httpClientWithRoundTripper(statusCode int, response string) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: statusCode,
				Body:       ioutil.NopCloser(bytes.NewBufferString(response)),
			}
		}),
	}
}

// in mock can create httpClientWithRoundTripper as a http.client and use above client and make client.Do call to mock up the
// round trip function
