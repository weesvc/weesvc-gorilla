package api

import (
	"net/http"
	"strings"
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

func TestAddressForRequestProxied(t *testing.T) {
	testCases := []struct {
		remoteAddr string
		expected   string
	}{
		{
			remoteAddr: "12.34.56.78:8080,23.45.67.89:8080",
			expected:   "23.45.67.89",
		},
	}
	fixture := API{Config: &Config{ProxyCount: 1}}
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
	// For convenience, we're setting the address(es) as the forwarded for as well as remoteAddr
	header := http.Header{}
	remote := ""
	if strings.Contains(address, ",") {
		// List of addresses...treat as proxied
		header.Add("X-Forwarded-For", address)
	} else {
		remote = address
	}
	return http.Request{RemoteAddr: remote, Header: header}
}
