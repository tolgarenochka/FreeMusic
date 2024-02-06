package main

import (
	"context"
	"fmt"

	"FreeMusic/internal/config"
	"FreeMusic/internal/models"
	"FreeMusic/internal/repository/mongodb"

	"github.com/mailru/easyjson/buffer"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const configPath = "./configs/local_config.json"

// main ...
func main() {
	config, err := config.InitConfig(configPath)
	if err != nil {
		logrus.Fatalf("error initializing configs: %v", err)
	}

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d",
		config.DbFilesUsername, config.DbFilesPassword,
		config.DbFilesHost, config.DbFilesPort))
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logrus.Fatalf("error getting client: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		logrus.Fatalf("error pinging client: %v", err)
	}

	defer client.Disconnect(context.Background())

	// Access the database and collection
	database := client.Database(config.DbFilesName)

	mongoFileStorage, err := mongodb.NewMongoFileStorage(config)
	if err != nil {
		logrus.Fatalf("error getting NewMongoFileStorage: %v", err)
	}

	a := buffer.Buffer{Buf: make([]byte, 12345)}

	req := models.UploadFileRequest{
		File:          a.ReadCloser(),
		FileName:      "init_file",
		FileExtension: "txt",
		UserID:        0,
	}

	_, err = mongoFileStorage.UploadFile(context.Background(), req)
	if err != nil {
		logrus.Fatalf("error upload test file: %v", err)
	}

	// Define the index model
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"user_id", 1},
			{"file_name", 1},
		},
		Options: options.Index().SetName("unique_files_index").SetUnique(true),
	}

	collection := database.Collection(config.DBFileCollectionName)

	// Create the index
	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		logrus.Fatalf("error create index: %v", err)
	}

	logrus.Println("Index created successfully.")
}
