package banner_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"kl-assignment/app/banner"
	"kl-assignment/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) GetAllBanners(ctx context.Context) ([]banner.BannerResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]banner.BannerResponse), args.Error(1)
}

func (m *mockService) GetBannersByUserId(ctx context.Context, userId string) ([]banner.BannerResponse, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]banner.BannerResponse), args.Error(1)
}

func setupBannerRouter(h banner.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorResponseMiddleware())

	api := r.Group("/api/banner")
	api.GET("/", h.GetAllBanners)
	api.GET("/user/", h.GetBannersByUserId)
	api.GET("/:user_id", h.GetBannersByUserId)
	return r
}

func TestHandler_GetAllBanners_Success(t *testing.T) {
	svc := new(mockService)
	h := banner.NewHandler(svc)
	r := setupBannerRouter(h)

	expected := []banner.BannerResponse{{BannerId: "b-1"}}
	svc.On("GetAllBanners", mock.Anything).Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/banner/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHandler_GetAllBanners_ServiceError(t *testing.T) {
	svc := new(mockService)
	h := banner.NewHandler(svc)
	r := setupBannerRouter(h)

	svc.On("GetAllBanners", mock.Anything).Return(nil, errors.New("internal server error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/banner/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestHandler_GetBannersByUserId_Success(t *testing.T) {
	svc := new(mockService)
	h := banner.NewHandler(svc)
	r := setupBannerRouter(h)

	userId := "user-123"
	expected := []banner.BannerResponse{{BannerId: "b-2"}}
	svc.On("GetBannersByUserId", mock.Anything, userId).Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/banner/"+userId, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHandler_GetBannersByUserId_EmptyUserId(t *testing.T) {
	svc := new(mockService)
	h := banner.NewHandler(svc)
	r := setupBannerRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/banner/user/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_GetBannersByUserId_ServiceError(t *testing.T) {
	svc := new(mockService)
	h := banner.NewHandler(svc)
	r := setupBannerRouter(h)

	userId := "user-123"
	svc.On("GetBannersByUserId", mock.Anything, userId).Return(nil, errors.New("internal server error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/banner/"+userId, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}
