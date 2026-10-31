package transaction

import (
	"kl-assignment/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetTransactionsByUserId(c *gin.Context)
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) GetTransactionsByUserId(c *gin.Context) {
	userId := c.Param("user_id")
	if userId == "" {
		c.Error(response.NewBadRequestError("User ID is required"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.service.GetTransactionsByUserId(c.Request.Context(), userId, page, limit)
	if err != nil {
		c.Error(response.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}
