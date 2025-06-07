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
	
	// Initialize Cloudinary if using it
	if cfg.Provider == "cloudinary" && cfg.CloudName != "" && cfg.APIKey != "" && cfg.APISecret != "" {
		cld, err := cloudinary.NewFromParams(cfg.CloudName, cfg.APIKey, cfg.APISecret)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Cloudinary: %v", err)
		}
		service.cloudinary = cld
	}
	
	return service, nil
}

// UploadFile uploads a file to the storage provider
func (s *StorageService) UploadFile(file *multipart.FileHeader, folder string) (string, error) {
	// Use Cloudinary for uploads
	if s.config.Provider == "cloudinary" && s.cloudinary != nil {
		return s.uploadToCloudinary(file, folder)
	}
	
	// Default fallback to local storage
	return s.uploadToLocal(file, folder)
}

// uploadToCloudinary uploads a file to Cloudinary
func (s *StorageService) uploadToCloudinary(file *multipart.FileHeader, folder string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()
	
	// Set upload folder
	uploadFolder := s.config.UploadFolder
	if folder != "" {
		uploadFolder = path.Join(uploadFolder, folder)
	}
	
	// Upload to Cloudinary
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

// uploadToLocal uploads a file to local storage (fallback)
func (s *StorageService) uploadToLocal(file *multipart.FileHeader, folder string) (string, error) {
	// This is a simplified mock implementation for local storage
	// In a real implementation, you would save the file to disk
	return fmt.Sprintf("/uploads/%s/%s", folder, file.Filename), nil
}

// GetFileURL returns the URL for a file
func (s *StorageService) GetFileURL(fileName string, folder string) string {
	if s.config.Provider == "cloudinary" {
		uploadFolder := s.config.UploadFolder
		if folder != "" {
			uploadFolder = path.Join(uploadFolder, folder)
		}
		return fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/%s/%s", 
			s.config.CloudName, uploadFolder, fileName)
	}
	
	// Local storage fallback
	return fmt.Sprintf("/uploads/%s/%s", folder, fileName)
}

// DeleteFile deletes a file from storage
func (s *StorageService) DeleteFile(fileURL string) error {
	if s.config.Provider == "cloudinary" && s.cloudinary != nil {
		// Extract public ID from URL
		parts := strings.Split(fileURL, "/upload/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid Cloudinary URL format")
		}
		
		publicIDWithExt := parts[1]
		publicID := strings.TrimSuffix(publicIDWithExt, path.Ext(publicIDWithExt))
		
		// Delete from Cloudinary
		ctx := context.Background()
		_, err := s.cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
		if err != nil {
			return fmt.Errorf("failed to delete from Cloudinary: %v", err)
		}
		
		return nil
	}
	
	// Local storage fallback (mock implementation)
	return nil
}
