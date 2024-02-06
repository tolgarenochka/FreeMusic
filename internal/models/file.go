package models

import (
	"bytes"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// DownloadFileRequest ...
type DownloadFileRequest struct {
	FileName string `json:"filename"`
	UserID   uint64 `json:"user_id"`
}

// DownloadFileResponse ...
type DownloadFileResponse struct {
	FileInfo   *FileInfo              `json:"file_info"`
	FileStream *gridfs.DownloadStream `json:"file_stream"`
}

// DownloadAudioImageFileResponse ...
type DownloadAudioImageFileResponse struct {
	FileInfo *FileInfo     `json:"file_info"`
	FileBody *bytes.Reader `json:"file_body"`
}

// UploadFileRequest ...
type UploadFileRequest struct {
	File          io.ReadCloser `json:"file"`
	FileName      string        `json:"filename"`
	FileImage     []byte        `json:"file_image"`
	FileExtension string        `json:"file_extension"`
	UserID        uint64        `json:"user_id"`
	Artist        string        `json:"artist,omitempty"`
	Duration      string        `json:"duration"`
}

// DropFileRequest ...
type DropFileRequest struct {
	FileName string `json:"filename"`
	UserID   uint64 `json:"user_id"`
}

// UploadFileResponse ...
type UploadFileResponse struct {
	IDHex string `json:"id_hex"`
}

// InfoAboutMusicFile ...
type InfoAboutMusicFile struct {
	Artist   string `json:"artist"`
	FileName string `json:"fileName"`
	Duration string `json:"duration"`
	Tag      string `json:"tag"`
}

// GetAllMusicFilesInfoResponse ...
type GetAllMusicFilesInfoResponse struct {
	InfoAboutMusicFile []InfoAboutMusicFile `json:"info_about_music_file"`
}

// FileInfo ...
type FileInfo struct {
	ID             primitive.ObjectID `bson:"_id"`
	UserID         uint64             `bson:"user_id"`
	FileIDHex      string             `bson:"file_id_hex"`
	FileImageIDHex string             `bson:"file_image_id_hex"`
	FileExtension  string             `bson:"file_extension"`
	FileName       string             `bson:"file_name"`
	Artist         string             `bson:"artist,omitempty"`
	Duration       string             `bson:"duration"`
}
