package storage

import (
	"fmt"
	"io"
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
	defer dst.Close()

	if _, err := io.Copy(dst, content); err != nil {
		return "", fmt.Errorf("failed to save content: %w", err)
	}

	return fullPath, nil
}
