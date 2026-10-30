package account

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetAccountByUserId(ctx context.Context, userId string) ([]AccountWithDetail, error)
	GetAccountFlags(ctx context.Context, accountIds []string) ([]AccountFlag, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAccountByUserId(ctx context.Context, userId string) ([]AccountWithDetail, error) {
	var results []AccountWithDetail
	selectField := []string{
		"accounts.account_id",
		"accounts.user_id",
		"accounts.type",
		"accounts.currency",
		"accounts.account_number",
		"accounts.issuer",
		"account_balances.amount",
		"account_details.color",
		"account_details.is_main_account",
		"account_details.progress",
	}
	err := r.db.Debug().WithContext(ctx).Table("accounts").
		Select(selectField).
		Joins("LEFT JOIN account_balances ON accounts.account_id = account_balances.account_id").
		Joins("LEFT JOIN account_details ON accounts.account_id = account_details.account_id").
		Where("accounts.user_id = ?", userId).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *repository) GetAccountFlags(ctx context.Context, accountIds []string) ([]AccountFlag, error) {
	var flags []AccountFlag
	err := r.db.Debug().WithContext(ctx).
		Where("account_id IN ?", accountIds).
		Find(&flags).Error
	return flags, err
}
