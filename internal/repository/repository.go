package repository

import (
	"FreeMusic/internal/models"

	"golang.org/x/net/context"
)

// FileStorage ...
type FileStorage interface {
	UploadFile(ctx context.Context, req models.UploadFileRequest) (*models.UploadFileResponse, error)
	DownloadFile(ctx context.Context, req models.DownloadFileRequest, fileExtension models.FileExtension) (*models.DownloadFileResponse, error)
	DownloadAudioImageFile(ctx context.Context, req models.DownloadFileRequest) (*models.DownloadAudioImageFileResponse, error)
	GetAllMusicFilesInfo(ctx context.Context, userID uint64) (*models.GetAllMusicFilesInfoResponse, error)
	DropFile(ctx context.Context, request models.DropFileRequest) error
}

// Repository ...
type Repository struct {
	FileStorage *FileStorage
}

// NewRepository ...
func NewRepository(fileStorage FileStorage) *Repository {
	return &Repository{
		FileStorage: &fileStorage,
	}
}
