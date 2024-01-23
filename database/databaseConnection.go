package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {
	//load env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	//Get mongo connection string
	MongoDB := os.Getenv("MONGO_URL")

	//create client options
	clientOptions := options.Client().ApplyURI(MongoDB)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//connect to mongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	//return mongo client
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collection string) *mongo.Collection {
	col := client.Database("Auth").Collection("collection")
	return col
}
