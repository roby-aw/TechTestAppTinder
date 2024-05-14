package utils

import (
	"context"
	"fmt"

	"roby-backend-golang/config"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DatabaseDriver string

const (
	MongoDB DatabaseDriver = "mongodb"
)

type DatabaseConnection struct {
	Driver DatabaseDriver

	MongoDB     *mongo.Database
	mongoClient *mongo.Client

	AwsS3 *session.Session

	Redis *redis.Client
}

func NewConnectionDatabase(config *config.AppConfig) *DatabaseConnection {
	var db DatabaseConnection
	db.mongoClient = newMongodb(config)
	db.MongoDB = db.mongoClient.Database(config.Database.DBNAME)
	db.Redis = NewRedisClient(config.Database.REDIS_HOST, config.Database.REDIS_PASS)
	db.AwsS3 = InitAwss3(config)

	return &db
}

func newMongodb(config *config.AppConfig) *mongo.Client {
	url := fmt.Sprintf("mongodb://%s:%s@%s%s", config.Database.DBUSER, config.Database.DBPASS, config.Database.DBURL, config.Database.DBNAME)

	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		panic(err)
	}
	return client
}

func (db *DatabaseConnection) CloseConnection() {
	db.mongoClient.Disconnect(context.Background())
}

func NewRedisClient(host, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Redis!")

	return client
}
