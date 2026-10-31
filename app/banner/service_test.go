package banner_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"kl-assignment/app/banner"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetAllBanners(ctx context.Context) ([]banner.Banner, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]banner.Banner), args.Error(1)
}

func (m *mockRepo) GetBannersByUserId(ctx context.Context, userId string) ([]banner.Banner, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]banner.Banner), args.Error(1)
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestService_GetAllBanners_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := banner.NewService(repo, newTestLogger())

	title := "Promo"
	banners := []banner.Banner{
		{
			BannerId:    "b-1",
			Title:       &title,
			Description: "Desc",
		},
	}
	repo.On("GetAllBanners", ctx).Return(banners, nil)

	res, err := svc.GetAllBanners(ctx)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "b-1", res[0].BannerId)
	assert.Equal(t, "Promo", *res[0].Title)
	repo.AssertExpectations(t)
}

func TestService_GetAllBanners_RepoError(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := banner.NewService(repo, newTestLogger())

	repo.On("GetAllBanners", ctx).Return(nil, errors.New("db error"))

	res, err := svc.GetAllBanners(ctx)

	assert.Error(t, err)
	assert.Nil(t, res)
	repo.AssertExpectations(t)
}

func TestService_GetBannersByUserId_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := banner.NewService(repo, newTestLogger())

	userId := "user-1"
	title := "User Promo"
	banners := []banner.Banner{
		{
			BannerId:    "b-2",
			Title:       &title,
			Description: "User Desc",
		},
	}
	repo.On("GetBannersByUserId", ctx, userId).Return(banners, nil)

	res, err := svc.GetBannersByUserId(ctx, userId)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "b-2", res[0].BannerId)
	repo.AssertExpectations(t)
}

func TestService_GetBannersByUserId_RepoError(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := banner.NewService(repo, newTestLogger())

	userId := "user-1"
	repo.On("GetBannersByUserId", ctx, userId).Return(nil, errors.New("db error"))

	res, err := svc.GetBannersByUserId(ctx, userId)

	assert.Error(t, err)
	assert.Nil(t, res)
	repo.AssertExpectations(t)
}
