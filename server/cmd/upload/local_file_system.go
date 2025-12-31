package main

import (
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrPathTraversal = errors.New("path traversal detected")
	ErrInvalidPath   = errors.New("invalid file path")
)

// Local filesystem
type localFileStorage struct {
	baseURL string
	rootDir string
}

func newLocalFileStorage(baseURL, rootDir string) *localFileStorage {
	return &localFileStorage{
		baseURL: baseURL,
		rootDir: rootDir,
	}
}

func (s *localFileStorage) resolve(rel string) (string, error) {
	rel = strings.TrimPrefix(rel, "/")
	joined := filepath.Join(s.rootDir, rel)
	cleaned := filepath.Clean(joined)
	abs, err := filepath.Abs(cleaned)
	if err != nil {
		return "", ErrInvalidPath
	}
	rootWithSep := s.rootDir
	if !strings.HasSuffix(rootWithSep, string(os.PathSeparator)) {
		rootWithSep += string(os.PathSeparator)
	}
	if abs != s.rootDir && !strings.HasPrefix(abs, rootWithSep) {
		return "", ErrPathTraversal
	}
	return url.JoinPath(rootWithSep, rel)
}

func (s *localFileStorage) SaveStream(rel string, reader io.Reader, mimeType string) (string, int, error) {
	dstPath, err := s.resolve(rel)
	if err != nil {
		return "", 0, err
	}

	// ensure parent exists
	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return "", 0, err
	}

	tmpPath := filepath.Base(dstPath) + ".tmp"
	if err := writeTempFile(tmpPath, reader); err != nil {
		return "", 0, err
	}

	if err := os.Rename(tmpPath, dstPath); err != nil {
		return "", 0, err
	}
	os.Chmod(dstPath, 0o644)

	url, err := url.JoinPath(s.baseURL, rel)
	if err != nil {
		return "", 0, err
	}
	info, err := os.Stat(dstPath)
	if err != nil {
		return "", 0, err
	}
	return url, int(info.Size()), nil
}

func (s *localFileStorage) SaveFile(rel string, data []byte, mimeType string) (string, error) {
	dstPath, err := s.resolve(rel)
	if err != nil {
		return "", err
	}
	// ensure parent exists
	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return "", err
	}

	if err := os.WriteFile(dstPath, data, 0o644); err != nil {
		return "", err
	}

	return url.JoinPath(s.baseURL, rel)
}

func (s *localFileStorage) SaveImage(dstPath string, data []byte, mimeType string) (string, error) {
	// TODO: reduce image quality base on mimeType and create thumbnails
	return s.SaveFile(dstPath, data, mimeType)
}

func (s *localFileStorage) Delete(rel string) error {
	dstPath, err := s.resolve(rel)
	if err != nil {
		return err
	}
	return os.Remove(dstPath)
}

func writeTempFile(tmpPath string, reader io.Reader) error {
	f, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(f, reader)
	closeErr := f.Close()
	if copyErr != nil {
		os.Remove(tmpPath)
		return copyErr
	}
	if closeErr != nil {
		os.Remove(tmpPath)
		return closeErr
	}
	return nil
}
