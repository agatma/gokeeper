package main

import (
	"context"
	"gokeeper/internal/client/core/app"
	"gokeeper/internal/client/core/config"
	"log"
	_ "net/http/pprof"
)

func main() {
	client := app.NewClient(config.NewConfig())
	if err := client.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
