package main

import (
	"kl-assignment/config"
	"kl-assignment/database"
	"kl-assignment/router"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err := database.ConnectMysql(cfg.MysqlConfig)
	if err != nil {
		panic(err)
	}

	r := gin.New()
	router.Register(r, cfg, db)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
