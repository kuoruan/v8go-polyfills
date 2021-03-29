package polyfills

import (
	"go.kuoruan.net/v8go-polyfills/base64"
	"go.kuoruan.net/v8go-polyfills/console"
	"go.kuoruan.net/v8go-polyfills/fetch"
	"go.kuoruan.net/v8go-polyfills/internal"
	"go.kuoruan.net/v8go-polyfills/url"

	"rogchap.com/v8go"
)

func InjectToGlobalObject(iso *v8go.Isolate, global *v8go.ObjectTemplate, opt ...interface{}) error {
	var fetchOpts []fetch.Option
	var consoleOpts []console.Option

	for _, o := range opt {
		switch t := o.(type) {
		case fetch.Option:
			fetchOpts = append(fetchOpts, t)
		case console.Option:
			consoleOpts = append(consoleOpts, t)
		}
	}

	if err := fetch.InjectTo(iso, global, fetchOpts...); err != nil {
		return err
	}

	if err := console.InjectTo(iso, global, consoleOpts...); err != nil {
		return err
	}

	if err := base64.InjectTo(iso, global); err != nil {
		return err
	}

	return nil
}

func InjectToContext(ctx *v8go.Context) error {
	for _, p := range []func(*v8go.Context) error{
		url.InjectTo,
	} {
		if err := p(ctx); err != nil {
			return err
		}
	}

	return nil
}

func Version() string {
	return internal.Version
}
