package user_test

import (
	"context"
	"regexp"
	"testing"

	"kl-assignment/app/user"

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

func TestRepository_GetUserWithGreeting_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := user.NewRepository(db)

	userId := "user-123"
	rows := sqlmock.NewRows([]string{"user_id", "name", "greeting"}).
		AddRow(userId, "Clare", "Have a nice day")

	expectedSQL := regexp.QuoteMeta("SELECT users.user_id,users.name,user_greetings.greeting FROM `users` LEFT JOIN user_greetings ON users.user_id = user_greetings.user_id WHERE users.user_id = ?")
	mock.ExpectQuery(expectedSQL).WithArgs(userId).WillReturnRows(rows)

	res, err := repo.GetUserWithGreeting(context.Background(), userId)

	assert.NoError(t, err)
	assert.Equal(t, userId, res.UserId)
	assert.Equal(t, "Clare", res.Name)
	assert.Equal(t, "Have a nice day", res.Greeting)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetUserWithGreeting_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := user.NewRepository(db)

	userId := "non-existing"
	rows := sqlmock.NewRows([]string{"user_id", "name", "greeting"})

	expectedSQL := regexp.QuoteMeta("SELECT users.user_id,users.name,user_greetings.greeting FROM `users` LEFT JOIN user_greetings ON users.user_id = user_greetings.user_id WHERE users.user_id = ?")
	mock.ExpectQuery(expectedSQL).WithArgs(userId).WillReturnRows(rows)

	_, err := repo.GetUserWithGreeting(context.Background(), userId)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}
