package transaction_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"kl-assignment/app/transaction"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetTransactionsByUserId(ctx context.Context, userId string, limit, offset int) ([]transaction.Transaction, int64, error) {
	args := m.Called(ctx, userId, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]transaction.Transaction), args.Get(1).(int64), args.Error(2)
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestService_GetTransactionsByUserId_Success(t *testing.T) {
	repo := new(mockRepo)
	svc := transaction.NewService(repo, newTestLogger())

	name := "Coffee Shop"
	image := "coffee.png"
	isBank := true

	txList := []transaction.Transaction{
		{
			TransactionId: "tx-001",
			UserId:        &name,
			Name:          &name,
			Image:         &image,
			IsBank:        &isBank,
		},
		{
			TransactionId: "tx-002",
			UserId:        nil,
			Name:          nil,
			Image:         nil,
			IsBank:        nil,
		},
	}

	repo.On("GetTransactionsByUserId", mock.Anything, "usr-001", 10, 0).Return(txList, int64(25), nil)

	res, err := svc.GetTransactionsByUserId(context.Background(), "usr-001", 1, 10)

	assert.NoError(t, err)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, 1, res.Page)
	assert.Equal(t, 10, res.Limit)
	assert.Equal(t, int64(25), res.TotalRows)
	assert.Equal(t, 3, res.TotalPages)

	// First item with values
	assert.Equal(t, "tx-001", res.Data[0].TransactionId)
	assert.Equal(t, "Coffee Shop", res.Data[0].Name)
	assert.Equal(t, "coffee.png", res.Data[0].Image)
	assert.True(t, res.Data[0].IsBank)

	// Second item with nils handled
	assert.Equal(t, "tx-002", res.Data[1].TransactionId)
	assert.Equal(t, "", res.Data[1].Name)
	assert.Equal(t, "", res.Data[1].Image)
	assert.False(t, res.Data[1].IsBank)

	repo.AssertExpectations(t)
}

func TestService_GetTransactionsByUserId_ClampingPageLimit(t *testing.T) {
	repo := new(mockRepo)
	svc := transaction.NewService(repo, newTestLogger())

	repo.On("GetTransactionsByUserId", mock.Anything, "usr-001", 10, 0).Return([]transaction.Transaction{}, int64(0), nil)

	res, err := svc.GetTransactionsByUserId(context.Background(), "usr-001", 0, -5)

	assert.NoError(t, err)
	assert.Equal(t, 1, res.Page)
	assert.Equal(t, 10, res.Limit)
	assert.Equal(t, 0, res.TotalPages)
	repo.AssertExpectations(t)
}

func TestService_GetTransactionsByUserId_RepoError(t *testing.T) {
	repo := new(mockRepo)
	svc := transaction.NewService(repo, newTestLogger())

	repo.On("GetTransactionsByUserId", mock.Anything, "usr-001", 10, 0).Return(nil, int64(0), errors.New("database timeout"))

	res, err := svc.GetTransactionsByUserId(context.Background(), "usr-001", 1, 10)

	assert.Error(t, err)
	assert.Equal(t, "database timeout", err.Error())
	assert.Empty(t, res.Data)
	repo.AssertExpectations(t)
}
