package debitcard_test

import (
	"testing"

	debitcard "kl-assignment/app/debit_card"

	"github.com/stretchr/testify/assert"
)

func TestApp_New(t *testing.T) {
	db, _ := setupMockDB(t)
	logger := newTestLogger()

	h := debitcard.New(db, logger)
	assert.NotNil(t, h)
}
