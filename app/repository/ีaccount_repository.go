package repository

import (
	"kl-assignment/app/model"

	"gorm.io/gorm"
)

type AccountRepository interface {
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (r *accountRepository) GetAccount(accountID string) (res model.Account, err error) {
	err = r.db.Where("account_id = ?", accountID).
		Find(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}
