package main

import (
	"github.com/daheige/go-ddd-api/config"
	"github.com/daheige/go-ddd-api/interfaces"
)

func init() {
	config.InitDB()
}

func main() {
	defer config.CloseDB()
	interfaces.Run(8000)
}
