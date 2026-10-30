package user

import (
	"kl-assignment/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetUserProfile(c *gin.Context)
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) GetUserProfile(c *gin.Context) {
	userId := c.Param("user_id")

	if userId == "" {
		c.Error(response.NewBadRequestError("INVALID_USER_ID", "User ID is required"))
		return
	}

	result, err := h.service.GetUserProfile(c.Request.Context(), userId)
	if err != nil {
		c.Error(response.NewBadRequestError("INVALID_REQUEST", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.MappingStandardResponse(response.SUCCESS, result))
}
