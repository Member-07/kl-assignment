package banner

import (
	"context"
	"log/slog"
)

type Service interface {
	GetAllBanners(ctx context.Context) ([]BannerResponse, error)
	GetBannersByUserId(ctx context.Context, userId string) ([]BannerResponse, error)
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) GetAllBanners(ctx context.Context) ([]BannerResponse, error) {
	bannerList, err := s.repo.GetAllBanners(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]BannerResponse, 0, len(bannerList))
	for _, banner := range bannerList {
		results = append(results, BannerResponse{
			BannerId:    banner.BannerId,
			Title:       banner.Title,
			Description: banner.Description,
			Image:       banner.Image,
		})
	}

	return results, nil
}

func (s *service) GetBannersByUserId(ctx context.Context, userId string) ([]BannerResponse, error) {
	bannerList, err := s.repo.GetBannersByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	results := make([]BannerResponse, 0, len(bannerList))
	for _, banner := range bannerList {
		results = append(results, BannerResponse{
			BannerId:    banner.BannerId,
			Title:       banner.Title,
			Description: banner.Description,
			Image:       banner.Image,
		})
	}

	return results, nil
}
