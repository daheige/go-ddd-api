package providers

import (
	"fmt"

	"github.com/daheige/go-ddd-api/internal/api"
	"github.com/daheige/go-ddd-api/internal/api/news"
	"github.com/daheige/go-ddd-api/internal/api/topics"
	"github.com/daheige/go-ddd-api/internal/application"
	"github.com/daheige/go-ddd-api/internal/infras/config"
	"github.com/daheige/go-ddd-api/internal/infras/migration"
	"github.com/daheige/go-ddd-api/internal/infras/monitor"
	"github.com/daheige/go-ddd-api/internal/infras/persistence"
	"github.com/go-god/gdi"
	"github.com/go-god/gdi/factory"
	"github.com/go-god/logger"
	"go.uber.org/zap"
)

// Inject dependency injection
func Inject(app *api.NewsHandler) error {
	loggerInject() // inject logger
	conf := config.NewConfig()
	di := factory.CreateDI(factory.FbInject) // create a di container
	err := di.Provide(
		&gdi.Object{Value: app},
		&gdi.Object{Value: conf.AppConfig()}, // app section inject
		&gdi.Object{Value: conf.InitDB()},    // db inject
		&gdi.Object{Value: &migration.MigrateAction{}},
		&gdi.Object{Value: &news.NewsHandler{}},
		&gdi.Object{Value: &topics.TopicHandler{}},
		&gdi.Object{Value: &application.TopicService{}},
		&gdi.Object{Value: &application.NewsService{}},
		&gdi.Object{Value: &persistence.NewsRepositoryImpl{}},
		&gdi.Object{Value: &persistence.TopicRepositoryImpl{}},
		&gdi.Object{Value: &api.RouterHandler{}},
		&gdi.Object{Value: &monitor.PrometheusMonitor{}},
	)
	if err != nil {
		return fmt.Errorf("provide error:%s", err.Error())
	}

	// invoke object
	err = di.Invoke()
	if err != nil {
		return fmt.Errorf("invoke error:%s", err.Error())
	}

	return nil
}

func loggerInject() {
	opts := []logger.Option{
		logger.WithLogDir("./logs"),          // log dir
		logger.WithLogFilename("go-app.log"), // default zap.log
		logger.WithStdout(false),             // In the common production environment, do not output it to stdout
		logger.WithJsonFormat(true),          // json formatting
		logger.WithAddCaller(true),           // Print line number
		logger.WithEnableColor(false),        // Whether logs are dyed. By default, logs are not dyed

		// Set the lowest level of log printing. If this parameter is not set, the default level is info
		logger.WithLogLevel(zap.DebugLevel),
		logger.WithMaxAge(3),       // The maximum storage time is 3 days
		logger.WithMaxSize(20),     // Each log file has a maximum of 20MB
		logger.WithCompress(false), // Log no compression
	}

	// Generates the default log handle object
	logger.Default(opts...)
}
