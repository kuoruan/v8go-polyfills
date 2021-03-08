package internal

type Polyfill int

const (
	PolyfillFetch Polyfill = iota + 1
)

type Option interface {
	Polyfill() Polyfill
}
