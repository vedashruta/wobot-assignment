package mongo

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(fileName string) (client *mongo.Client, err error) {
	err = godotenv.Load(fileName)
	if err != nil {
		return
	}
	uri := os.Getenv("MONGO_URI")
	timeoutStr := os.Getenv("MONGO_TIMEOUT")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	clientOptions := options.Client().ApplyURI(uri)
	tempClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}
	err = tempClient.Ping(ctx, nil)
	if err != nil {
		return
	}
	client = tempClient
	return
}
