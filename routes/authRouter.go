package routes

import (
	"go-gin-jwt/controllers"

	"github.com/gin-gonic/gin"
)

var auth = controllers.NewAuth()

func AuthRoutes(r *gin.Engine) {
	r.POST("/auth/signup", auth.SignUp)
	r.POST("/auth/login", auth.Login)
}
