package middleware

import (
	"go-gin-jwt/helper"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	signedToken := c.Request.Header.Get("token")
	if signedToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token missing in headers",
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.ParseWithClaims(signedToken, &helper.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(*helper.Claims); ok && token.Valid {
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()

}
