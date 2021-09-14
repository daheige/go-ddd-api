package config

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// DBConfig db config
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Dbname   string
}

// Init read config and init gorm db
func Init() *gorm.DB {
	dbConf := &DBConfig{}
	dbConf.read()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbConf.User,
		dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Dbname)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalln("db open error: ", err)
	}

	return db
}

// read and parse the configuration file
func (c *DBConfig) read() {
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Fatal(err)
	}
}
