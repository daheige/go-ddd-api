package main

import (
	"flag"
	"log"

	"github.com/go-god/gdi"
	"github.com/go-god/gdi/factory"

	"github.com/daheige/go-ddd-api/api"
	"github.com/daheige/go-ddd-api/api/news"
	"github.com/daheige/go-ddd-api/api/topics"
	"github.com/daheige/go-ddd-api/internal/application"
	"github.com/daheige/go-ddd-api/internal/infras/config"
	"github.com/daheige/go-ddd-api/internal/infras/migration"
	"github.com/daheige/go-ddd-api/internal/infras/persistence"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8000, "app run port,eg:8000")
	flag.Parse()
}

func main() {
	var app api.AppService
	di := factory.CreateDI(factory.FbInject) // create a di container
	err := di.Provide(
		&gdi.Object{Value: &app},
		&gdi.Object{Value: config.Init()},
		&gdi.Object{Value: &migration.MigrateAction{}},
		&gdi.Object{Value: &news.NewsHandler{}},
		&gdi.Object{Value: &topics.TopicHandler{}},
		&gdi.Object{Value: &application.TopicService{}},
		&gdi.Object{Value: &application.NewsService{}},
		&gdi.Object{Value: &persistence.NewsRepositoryImpl{}},
		&gdi.Object{Value: &persistence.TopicRepositoryImpl{}},
	)

	if err != nil {
		log.Fatalln("provide error: ", err)
	}

	err = di.Invoke()
	if err != nil {
		log.Fatalln("invoke error: ", err)
	}

	// app service run
	app.Run(port)
}
