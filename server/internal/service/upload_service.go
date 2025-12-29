package service

import "io"

var UploadService = newUploadService()

func newUploadService() *uploadService {
	return &uploadService{}
}

type uploadService struct {
}

func (s *uploadService) saveStream(r io.Reader) error {
	return nil
}
