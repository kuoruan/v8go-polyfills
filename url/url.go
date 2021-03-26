package url

import (
	_ "embed"
	"errors"

	"rogchap.com/v8go"
)

//go:embed bundle.js
var urlPolyfill string

func Inject(ctx *v8go.Context) error {
	if ctx == nil {
		return errors.New("ctx is required")
	}

	_, err := ctx.RunScript(urlPolyfill, "url-polyfill.js")
	return err
}
