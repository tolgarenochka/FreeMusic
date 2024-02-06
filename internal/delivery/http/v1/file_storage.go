package v1

import (
	"FreeMusic/internal/metrics"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	appError "FreeMusic/internal/app_errors"
	"FreeMusic/internal/models"
	"FreeMusic/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// uploadFile @Summary UploadFile
// @Tags FileStorage
// @Description upload file
// @Accept multipart/form-data
// @Produce  json
//
// @Param Authorization header string true "Auth header"
// @Param filename formData string true "file name"
// @Param body formData file true "File to upload"
//
// @Success 200 {object} models.UploadFileResponse
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/file/upload [post]
func (h *Handler) uploadFile(c *gin.Context) {
	metrics.HttpRequestsRPS.WithLabelValues("/upload").Observe(1)

	userID, err := getUserId(c)
	if err != nil {
		logrus.Errorf("uploadFile err: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't get userID"})
		return
	}

	fileHeader, err := c.FormFile("body")
	if err != nil {
		logrus.Errorf("uploadFile err: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't file"})
		return
	}

	fileExtension := filepath.Ext(fileHeader.Filename)

	filename := c.PostForm("filename")
	if len(filename) == 0 {
		logrus.Errorf("uploadFile err: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't get filename"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		logrus.Errorf("uploadFile err: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't open file"})
		return
	}
	defer file.Close()

	req := models.UploadFileRequest{
		File:          file,
		FileName:      filename,
		FileExtension: fileExtension,
		UserID:        userID,
	}

	if fileExtension == models.MP3 {
		req.FileImage, req.Duration, req.Artist, err = utils.HandleMP3File(file)
		if err != nil {
			logrus.Errorf("uploadFile err: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't handle mp3 file"})
			return
		}
	}

	resp, err := h.services.UploadFile(req)
	if err != nil {
		logrus.Errorf("UploadFile err: %v", err)

		var appError *appError.DuplicateFound
		if errors.As(err, &appError) {
			c.AbortWithStatusJSON(http.StatusConflict, errorResponse{"file with that name already exists"})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't save file"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// getAllMusicFilesInfo @Summary getAllMusicFilesInfo
// @Tags FileStorage
// @Description get all music files info (name, artist, duration, tag)
//
// @Produce  json
//
// @Param Authorization header string true "Auth header"
//
// @Success 200 {object} models.GetAllMusicFilesInfoResponse
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/file/get-all-music [get]
func (h *Handler) getAllMusicFilesInfo(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		logrus.Errorf("getAllMusicFilesInfo err: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't get userID"})
		return
	}

	resp, err := h.services.GetAllMusicFilesInfo(userID)
	if err != nil {
		logrus.Errorf("getAllMusicFilesInfo err: %v", err)

		var appError *appError.NoData
		if errors.As(err, &appError) {
			c.AbortWithStatusJSON(http.StatusNotFound, errorResponse{"no data"})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't download file"})
		}

		return
	}
	if err != nil {
		logrus.Errorf("getAllMusicFilesInfo err: %v", err)

		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't get music"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// downloadFile @Summary DownloadFile
// @Tags FileStorage
// @Description download file
// @Accept multipart/form-data
// @Produce  json
//
// @Param Authorization header string true "Auth header"
// @Param filename formData string true "file name"
//
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/file/download [post]
func (h *Handler) downloadFile(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		logrus.Errorf("downloadFile err: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't get userID"})
		return
	}

	filename := c.PostForm("filename")
	if len(filename) == 0 {
		logrus.Errorf("downloadFile err: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't get filename"})
		return
	}

	req := models.DownloadFileRequest{
		FileName: filename,
		UserID:   userID,
	}

	resp, err := h.services.DownloadFile(req, models.Any)
	if err != nil {
		logrus.Errorf("downloadFile err: %v", err)

		var appError *appError.FileNotFound
		if errors.As(err, &appError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"file is not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't download file"})
		}

		return
	}
	defer resp.FileStream.Close()

	fullFileName := resp.FileInfo.FileName + resp.FileInfo.FileExtension
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fullFileName))
	c.Writer.Header().Set("Content-Type", "application/octet-stream")

	_, err = io.Copy(c.Writer, resp.FileStream)
	if err != nil {
		logrus.Errorf("downloadFile err: can't copy file to resp stream, err: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't download file"})
		return
	}

	c.Writer.WriteHeader(200)
}

// downloadAudioImage @Summary DownloadAudioImageFile
// @Tags FileStorage
// @Description download audio image
// @Accept multipart/form-data
// @Produce  json
//
// @Param Authorization header string true "Auth header"
// @Param filename formData string true "file name"
//
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/file/download-audio-image [post]
func (h *Handler) downloadAudioImage(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		logrus.Errorf("downloadAudioImage err: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't get userID"})
		return
	}

	filename := c.PostForm("filename")
	if len(filename) == 0 {
		logrus.Errorf("downloadAudioImage err: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't get filename"})
		return
	}

	req := models.DownloadFileRequest{
		FileName: filename,
		UserID:   userID,
	}

	resp, err := h.services.DownloadAudioImageFile(req)
	if err != nil {
		logrus.Errorf("downloadAudioImage err: %v", err)

		var appError *appError.FileNotFound
		if errors.As(err, &appError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"file is not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't download file"})
		}

		return
	}

	fullFileName := resp.FileInfo.FileName + ".jpg"
	c.Header("Content-Type", "image/jpeg")
	http.ServeContent(c.Writer, c.Request, fullFileName, time.Now(), resp.FileBody)
}

// downloadAudio @Summary DownloadAudio
// @Tags FileStorage
// @Description stream audio file
// @Accept multipart/form-data
// @Produce  json
//
// @Param Authorization header string true "Auth header"
// @Param filename formData string true "file name"
//
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/file/download-audio [post]
func (h *Handler) downloadAudio(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		logrus.Errorf("downloadAudio err: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't get userID"})
		return
	}

	filename := c.PostForm("filename")
	if len(filename) == 0 {
		logrus.Errorf("downloadAudio err: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't get filename"})
		return
	}

	req := models.DownloadFileRequest{
		FileName: filename,
		UserID:   userID,
	}

	resp, err := h.services.DownloadFile(req, models.MP3)
	if err != nil {
		logrus.Errorf("downloadAudio err: %v", err)

		var appError *appError.FileNotFound
		if errors.As(err, &appError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"file is not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't download file"})
		}

		return
	}
	defer resp.FileStream.Close()

	c.Header("Content-Type", "audio/mpeg")
	c.Header("filename", resp.FileInfo.FileName)
	c.Header("artist", resp.FileInfo.Artist)
	c.Writer.WriteHeader(http.StatusOK)

	bufferSize := 1024
	buffer := make([]byte, bufferSize)
	for {
		bytesRead, err := resp.FileStream.Read(buffer)
		if err != nil && err.Error() == "EOF" {
			logrus.Infof("downloadAudio succesful")
			break
		}
		if err != nil {
			logrus.Errorf("downloadAudio err: %v", err)
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		c.Writer.Write(buffer[:bytesRead])
		c.Writer.Flush()
	}
}

// dropFile @Summary DropFile
// @Tags FileStorage
// @Description drop file
// @Accept multipart/form-data
// @Produce  json
//
// @Param Authorization header string true "Auth header"
// @Param filename formData string true "file name"
//
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/file/drop [delete]
func (h *Handler) dropFile(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		logrus.Errorf("dropFile err: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't get userID"})
		return
	}

	filename := c.PostForm("filename")
	if len(filename) == 0 {
		logrus.Errorf("dropFile err: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't get filename"})
		return
	}

	req := models.DropFileRequest{
		FileName: filename,
		UserID:   userID,
	}

	err = h.services.DropFile(req)
	if err != nil {
		logrus.Errorf("dropFile err: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"can't drop file"})
		return
	}

	c.JSON(http.StatusOK, okResponse{"ok"})
}
