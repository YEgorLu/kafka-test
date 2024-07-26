package main

import (
	_ "github.com/YEgorLu/kafka-test/api/rest"
	"github.com/YEgorLu/kafka-test/internal/config"
	"github.com/YEgorLu/kafka-test/internal/db"
	"github.com/YEgorLu/kafka-test/internal/logger"
	"github.com/YEgorLu/kafka-test/internal/server"
)

// @title Kafka Test Swagger API
// @version 1.0
// @description Swagger API for Kafka Test.
// @termsOfService http://swagger.io/terms/

// @license.name MIT
// @license.url https://github.com/YEgorLu/kafka-test/blob/master/LICENSE

// @BasePath /
func main() {
	l := logger.Get()
	serverConfig := server.ServerConfig{
		Port: config.App.Port,
		Addr: config.App.ServerHostAddr,
		Log:  l,
	}
	defer db.CloseAll()
	err := server.
		NewServer(&serverConfig).
		Configure().
		WithSwagger().
		Run()
	if err != nil {
		l.Error(err)
	}
	l.Info("Program closed")
}
