package account

import (
	"kl-assignment/pkg/response"
	"kl-assignment/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetAccountsByUserId(c *gin.Context)
	SetMainAccount(c *gin.Context)
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) GetAccountsByUserId(c *gin.Context) {
	userId := c.Param("user_id")
	if userId == "" {
		c.Error(response.NewBadRequestError("User ID is required"))
		return
	}

	result, err := h.service.GetAccountsByUserId(c.Request.Context(), userId)
	if err != nil {
		c.Error(response.NewBadRequestError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.MappingStandardResponse(response.SUCCESS, result))
}

func (h *handler) SetMainAccount(c *gin.Context) {
	var req SetMainAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(response.NewBadRequestError("Invalid request body"))
		return
	}

	if err := validation.Struct(req); err != nil {
		c.Error(response.NewBadRequestError(err.Error()))
		return
	}

	if err := h.service.SetMainAccount(c.Request.Context(), req); err != nil {
		c.Error(response.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.MappingStandardResponse(response.SUCCESS, ""))
}
