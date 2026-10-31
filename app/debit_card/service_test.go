package debitcard_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	debitcard "kl-assignment/app/debit_card"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetCardsByUserId(ctx context.Context, userId string) ([]debitcard.DebitCard, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]debitcard.DebitCard), args.Error(1)
}

func (m *mockRepo) UpdateCardStatus(ctx context.Context, req debitcard.DebitCardStatus) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestService_GetCardsByUserId_Success(t *testing.T) {
	repo := new(mockRepo)
	svc := debitcard.NewService(repo, newTestLogger())

	cards := []debitcard.DebitCard{
		{
			CardID:      "card-001",
			UserID:      "usr-001",
			Name:        "Platinum Card",
			Issuer:      "VISA",
			Number:      "4111111111111111",
			Color:       "#000000",
			BorderColor: "#FFFFFF",
			Status:      "Active",
		},
	}

	repo.On("GetCardsByUserId", mock.Anything, "usr-001").Return(cards, nil)

	res, err := svc.GetCardsByUserId(context.Background(), "usr-001")

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "card-001", res[0].CardId)
	assert.Equal(t, "Platinum Card", res[0].Name)
	assert.Equal(t, "Active", res[0].Status)
	assert.Equal(t, "VISA", res[0].Detail.Issuer)
	assert.Equal(t, "4111111111111111", res[0].Detail.Number)
	assert.Equal(t, "#000000", res[0].Design.Color)
	assert.Equal(t, "#FFFFFF", res[0].Design.BorderColor)
	repo.AssertExpectations(t)
}

func TestService_GetCardsByUserId_RepoError(t *testing.T) {
	repo := new(mockRepo)
	svc := debitcard.NewService(repo, newTestLogger())

	repo.On("GetCardsByUserId", mock.Anything, "usr-001").Return(nil, errors.New("db error"))

	res, err := svc.GetCardsByUserId(context.Background(), "usr-001")

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "db error", err.Error())
	repo.AssertExpectations(t)
}

func TestService_UpdateCardStatus_Success(t *testing.T) {
	repo := new(mockRepo)
	svc := debitcard.NewService(repo, newTestLogger())

	req := debitcard.UpdateCardStatusRequest{
		UserId: "usr-001",
		CardId: "card-001",
		Status: "InActive",
	}

	expectedStatus := debitcard.DebitCardStatus{
		UserId: "usr-001",
		CardId: "card-001",
		Status: "InActive",
	}

	repo.On("UpdateCardStatus", mock.Anything, expectedStatus).Return(nil)

	err := svc.UpdateCardStatus(context.Background(), req)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestService_UpdateCardStatus_RepoError(t *testing.T) {
	repo := new(mockRepo)
	svc := debitcard.NewService(repo, newTestLogger())

	req := debitcard.UpdateCardStatusRequest{
		UserId: "usr-001",
		CardId: "card-001",
		Status: "Active",
	}

	expectedStatus := debitcard.DebitCardStatus{
		UserId: "usr-001",
		CardId: "card-001",
		Status: "Active",
	}

	repo.On("UpdateCardStatus", mock.Anything, expectedStatus).Return(errors.New("update failed"))

	err := svc.UpdateCardStatus(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "update failed", err.Error())
	repo.AssertExpectations(t)
}
