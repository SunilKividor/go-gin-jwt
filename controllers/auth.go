package controllers

import (
	"context"
	"go-gin-jwt/database"
	"go-gin-jwt/helper"
	models "go-gin-jwt/models/auth"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

func verifyPassword(required string, provided string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(required), []byte(provided))
	return err == nil
}

func (auth *auth) SignUp(c *gin.Context) {

	//Get the email and passoword from body
	var body models.RequestBody
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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

func (auth *auth) Login(c *gin.Context) {

	//get login credentials
	var body models.RequestBody
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error reading the body",
		})
		return
	}

	//find the user in database
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	collection := database.OpenCollection(database.Client, "users")
	filter := bson.D{{Key: "username", Value: body.Username}}
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no such user exist",
		})
		return
	}

	//hash and check the password
	if !verifyPassword(user.Password, body.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect password",
		})
		return
	}

	//send access and refresh tokens
	res := models.ResponseBody{
		Username:     user.Username,
		AccesToken:   user.AccesToken,
		RefreshToken: user.RefreshToken,
	}

	c.JSON(http.StatusOK, res)
}
