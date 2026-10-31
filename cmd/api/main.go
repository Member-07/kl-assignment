package main

import (
	"kl-assignment/pkg/config"
	"kl-assignment/pkg/database"
	"kl-assignment/pkg/logger"
	"kl-assignment/router"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := logger.New(cfg.LogLevel)

	logger.Info("Starting API server")

	db, err := database.ConnectMysql(cfg.MysqlConfig)
	if err != nil {
		panic(err)
	}

	r := gin.New()
	router.Register(r, cfg, db, logger)

	r.Run()
}
