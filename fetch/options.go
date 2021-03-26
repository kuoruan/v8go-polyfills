package fetch

import (
	"net/http"
)

type Option interface {
	apply(ft *fetcher)
}

type funcOption func(ft *fetcher)

func (f funcOption) apply(ft *fetcher) {
	f(ft)
}

func (f funcOption) Polyfill() string {
	return "fetch"
}

func WithLocalHandler(handler http.Handler) Option {
	return funcOption(func(ft *fetcher) {
		ft.LocalHandler = handler
	})
}
