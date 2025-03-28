package service

import (
	"bbs-go/internal/model"
	"bbs-go/internal/repository"

	"bbs-go/pkg/web/params"
	"bbs-go/sqls"
)

var EmailCodeService = newEmailCodeService()

func newEmailCodeService() *emailCodeService {
	return &emailCodeService{}
}

type emailCodeService struct {
}

func (s *emailCodeService) Get(id int64) *model.EmailCode {
	return repository.EmailCodeRepository.Get(sqls.DB(), id)
}

func (s *emailCodeService) Take(where ...interface{}) *model.EmailCode {
	return repository.EmailCodeRepository.Take(sqls.DB(), where...)
}

func (s *emailCodeService) Find(cnd *sqls.Cnd) []model.EmailCode {
	return repository.EmailCodeRepository.Find(sqls.DB(), cnd)
}

func (s *emailCodeService) FindOne(cnd *sqls.Cnd) *model.EmailCode {
	return repository.EmailCodeRepository.FindOne(sqls.DB(), cnd)
}

func (s *emailCodeService) FindPageByParams(params *params.QueryParams) (list []model.EmailCode, paging *sqls.Paging) {
	return repository.EmailCodeRepository.FindPageByParams(sqls.DB(), params)
}

func (s *emailCodeService) FindPageByCnd(cnd *sqls.Cnd) (list []model.EmailCode, paging *sqls.Paging) {
	return repository.EmailCodeRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *emailCodeService) Count(cnd *sqls.Cnd) int64 {
	return repository.EmailCodeRepository.Count(sqls.DB(), cnd)
}

func (s *emailCodeService) Create(t *model.EmailCode) error {
	return repository.EmailCodeRepository.Create(sqls.DB(), t)
}

func (s *emailCodeService) Update(t *model.EmailCode) error {
	return repository.EmailCodeRepository.Update(sqls.DB(), t)
}

func (s *emailCodeService) Updates(id int64, columns map[string]interface{}) error {
	return repository.EmailCodeRepository.Updates(sqls.DB(), id, columns)
}

func (s *emailCodeService) UpdateColumn(id int64, name string, value interface{}) error {
	return repository.EmailCodeRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *emailCodeService) Delete(id int64) {
	repository.EmailCodeRepository.Delete(sqls.DB(), id)
}
