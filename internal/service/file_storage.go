package service

import (
	"context"

	"FreeMusic/internal/models"
	"FreeMusic/internal/repository"

	"github.com/pkg/errors"
)

// fileManager ...
type fileManager struct {
	repo repository.FileStorage
}

// NewFileManager ...
func NewFileManager(repo *repository.FileStorage) *fileManager {
	return &fileManager{
		repo: *repo,
	}
}

// UploadFile ...
func (f *fileManager) UploadFile(req models.UploadFileRequest) (*models.UploadFileResponse, error) {
	resp, err := f.repo.UploadFile(context.Background(), req)
	if err != nil {
		return nil, errors.Wrap(err, "UploadFile error")
	}

	return resp, nil
}

// DownloadFile ...
func (f *fileManager) DownloadFile(req models.DownloadFileRequest, fileExtension models.FileExtension) (*models.DownloadFileResponse, error) {
	resp, err := f.repo.DownloadFile(context.Background(), req, fileExtension)
	if err != nil {
		return nil, errors.Wrap(err, "DownloadFile error")
	}

	return resp, nil
}

// DownloadAudioImageFile ...
func (f *fileManager) DownloadAudioImageFile(req models.DownloadFileRequest) (*models.DownloadAudioImageFileResponse, error) {
	resp, err := f.repo.DownloadAudioImageFile(context.Background(), req)
	if err != nil {
		return nil, errors.Wrap(err, "DownloadAudioImageFile error")
	}

	return resp, nil
}

// GetAllMusicFilesInfo ...
func (f *fileManager) GetAllMusicFilesInfo(userID uint64) (*models.GetAllMusicFilesInfoResponse, error) {
	resp, err := f.repo.GetAllMusicFilesInfo(context.Background(), userID)
	if err != nil {
		return nil, errors.Wrap(err, "DownloadFile error")
	}

	return resp, nil
}

// DropFile ...
func (f *fileManager) DropFile(req models.DropFileRequest) error {
	err := f.repo.DropFile(context.Background(), req)
	if err != nil {
		return errors.Wrap(err, "DropFile error")
	}

	return nil
}
