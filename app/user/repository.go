package user

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetUserWithGreeting(ctx context.Context, userID string) (UserWithGreeting, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetUserWithGreeting(ctx context.Context, userID string) (UserWithGreeting, error) {
	var result UserWithGreeting

	selectField := []string{
		"users.user_id",
		"users.name",
		"user_greetings.greeting",
	}
	res := r.db.Debug().WithContext(ctx).Table("users").
		Select(selectField).
		Joins("LEFT JOIN user_greetings ON users.user_id = user_greetings.user_id").
		Where("users.user_id = ?", userID).
		Scan(&result)

	if res.Error != nil {
		return result, res.Error
	}

	if res.RowsAffected == 0 {
		return result, gorm.ErrRecordNotFound
	}

	return result, nil
}
