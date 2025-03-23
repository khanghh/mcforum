package services

import (
	"bbs-go/internal/models/constants"

	"bbs-go/sqls"
	"bbs-go/web/params"

	"bbs-go/internal/models"
	"bbs-go/internal/repositories"
)

var ForumService = newForumService()

func newForumService() *forumService {
	return &forumService{}
}

type forumService struct {
}

func (s *forumService) Get(id int64) *models.Forum {
	return repositories.ForumRepository.Get(sqls.DB(), id)
}

func (s *forumService) Take(where ...interface{}) *models.Forum {
	return repositories.ForumRepository.Take(sqls.DB(), where...)
}

func (s *forumService) Find(cnd *sqls.Cnd) []models.Forum {
	return repositories.ForumRepository.Find(sqls.DB(), cnd)
}

func (s *forumService) FindOne(cnd *sqls.Cnd) *models.Forum {
	return repositories.ForumRepository.FindOne(sqls.DB(), cnd)
}

func (s *forumService) FindPageByParams(params *params.QueryParams) (list []models.Forum, paging *sqls.Paging) {
	return repositories.ForumRepository.FindPageByParams(sqls.DB(), params)
}

func (s *forumService) FindPageByCnd(cnd *sqls.Cnd) (list []models.Forum, paging *sqls.Paging) {
	return repositories.ForumRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *forumService) Create(t *models.Forum) error {
	return repositories.ForumRepository.Create(sqls.DB(), t)
}

func (s *forumService) Update(t *models.Forum) error {
	return repositories.ForumRepository.Update(sqls.DB(), t)
}

func (s *forumService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ForumRepository.Updates(sqls.DB(), id, columns)
}

func (s *forumService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ForumRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *forumService) Delete(id int64) {
	repositories.ForumRepository.Delete(sqls.DB(), id)
}

func (s *forumService) GetAll() []models.Forum {
	return repositories.ForumRepository.Find(sqls.DB(), sqls.NewCnd().Eq("status", constants.StatusOK).Asc("sort_no").Desc("id"))
}
