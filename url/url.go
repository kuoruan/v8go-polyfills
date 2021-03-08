package url

import (
	_ "embed"
	"errors"

	"go.kuoruan.net/v8go-polyfills/internal"
)

//go:embed bundle.js
var urlPolyfill string

func Inject(ctx internal.Context) error {
	if ctx == nil {
		return errors.New("ctx is required")
	}

	_, err := ctx.RunScript(urlPolyfill, "url-polyfill.js")
	return err
}
