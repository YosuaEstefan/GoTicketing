// middleware/jwt_auth.go (yang diperbarui)
package middleware

import (
	"net/http"
	"strings"
	"ticket/utils"

	"github.com/gin-gonic/gin"

)

// JWTMiddleware interface mendefinisikan fungsi untuk autentikasi JWT
type JWTMiddleware interface {
	JWTAuth() gin.HandlerFunc
}

type jwtMiddleware struct {
	jwtUtil utils.JWTUtil
}

// NewJWTMiddleware menciptakan instance baru middleware JWT
func NewJWTMiddleware(jwtUtil utils.JWTUtil) JWTMiddleware {
	return &jwtMiddleware{
		jwtUtil: jwtUtil,
	}
}

func (m *jwtMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Header Authorization diperlukan"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token auth tidak valid"})
			c.Abort()
			return
		}

		token, err := m.jwtUtil.ValidateToken(parts[1])
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token auth tidak valid"})
			c.Abort()
			return
		}

		// Mengekstrak ID pengguna dan peran dari token
		userID, err := m.jwtUtil.GetUserIDFromToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Info pengguna dalam token tidak valid"})
			c.Abort()
			return
		}

		role, err := m.jwtUtil.GetUserRoleFromToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Info pengguna dalam token tidak valid"})
			c.Abort()
			return
		}

		// Menetapkan ID pengguna dan peran dalam konteks
		c.Set("userID", userID)
		c.Set("userRole", role)

		c.Next()
	}
}
