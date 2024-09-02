package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// mongoUri = "mongodb://localhost:27017"
	mongoUri = os.Getenv("MONGO_URI")
	client   *mongo.Client
	database *mongo.Database
	dbName   = "knowledge_db" // 替换为实际的数据库名称
)

// mongo实现本地连接的方法
func InitDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	database = client.Database(dbName)
}

// docker 实现
func InitDBDocker() {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://root:123456789@localhost:27018/%s?authSource=admin", dbName))

	var err error
	client, err = mongo.NewClient(clientOptions)
	if err != nil {
		fmt.Println("Failed to create new MongoDB client:", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Failed to ping MongoDB:", err)
		return
	}

	database = client.Database(dbName)
	fmt.Println("Connected to MongoDB")
	// 以上是使用docker 来作为mongo的容器
}

func InitDB_docker() {
	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}

	clientOptions := options.Client().ApplyURI(mongoUri)

	var err error
	client, err = mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create new MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	database = client.Database(dbName)
	fmt.Println("Connected to MongoDB")
}

func GetCollection(collectionName string) *mongo.Collection {
	if client == nil {
		log.Fatal("MongoDB client is not initialized")
	}
	return client.Database(dbName).Collection(collectionName)
}
