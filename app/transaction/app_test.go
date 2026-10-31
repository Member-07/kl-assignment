package transaction_test

import (
	"testing"

	"kl-assignment/app/transaction"

	"github.com/stretchr/testify/assert"
)

func TestApp_New(t *testing.T) {
	db, _ := setupMockDB(t)
	logger := newTestLogger()

	h := transaction.New(db, logger)
	assert.NotNil(t, h)
}
