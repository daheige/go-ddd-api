package config

import (
	"fmt"
	"log"

	"github.com/daheige/go-ddd-api/internal/infras/setting"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// DBConf database config
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Dbname   string
}

// AppConfig app section
type AppConfig struct {
	Port     int
	AppName  string
	AppEnv   string
	AppDebug bool
}

// configImpl config
type configImpl struct {
	DB  DBConfig
	App AppConfig
}

// New load config
func NewConfig() *configImpl {
	s := &configImpl{}
	s.load()

	return s
}

// read and parse the configuration file
func (s *configImpl) load() {
	conf := setting.New(setting.WithConfigFile("./config.yaml"))
	if err := conf.Load(); err != nil {
		log.Fatalf("read config file err:%s\n", err.Error())
	}

	if err := conf.ReadSection("app", &s.App); err != nil {
		log.Fatalf("read config file err:%s", err.Error())
	}

	if err := conf.ReadSection("db", &s.DB); err != nil {
		log.Fatalf("read config file err:%s", err.Error())
	}
}

// InitDB init gorm db
func (s *configImpl) InitDB() *gorm.DB {
	dbConf := s.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbConf.User,
		dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Dbname)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalln("db open error: ", err)
	}

	return db
}

// AppConfig returns app config
func (s *configImpl) AppConfig() *AppConfig {
	return &s.App
}
