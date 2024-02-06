package config

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"os"
)

// Config ...
type Config struct {
	AppHost           string `json:"app_host,omitempty"`
	AppPort           uint64 `json:"app_port,omitempty"`
	AppMaxHeaderBytes int    `json:"app_max_header_bytes"`
	AppMaxFileSize    int    `json:"app_max_file_size"`
	AppReadTimeout    int    `json:"app_read_timeout"`
	AppWriteTimeout   int    `json:"app_write_timeout"`

	DbFilesHost          string `json:"db_files_host"`
	DbFilesPort          int    `json:"db_files_port"`
	DbFilesUsername      string `json:"db_files_username"`
	DbFilesPassword      string `json:"db_files_password"`
	DbFilesName          string `json:"db_files_name"`
	DBFileCollectionName string `json:"db_file_collection_name"`
}

// Validate ...
func (s *Config) Validate() error {
	if len(s.AppHost) == 0 {
		return errors.New("AppHost is empty")
	}

	if s.AppReadTimeout <= 0 {
		return errors.New("AppReadTimeout <= 0")
	}

	if s.AppWriteTimeout <= 0 {
		return errors.New("AppWriteTimeout <= 0")
	}

	if s.AppMaxHeaderBytes <= 0 {
		return errors.New("AppMaxHeaderBytes <= 0")
	}

	if s.AppPort <= 0 {
		return errors.New("AppPort <= 0")
	}

	if s.AppMaxFileSize <= 0 {
		return errors.New("AppMaxFileSize <= 0")
	}

	if len(s.DbFilesPassword) == 0 {
		return errors.New("DbFilesPassword is empty")
	}

	if len(s.DbFilesHost) == 0 {
		return errors.New("DbFilesHost is empty")
	}

	if len(s.DbFilesPassword) == 0 {
		return errors.New("DbFilesPassword is empty")
	}

	if s.DbFilesPort <= 0 {
		return errors.New("DbFilesPort <= 0")
	}

	if len(s.DbFilesName) == 0 {
		return errors.New("DbFilesName is empty")
	}

	if len(s.DBFileCollectionName) == 0 {
		return errors.New("DBFileCollectionName is empty")
	}

	return nil
}

// InitConfig ...
func InitConfig(path string) (*Config, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		logrus.Errorf("Error reading JSON file: %v", err)
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(fileContent, &config); err != nil {
		logrus.Errorf("Error unmarshaling JSON: %v", err)
		return nil, err
	}

	if err := config.Validate(); err != nil {
		logrus.Errorf("Error unmarshaling JSON: %v", err)
		return nil, err
	}

	return &config, nil
}
