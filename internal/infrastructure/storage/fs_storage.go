package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/antonioparicio/numismaticapp/internal/domain"
)

type LocalStorage struct {
	BaseDir string
}

func NewLocalStorage(baseDir string) *LocalStorage {
	return &LocalStorage{BaseDir: baseDir}
}

func (s *LocalStorage) Save(coinID, filename string, data []byte) (string, error) {
	dir := filepath.Join(s.BaseDir, "coins", coinID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return path, nil
}

func (s *LocalStorage) Load(coinID, filename string) ([]byte, error) {
	path := s.GetPath(coinID, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil
}

func (s *LocalStorage) Exists(coinID, filename string) bool {
	path := s.GetPath(coinID, filename)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (s *LocalStorage) GetPath(coinID, filename string) string {
	return filepath.Join(s.BaseDir, "coins", coinID, filename)
}

var _ domain.ImageStorage = (*LocalStorage)(nil)
