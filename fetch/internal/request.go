package internal

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	RequestRedirectError  = "error"
	RequestRedirectFollow = "follow"
	RequestRedirectManual = "manual"
)

/*
 RequestInit is the fetch API defined object.
 Only supports raw request now.
*/
type RequestInit struct {
	Body     string            `json:"body"`
	Headers  map[string]string `json:"headers"`
	Method   string            `json:"method"`
	Redirect string            `json:"redirect"`
}

/*
 Request is the request object used by fetch
*/
type Request struct {
	Body     string
	Method   string
	Redirect string

	Header     http.Header
	URL        *url.URL
	RemoteAddr string
}

/*
 parse and check the request URL, return *url.URL
*/
func ParseRequestURL(rawURL string) (*url.URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("url '%s' is not valid, %w", rawURL, err)
	}

	/**
	 * Check the scheme, we only support http and https at this time
	 */
	switch u.Scheme {
	case "http", "https":
	case "": // then scheme is empty, it's a local request
		if !strings.HasPrefix(u.Path, "/") {
			return nil, fmt.Errorf("unsupported relatve path %s", u.Path)
		}
	default:
		return nil, fmt.Errorf("unsupported scheme %s", u.Scheme)
	}

	return u, nil
}
