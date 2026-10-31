package banner

import (
	"kl-assignment/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetAllBanners(c *gin.Context)
	GetBannersByUserId(c *gin.Context)
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) GetAllBanners(c *gin.Context) {
	result, err := h.service.GetAllBanners(c.Request.Context())
	if err != nil {
		c.Error(response.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *handler) GetBannersByUserId(c *gin.Context) {
	userId := c.Param("user_id")
	if userId == "" {
		c.Error(response.NewBadRequestError("User ID is required"))
		return
	}

	result, err := h.service.GetBannersByUserId(c.Request.Context(), userId)
	if err != nil {
		c.Error(response.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}
