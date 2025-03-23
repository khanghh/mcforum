package service

import (
	"bbs-go/sqls"
	"bbs-go/web/params"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"
)

var LinkService = newLinkService()

func newLinkService() *linkService {
	return &linkService{}
}

type linkService struct {
}

func (s *linkService) Get(id int64) *model.Link {
	return repository.LinkRepository.Get(sqls.DB(), id)
}

func (s *linkService) Take(where ...interface{}) *model.Link {
	return repository.LinkRepository.Take(sqls.DB(), where...)
}

func (s *linkService) Find(cnd *sqls.Cnd) []model.Link {
	return repository.LinkRepository.Find(sqls.DB(), cnd)
}

func (s *linkService) FindOne(cnd *sqls.Cnd) *model.Link {
	return repository.LinkRepository.FindOne(sqls.DB(), cnd)
}

func (s *linkService) FindPageByParams(params *params.QueryParams) (list []model.Link, paging *sqls.Paging) {
	return repository.LinkRepository.FindPageByParams(sqls.DB(), params)
}

func (s *linkService) FindPageByCnd(cnd *sqls.Cnd) (list []model.Link, paging *sqls.Paging) {
	return repository.LinkRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *linkService) Create(t *model.Link) error {
	return repository.LinkRepository.Create(sqls.DB(), t)
}

func (s *linkService) Update(t *model.Link) error {
	return repository.LinkRepository.Update(sqls.DB(), t)
}

func (s *linkService) Updates(id int64, columns map[string]interface{}) error {
	return repository.LinkRepository.Updates(sqls.DB(), id, columns)
}

func (s *linkService) UpdateColumn(id int64, name string, value interface{}) error {
	return repository.LinkRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *linkService) Delete(id int64) {
	repository.LinkRepository.Delete(sqls.DB(), id)
}
