package router

import (
	"kl-assignment/app/health"
	"kl-assignment/app/user"
	"kl-assignment/pkg/config"
	"kl-assignment/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(r *gin.Engine, cfg config.Config, db *gorm.DB) {
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.ErrorResponseMiddleware())

	api := r.Group("/api")

	// health
	healthHandler := health.New()
	api.GET("/health", healthHandler.Check)

	// User
	userHandler := user.New(db)
	user := api.Group("/user")
	user.GET("/:user_id", userHandler.GetUserProfile)

}
