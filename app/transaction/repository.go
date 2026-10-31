package transaction

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetTransactionsByUserId(ctx context.Context, userId string, limit, offset int) ([]Transaction, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetTransactionsByUserId(ctx context.Context, userId string, limit, offset int) ([]Transaction, int64, error) {
	var transactions []Transaction
	var totalRows int64

	dbQuery := r.db.Debug().Model(&Transaction{}).Where("user_id = ?", userId)

	if err := dbQuery.WithContext(ctx).Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	err := dbQuery.WithContext(ctx).Offset(offset).
		Limit(limit).
		Find(&transactions).Error

	if err != nil {
		return nil, 0, err
	}

	return transactions, totalRows, nil
}
