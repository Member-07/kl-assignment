package debitcard

import (
	"kl-assignment/pkg/response"
	"kl-assignment/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetCardsByUserId(c *gin.Context)
	UpdateCardStatus(c *gin.Context)
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) GetCardsByUserId(c *gin.Context) {
	userId := c.Param("user_id")
	if userId == "" {
		c.Error(response.NewBadRequestError("User ID is required"))
		return
	}

	result, err := h.service.GetCardsByUserId(c.Request.Context(), userId)
	if err != nil {
		c.Error(response.NewBadRequestError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.MappingStandardResponse(response.SUCCESS, result))

}

func (h *handler) UpdateCardStatus(c *gin.Context) {
	var req UpdateCardStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(response.NewBadRequestError("Invalid request body"))
		return
	}

	if err := validation.Struct(req); err != nil {
		c.Error(response.NewBadRequestError(err.Error()))
		return
	}

	err := h.service.UpdateCardStatus(c.Request.Context(), req)
	if err != nil {
		c.Error(response.NewBadRequestError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.MappingStandardResponse(response.SUCCESS, ""))
}
