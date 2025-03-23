package service

import (
	"bbs-go/internal/model/constants"
	"bbs-go/internal/pkg/event"
	"errors"

	"bbs-go/common/dates"
	"bbs-go/sqls"
	"bbs-go/web/params"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"
)

var FavoriteService = newFavoriteService()

func newFavoriteService() *favoriteService {
	return &favoriteService{}
}

type favoriteService struct {
}

func (s *favoriteService) Get(id int64) *model.Favorite {
	return repository.FavoriteRepository.Get(sqls.DB(), id)
}

func (s *favoriteService) Take(where ...interface{}) *model.Favorite {
	return repository.FavoriteRepository.Take(sqls.DB(), where...)
}

func (s *favoriteService) Find(cnd *sqls.Cnd) []model.Favorite {
	return repository.FavoriteRepository.Find(sqls.DB(), cnd)
}

func (s *favoriteService) FindOne(cnd *sqls.Cnd) *model.Favorite {
	return repository.FavoriteRepository.FindOne(sqls.DB(), cnd)
}

func (s *favoriteService) FindPageByParams(params *params.QueryParams) (list []model.Favorite, paging *sqls.Paging) {
	return repository.FavoriteRepository.FindPageByParams(sqls.DB(), params)
}

func (s *favoriteService) FindPageByCnd(cnd *sqls.Cnd) (list []model.Favorite, paging *sqls.Paging) {
	return repository.FavoriteRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *favoriteService) Create(t *model.Favorite) error {
	return repository.FavoriteRepository.Create(sqls.DB(), t)
}

func (s *favoriteService) Update(t *model.Favorite) error {
	return repository.FavoriteRepository.Update(sqls.DB(), t)
}

func (s *favoriteService) Updates(id int64, columns map[string]interface{}) error {
	return repository.FavoriteRepository.Updates(sqls.DB(), id, columns)
}

func (s *favoriteService) UpdateColumn(id int64, name string, value interface{}) error {
	return repository.FavoriteRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *favoriteService) Delete(id int64) {
	repository.FavoriteRepository.Delete(sqls.DB(), id)
}

func (s *favoriteService) IsFavorited(userId int64, entityType string, entityId int64) bool {
	return repository.FavoriteRepository.Take(sqls.DB(), "user_id = ? and entity_type = ? and entity_id = ?",
		userId, entityType, entityId) != nil
}

func (s *favoriteService) GetBy(userId int64, entityType string, entityId int64) *model.Favorite {
	return repository.FavoriteRepository.Take(sqls.DB(), "user_id = ? and entity_type = ? and entity_id = ?",
		userId, entityType, entityId)
}

// AddArticleFavorite 收藏文章
func (s *favoriteService) AddArticleFavorite(userId, articleId int64) error {
	article := repository.ArticleRepository.Get(sqls.DB(), articleId)
	if article == nil || article.Status != constants.StatusOK {
		return errors.New("收藏的文章不存在")
	}
	return s.addFavorite(userId, constants.EntityArticle, articleId)
}

// AddTopicFavorite 收藏主题
func (s *favoriteService) AddTopicFavorite(userId, topicId int64) error {
	topic := repository.TopicRepository.Get(sqls.DB(), topicId)
	if topic == nil || topic.Status != constants.StatusOK {
		return errors.New("收藏的话题不存在")
	}
	return s.addFavorite(userId, constants.EntityTopic, topicId)
}

func (s *favoriteService) addFavorite(userId int64, entityType string, entityId int64) error {
	if s.IsFavorited(userId, entityType, entityId) { // 已经收藏
		return nil
	}
	if err := repository.FavoriteRepository.Create(sqls.DB(), &model.Favorite{
		UserId:     userId,
		EntityType: entityType,
		EntityId:   entityId,
		CreateTime: dates.NowTimestamp(),
	}); err != nil {
		return err
	}

	// 发送事件
	event.Send(event.UserFavoriteEvent{
		UserId:     userId,
		EntityId:   entityId,
		EntityType: entityType,
	})
	return nil
}
