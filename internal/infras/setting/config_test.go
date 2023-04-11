package setting

import (
	"log"
	"testing"
)

func TestConfigRepoImpl_Load(t *testing.T) {
	conf := New(WithConfigFile("./test.yaml"))
	if err := conf.Load(); err != nil {
		log.Fatalln("load config err: ", err)
	}

	log.Println(conf.IsSet("app"))

	var app int64
	conf.ReadSection("app", &app)
	log.Println(app)
}
