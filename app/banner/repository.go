package banner

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetAllBanners(ctx context.Context) ([]Banner, error)
	GetBannersByUserId(ctx context.Context, userId string) ([]Banner, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAllBanners(ctx context.Context) ([]Banner, error) {
	var banners []Banner
	err := r.db.Debug().WithContext(ctx).
		Where("user_id IS NULL").
		Find(&banners).Error
	return banners, err
}

func (r *repository) GetBannersByUserId(ctx context.Context, userId string) ([]Banner, error) {
	var banners []Banner
	err := r.db.Debug().WithContext(ctx).Where("user_id = ? OR user_id IS NULL", userId).Find(&banners).Error
	return banners, err
}
