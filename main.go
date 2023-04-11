package main

import (
	"log"

	"github.com/daheige/go-ddd-api/internal/api"
	"github.com/daheige/go-ddd-api/internal/api/news"
	"github.com/daheige/go-ddd-api/internal/api/topics"
	"github.com/daheige/go-ddd-api/internal/application"
	"github.com/daheige/go-ddd-api/internal/infras/config"
	"github.com/daheige/go-ddd-api/internal/infras/migration"
	"github.com/daheige/go-ddd-api/internal/infras/persistence"
	"github.com/go-god/gdi"
	"github.com/go-god/gdi/factory"
)

func main() {
	var (
		app  api.NewsService
		conf = config.NewConfig()
	)

	di := factory.CreateDI(factory.FbInject) // create a di container
	err := di.Provide(
		&gdi.Object{Value: &app},
		&gdi.Object{Value: conf.AppConfig()}, // app section inject
		&gdi.Object{Value: conf.InitDB()},    // db inject
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
	app.Run()
}
