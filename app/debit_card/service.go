package debitcard

import (
	"context"
	"log/slog"
)

type Service interface {
	GetCardsByUserId(ctx context.Context, userId string) ([]DebitCardResponse, error)
	UpdateCardStatus(ctx context.Context, req UpdateCardStatusRequest) error
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) GetCardsByUserId(ctx context.Context, userId string) ([]DebitCardResponse, error) {
	cardList, err := s.repo.GetCardsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	results := make([]DebitCardResponse, 0, len(cardList))
	for _, card := range cardList {
		results = append(results, DebitCardResponse{
			CardId: card.CardID,
			Name:   card.Name,
			Status: card.Status,
			Detail: CardDetailResponse{
				Issuer: card.Issuer,
				Number: card.Number,
			},
			Design: CardDesignResponse{
				Color:       card.Color,
				BorderColor: card.BorderColor,
			},
		})
	}

	return results, nil
}

func (s *service) UpdateCardStatus(ctx context.Context, req UpdateCardStatusRequest) error {
	err := s.repo.UpdateCardStatus(ctx, DebitCardStatus{
		UserId: req.UserId,
		CardId: req.CardId,
		Status: req.Status,
	})

	if err != nil {
		return err
	}
	return nil
}
