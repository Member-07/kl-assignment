package user

import (
	"context"
	"log/slog"
)

type Service interface {
	GetUserProfile(ctx context.Context, userId string) (UserResponse, error)
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) GetUserProfile(ctx context.Context, userId string) (UserResponse, error) {
	res, err := s.repo.GetUserWithGreeting(ctx, userId)
	if err != nil {
		s.logger.Error("GetUserWithGreeting", "error", err)
		return UserResponse{}, err
	}

	return UserResponse{
		UserId:      res.UserId,
		Name:        res.Name,
		GreetingMsg: res.Greeting,
	}, nil
}
