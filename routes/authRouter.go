package routes

import (
	"go-gin-jwt/controllers"
	"go-gin-jwt/middleware"

	"github.com/gin-gonic/gin"
)

var auth = controllers.NewAuth()

func AuthRoutes(r *gin.Engine) {
	r.POST("/auth/signup", auth.SignUp)
	r.POST("/auth/login", auth.Login)
	r.GET("/auth/refresh", auth.RefreshAccessToken)
	authGroup := r.Group("/user")
	authGroup.Use(middleware.RequireAuth)
	{
		authGroup.GET("/getall", auth.GetAllUsers)
	}
}
