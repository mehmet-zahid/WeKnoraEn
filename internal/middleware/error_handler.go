package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Tencent/WeKnora/internal/errors"
)

// ErrorHandler is a middleware that handles application errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there are errors
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last().Err

			// Check if it's an application error
			if appErr, ok := errors.IsAppError(err); ok {
				// Return application error
				c.JSON(appErr.HTTPCode, gin.H{
					"success": false,
					"error": gin.H{
						"code":    appErr.Code,
						"message": appErr.Message,
						"details": appErr.Details,
					},
				})
				return
			}

			// Handle other types of errors
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    errors.ErrInternalServer,
					"message": "Internal server error",
				},
			})
		}
	}
}
