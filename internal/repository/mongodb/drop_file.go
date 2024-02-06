package mongodb

import (
	"context"

	"FreeMusic/internal/models"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// DropFile ...
func (m *mongoFileStorage) DropFile(ctx context.Context, req models.DropFileRequest) error {
	db := m.client.Database(m.databaseName)

	fileInfo, err := m.dropFileInfo(ctx, db, req)
	if err != nil {
		return errors.Wrap(err, "DropFile error")
	}

	dropFileErr := m.dropFileInFileStorage(ctx, db, fileInfo.FileIDHex)
	if dropFileErr != nil {
		err = errors.Wrap(dropFileErr, "DropFile error")
	}

	dropFileImageErr := m.dropFileInFileStorage(ctx, db, fileInfo.FileImageIDHex)
	if dropFileImageErr != nil {
		err = errors.Wrap(err, "DropFile error")
	}

	return nil
}

// dropFileInfo ...
func (m *mongoFileStorage) dropFileInfo(ctx context.Context, db *mongo.Database, req models.DropFileRequest) (*models.FileInfo, error) {
	filter := bson.M{
		"file_name": req.FileName,
		"user_id":   req.UserID,
	}
	var deletedDocument models.FileInfo

	collection := db.Collection(m.fileCollectionName)

	err := collection.FindOneAndDelete(ctx, filter).Decode(&deletedDocument)
	if err != nil {
		return nil, errors.Wrap(err, "dropFileInfo: can't drop file info")
	}

	return &deletedDocument, nil
}

// dropFileInFileStorage ...
func (m *mongoFileStorage) dropFileInFileStorage(ctx context.Context, db *mongo.Database, fileIDHex string) error {
	fs, err := gridfs.NewBucket(db)
	if err != nil {
		return errors.Wrap(err, "dropFileInFileStorage: can't get bucket")
	}

	fileID, err := primitive.ObjectIDFromHex(fileIDHex)
	if err != nil {
		return errors.Wrap(err, "dropFileInFileStorage: can't convert FileIDHex to primitive.ObjectIDFromHex")
	}

	err = fs.Delete(fileID)
	if err != nil && err.Error() != "file with given parameters not found" {
		return errors.Wrap(err, "dropFileInFileStorage: can't delete file by fileIDHex")
	}

	return nil
}
