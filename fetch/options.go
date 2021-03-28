package fetch

import (
	"net/http"
	"net/url"
)

type UserAgentProvider interface {
	GetUserAgent(u *url.URL) string
}

type UserAgentProviderFunc func(u *url.URL) string

func (f UserAgentProviderFunc) GetUserAgent(u *url.URL) string {
	return f(u)
}

type Option interface {
	apply(ft *fetcher)
}

type optionFunc func(ft *fetcher)

func (f optionFunc) apply(ft *fetcher) {
	f(ft)
}

func WithLocalHandler(handler http.Handler) Option {
	return optionFunc(func(ft *fetcher) {
		ft.LocalHandler = handler
	})
}

func WithUserAgentProvider(provider UserAgentProvider) Option {
	return optionFunc(func(ft *fetcher) {
		ft.UserAgentProvider = provider
	})
}

func WithAddrLocal(addr string) Option {
	return optionFunc(func(ft *fetcher) {
		ft.AddrLocal = addr
	})
}
