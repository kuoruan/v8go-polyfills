package internal

import "rogchap.com/v8go"

type Context interface {
	Isolate() (*v8go.Isolate, error)
	Global() *v8go.Object
	RunScript(source string, origin string) (*v8go.Value, error)
}
