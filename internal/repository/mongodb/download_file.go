package mongodb

import (
	"context"

	"FreeMusic/internal/app_errors"
	"FreeMusic/internal/models"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// DownloadFile ...
func (m *mongoFileStorage) DownloadFile(ctx context.Context, req models.DownloadFileRequest, fileExtension models.FileExtension) (*models.DownloadFileResponse, error) {
	db := m.client.Database(m.databaseName)

	fileInfo, err := findIDHexByFileNameAndUserID(ctx, m, db, req, fileExtension)
	if err != nil {
		return nil, errors.Wrap(err, "DownloadFile error")
	}

	if fileInfo == nil {
		return nil, &app_errors.FileNotFound{
			Message: "file not found",
		}
	}

	fileStream, err := getFileStreamByFileIDHex(db, fileInfo.FileIDHex)
	if err != nil {
		return nil, errors.Wrap(err, "DownloadFile error")
	}

	if fileStream == nil {
		return nil, &app_errors.FileNotFound{
			Message: "file not found",
		}
	}

	var resp models.DownloadFileResponse
	resp.FileInfo = fileInfo
	resp.FileStream = fileStream

	return &resp, nil
}

// findIDHexByFileNameAndUserID ...
func findIDHexByFileNameAndUserID(ctx context.Context, m *mongoFileStorage, db *mongo.Database, req models.DownloadFileRequest, fileExtension models.FileExtension) (*models.FileInfo, error) {
	collection := db.Collection(m.fileCollectionName)
	var filter primitive.M
	if fileExtension == models.Any {
		filter = bson.M{
			"file_name": req.FileName,
			"user_id":   req.UserID,
		}
	} else {
		filter = bson.M{
			"file_name":      req.FileName,
			"user_id":        req.UserID,
			"file_extension": fileExtension,
		}
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "findIDHexByFileNameAndUserID: can't get cursor on collections")
	}
	defer cursor.Close(ctx)

	var fileInfo models.FileInfo
	var isGetDataFromDB bool
	for cursor.Next(ctx) {
		isGetDataFromDB = true
		if err := cursor.Decode(&fileInfo); err != nil {
			return nil, errors.Wrap(err, "findIDHexByFileNameAndUserID: can't decode file from db")
		}
		break
	}

	if err = cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "findIDHexByFileNameAndUserID: cursor error")
	}

	if !isGetDataFromDB {
		return nil, errors.Wrap(err, "findIDHexByFileNameAndUserID: can't get data from db")
	}

	if fileInfo.UserID == 0 {
		return nil, errors.Wrap(err, "findIDHexByFileNameAndUserID: can't bad data from collections 'files'")
	}

	return &fileInfo, nil
}

// getFileStreamByFileIDHex ...
func getFileStreamByFileIDHex(db *mongo.Database, fileIDHex string) (*gridfs.DownloadStream, error) {
	if fileIDHex == "" {
		return nil, errors.Wrap(nil, "getFileStreamByFileIDHex: get empty fileIDHex")
	}

	fs, err := gridfs.NewBucket(db)
	if err != nil {
		return nil, errors.Wrap(err, "getFileStreamByFileIDHex: can't get bucket")
	}

	fileID, err := primitive.ObjectIDFromHex(fileIDHex)
	if err != nil {
		return nil, errors.Wrap(err, "getFileStreamByFileIDHex: can't convert FileIDHex to primitive.ObjectIDFromHex")
	}

	fileStream, err := fs.OpenDownloadStream(fileID)
	if err != nil {
		return nil, errors.Wrap(err, "getFileStreamByFileIDHex: can't get file stream")
	}

	return fileStream, nil
}
