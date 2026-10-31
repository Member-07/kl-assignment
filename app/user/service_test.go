package user_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"kl-assignment/app/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetUserWithGreeting(ctx context.Context, userID string) (user.UserWithGreeting, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return user.UserWithGreeting{}, args.Error(1)
	}
	return args.Get(0).(user.UserWithGreeting), args.Error(1)
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestService_GetUserProfile_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := user.NewService(repo, newTestLogger())

	userId := "user-123"
	userGreeting := user.UserWithGreeting{
		UserId:   userId,
		Name:     "Clare",
		Greeting: "Good morning",
	}
	repo.On("GetUserWithGreeting", ctx, userId).Return(userGreeting, nil)

	res, err := svc.GetUserProfile(ctx, userId)

	assert.NoError(t, err)
	assert.Equal(t, userId, res.UserId)
	assert.Equal(t, "Clare", res.Name)
	assert.Equal(t, "Good morning", res.GreetingMsg)
	repo.AssertExpectations(t)
}

func TestService_GetUserProfile_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := user.NewService(repo, newTestLogger())

	userId := "non-existing"
	repo.On("GetUserWithGreeting", ctx, userId).Return(nil, gorm.ErrRecordNotFound)

	res, err := svc.GetUserProfile(ctx, userId)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Empty(t, res.UserId)
	repo.AssertExpectations(t)
}

func TestService_GetUserProfile_RepoError(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := user.NewService(repo, newTestLogger())

	userId := "user-123"
	repo.On("GetUserWithGreeting", ctx, userId).Return(nil, errors.New("db error"))

	res, err := svc.GetUserProfile(ctx, userId)

	assert.Error(t, err)
	assert.Empty(t, res.UserId)
	repo.AssertExpectations(t)
}
