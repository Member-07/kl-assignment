package banner_test

import (
	"context"
	"regexp"
	"testing"

	"kl-assignment/app/banner"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestRepository_GetAllBanners_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := banner.NewRepository(db)

	rows := sqlmock.NewRows([]string{"banner_id", "user_id", "title", "description", "image"}).
		AddRow("b-1", nil, "Title 1", "Desc 1", "img1.png")

	expectedSQL := regexp.QuoteMeta("SELECT * FROM `banners` WHERE user_id IS NULL")
	mock.ExpectQuery(expectedSQL).WillReturnRows(rows)

	res, err := repo.GetAllBanners(context.Background())

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "b-1", res[0].BannerId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetBannersByUserId_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := banner.NewRepository(db)

	userId := "user-1"
	rows := sqlmock.NewRows([]string{"banner_id", "user_id", "title", "description", "image"}).
		AddRow("b-2", userId, "Title 2", "Desc 2", "img2.png")

	expectedSQL := regexp.QuoteMeta("SELECT * FROM `banners` WHERE user_id = ? OR user_id IS NULL")
	mock.ExpectQuery(expectedSQL).WithArgs(userId).WillReturnRows(rows)

	res, err := repo.GetBannersByUserId(context.Background(), userId)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "b-2", res[0].BannerId)
	assert.NoError(t, mock.ExpectationsWereMet())
}
