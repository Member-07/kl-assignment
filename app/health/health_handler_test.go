package health_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"kl-assignment/app/health"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestApp_New(t *testing.T) {
	h := health.New()
	assert.NotNil(t, h)
}

func TestHealthHandler_Check(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := health.NewHandler()
	r.GET("/health", h.Check)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}
