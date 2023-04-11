package main

import (
	"log"

	"github.com/daheige/go-ddd-api/internal/api"
	"github.com/daheige/go-ddd-api/internal/providers"
)

func main() {
	var app api.NewsHandler
	// dependency injection
	if err := providers.Inject(&app); err != nil {
		log.Fatalf("provier error:%s\n", err)
	}

	// app service run
	app.Run()
}
