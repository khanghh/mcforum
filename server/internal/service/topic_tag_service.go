package service

import (
	"bbs-go/internal/model/constants"

	"bbs-go/sqls"
	"bbs-go/web/params"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"
)

var TopicTagService = newTopicTagService()

func newTopicTagService() *topicTagService {
	return &topicTagService{}
}

type topicTagService struct {
}

func (s *topicTagService) Get(id int64) *model.TopicTag {
	return repository.TopicTagRepository.Get(sqls.DB(), id)
}

func (s *topicTagService) Take(where ...interface{}) *model.TopicTag {
	return repository.TopicTagRepository.Take(sqls.DB(), where...)
}

func (s *topicTagService) Find(cnd *sqls.Cnd) []model.TopicTag {
	return repository.TopicTagRepository.Find(sqls.DB(), cnd)
}

func (s *topicTagService) FindOne(cnd *sqls.Cnd) *model.TopicTag {
	return repository.TopicTagRepository.FindOne(sqls.DB(), cnd)
}

func (s *topicTagService) FindPageByParams(params *params.QueryParams) (list []model.TopicTag, paging *sqls.Paging) {
	return repository.TopicTagRepository.FindPageByParams(sqls.DB(), params)
}

func (s *topicTagService) FindPageByCnd(cnd *sqls.Cnd) (list []model.TopicTag, paging *sqls.Paging) {
	return repository.TopicTagRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *topicTagService) Create(t *model.TopicTag) error {
	return repository.TopicTagRepository.Create(sqls.DB(), t)
}

func (s *topicTagService) Update(t *model.TopicTag) error {
	return repository.TopicTagRepository.Update(sqls.DB(), t)
}

func (s *topicTagService) Updates(id int64, columns map[string]interface{}) error {
	return repository.TopicTagRepository.Updates(sqls.DB(), id, columns)
}

func (s *topicTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return repository.TopicTagRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *topicTagService) DeleteByTopicId(topicId int64) {
	sqls.DB().Model(model.TopicTag{}).Where("topic_id = ?", topicId).UpdateColumn("status", constants.StatusDeleted)
}

func (s *topicTagService) UndeleteByTopicId(topicId int64) {
	sqls.DB().Model(model.TopicTag{}).Where("topic_id = ?", topicId).UpdateColumn("status", constants.StatusOK)
}
