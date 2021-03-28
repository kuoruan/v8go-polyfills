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
