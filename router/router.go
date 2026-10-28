package router

import (
	"kl-assignment/config"
	"kl-assignment/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(r *gin.Engine, cfg config.Config, db *gorm.DB) {
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	// api := r.Group("/api")
	// auth := api.Group("/auth")
	// auth.POST("/register", authHandler.Register)
	// auth.POST("/login", authHandler.Login)

	// user := api.Group("/users")

	// user.Use(middleware.Auth(jwt))
	// {
	// 	user.POST("/", userHandler.Create)
	// 	user.GET("", userHandler.List)
	// 	user.GET("/:id", userHandler.GetByID)
	// 	user.PUT("/:id", userHandler.Update)
	// 	user.DELETE("/:id", userHandler.Delete)
	// }

}
