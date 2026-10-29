package user

import (
	"context"
)

type Service interface {
	GetUserProfile(ctx context.Context, userId string) (UserResponse, error)
}

type service struct {
	repo Repository
}

func NewUserService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetUserProfile(ctx context.Context, userId string) (UserResponse, error) {
	res, err := s.repo.GetUserWithGreeting(ctx, userId)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{
		UserId:      res.UserId,
		Name:        res.Name,
		GreetingMsg: res.Greeting,
	}, nil
}
