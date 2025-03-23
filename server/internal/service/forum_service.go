package service

import (
	"bbs-go/internal/model/constants"

	"bbs-go/sqls"
	"bbs-go/web/params"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"
)

var ForumService = newForumService()

func newForumService() *forumService {
	return &forumService{}
}

type forumService struct {
}

func (s *forumService) Get(id int64) *model.Forum {
	return repository.ForumRepository.Get(sqls.DB(), id)
}

func (s *forumService) Take(where ...interface{}) *model.Forum {
	return repository.ForumRepository.Take(sqls.DB(), where...)
}

func (s *forumService) Find(cnd *sqls.Cnd) []model.Forum {
	return repository.ForumRepository.Find(sqls.DB(), cnd)
}

func (s *forumService) FindOne(cnd *sqls.Cnd) *model.Forum {
	return repository.ForumRepository.FindOne(sqls.DB(), cnd)
}

func (s *forumService) FindPageByParams(params *params.QueryParams) (list []model.Forum, paging *sqls.Paging) {
	return repository.ForumRepository.FindPageByParams(sqls.DB(), params)
}

func (s *forumService) FindPageByCnd(cnd *sqls.Cnd) (list []model.Forum, paging *sqls.Paging) {
	return repository.ForumRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *forumService) Create(t *model.Forum) error {
	return repository.ForumRepository.Create(sqls.DB(), t)
}

func (s *forumService) Update(t *model.Forum) error {
	return repository.ForumRepository.Update(sqls.DB(), t)
}

func (s *forumService) Updates(id int64, columns map[string]interface{}) error {
	return repository.ForumRepository.Updates(sqls.DB(), id, columns)
}

func (s *forumService) UpdateColumn(id int64, name string, value interface{}) error {
	return repository.ForumRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *forumService) Delete(id int64) {
	repository.ForumRepository.Delete(sqls.DB(), id)
}

func (s *forumService) GetAll() []model.Forum {
	return repository.ForumRepository.Find(sqls.DB(), sqls.NewCnd().Eq("status", constants.StatusOK).Asc("sort_no").Desc("id"))
}
