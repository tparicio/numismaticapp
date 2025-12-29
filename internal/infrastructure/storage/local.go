package storage

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type LocalFileStorage struct {
	BaseDir string
}

func NewLocalFileStorage(baseDir string) *LocalFileStorage {
	return &LocalFileStorage{BaseDir: baseDir}
}

func (s *LocalFileStorage) EnsureDir(coinID uuid.UUID) (string, error) {
	path := filepath.Join(s.BaseDir, "coins", coinID.String())
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}
	return path, nil
}

// DeleteCoinDirectory removes the entire directory for a coin
func (s *LocalFileStorage) DeleteCoinDirectory(coinID uuid.UUID) error {
	dir := filepath.Join(s.BaseDir, "coins", coinID.String())

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Directory doesn't exist, nothing to delete
		return nil
	}

	// Remove the entire coin directory
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to delete coin directory: %w", err)
	}

	return nil
}

func (s *LocalFileStorage) SaveFile(coinID uuid.UUID, filename string, content io.Reader) (string, error) {
	dir, err := s.EnsureDir(coinID)
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(dir, filename)
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if err := dst.Close(); err != nil {
			slog.Error("Failed to close destination file", "path", fullPath, "error", err)
		}
	}()

	if _, err := io.Copy(dst, content); err != nil {
		return "", fmt.Errorf("failed to save content: %w", err)
	}

	return fullPath, nil
}

func (s *LocalFileStorage) SaveGroupFile(groupID int, filename string, content io.Reader) (string, error) {
	dir := filepath.Join(s.BaseDir, "groups", fmt.Sprintf("%d", groupID))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	fullPath := filepath.Join(dir, filename)
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if err := dst.Close(); err != nil {
			slog.Error("Failed to close destination file", "path", fullPath, "error", err)
		}
	}()

	if _, err := io.Copy(dst, content); err != nil {
		return "", fmt.Errorf("failed to save content: %w", err)
	}

	return fullPath, nil
}
