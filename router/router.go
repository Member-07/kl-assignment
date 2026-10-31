package router

import (
	"kl-assignment/app/account"
	"kl-assignment/app/banner"
	debitcard "kl-assignment/app/debit_card"
	"kl-assignment/app/health"
	"kl-assignment/app/transaction"
	"kl-assignment/app/user"
	"kl-assignment/pkg/config"
	"kl-assignment/pkg/middleware"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(r *gin.Engine, cfg config.Config, db *gorm.DB, logger *slog.Logger) {
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.ErrorResponseMiddleware())

	api := r.Group("/api")

	// health
	healthHandler := health.New()
	api.GET("/health", healthHandler.Check)

	// user
	userHandler := user.New(db, logger)
	user := api.Group("/user")
	user.GET("/:user_id", userHandler.GetUserProfile)

	// account
	accountHandler := account.New(db, logger)
	account := api.Group("/account")
	account.GET("/:user_id", accountHandler.GetAccountsByUserId)

	// debitcard
	debitcardHandler := debitcard.New(db, logger)
	debitcard := api.Group("/debit-card")
	debitcard.GET("/:user_id", debitcardHandler.GetCardsByUserId)
	debitcard.PUT("/status", debitcardHandler.UpdateCardStatus)

	// banner
	bannerHandler := banner.New(db, logger)
	banner := api.Group("/banner")
	banner.GET("/", bannerHandler.GetAllBanners)
	banner.GET("/:user_id", bannerHandler.GetBannersByUserId)

	// transaction
	transactionHandler := transaction.New(db, logger)
	transaction := api.Group("/transaction")
	transaction.GET("/:user_id", transactionHandler.GetTransactionsByUserId)
}
