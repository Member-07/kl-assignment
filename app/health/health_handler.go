package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Check(c *gin.Context)
}

type handler struct {
}

func NewHandler() Handler {
	return &handler{}
}

func (h *handler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().UTC(),
	})
}
