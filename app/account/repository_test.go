package account_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"kl-assignment/app/account"

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

func TestRepository_GetAccountByUserId_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := account.NewRepository(db)

	userId := "user-1"
	rows := sqlmock.NewRows([]string{
		"account_id", "user_id", "type", "currency", "account_number", "issuer",
		"amount", "color", "is_main_account", "progress",
	}).AddRow("acc-1", userId, "Saving", "THB", "123456", "TestBank", 100.50, "#fff", 1, 50)

	expectedSQL := regexp.QuoteMeta("SELECT accounts.account_id,accounts.user_id,accounts.type,accounts.currency,accounts.account_number,accounts.issuer,account_balances.amount,account_details.color,account_details.is_main_account,account_details.progress FROM `accounts` LEFT JOIN account_balances ON accounts.account_id = account_balances.account_id LEFT JOIN account_details ON accounts.account_id = account_details.account_id WHERE accounts.user_id = ?")
	mock.ExpectQuery(expectedSQL).WithArgs(userId).WillReturnRows(rows)

	res, err := repo.GetAccountByUserId(context.Background(), userId)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "acc-1", res[0].AccountId)
	assert.True(t, res[0].IsMainAccount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAccountByUserId_QueryError(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := account.NewRepository(db)

	userId := "user-1"
	expectedSQL := regexp.QuoteMeta("SELECT accounts.account_id,accounts.user_id,accounts.type,accounts.currency,accounts.account_number,accounts.issuer,account_balances.amount,account_details.color,account_details.is_main_account,account_details.progress FROM `accounts` LEFT JOIN account_balances ON accounts.account_id = account_balances.account_id LEFT JOIN account_details ON accounts.account_id = account_details.account_id WHERE accounts.user_id = ?")
	mock.ExpectQuery(expectedSQL).WithArgs(userId).WillReturnError(errors.New("db error"))

	res, err := repo.GetAccountByUserId(context.Background(), userId)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAccountFlags_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := account.NewRepository(db)

	accountIds := []string{"acc-1"}
	rows := sqlmock.NewRows([]string{"flag_id", "account_id", "user_id", "flag_type", "flag_value", "created_at", "updated_at"}).
		AddRow(1, "acc-1", "user-1", "lock", "Disbursement", time.Now(), time.Now())

	expectedSQL := regexp.QuoteMeta("SELECT * FROM `account_flags` WHERE account_id IN (?)")
	mock.ExpectQuery(expectedSQL).WithArgs("acc-1").WillReturnRows(rows)

	res, err := repo.GetAccountFlags(context.Background(), accountIds)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "lock", res[0].FlagType)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAccountFlags_QueryError(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := account.NewRepository(db)

	accountIds := []string{"acc-1"}
	expectedSQL := regexp.QuoteMeta("SELECT * FROM `account_flags` WHERE account_id IN (?)")
	mock.ExpectQuery(expectedSQL).WithArgs("acc-1").WillReturnError(errors.New("db error"))

	res, err := repo.GetAccountFlags(context.Background(), accountIds)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_SetMainAccount_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := account.NewRepository(db)

	userId := "user-1"
	accountId := "acc-1"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `account_details` SET `is_main_account`=? WHERE user_id = ?")).
		WithArgs(false, userId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta("UPDATE `account_details` SET `is_main_account`=? WHERE account_id = ? AND user_id = ?")).
		WithArgs(true, accountId, userId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.SetMainAccount(context.Background(), userId, accountId)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_SetMainAccount_FirstUpdateError(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := account.NewRepository(db)

	userId := "user-1"
	accountId := "acc-1"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `account_details` SET `is_main_account`=? WHERE user_id = ?")).
		WithArgs(false, userId).
		WillReturnError(errors.New("first update error"))
	mock.ExpectRollback()

	err := repo.SetMainAccount(context.Background(), userId, accountId)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_SetMainAccount_SecondUpdateError(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := account.NewRepository(db)

	userId := "user-1"
	accountId := "acc-1"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `account_details` SET `is_main_account`=? WHERE user_id = ?")).
		WithArgs(false, userId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta("UPDATE `account_details` SET `is_main_account`=? WHERE account_id = ? AND user_id = ?")).
		WithArgs(true, accountId, userId).
		WillReturnError(errors.New("second update error"))
	mock.ExpectRollback()

	err := repo.SetMainAccount(context.Background(), userId, accountId)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
