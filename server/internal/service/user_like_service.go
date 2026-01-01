package service

import (
	"bbs-go/internal/errs"
	"bbs-go/internal/event"
	"bbs-go/internal/model/constants"
	"errors"

	"bbs-go/common/dates"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"
)

var UserLikeService = newUserLikeService()

func newUserLikeService() *userLikeService {
	return &userLikeService{}
}

type userLikeService struct {
}

func (s *userLikeService) Get(id int64) *model.UserLike {
	return repository.UserLikeRepository.Get(sqls.DB(), id)
}

func (s *userLikeService) Take(where ...interface{}) *model.UserLike {
	return repository.UserLikeRepository.Take(sqls.DB(), where...)
}

func (s *userLikeService) Find(cnd *sqls.Cnd) []model.UserLike {
	return repository.UserLikeRepository.Find(sqls.DB(), cnd)
}

func (s *userLikeService) FindOne(cnd *sqls.Cnd) *model.UserLike {
	return repository.UserLikeRepository.FindOne(sqls.DB(), cnd)
}

func (s *userLikeService) FindPageByParams(params *params.QueryParams) (list []model.UserLike, paging *sqls.Paging) {
	return repository.UserLikeRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userLikeService) FindPageByCnd(cnd *sqls.Cnd) (list []model.UserLike, paging *sqls.Paging) {
	return repository.UserLikeRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userLikeService) Create(t *model.UserLike) error {
	return repository.UserLikeRepository.Create(sqls.DB(), t)
}

func (s *userLikeService) Update(t *model.UserLike) error {
	return repository.UserLikeRepository.Update(sqls.DB(), t)
}

func (s *userLikeService) Updates(id int64, columns map[string]interface{}) error {
	return repository.UserLikeRepository.Updates(sqls.DB(), id, columns)
}

func (s *userLikeService) UpdateColumn(id int64, name string, value interface{}) error {
	return repository.UserLikeRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *userLikeService) Delete(id int64) {
	repository.UserLikeRepository.Delete(sqls.DB(), id)
}

// Count the number of likes
func (s *userLikeService) Count(entityType string, entityId int64) int64 {
	var count int64 = 0
	sqls.DB().Model(&model.UserLike{}).
		Where("entity_id = ?", entityId).
		Where("entity_type = ?", entityType).
		Where("status = ?", constants.StatusActive).
		Count(&count)
	return count
}

// Recent recent likes
func (s *userLikeService) Recent(entityType string, entityId int64, count int) []model.UserLike {
	return s.Find(sqls.NewCnd().
		Eq("entity_id", entityId).
		Eq("entity_type", entityType).
		Eq("status", constants.StatusActive).
		Desc("id").Limit(count),
	)
}

// IsLiked Check whether liked
func (s *userLikeService) IsLiked(userId int64, entityType string, entityId int64) bool {
	isLiked, _, _ := s.isLiked(sqls.DB(), userId, entityType, entityId)
	return isLiked
}

// GetUserLikes filter the given list of enitity ids to get the items liked by user
func (s *userLikeService) GetUserLikes(userId int64, entityType string, entityIds []int64) (likedEntityIds []int64) {
	list := repository.UserLikeRepository.Find(
		sqls.DB(),
		sqls.NewCnd().
			Eq("user_id", userId).
			Eq("entity_type", entityType).
			In("entity_id", entityIds).
			Eq("status", constants.StatusActive),
	)
	for _, like := range list {
		likedEntityIds = append(likedEntityIds, like.EntityID)
	}
	return
}

// TopicLike Like a topic
func (s *userLikeService) TopicLike(userId int64, topicId int64) error {
	topic := repository.TopicRepository.Get(sqls.DB(), topicId)
	if topic == nil || topic.Status != constants.StatusActive {
		return errs.ErrTopicNotFound
	}

	likeCreated := false
	if err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		isLiked, likeExists, err := s.isLiked(tx, userId, constants.EntityTopic, topicId)
		if err != nil || isLiked {
			return nil
		}

		if err = s.like(tx, userId, constants.EntityTopic, topicId); err != nil {
			return err
		}
		likeCreated = !likeExists
		// update like count
		return tx.Model(model.Topic{}).Where("id = ?", topicId).
			UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
	}); err != nil {
		return err
	}

	// send event for new like only
	if likeCreated {
		event.Send(event.UserLikeEvent{
			UserId:     userId,
			EntityId:   topicId,
			EntityType: constants.EntityTopic,
		})
	}

	return nil
}

func (s *userLikeService) TopicUnLike(userId int64, topicId int64) error {
	if err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		isLiked, _, err := s.isLiked(tx, userId, constants.EntityTopic, topicId)
		if err != nil || !isLiked {
			return nil
		}

		if err := s.unlike(tx, userId, constants.EntityTopic, topicId); err != nil {
			return err
		}

		// update like count
		return repository.TopicRepository.UpdateColumn(tx, topicId, "like_count", gorm.Expr("like_count - 1"))
	}); err != nil {
		return err
	}

	// send event
	// event.Send(event.UserUnLikeEvent{
	// 	UserId:     userId,
	// 	EntityId:   topicId,
	// 	EntityType: constants.EntityTopic,
	// })

	return nil
}

// CommentLike topic like
func (s *userLikeService) CommentLike(userId int64, commentId int64) error {
	comment := repository.CommentRepository.Get(sqls.DB(), commentId)
	if comment == nil || comment.Status != constants.StatusActive {
		return errs.ErrCommentNotFound
	}

	likeCreated := false
	if err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		isLiked, likeExists, err := s.isLiked(tx, userId, constants.EntityComment, commentId)
		if err != nil || isLiked {
			return nil
		}

		if err = s.like(tx, userId, constants.EntityComment, commentId); err != nil {
			return err
		}
		likeCreated = !likeExists
		// update like count
		return repository.CommentRepository.UpdateColumn(tx, commentId, "like_count", gorm.Expr("like_count + 1"))
	}); err != nil {
		return err
	}

	// send event for new like only
	if likeCreated {
		event.Send(event.UserLikeEvent{
			UserId:     userId,
			EntityId:   commentId,
			EntityType: constants.EntityComment,
		})
	}

	return nil
}

// CommentLike Like a comment
func (s *userLikeService) CommentUnLike(userId int64, commentId int64) error {
	comment := repository.CommentRepository.Get(sqls.DB(), commentId)
	if comment == nil || comment.Status != constants.StatusActive {
		return errors.New("Comment does not exist")
	}

	if err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		isLiked, _, err := s.isLiked(tx, userId, constants.EntityComment, commentId)
		if err != nil || !isLiked {
			return nil
		}

		if err := s.unlike(tx, userId, constants.EntityComment, commentId); err != nil {
			return err
		}
		// update like count
		return repository.CommentRepository.UpdateColumn(tx, commentId, "like_count", gorm.Expr("like_count - 1"))
	}); err != nil {
		return err
	}

	// send event
	event.Send(event.UserUnLikeEvent{
		UserId:     userId,
		EntityId:   commentId,
		EntityType: constants.EntityComment,
	})

	return nil
}

func (r *userLikeService) isLiked(db *gorm.DB, userId int64, entityType string, entityId int64) (bool, bool, error) {
	var existing model.UserLike
	err := db.Model(&model.UserLike{}).
		Where("user_id = ? AND entity_id = ? AND entity_type = ?", userId, entityId, entityType).
		Take(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, false, nil
		}
		return false, false, err
	}
	return existing.Status == constants.StatusActive, true, nil
}

// like return true if a new like record is created
func (s *userLikeService) like(tx *gorm.DB, userID int64, entityType string, entityID int64) error {
	like := model.UserLike{
		UserID:     userID,
		EntityType: entityType,
		EntityID:   entityID,
		Status:     constants.StatusActive,
		CreateTime: dates.NowTimestamp(),
	}

	// Upsert (Insert or Update)
	ret := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "entity_id"}, {Name: "entity_type"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"status": constants.StatusActive}),
	}).Create(&like)
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

func (s userLikeService) unlike(tx *gorm.DB, userId int64, entityType string, entityId int64) error {
	return tx.Model(&model.UserLike{}).
		Where("user_id = ? AND entity_id = ? AND entity_type = ?", userId, entityId, entityType).
		Update("status", constants.StatusDeleted).Error
}
