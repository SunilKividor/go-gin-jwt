package controllers

import (
	"context"
	"go-gin-jwt/database"
	"go-gin-jwt/helper"
	"go-gin-jwt/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type auth struct{}

func NewAuth() *auth {
	return &auth{}
}

func hashpassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (auth *auth) SignUp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Get the email and passoword from body
	var body models.RequestBody
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//hash the password
	hash, err := hashpassword(body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	//generate all tokens
	accessToken, refeshToken, err := helper.GenerateAllTokens(body.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create tokens",
		})
		return
	}

	//create the user
	user := &models.User{
		Username:     body.Username,
		Password:     hash,
		AccesToken:   accessToken,
		RefreshToken: refeshToken,
	}
	//save user in databaase
	collection := database.OpenCollection(database.Client, "users")
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to signup user",
		})
		return
	}
	res := models.ResponseBody{
		Username:     user.Username,
		AccesToken:   user.AccesToken,
		RefreshToken: user.RefreshToken,
	}
	//respond
	c.JSON(http.StatusOK, res)
}
