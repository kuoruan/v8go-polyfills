package url

import (
	"errors"

	"rogchap.com/v8go"
)

func InjectTo(ctx *v8go.Context) error {
	if ctx == nil {
		return errors.New("v8go-polyfills/url: ctx is required")
	}

	_, err := ctx.RunScript(urlPolyfill, "url-polyfill.js")
	return err
}
