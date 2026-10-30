package debitcard

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetCardsByUserId(ctx context.Context, userId string) ([]DebitCard, error)
	UpdateCardStatus(ctx context.Context, req DebitCardStatus) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetCardsByUserId(ctx context.Context, userId string) ([]DebitCard, error) {
	var results []DebitCard
	selectField := []string{
		"debit_cards.card_id",
		"debit_cards.user_id",
		"debit_cards.name",
		"debit_card_details.issuer",
		"debit_card_details.number",
		"debit_card_design.color",
		"debit_card_design.border_color",
		"debit_card_status.status",
	}

	err := r.db.Debug().WithContext(ctx).Table("debit_cards").
		Select(selectField).
		Joins("LEFT JOIN debit_card_details ON debit_cards.card_id = debit_card_details.card_id").
		Joins("LEFT JOIN debit_card_design ON debit_cards.card_id = debit_card_design.card_id").
		Joins("LEFT JOIN debit_card_status ON debit_cards.card_id = debit_card_status.card_id").
		Where("debit_cards.user_id = ?", userId).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, err
}

func (r *repository) UpdateCardStatus(ctx context.Context, req DebitCardStatus) error {
	return r.db.Debug().WithContext(ctx).Model(&DebitCardStatus{}).
		Where("card_id = ?", req.CardId).
		Where("user_id = ?", req.UserId).
		Update("status", req.Status).Error
}
