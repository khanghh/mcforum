package service

import (
	"bbs-go/internal/model/constants"

	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"
)

var ArticleTagService = newArticleTagService()

func newArticleTagService() *articleTagService {
	return &articleTagService{}
}

type articleTagService struct {
}

func (s *articleTagService) Get(id int64) *model.ArticleTag {
	return repository.ArticleTagRepository.Get(sqls.DB(), id)
}

func (s *articleTagService) Take(where ...interface{}) *model.ArticleTag {
	return repository.ArticleTagRepository.Take(sqls.DB(), where...)
}

func (s *articleTagService) Find(cnd *sqls.Cnd) []model.ArticleTag {
	return repository.ArticleTagRepository.Find(sqls.DB(), cnd)
}

func (s *articleTagService) FindPageByParams(params *params.QueryParams) (list []model.ArticleTag, paging *sqls.Paging) {
	return repository.ArticleTagRepository.FindPageByParams(sqls.DB(), params)
}

func (s *articleTagService) FindPageByCnd(cnd *sqls.Cnd) (list []model.ArticleTag, paging *sqls.Paging) {
	return repository.ArticleTagRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *articleTagService) Create(t *model.ArticleTag) error {
	return repository.ArticleTagRepository.Create(sqls.DB(), t)
}

func (s *articleTagService) Update(t *model.ArticleTag) error {
	return repository.ArticleTagRepository.Update(sqls.DB(), t)
}

func (s *articleTagService) Updates(id int64, columns map[string]interface{}) error {
	return repository.ArticleTagRepository.Updates(sqls.DB(), id, columns)
}

func (s *articleTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return repository.ArticleTagRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *articleTagService) DeleteByArticleId(articleId int64) {
	sqls.DB().Model(model.ArticleTag{}).Where("article_id = ?", articleId).UpdateColumn("status", constants.StatusDeleted)
}
