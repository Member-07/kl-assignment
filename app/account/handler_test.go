package account_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"kl-assignment/app/account"
	"kl-assignment/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) GetAccountsByUserId(ctx context.Context, userId string) ([]account.AccountResponse, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]account.AccountResponse), args.Error(1)
}

func (m *mockService) SetMainAccount(ctx context.Context, req account.SetMainAccountRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func setupTestRouter(h account.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorResponseMiddleware())

	api := r.Group("/api/account")
	api.GET("", h.GetAccountsByUserId)
	api.GET("/:user_id", h.GetAccountsByUserId)
	api.PUT("/main", h.SetMainAccount)
	return r
}

func TestHandler_GetAccountsByUserId_Success(t *testing.T) {
	svc := new(mockService)
	h := account.NewHandler(svc)
	r := setupTestRouter(h)

	userId := "user-123"
	expected := []account.AccountResponse{
		{AccountId: "acc-1", UserId: userId},
	}
	svc.On("GetAccountsByUserId", mock.Anything, userId).Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/account/"+userId, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHandler_GetAccountsByUserId_EmptyUserId(t *testing.T) {
	svc := new(mockService)
	h := account.NewHandler(svc)
	r := setupTestRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/account", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_GetAccountsByUserId_ServiceError(t *testing.T) {
	svc := new(mockService)
	h := account.NewHandler(svc)
	r := setupTestRouter(h)

	userId := "user-123"
	svc.On("GetAccountsByUserId", mock.Anything, userId).Return(nil, errors.New("something went wrong"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/account/"+userId, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertExpectations(t)
}

func TestHandler_SetMainAccount_Success(t *testing.T) {
	svc := new(mockService)
	h := account.NewHandler(svc)
	r := setupTestRouter(h)

	bodyReq := account.SetMainAccountRequest{
		UserId:    "user-1",
		AccountId: "acc-1",
	}
	svc.On("SetMainAccount", mock.Anything, bodyReq).Return(nil)

	payload, _ := json.Marshal(bodyReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/account/main", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHandler_SetMainAccount_InvalidBody(t *testing.T) {
	svc := new(mockService)
	h := account.NewHandler(svc)
	r := setupTestRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/account/main", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_SetMainAccount_ValidationError(t *testing.T) {
	svc := new(mockService)
	h := account.NewHandler(svc)
	r := setupTestRouter(h)

	bodyReq := account.SetMainAccountRequest{
		UserId: "user-1",
	}

	payload, _ := json.Marshal(bodyReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/account/main", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_SetMainAccount_ServiceError(t *testing.T) {
	svc := new(mockService)
	h := account.NewHandler(svc)
	r := setupTestRouter(h)

	bodyReq := account.SetMainAccountRequest{
		UserId:    "user-1",
		AccountId: "acc-1",
	}
	svc.On("SetMainAccount", mock.Anything, bodyReq).Return(errors.New("db error"))

	payload, _ := json.Marshal(bodyReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/account/main", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}
