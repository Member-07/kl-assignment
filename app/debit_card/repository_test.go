package debitcard_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	debitcard "kl-assignment/app/debit_card"

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

func TestRepository_GetCardsByUserId_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := debitcard.NewRepository(db)

	rows := sqlmock.NewRows([]string{
		"card_id", "user_id", "name", "issuer", "number", "color", "border_color", "status",
	}).AddRow("card-1", "usr-1", "Card 1", "VISA", "4111", "#000", "#FFF", "Active")

	expectedSQL := regexp.QuoteMeta("SELECT debit_cards.card_id,debit_cards.user_id,debit_cards.name,debit_card_details.issuer,debit_card_details.number,debit_card_design.color,debit_card_design.border_color,debit_card_status.status FROM `debit_cards` LEFT JOIN debit_card_details ON debit_cards.card_id = debit_card_details.card_id LEFT JOIN debit_card_design ON debit_cards.card_id = debit_card_design.card_id LEFT JOIN debit_card_status ON debit_cards.card_id = debit_card_status.card_id WHERE debit_cards.user_id = ?")
	mock.ExpectQuery(expectedSQL).WithArgs("usr-1").WillReturnRows(rows)

	res, err := repo.GetCardsByUserId(context.Background(), "usr-1")

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "card-1", res[0].CardID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetCardsByUserId_Error(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := debitcard.NewRepository(db)

	expectedSQL := regexp.QuoteMeta("SELECT debit_cards.card_id,debit_cards.user_id,debit_cards.name,debit_card_details.issuer,debit_card_details.number,debit_card_design.color,debit_card_design.border_color,debit_card_status.status FROM `debit_cards` LEFT JOIN debit_card_details ON debit_cards.card_id = debit_card_details.card_id LEFT JOIN debit_card_design ON debit_cards.card_id = debit_card_design.card_id LEFT JOIN debit_card_status ON debit_cards.card_id = debit_card_status.card_id WHERE debit_cards.user_id = ?")
	mock.ExpectQuery(expectedSQL).WithArgs("usr-1").WillReturnError(errors.New("query failed"))

	res, err := repo.GetCardsByUserId(context.Background(), "usr-1")

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "query failed", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_UpdateCardStatus_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := debitcard.NewRepository(db)

	req := debitcard.DebitCardStatus{
		CardId: "card-1",
		UserId: "usr-1",
		Status: "InActive",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `debit_card_status` SET `status`=? WHERE card_id = ? AND user_id = ?")).
		WithArgs("InActive", "card-1", "usr-1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateCardStatus(context.Background(), req)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_UpdateCardStatus_Error(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := debitcard.NewRepository(db)

	req := debitcard.DebitCardStatus{
		CardId: "card-1",
		UserId: "usr-1",
		Status: "InActive",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `debit_card_status` SET `status`=? WHERE card_id = ? AND user_id = ?")).
		WithArgs("InActive", "card-1", "usr-1").
		WillReturnError(errors.New("update error"))
	mock.ExpectRollback()

	err := repo.UpdateCardStatus(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}
