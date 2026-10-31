package banner_test

import (
	"testing"

	"kl-assignment/app/banner"

	"github.com/stretchr/testify/assert"
)

func TestApp_New(t *testing.T) {
	db, _ := setupMockDB(t)
	logger := newTestLogger()

	h := banner.New(db, logger)
	assert.NotNil(t, h)
}
