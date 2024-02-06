package mongodb

import (
	"context"
	"fmt"

	"FreeMusic/internal/config"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoFileStorage ...
type mongoFileStorage struct {
	client             *mongo.Client
	databaseName       string
	fileCollectionName string
}

// NewMongoFileStorage ...
func NewMongoFileStorage(config *config.Config) (*mongoFileStorage, error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d",
		config.DbFilesUsername, config.DbFilesPassword,
		config.DbFilesHost, config.DbFilesPort))
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, errors.Wrap(err, "NewMongoFileStorage error")
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewMongoFileStorage error")
	}

	return &mongoFileStorage{
		client:             client,
		databaseName:       config.DbFilesName,
		fileCollectionName: config.DBFileCollectionName,
	}, nil
}

// Disconnect ...
func (m *mongoFileStorage) Disconnect(ctx context.Context) {
	m.client.Disconnect(ctx)
}
