package account

import (
	"context"
	"log/slog"
)

type Service interface {
	GetAccountsByUserId(ctx context.Context, userId string) ([]AccountResponse, error)
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) GetAccountsByUserId(ctx context.Context, userId string) ([]AccountResponse, error) {
	accountList, err := s.repo.GetAccountByUserId(ctx, userId)
	if err != nil {
		s.logger.Error("GetAccountByUserId", "error", err)
		return nil, err
	}

	results := make([]AccountResponse, 0, len(accountList))
	accountIds := make([]string, 0, len(accountList))

	for _, account := range accountList {
		accountIds = append(accountIds, account.AccountId)
	}

	accountFlagList, err := s.repo.GetAccountFlags(ctx, accountIds)
	if err != nil {
		s.logger.Error("GetAccountFlags", "error", err)
		return nil, err
	}

	mapAccountFlag := make(map[string][]AccountFlagResponse)
	for _, flag := range accountFlagList {
		mapAccountFlag[flag.AccountId] = append(mapAccountFlag[flag.AccountId], AccountFlagResponse{
			FlagType:  flag.FlagType,
			FlagValue: flag.FlagValue,
		})
	}

	for _, account := range accountList {
		accFlags, ok := mapAccountFlag[account.AccountId]
		if !ok {
			accFlags = []AccountFlagResponse{}
		}

		results = append(results, mappingAccountResponse(account, accFlags))
	}

	return results, nil
}

func mappingAccountResponse(account AccountWithDetail, accFlags []AccountFlagResponse) AccountResponse {
	return AccountResponse{
		AccountId:     account.AccountId,
		AccountNumber: account.AccountNumber,
		Amount:        account.Amount,
		Currency:      account.Currency,
		UserId:        account.UserId,
		Type:          account.Type,
		Issuer:        account.Issuer,
		Color:         account.Color,
		IsMainAccount: account.IsMainAccount,
		Progress:      account.Progress,
		Flags:         accFlags,
	}
}
