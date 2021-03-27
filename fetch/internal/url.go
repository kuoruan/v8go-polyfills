package internal

import (
	"fmt"
	"net/url"
	"strings"
)

func ParseURL(rawURL string) (*url.URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("url '%s' is not valid, %w", rawURL, err)
	}

	switch u.Scheme {
	case "http", "https":
	case "":
		if !strings.HasPrefix(u.Path, "/") {
			return nil, fmt.Errorf("unsupported relatve path %s", u.Path)
		}
	default:
		return nil, fmt.Errorf("unsupported scheme %s", u.Scheme)
	}

	return u, nil
}
