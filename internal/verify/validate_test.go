// Copyright © 2026 Ping Identity Corporation

package verify

import (
	"testing"
)

func TestIsURLWithHTTPS(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		value   string
		wantErr bool
	}{
		// Positive cases
		{name: "plain domain", value: "https://www.my-test-domain.com", wantErr: false},
		{name: "with port", value: "https://www.my-test-domain.com:3000", wantErr: false},
		{name: "with path", value: "https://my-test-domain.com/path", wantErr: false},
		{name: "with path and trailing slash", value: "https://my-test-domain.com/path/", wantErr: false},
		{name: "with query", value: "https://my-test-domain.com/path?query=value", wantErr: false},
		{name: "with query and fragment", value: "https://my-test-domain.com/path?query=value&query2=value2#fragment", wantErr: false},
		{name: "percent-encoded path", value: "https://example.org/a/path/with/a/strange%40character", wantErr: false},
		{name: "encoded percent in path", value: "https://example.org/a/path/with/a/strange%2540character", wantErr: false},
		{name: "dotted path segments", value: "https://my-test-domain.com/path-1.2.3/path-1.2.3/path-1.2.3/", wantErr: false},
		{name: "idn host", value: "https://例え.jp/path", wantErr: false},
		{name: "uppercase scheme", value: "HTTPS://www.my-test-domain.com", wantErr: false},
		{name: "sub domains", value: "https://test-1234.www.subdomain.my-test-domain.co.uk", wantErr: false},
		{name: "ipv4 host", value: "https://192.168.1.1/path", wantErr: false},
		{name: "ipv6 host", value: "https://[::1]:8443/path", wantErr: false},
		{name: "localhost", value: "https://localhost/path", wantErr: false},
		{name: "localhost with port", value: "https://localhost:3000/path", wantErr: false},

		// Negative cases
		{name: "http scheme", value: "http://www.my-test-domain.com", wantErr: true},
		{name: "ftp scheme", value: "ftp://www.my-test-domain.com", wantErr: true},
		{name: "no scheme", value: "www.my-test-domain.com", wantErr: true},
		{name: "no host", value: "https:///path", wantErr: true},
		{name: "no host, port only", value: "https://:3000/path", wantErr: true},
		{name: "no host, userinfo and port only", value: "https://user@:3000/path", wantErr: true},
		{name: "empty", value: "", wantErr: true},
		{name: "invalid port", value: "https://example.com:port/path", wantErr: true},
		{name: "invalid percent escape", value: "https://example.com/%zz", wantErr: true},
		{name: "truncated percent escape", value: "https://example.com/%2", wantErr: true},
		{name: "percent escape in host", value: "https://exa%mple.com/", wantErr: true},
		{name: "control character", value: "https://example.com/path\tx?q", wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := IsURLWithHTTPS(tc.value)
			if (err != nil) != tc.wantErr {
				t.Errorf("IsURLWithHTTPS(%q) error = %v, wantErr %v", tc.value, err, tc.wantErr)
			}
		})
	}
}

func TestIsURLWithHTTPorHTTPS(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		value   string
		wantErr bool
	}{
		// Positive cases
		{name: "https scheme", value: "https://www.my-test-domain.com", wantErr: false},
		{name: "http scheme", value: "http://www.my-test-domain.com", wantErr: false},
		{name: "http with port and path", value: "http://www.my-test-domain.com:3000/path", wantErr: false},
		{name: "percent-encoded path", value: "http://example.org/a/path/with/a/strange%40character", wantErr: false},
		{name: "localhost", value: "http://localhost:3000/path?query=value", wantErr: false},
		{name: "percent-encoded path with dotted segments", value: "http://my-test-domain.com/path-1.2.3/with-%25-encoding/", wantErr: false},

		// Negative cases
		{name: "ftp scheme", value: "ftp://www.my-test-domain.com", wantErr: true},
		{name: "no scheme", value: "www.my-test-domain.com", wantErr: true},
		{name: "no host", value: "http:///path", wantErr: true},
		{name: "empty", value: "", wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := IsURLWithHTTPorHTTPS(tc.value)
			if (err != nil) != tc.wantErr {
				t.Errorf("IsURLWithHTTPorHTTPS(%q) error = %v, wantErr %v", tc.value, err, tc.wantErr)
			}
		})
	}
}
