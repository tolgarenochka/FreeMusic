package mongodb

import (
	"bytes"
	"io"

	appError "FreeMusic/internal/app_errors"
	"FreeMusic/internal/models"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"golang.org/x/net/context"
)

// UploadFile ...
func (m *mongoFileStorage) UploadFile(ctx context.Context, req models.UploadFileRequest) (*models.UploadFileResponse, error) {
	db := m.client.Database(m.databaseName)

	fileIDHex, err := m.saveFileInFileStorage(db, req)
	if err != nil {
		return nil, errors.Wrap(err, "UploadFile error")
	}
	defer m.dropFileIfError(ctx, db, fileIDHex, err)

	fileImageIDHex, err := m.saveFileImageInFileStorage(db, req)
	if err != nil {

		return nil, errors.Wrap(err, "UploadFile error")
	}
	defer m.dropFileIfError(ctx, db, fileImageIDHex, err)

	resultID, err := m.saveInfoAboutFile(ctx, db, req, fileIDHex, fileImageIDHex)
	if err != nil {
		return nil, errors.Wrap(err, "UploadFile error")
	}

	return &models.UploadFileResponse{IDHex: resultID}, nil
}

// dropFileIfError ...
func (m *mongoFileStorage) dropFileIfError(ctx context.Context, db *mongo.Database, IDHex string, err error) {
	if err != nil {
		m.dropFileInFileStorage(ctx, db, IDHex)
	}
}

// saveFileImageInFileStorage ...
func (m *mongoFileStorage) saveFileImageInFileStorage(db *mongo.Database, req models.UploadFileRequest) (string, error) {
	fs, err := gridfs.NewBucket(db)
	if err != nil {
		return "", errors.Wrap(err, "saveFileInFileStorage: can't get bucket")
	}

	uploadStream, err := fs.OpenUploadStream(req.FileName + "_image")
	if err != nil {
		return "", errors.Wrap(err, "saveFileImageInFileStorage: can't get upload stream")
	}
	defer uploadStream.Close()

	fileImageReader := bytes.NewReader(req.FileImage)
	_, err = io.Copy(uploadStream, fileImageReader)
	if err != nil {
		return "", errors.Wrap(err, "saveFileImageInFileStorage: can't save file")
	}

	return uploadStream.FileID.(primitive.ObjectID).Hex(), nil
}

// saveFileInFileStorage ...
func (m *mongoFileStorage) saveFileInFileStorage(db *mongo.Database, req models.UploadFileRequest) (string, error) {
	fs, err := gridfs.NewBucket(db)
	if err != nil {
		return "", errors.Wrap(err, "saveFileInFileStorage: can't get bucket")
	}

	uploadStream, err := fs.OpenUploadStream(req.FileName)
	if err != nil {
		return "", errors.Wrap(err, "saveFileInFileStorage: can't get upload stream")
	}
	defer uploadStream.Close()

	_, err = io.Copy(uploadStream, req.File)
	if err != nil {
		return "", errors.Wrap(err, "saveFileInFileStorage: can't save file")
	}

	return uploadStream.FileID.(primitive.ObjectID).Hex(), nil
}

// saveInfoAboutFile ...
func (m *mongoFileStorage) saveInfoAboutFile(ctx context.Context, db *mongo.Database, req models.UploadFileRequest, fileIDHex, fileImageIDHex string) (string, error) {
	document := bson.M{
		"user_id":           req.UserID,
		"file_id_hex":       fileIDHex,
		"file_extension":    req.FileExtension,
		"file_name":         req.FileName,
		"file_image_id_hex": fileImageIDHex,
	}

	if len(req.Artist) != 0 {
		document["artist"] = req.Artist
	}

	if len(req.Duration) != 0 {
		document["duration"] = req.Duration
	}

	collection := db.Collection(m.fileCollectionName)

	result, err := collection.InsertOne(ctx, document)

	if err != nil {
		if writeException, ok := err.(mongo.WriteException); ok {
			for _, writeError := range writeException.WriteErrors {
				if writeError.Code == 11000 {
					// Код ошибки 11000 обозначает дубликацию
					return "", errors.Wrap(&appError.DuplicateFound{
						Message: "file with that name already exists",
					}, "saveInfoAboutFile: duplicate found")
				}
			}
		}
		return "", errors.Wrap(err, "saveInfoAboutFile: can't save info about file")
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
