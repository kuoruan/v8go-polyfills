package polyfills

import (
	"go.kuoruan.net/v8go-polyfills/fetch"

	"rogchap.com/v8go"
)

type PolyfillOption interface {
	Polyfill() string
}

func InjectAll(ctx *v8go.Context, opt ...PolyfillOption) error {
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
