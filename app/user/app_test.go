package user_test

import (
	"testing"

	"kl-assignment/app/user"

	"github.com/stretchr/testify/assert"
)

func TestApp_New(t *testing.T) {
	db, _ := setupMockDB(t)
	logger := newTestLogger()

	h := user.New(db, logger)
	assert.NotNil(t, h)
}
