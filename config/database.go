package config

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
)

type AppConfig struct {
	DB *gorm.DB
}

var AppConf = AppConfig{}

// ConfigDB db seting
type ConfigDB struct {
	User     string
	Password string
	Host     string
	Port     string
	Dbname   string
}

// InitDB init gorm db
func InitDB() {
	dbConf := ConfigDB{}
	dbConf.Read()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbConf.User,
		dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Dbname)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalln("db open error: ", err)
	}

	AppConf.DB = db
}

// CloseDB close db instance.
func CloseDB() {
	AppConf.DB.Close()
}

// Read and parse the configuration file
func (c *ConfigDB) Read() {
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Fatal(err)
	}
}
