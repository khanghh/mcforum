package service

import (
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"
)

var UserScoreLogService = newUserScoreLogService()

func newUserScoreLogService() *userScoreLogService {
	return &userScoreLogService{}
}

type userScoreLogService struct {
}

func (s *userScoreLogService) Get(id int64) *model.UserScoreLog {
	return repository.UserScoreLogRepository.Get(sqls.DB(), id)
}

func (s *userScoreLogService) Take(where ...interface{}) *model.UserScoreLog {
	return repository.UserScoreLogRepository.Take(sqls.DB(), where...)
}

func (s *userScoreLogService) Find(cnd *sqls.Cnd) []model.UserScoreLog {
	return repository.UserScoreLogRepository.Find(sqls.DB(), cnd)
}

func (s *userScoreLogService) FindOne(cnd *sqls.Cnd) *model.UserScoreLog {
	return repository.UserScoreLogRepository.FindOne(sqls.DB(), cnd)
}

func (s *userScoreLogService) FindPageByParams(params *params.QueryParams) (list []model.UserScoreLog, paging *sqls.Paging) {
	return repository.UserScoreLogRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userScoreLogService) FindPageByCnd(cnd *sqls.Cnd) (list []model.UserScoreLog, paging *sqls.Paging) {
	return repository.UserScoreLogRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userScoreLogService) Create(t *model.UserScoreLog) error {
	return repository.UserScoreLogRepository.Create(sqls.DB(), t)
}

func (s *userScoreLogService) Update(t *model.UserScoreLog) error {
	return repository.UserScoreLogRepository.Update(sqls.DB(), t)
}

func (s *userScoreLogService) Updates(id int64, columns map[string]interface{}) error {
	return repository.UserScoreLogRepository.Updates(sqls.DB(), id, columns)
}

func (s *userScoreLogService) UpdateColumn(id int64, name string, value interface{}) error {
	return repository.UserScoreLogRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *userScoreLogService) Delete(id int64) {
	repository.UserScoreLogRepository.Delete(sqls.DB(), id)
}
