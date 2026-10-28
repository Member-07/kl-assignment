package repository

import "gorm.io/gorm"

type BannerRepository interface {
}

type bannerRepository struct {
	db *gorm.DB
}

func NewBannerRepository(db *gorm.DB) BannerRepository {
	return &bannerRepository{
		db: db,
	}
}

func (r *bannerRepository) GetBanner() {

}
