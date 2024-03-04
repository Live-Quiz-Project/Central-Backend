package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AdminRequiredAuthentication(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	accessTokenString := header[len("Bearer "):]
	if accessTokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	accessToken, err := jwt.Parse(accessTokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return []byte(os.Getenv("ADMIN_ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if claims, ok := accessToken.Claims.(jwt.MapClaims); ok && accessToken.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("admin", claims["admin"])

		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
