package debitcard_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	debitcard "kl-assignment/app/debit_card"
	"kl-assignment/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) GetCardsByUserId(ctx context.Context, userId string) ([]debitcard.DebitCardResponse, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]debitcard.DebitCardResponse), args.Error(1)
}

func (m *mockService) UpdateCardStatus(ctx context.Context, req debitcard.UpdateCardStatusRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func setupDebitCardRouter(h debitcard.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorResponseMiddleware())

	api := r.Group("/api/debit-card")
	api.GET("/", h.GetCardsByUserId)
	api.GET("/:user_id", h.GetCardsByUserId)
	api.POST("/status", h.UpdateCardStatus)
	return r
}

func TestHandler_GetCardsByUserId_Success(t *testing.T) {
	svc := new(mockService)
	h := debitcard.NewHandler(svc)
	r := setupDebitCardRouter(h)

	expected := []debitcard.DebitCardResponse{
		{
			CardId: "card-001",
			Name:   "Card 1",
			Status: "Active",
		},
	}
	svc.On("GetCardsByUserId", mock.Anything, "usr-001").Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/debit-card/usr-001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHandler_GetCardsByUserId_MissingUserId(t *testing.T) {
	svc := new(mockService)
	h := debitcard.NewHandler(svc)
	r := setupDebitCardRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/debit-card/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "User ID is required")
}

func TestHandler_GetCardsByUserId_ServiceError(t *testing.T) {
	svc := new(mockService)
	h := debitcard.NewHandler(svc)
	r := setupDebitCardRouter(h)

	svc.On("GetCardsByUserId", mock.Anything, "usr-001").Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/debit-card/usr-001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "not found")
	svc.AssertExpectations(t)
}

func TestHandler_UpdateCardStatus_Success(t *testing.T) {
	svc := new(mockService)
	h := debitcard.NewHandler(svc)
	r := setupDebitCardRouter(h)

	bodyReq := debitcard.UpdateCardStatusRequest{
		UserId: "usr-001",
		CardId: "card-001",
		Status: "Active",
	}
	svc.On("UpdateCardStatus", mock.Anything, bodyReq).Return(nil)

	jsonBody, _ := json.Marshal(bodyReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/debit-card/status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHandler_UpdateCardStatus_InvalidJSON(t *testing.T) {
	svc := new(mockService)
	h := debitcard.NewHandler(svc)
	r := setupDebitCardRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/debit-card/status", bytes.NewBufferString("invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestHandler_UpdateCardStatus_ValidationError(t *testing.T) {
	svc := new(mockService)
	h := debitcard.NewHandler(svc)
	r := setupDebitCardRouter(h)

	bodyReq := debitcard.UpdateCardStatusRequest{
		UserId: "usr-001",
		CardId: "card-001",
		Status: "InvalidStatus",
	}

	jsonBody, _ := json.Marshal(bodyReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/debit-card/status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_UpdateCardStatus_ServiceError(t *testing.T) {
	svc := new(mockService)
	h := debitcard.NewHandler(svc)
	r := setupDebitCardRouter(h)

	bodyReq := debitcard.UpdateCardStatusRequest{
		UserId: "usr-001",
		CardId: "card-001",
		Status: "InActive",
	}
	svc.On("UpdateCardStatus", mock.Anything, bodyReq).Return(errors.New("update error"))

	jsonBody, _ := json.Marshal(bodyReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/debit-card/status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "update error")
	svc.AssertExpectations(t)
}
