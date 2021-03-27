package url

import (
	_ "embed"
)

//go:embed bundle.js
var urlPolyfill string
