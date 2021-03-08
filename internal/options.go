package internal

type Polyfill int

const (
	PolyfillFetch Polyfill = iota + 1
	PolyfillConsole
)

type Option interface {
	Polyfill() Polyfill
}
