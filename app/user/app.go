package user

import (
	"log/slog"

	"gorm.io/gorm"
)

func New(db *gorm.DB, logger *slog.Logger) Handler {
	repo := NewRepository(db)
	service := NewService(repo, logger)

	return NewHandler(service)
}
