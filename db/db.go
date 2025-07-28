package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"sync"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
const database= "db"

var (
	client     *mongo.Client
	clientOnce sync.Once
)

func getClient() (*mongo.Client, error) {
	var err error
	clientOnce.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

	
		uri := os.Getenv("mongo_url")
		
		clientOptions := options.Client().ApplyURI(uri)
		
		client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			return
		}
		
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			return
		}
		
		fmt.Println("Connected to MongoDB!")
	})
	
	return client, err
}

func GetUsersCollection() (*mongo.Collection) {
	client, err := getClient()
	if err != nil {
		return nil
	}
	return client.Database(database).Collection("users")
}

func GetTasksCollection() (*mongo.Collection) {
	client, err := getClient()
	if err != nil {
		return nil
	}
	return client.Database(database).Collection("tasks")
}