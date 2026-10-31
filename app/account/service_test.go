package account_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"kl-assignment/app/account"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetAccountByUserId(ctx context.Context, userId string) ([]account.AccountWithDetail, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]account.AccountWithDetail), args.Error(1)
}

func (m *mockRepo) GetAccountFlags(ctx context.Context, accountIds []string) ([]account.AccountFlag, error) {
	args := m.Called(ctx, accountIds)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]account.AccountFlag), args.Error(1)
}

func (m *mockRepo) SetMainAccount(ctx context.Context, userId string, accountId string) error {
	args := m.Called(ctx, userId, accountId)
	return args.Error(0)
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestService_GetAccountsByUserId_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := account.NewService(repo, newTestLogger())

	userId := "user-1"
	accounts := []account.AccountWithDetail{
		{
			AccountId:     "acc-1",
			AccountNumber: "123456",
			Amount:        decimal.NewFromFloat(100.50),
			Currency:      "THB",
			UserId:        userId,
			Type:          "Saving",
			Issuer:        "TestBank",
			Color:         "#fff",
			IsMainAccount: true,
			Progress:      50,
		},
	}
	flags := []account.AccountFlag{
		{
			AccountId: "acc-1",
			FlagType:  "lock",
			FlagValue: "Disbursement",
		},
	}

	repo.On("GetAccountByUserId", ctx, userId).Return(accounts, nil)
	repo.On("GetAccountFlags", ctx, []string{"acc-1"}).Return(flags, nil)

	res, err := svc.GetAccountsByUserId(ctx, userId)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "acc-1", res[0].AccountId)
	assert.Len(t, res[0].Flags, 1)
	assert.Equal(t, "lock", res[0].Flags[0].FlagType)
	repo.AssertExpectations(t)
}

func TestService_GetAccountsByUserId_NoFlags(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := account.NewService(repo, newTestLogger())

	userId := "user-1"
	accounts := []account.AccountWithDetail{
		{AccountId: "acc-1"},
	}

	repo.On("GetAccountByUserId", ctx, userId).Return(accounts, nil)
	repo.On("GetAccountFlags", ctx, []string{"acc-1"}).Return([]account.AccountFlag{}, nil)

	res, err := svc.GetAccountsByUserId(ctx, userId)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Empty(t, res[0].Flags)
	repo.AssertExpectations(t)
}

func TestService_GetAccountsByUserId_RepoError(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := account.NewService(repo, newTestLogger())

	userId := "user-1"
	repo.On("GetAccountByUserId", ctx, userId).Return(nil, errors.New("db error"))

	res, err := svc.GetAccountsByUserId(ctx, userId)

	assert.Error(t, err)
	assert.Nil(t, res)
	repo.AssertExpectations(t)
}

func TestService_GetAccountsByUserId_FlagsRepoError(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := account.NewService(repo, newTestLogger())

	userId := "user-1"
	accounts := []account.AccountWithDetail{{AccountId: "acc-1"}}

	repo.On("GetAccountByUserId", ctx, userId).Return(accounts, nil)
	repo.On("GetAccountFlags", ctx, []string{"acc-1"}).Return(nil, errors.New("flag db error"))

	res, err := svc.GetAccountsByUserId(ctx, userId)

	assert.Error(t, err)
	assert.Nil(t, res)
	repo.AssertExpectations(t)
}

func TestService_SetMainAccount_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := account.NewService(repo, newTestLogger())

	req := account.SetMainAccountRequest{
		UserId:    "user-1",
		AccountId: "acc-1",
	}

	repo.On("SetMainAccount", ctx, req.UserId, req.AccountId).Return(nil)

	err := svc.SetMainAccount(ctx, req)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestService_SetMainAccount_RepoError(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := account.NewService(repo, newTestLogger())

	req := account.SetMainAccountRequest{
		UserId:    "user-1",
		AccountId: "acc-1",
	}

	repo.On("SetMainAccount", ctx, req.UserId, req.AccountId).Return(errors.New("db error"))

	err := svc.SetMainAccount(ctx, req)

	assert.Error(t, err)
	repo.AssertExpectations(t)
}
