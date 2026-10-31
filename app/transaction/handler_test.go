package transaction_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"kl-assignment/app/transaction"
	"kl-assignment/pkg/middleware"
	"kl-assignment/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) GetTransactionsByUserId(ctx context.Context, userId string, page, limit int) (transaction.TransactionPaginationResponse, error) {
	args := m.Called(ctx, userId, page, limit)
	if args.Get(0) == nil {
		return transaction.TransactionPaginationResponse{}, args.Error(1)
	}
	return args.Get(0).(transaction.TransactionPaginationResponse), args.Error(1)
}

func setupTransactionRouter(h transaction.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorResponseMiddleware())

	api := r.Group("/api/transaction")
	api.GET("/", h.GetTransactionsByUserId)
	api.GET("/:user_id", h.GetTransactionsByUserId)
	return r
}

func TestHandler_GetTransactionsByUserId_Success(t *testing.T) {
	svc := new(mockService)
	h := transaction.NewHandler(svc)
	r := setupTransactionRouter(h)

	expected := transaction.TransactionPaginationResponse{
		Data: []transaction.TransactionResponse{
			{
				TransactionId: "tx-1",
				Name:          "Salary",
				IsBank:        true,
			},
		},
		Pagination: response.Pagination{
			Page:       2,
			Limit:      5,
			TotalRows:  12,
			TotalPages: 3,
		},
	}
	svc.On("GetTransactionsByUserId", mock.Anything, "usr-001", 2, 5).Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/transaction/usr-001?page=2&limit=5", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "tx-1")
	assert.Contains(t, w.Body.String(), "Salary")
	svc.AssertExpectations(t)
}

func TestHandler_GetTransactionsByUserId_DefaultQuery(t *testing.T) {
	svc := new(mockService)
	h := transaction.NewHandler(svc)
	r := setupTransactionRouter(h)

	expected := transaction.TransactionPaginationResponse{
		Data: []transaction.TransactionResponse{},
		Pagination: response.Pagination{
			Page:       1,
			Limit:      10,
			TotalRows:  0,
			TotalPages: 0,
		},
	}
	svc.On("GetTransactionsByUserId", mock.Anything, "usr-001", 1, 10).Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/transaction/usr-001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHandler_GetTransactionsByUserId_MissingUserId(t *testing.T) {
	svc := new(mockService)
	h := transaction.NewHandler(svc)
	r := setupTransactionRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/transaction/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "User ID is required")
}

func TestHandler_GetTransactionsByUserId_ServiceError(t *testing.T) {
	svc := new(mockService)
	h := transaction.NewHandler(svc)
	r := setupTransactionRouter(h)

	svc.On("GetTransactionsByUserId", mock.Anything, "usr-001", 1, 10).Return(nil, errors.New("db error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/transaction/usr-001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "db error")
	svc.AssertExpectations(t)
}
