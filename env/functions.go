package env

import (
	"os"
	"server/middlewares/jwt"
	"strconv"
	"time"

	mongoDB "go.mongodb.org/mongo-driver/mongo"

	"server/services/mongo"

	"github.com/joho/godotenv"
)

var (
	Name                string
	Port                string
	Storage             string
	Timeout             time.Duration
	DefaultStorageQuota int
)

var (
	FilesCollection *mongoDB.Collection
	UsersCollection *mongoDB.Collection
)

var (
	MongoClient *mongoDB.Client
)

func LoadEnv() (err error) {
	fileName := "env/.env"
	err = godotenv.Load(fileName)
	if err != nil {
		return
	}
	Port = os.Getenv("PORT")
	Storage = os.Getenv("STORAGE")
	timeoutStr := os.Getenv("SERVER_TIMEOUT")
	Timeout, err = time.ParseDuration(timeoutStr)
	if err != nil {
		return
	}
	tempStorageSize := os.Getenv("DEFAULT_STORAGE_QUOTA_MB")
	DefaultStorageQuota, err = strconv.Atoi(tempStorageSize)
	if err != nil {
		return
	}
	err = jwt.Init(fileName)
	if err != nil {
		return
	}
	err = os.MkdirAll(Storage, os.ModePerm)
	if err != nil {
		return
	}
	MongoClient, err = mongo.Connect(fileName)
	if err != nil {
		return
	}
	FilesCollection = MongoClient.Database("filenest").Collection("files")
	UsersCollection = MongoClient.Database("filenest").Collection("users")
	return
}
