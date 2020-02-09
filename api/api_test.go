package api

import (
	"net/http"
	"testing"
)

func TestAddressForRequest(t *testing.T) {
	testCases := []struct {
		remoteAddr string
		expected   string
	}{
		{
			remoteAddr: "localhost:8080",
			expected:   "localhost",
		},
		{
			remoteAddr: "[::1]:8080",
			expected:   "[::1]",
		},
	}
	fixture := API{Config: &Config{}}
	for _, tc := range testCases {
		tc := tc // pin
		t.Run(tc.remoteAddr, func(t *testing.T) {
			request := mockRequest(tc.remoteAddr)
			actual := fixture.addressForRequest(&request)
			if actual != tc.expected {
				t.Fatalf("expected %q, but got %q", tc.expected, actual)
			}
		})
	}
}

func mockRequest(address string) http.Request {
	return http.Request{RemoteAddr: address}
}
