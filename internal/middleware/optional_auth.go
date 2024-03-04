package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func LiveOptionalAuthentication(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" || header == "Bearer" {
		c.Set("uid", "NOT_HOST")
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
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
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
		c.Set("uid", claims["uid"])
		c.Set("name", claims["name"])
		c.Set("display_name", claims["display_name"])
		c.Set("display_emoji", claims["display_emoji"])
		c.Set("display_color", claims["display_color"])

		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
