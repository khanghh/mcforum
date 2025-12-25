package service

import (
	"bbs-go/internal/errs"
	"bbs-go/internal/event"
	"bbs-go/internal/model/constants"
	"errors"

	"bbs-go/common/dates"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

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

// AddArticleFavorite Favorite an article
func (s *favoriteService) AddArticleFavorite(userId, articleId int64) error {
	article := repository.ArticleRepository.Get(sqls.DB(), articleId)
	if article == nil || article.Status != constants.StatusActive {
		return errors.New("The article to favorite does not exist")
	}
	return s.addFavorite(userId, constants.EntityArticle, articleId)
}

// AddTopicFavorite Favorite a topic
func (s *favoriteService) AddTopicFavorite(userId, topicId int64) error {
	topic := repository.TopicRepository.Get(sqls.DB(), topicId)
	if topic == nil || topic.Status != constants.StatusActive {
		return errs.ErrTopicNotFound
	}
	return s.addFavorite(userId, constants.EntityTopic, topicId)
}

// RemoveTopicFavorite Remove topic favorite
func (s *favoriteService) RemoveTopicFavorite(userId, topicId int64) error {
	return s.removeFavorite(userId, constants.EntityTopic, topicId)
}

func (s *favoriteService) addFavorite(userId int64, entityType string, entityId int64) error {
	if s.IsFavorited(userId, entityType, entityId) { // already favorited
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

	// send event
	event.Send(event.UserFavoriteEvent{
		UserId:     userId,
		EntityId:   entityId,
		EntityType: entityType,
	})
	return nil
}

func (s *favoriteService) removeFavorite(userId int64, entityType string, entityId int64) error {
	tmp := s.GetBy(userId, entityType, entityId)
	if tmp != nil {
		repository.FavoriteRepository.Delete(sqls.DB(), tmp.Id)
		event.Send(event.UserUnfavoriteEvent{
			UserId:     userId,
			EntityId:   entityId,
			EntityType: entityType,
		})
	}
	return nil
}
