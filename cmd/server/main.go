package main

import (
	"gokeeper/internal/server/core/app"
	_ "net/http/pprof"
)

func main() {
	server, err := app.NewServer()
	if err != nil {
		panic(err)
	}
	server.Run()
}
