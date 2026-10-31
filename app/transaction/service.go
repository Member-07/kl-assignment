package transaction

import (
	"context"
	"log/slog"
	"math"
)

type Service interface {
	GetTransactionsByUserId(ctx context.Context, userId string, page, limit int) (TransactionPaginationResponse, error)
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) GetTransactionsByUserId(ctx context.Context, userId string, page, limit int) (TransactionPaginationResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	transactionList, totalRows, err := s.repo.GetTransactionsByUserId(ctx, userId, limit, offset)
	if err != nil {
		return TransactionPaginationResponse{}, err
	}

	transactionRes := make([]TransactionResponse, 0, len(transactionList))
	for _, transaction := range transactionList {
		name, image := "", ""
		if transaction.Name != nil {
			name = *transaction.Name
		}
		if transaction.Image != nil {
			image = *transaction.Image
		}

		isBank := false
		if transaction.IsBank != nil {
			isBank = *transaction.IsBank
		}

		transactionRes = append(transactionRes, TransactionResponse{
			TransactionId: transaction.TransactionId,
			Name:          name,
			Image:         image,
			IsBank:        isBank,
		})
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))

	return TransactionPaginationResponse{
		Data:       transactionRes,
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}, nil
}
