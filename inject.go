package polyfills

import "go.kuoruan.net/v8go-polyfills/fetch"

type Polyfill int

const (
	PolyfillFetch Polyfill = iota + 1
)

type Option interface {
	ForPolyfill() Polyfill
}

func InjectAll(ctx Context, opt ...Option) error {
	var fetchOpts []fetch.Option

	for _, o := range opt {
		switch t := o.(type) {
		case fetch.Option:
			fetchOpts = append(fetchOpts, t)
		}
	}

	if err := fetch.Inject(ctx, fetchOpts...); err != nil {
		return err
	}

	return nil
}
