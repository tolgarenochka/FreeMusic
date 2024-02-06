package mongodb

import (
	"bytes"
	"context"

	"FreeMusic/internal/app_errors"
	"FreeMusic/internal/models"

	"github.com/pkg/errors"
)

// DownloadAudioImageFile ...
func (m *mongoFileStorage) DownloadAudioImageFile(ctx context.Context, req models.DownloadFileRequest) (*models.DownloadAudioImageFileResponse, error) {
	db := m.client.Database(m.databaseName)

	fileInfo, err := findIDHexByFileNameAndUserID(ctx, m, db, req, models.MP3)
	if err != nil {
		return nil, errors.Wrap(err, "DownloadAudioImageFile error")
	}

	if fileInfo == nil {
		return nil, &app_errors.FileNotFound{
			Message: "file not found",
		}
	}

	fileStream, err := getFileStreamByFileIDHex(db, fileInfo.FileImageIDHex)
	if err != nil {
		return nil, errors.Wrap(err, "DownloadFile error")
	}

	if fileStream == nil {
		return nil, &app_errors.FileNotFound{
			Message: "file not found",
		}
	}

	var resp models.DownloadAudioImageFileResponse

	resp.FileInfo = fileInfo

	buffer := make([]byte, fileStream.GetFile().Length)
	_, err = fileStream.Read(buffer)
	if err != nil {
		return nil, errors.Wrap(err, "DownloadFile error")
	}
	resp.FileBody = bytes.NewReader(buffer)

	return &resp, nil
}
