package service

import (
	"bbs-go/internal/model"
	"bbs-go/internal/repository"
	"bbs-go/sqls"
)

var UploadService = newUploadService()

type uploadService struct {
}

func newUploadService() *uploadService {
	return &uploadService{}
}

func (s *uploadService) RecordUpload(upload *model.Upload) error {
	return repository.UploadRepository.Create(sqls.DB(), upload)
}
