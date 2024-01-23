package main

import (
	"go-gin-jwt/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes.AuthRoutes(r)

	r.Run(":8080")
}
