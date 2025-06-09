package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/nyumbanicare/internal/config"
)

// StorageService handles file uploads and retrievals
type StorageService struct {
	config     *config.StorageConfig
	cloudinary *cloudinary.Cloudinary
}

// NewStorageService creates a new storage service
func NewStorageService(cfg *config.StorageConfig) (*StorageService, error) {
	service := &StorageService{
		config: cfg,
	}

	if cfg.Provider == "cloudinary" && cfg.CloudName != "" && cfg.APIKey != "" && cfg.APISecret != "" {
		cld, err := cloudinary.NewFromParams(cfg.CloudName, cfg.APIKey, cfg.APISecret)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Cloudinary: %v", err)
		}
		service.cloudinary = cld
	}

	return service, nil
}

func (s *StorageService) UploadFile(file *multipart.FileHeader, folder string) (string, error) {
	if s.config.Provider == "cloudinary" && s.cloudinary != nil {
		return s.uploadToCloudinary(file, folder)
	}

	return "", fmt.Errorf("storage provider not configured or unsupported: %s", s.config.Provider)
}

func (s *StorageService) uploadToCloudinary(file *multipart.FileHeader, folder string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	uploadFolder := s.config.UploadFolder
	if folder != "" {
		uploadFolder = path.Join(uploadFolder, folder)
	}

	ctx := context.Background()
	uploadParams := uploader.UploadParams{
		PublicID: strings.TrimSuffix(file.Filename, path.Ext(file.Filename)),
		Folder:   uploadFolder,
	}

	uploadResult, err := s.cloudinary.Upload.Upload(ctx, src, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %v", err)
	}

	return uploadResult.SecureURL, nil
}

func (s *StorageService) GetFileURL(fileName string, folder string) string {
	if s.config.Provider == "cloudinary" {
		uploadFolder := s.config.UploadFolder
		if folder != "" {
			uploadFolder = path.Join(uploadFolder, folder)
		}
		return fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/%s/%s",
			s.config.CloudName, uploadFolder, fileName)
	}

	return ""
}

func (s *StorageService) DeleteFile(fileURL string) error {
	if s.config.Provider == "cloudinary" && s.cloudinary != nil {
		parts := strings.Split(fileURL, "/upload/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid Cloudinary URL format")
		}

		publicIDWithExt := parts[1]
		publicID := strings.TrimSuffix(publicIDWithExt, path.Ext(publicIDWithExt))

		ctx := context.Background()
		_, err := s.cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
		if err != nil {
			return fmt.Errorf("failed to delete from Cloudinary: %v", err)
		}

		return nil
	}

	return nil
}
