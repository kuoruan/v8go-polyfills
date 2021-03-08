package polyfills

import (
	"go.kuoruan.net/v8go-polyfills/fetch"
	"go.kuoruan.net/v8go-polyfills/internal"
)

func InjectAll(ctx internal.Context, opt ...internal.Option) error {
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
