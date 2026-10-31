package user_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"kl-assignment/app/user"
	"kl-assignment/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) GetUserProfile(ctx context.Context, userId string) (user.UserResponse, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return user.UserResponse{}, args.Error(1)
	}
	return args.Get(0).(user.UserResponse), args.Error(1)
}

func setupUserRouter(h user.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorResponseMiddleware())

	api := r.Group("/api/user")
	api.GET("", h.GetUserProfile)
	api.GET("/:user_id", h.GetUserProfile)
	return r
}

func TestHandler_GetUserProfile_Success(t *testing.T) {
	svc := new(mockService)
	h := user.NewHandler(svc)
	r := setupUserRouter(h)

	userId := "user-123"
	expected := user.UserResponse{
		UserId:      userId,
		Name:        "Clare",
		GreetingMsg: "Hello",
	}
	svc.On("GetUserProfile", mock.Anything, userId).Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/user/"+userId, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHandler_GetUserProfile_EmptyUserId(t *testing.T) {
	svc := new(mockService)
	h := user.NewHandler(svc)
	r := setupUserRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/user", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_GetUserProfile_ServiceError(t *testing.T) {
	svc := new(mockService)
	h := user.NewHandler(svc)
	r := setupUserRouter(h)

	userId := "user-123"
	svc.On("GetUserProfile", mock.Anything, userId).Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/user/"+userId, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertExpectations(t)
}
