package streams

import (
	_ "embed"
)

//go:embed bundle.js
var streamsPolyfill string
