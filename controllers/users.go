package controllers

import (
	"context"
	"go-gin-jwt/database"
	models "go-gin-jwt/models/auth"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (auth *auth) GetAllUsers(c *gin.Context) {
	var users []models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	collection := database.OpenCollection(database.Client, "users")
	cur, err := collection.Find(ctx, bson.D{})
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get all users",
		})
		return
	}
	err = cur.All(context.TODO(), &users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to decode all users",
		})
		return
	}

	c.JSON(http.StatusOK, users)
}
