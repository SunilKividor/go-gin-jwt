package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Username     string             `json:"username"`
	Password     string             `json:"password"`
	AccesToken   string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
}
