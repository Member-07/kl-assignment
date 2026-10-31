package transaction_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"kl-assignment/app/transaction"

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

func TestRepository_GetTransactionsByUserId_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := transaction.NewRepository(db)

	name := "Tx 1"
	img := "img.png"
	isBank := true

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	expectedCountSQL := regexp.QuoteMeta("SELECT count(*) FROM `transactions` WHERE user_id = ?")
	mock.ExpectQuery(expectedCountSQL).WithArgs("usr-1").WillReturnRows(countRows)

	dataRows := sqlmock.NewRows([]string{"transaction_id", "user_id", "name", "image", "isBank"}).
		AddRow("tx-1", "usr-1", &name, &img, &isBank)
	expectedFindSQL := regexp.QuoteMeta("SELECT * FROM `transactions` WHERE user_id = ? LIMIT ?")
	mock.ExpectQuery(expectedFindSQL).WithArgs("usr-1", 10).WillReturnRows(dataRows)

	res, total, err := repo.GetTransactionsByUserId(context.Background(), "usr-1", 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, res, 1)
	assert.Equal(t, "tx-1", res[0].TransactionId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetTransactionsByUserId_WithOffset(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := transaction.NewRepository(db)

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(15)
	expectedCountSQL := regexp.QuoteMeta("SELECT count(*) FROM `transactions` WHERE user_id = ?")
	mock.ExpectQuery(expectedCountSQL).WithArgs("usr-1").WillReturnRows(countRows)

	dataRows := sqlmock.NewRows([]string{"transaction_id", "user_id", "name", "image", "isBank"}).
		AddRow("tx-11", "usr-1", "Tx 11", "img.png", false)
	expectedFindSQL := regexp.QuoteMeta("SELECT * FROM `transactions` WHERE user_id = ? LIMIT ? OFFSET ?")
	mock.ExpectQuery(expectedFindSQL).WithArgs("usr-1", 10, 10).WillReturnRows(dataRows)

	res, total, err := repo.GetTransactionsByUserId(context.Background(), "usr-1", 10, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(15), total)
	assert.Len(t, res, 1)
	assert.Equal(t, "tx-11", res[0].TransactionId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetTransactionsByUserId_CountError(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := transaction.NewRepository(db)

	expectedCountSQL := regexp.QuoteMeta("SELECT count(*) FROM `transactions` WHERE user_id = ?")
	mock.ExpectQuery(expectedCountSQL).WithArgs("usr-1").WillReturnError(errors.New("count failed"))

	res, total, err := repo.GetTransactionsByUserId(context.Background(), "usr-1", 10, 0)

	assert.Error(t, err)
	assert.Equal(t, int64(0), total)
	assert.Nil(t, res)
	assert.Equal(t, "count failed", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetTransactionsByUserId_FindError(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := transaction.NewRepository(db)

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	expectedCountSQL := regexp.QuoteMeta("SELECT count(*) FROM `transactions` WHERE user_id = ?")
	mock.ExpectQuery(expectedCountSQL).WithArgs("usr-1").WillReturnRows(countRows)

	expectedFindSQL := regexp.QuoteMeta("SELECT * FROM `transactions` WHERE user_id = ? LIMIT ?")
	mock.ExpectQuery(expectedFindSQL).WithArgs("usr-1", 10).WillReturnError(errors.New("find failed"))

	res, total, err := repo.GetTransactionsByUserId(context.Background(), "usr-1", 10, 0)

	assert.Error(t, err)
	assert.Equal(t, int64(0), total)
	assert.Nil(t, res)
	assert.Equal(t, "find failed", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}
