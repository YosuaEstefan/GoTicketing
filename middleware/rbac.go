// middleware/rbac.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

)

func Authorize(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Check if user has required role
		authorized := false
		for _, role := range roles {
			if userRole == role {
				authorized = true
				break
			}
		}

		if !authorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
