// Copyright © 2025 Ping Identity Corporation

package verify

import (
	"fmt"
	"testing"
)

var (
	urlsWithoutProtocol = []string{
		"www.my-test-domain.%s",
		"www.my-test-domain.%s:3000",
		"test1234.my-test-domain.%s",
		"test1234.www.my-test-domain.%s",
		"test-1234.www.subdomain.my-test-domain.%s",
		"my-test-domain.%s",
		"my-test-domain.%s/",
		"my-test-domain.%s:3000/",
		"my-test-domain.%s/path",
		"my-test-domain.%s:3000/path",
		"my-test-domain.%s/path/",
		"my-test-domain.%s/path?query",
		"my-test-domain.%s:3000/path?query",
		"my-test-domain.%s/path?query=value",
		"my-test-domain.%s:3000/path?query=value&query2=value2",
		"my-test-domain.%s/path?query=value&query2=value2",
		"my-test-domain.%s/path?query=value&query2=value2#fragment",
		"my-test-domain.%s:3000/path?query=value&query2=value2#fragment",
		"my-test-domain.%s/path-1.2.3",
		"my-test-domain.%s/path-1.2.3/",
		"my-test-domain.%s:3000/path-1.2.3/",
		"my-test-domain.%s/path-1.2.3/path-1.2.3/path-1.2.3/path-1.2.3/path-1.2.3/path-1.2.3/path-1.2.3/",
	}

	tlds = []string{
		"com",
		"org",
		"net",
		"gov",
		"gov.uk",
		"edu",
		"mil",
		"co",
		"co.uk",
		"com.au",
		"com.br",
	}
)

func TestRegex_Localhost(t *testing.T) {

	localhostURLsWithoutProtocol := []string{
		"localhost",
		"localhost:3000",
		"localhost/",
		"localhost:3000/",
		"localhost/path",
		"localhost:3000/path",
		"localhost/path/",
		"localhost/path?query",
		"localhost:3000/path?query",
		"localhost/path?query=value",
		"localhost:3000/path?query=value&query2=value2",
		"localhost/path?query=value&query2=value2",
		"localhost/path?query=value&query2=value2#fragment",
		"localhost:3000/path?query=value&query2=value2#fragment",
		"localhost/path-1.2.3",
		"localhost/path-1.2.3/",
		"localhost:3000/path-1.2.3/",
		"localhost/path-1.2.3/path-1.2.3/path-1.2.3/path-1.2.3/path-1.2.3/path-1.2.3/path-1.2.3/",
	}
	protocols := []string{
		"https",
		"http",
	}

	for _, protocol := range protocols {
		for _, url := range localhostURLsWithoutProtocol {
			parsedUrl := fmt.Sprintf("%s://%s", protocol, url)
			t.Run(parsedUrl, func(t *testing.T) {
				if !IsURLWithHTTPorHTTPS.MatchString(parsedUrl) {
					t.Errorf("IsURLWithHTTPorHTTPS failed to match %q", parsedUrl)
				}
			})
		}
	}
}

func TestRegex_IsURLWithHTTPS_Positive(t *testing.T) {

	protocols := []string{
		"https",
	}

	for _, protocol := range protocols {
		for _, tld := range tlds {
			for _, url := range urlsWithoutProtocol {
				parsedUrl := fmt.Sprintf("%s://%s", protocol, fmt.Sprintf(url, tld))
				t.Run(parsedUrl, func(t *testing.T) {
					if !IsURLWithHTTPS.MatchString(parsedUrl) {
						t.Errorf("IsURLWithHTTPS failed to match %q", parsedUrl)
					}
				})
			}
		}
	}
}

func TestRegex_IsURLWithHTTPS_Negative(t *testing.T) {

	protocols := []string{
		"http",
		"ftp",
		"ftps",
		"ssh",
		"telnet",
		"smtp",
		"imap",
	}

	for _, protocol := range protocols {
		for _, tld := range tlds {
			for _, url := range urlsWithoutProtocol {
				parsedUrl := fmt.Sprintf("%s://%s", protocol, fmt.Sprintf(url, tld))
				t.Run(parsedUrl, func(t *testing.T) {
					if IsURLWithHTTPS.MatchString(parsedUrl) {
						t.Errorf("IsURLWithHTTPS erroneous match %q", parsedUrl)
					}
				})
			}
		}
	}
}

func TestRegex_IsURLWithHTTPorHTTPS_Positive(t *testing.T) {

	protocols := []string{
		"https",
		"http",
	}

	for _, protocol := range protocols {
		for _, tld := range tlds {
			for _, url := range urlsWithoutProtocol {
				parsedUrl := fmt.Sprintf("%s://%s", protocol, fmt.Sprintf(url, tld))
				t.Run(parsedUrl, func(t *testing.T) {
					if !IsURLWithHTTPorHTTPS.MatchString(parsedUrl) {
						t.Errorf("IsURLWithHTTPorHTTPS failed to match %q", parsedUrl)
					}
				})
			}
		}
	}
}

func TestRegex_IsURLWithHTTPorHTTPS_Negative(t *testing.T) {

	protocols := []string{
		"ftp",
		"ftps",
		"ssh",
		"telnet",
		"smtp",
		"imap",
	}

	for _, protocol := range protocols {
		for _, tld := range tlds {
			for _, url := range urlsWithoutProtocol {
				parsedUrl := fmt.Sprintf("%s://%s", protocol, fmt.Sprintf(url, tld))
				t.Run(parsedUrl, func(t *testing.T) {
					if IsURLWithHTTPorHTTPS.MatchString(parsedUrl) {
						t.Errorf("IsURLWithHTTPorHTTPS erroneous match %q", parsedUrl)
					}
				})
			}
		}
	}
}
