package account_test

import (
	"testing"

	"kl-assignment/app/account"

	"github.com/stretchr/testify/assert"
)

func TestApp_New(t *testing.T) {
	db, _ := setupMockDB(t)
	logger := newTestLogger()

	h := account.New(db, logger)
	assert.NotNil(t, h)
}
