package service

import (
	"bbs-go/internal/model"
	"bbs-go/internal/repository"

	"bbs-go/sqls"
	"bbs-go/web/params"
)

var UserReportService = newUserReportService()

func newUserReportService() *userReportService {
	return &userReportService{}
}

type userReportService struct {
}

func (s *userReportService) Get(id int64) *model.UserReport {
	return repository.UserReportRepository.Get(sqls.DB(), id)
}

func (s *userReportService) Take(where ...interface{}) *model.UserReport {
	return repository.UserReportRepository.Take(sqls.DB(), where...)
}

func (s *userReportService) Find(cnd *sqls.Cnd) []model.UserReport {
	return repository.UserReportRepository.Find(sqls.DB(), cnd)
}

func (s *userReportService) FindOne(cnd *sqls.Cnd) *model.UserReport {
	return repository.UserReportRepository.FindOne(sqls.DB(), cnd)
}

func (s *userReportService) FindPageByParams(params *params.QueryParams) (list []model.UserReport, paging *sqls.Paging) {
	return repository.UserReportRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userReportService) FindPageByCnd(cnd *sqls.Cnd) (list []model.UserReport, paging *sqls.Paging) {
	return repository.UserReportRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userReportService) Count(cnd *sqls.Cnd) int64 {
	return repository.UserReportRepository.Count(sqls.DB(), cnd)
}

func (s *userReportService) Create(t *model.UserReport) error {
	return repository.UserReportRepository.Create(sqls.DB(), t)
}

func (s *userReportService) Update(t *model.UserReport) error {
	return repository.UserReportRepository.Update(sqls.DB(), t)
}

func (s *userReportService) Updates(id int64, columns map[string]interface{}) error {
	return repository.UserReportRepository.Updates(sqls.DB(), id, columns)
}

func (s *userReportService) UpdateColumn(id int64, name string, value interface{}) error {
	return repository.UserReportRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *userReportService) Delete(id int64) {
	repository.UserReportRepository.Delete(sqls.DB(), id)
}
