package middleware

import (
	"kl-assignment/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if appErr, ok := err.(*response.AppError); ok {
				c.JSON(appErr.StatusCode, response.MappingStandardResponse(appErr.Message, appErr.Message))
				return
			}

			c.JSON(http.StatusInternalServerError, response.MappingStandardResponse("INTERNAL_SERVER_ERROR", err.Error()))
			return
		}

	}
}
