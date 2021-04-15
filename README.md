# Polyfills for [V8Go](https://github.com/rogchap/v8go)

## Install

```shell
go get -u go.kuoruan.net/v8go-polyfills
```

> This module uses Golang [embed](https://golang.org/pkg/embed/), so requires Go version 1.16

## Polyfill List

* base64: `atob` and `btoa`

* console: `console.log`

* fetch: `fetch`

* timers: `setTimeout`, `clearTimeout`, `setInterval` and `clearInterval`

* url: `URL` and `URLSearchParams`

## Usage

### fetch polyfill

```go
package main

import (
	"errors"
	"fmt"
	"time"

	"go.kuoruan.net/v8go-polyfills/fetch"
	"rogchap.com/v8go"
)

func main() {
	iso, _ := v8go.NewIsolate()
	global, _ := v8go.NewObjectTemplate(iso)

	if err := fetch.InjectTo(iso, global); err != nil {
		panic(err)
	}

	ctx, _ := v8go.NewContext(iso, global)

	val, err := ctx.RunScript("fetch('https://www.example.com').then(res => res.text())", "fetch.js")
	if err != nil {
		panic(err)
	}

	proms, err := val.AsPromise()
	if err != nil {
		panic(err)
	}
	done := make(chan bool, 1)

	go func() {
		for proms.State() == v8go.Pending {
			continue
		}
		done <- true
	}()

	select {
	case <-time.After(time.Second * 10):
		panic(errors.New("request timeout"))
	case <-done:
		html := proms.Result().String()
		fmt.Println(html)
	}
}
```
