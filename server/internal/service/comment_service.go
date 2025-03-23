package service

import (
	"bbs-go/internal/model/constants"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/iplocator"
	"errors"
	"log/slog"
	"strings"

	"bbs-go/common/dates"
	"bbs-go/common/jsons"
	"bbs-go/common/strs"
	"bbs-go/sqls"
	"bbs-go/web/params"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"
)

var CommentService = newCommentService()

func newCommentService() *commentService {
	return &commentService{}
}

type commentService struct {
}

func (s *commentService) Get(id int64) *model.Comment {
	return repository.CommentRepository.Get(sqls.DB(), id)
}

func (s *commentService) Take(where ...interface{}) *model.Comment {
	return repository.CommentRepository.Take(sqls.DB(), where...)
}

func (s *commentService) Find(cnd *sqls.Cnd) []model.Comment {
	return repository.CommentRepository.Find(sqls.DB(), cnd)
}

func (s *commentService) FindOne(cnd *sqls.Cnd) *model.Comment {
	return repository.CommentRepository.FindOne(sqls.DB(), cnd)
}

func (s *commentService) FindPageByParams(params *params.QueryParams) (list []model.Comment, paging *sqls.Paging) {
	return repository.CommentRepository.FindPageByParams(sqls.DB(), params)
}

func (s *commentService) FindPageByCnd(cnd *sqls.Cnd) (list []model.Comment, paging *sqls.Paging) {
	return repository.CommentRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *commentService) Count(cnd *sqls.Cnd) int64 {
	return repository.CommentRepository.Count(sqls.DB(), cnd)
}

func (s *commentService) Create(t *model.Comment) error {
	return repository.CommentRepository.Create(sqls.DB(), t)
}

func (s *commentService) Update(t *model.Comment) error {
	return repository.CommentRepository.Update(sqls.DB(), t)
}

func (s *commentService) Updates(id int64, columns map[string]interface{}) error {
	return repository.CommentRepository.Updates(sqls.DB(), id, columns)
}

func (s *commentService) UpdateColumn(id int64, name string, value interface{}) error {
	return repository.CommentRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *commentService) Delete(id int64) error {
	return repository.CommentRepository.UpdateColumn(sqls.DB(), id, "status", constants.StatusDeleted)
}

// Publish 发表评论
func (s *commentService) Publish(userId int64, form model.CreateCommentForm) (*model.Comment, error) {
	form.Content = strings.TrimSpace(form.Content)
	if strs.IsBlank(form.EntityType) {
		return nil, errors.New("参数非法")
	}
	if form.EntityId <= 0 {
		return nil, errors.New("参数非法")
	}
	if strs.IsBlank(form.Content) {
		return nil, errors.New("请输入评论内容")
	}

	comment := &model.Comment{
		UserId:      userId,
		EntityType:  form.EntityType,
		EntityId:    form.EntityId,
		Content:     form.Content,
		ContentType: constants.ContentTypeText,
		QuoteId:     form.QuoteId,
		Status:      constants.StatusOK,
		UserAgent:   form.UserAgent,
		Ip:          form.Ip,
		IpLocation:  iplocator.IpLocation(form.Ip),
		CreateTime:  dates.NowTimestamp(),
	}

	if len(form.ImageList) > 0 {
		imageListStr, err := jsons.ToStr(form.ImageList)
		if err == nil {
			comment.ImageList = imageListStr
		} else {
			slog.Error(err.Error(), slog.Any("err", err))
		}
	}

	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := repository.CommentRepository.Create(tx, comment); err != nil {
			return err
		}

		if form.EntityType == constants.EntityTopic {
			if err := TopicService.onComment(tx, form.EntityId, comment); err != nil {
				return err
			}
		} else if form.EntityType == constants.EntityComment { // 二级评论
			if err := s.onComment(tx, comment); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 用户跟帖计数
	UserService.IncrCommentCount(userId)
	// 获得积分
	UserService.IncrScoreForPostComment(comment)
	// 发送事件
	event.Send(event.CommentCreateEvent{
		UserId:    userId,
		CommentId: comment.Id,
	})

	return comment, nil
}

// onComment 评论被回复（二级评论）
func (s *commentService) onComment(tx *gorm.DB, comment *model.Comment) error {
	return repository.CommentRepository.UpdateColumn(tx, comment.EntityId, "comment_count", gorm.Expr("comment_count + 1"))
}

// // 统计数量
// func (s *commentService) Count(entityType string, entityId int64) int64 {
// 	var count int64 = 0
// 	sqls.DB().Model(&model.Comment{}).Where("entity_type = ? and entity_id = ?", entityType, entityId).Count(&count)
// 	return count
// }

// GetComments 列表
func (s *commentService) GetComments(entityType string, entityId int64, cursor int64) (comments []model.Comment, nextCursor int64, hasMore bool) {
	limit := 20
	cnd := sqls.NewCnd().Eq("entity_type", entityType).Eq("entity_id", entityId).Eq("status", constants.StatusOK).Desc("id").Limit(limit)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	comments = repository.CommentRepository.Find(sqls.DB(), cnd)
	if len(comments) > 0 {
		nextCursor = comments[len(comments)-1].Id
		hasMore = len(comments) >= limit
	} else {
		nextCursor = cursor
	}
	return
}

// GetReplies 二级回复列表
func (s *commentService) GetReplies(commentId int64, cursor int64, limit int) (comments []model.Comment, nextCursor int64, hasMore bool) {
	cnd := sqls.NewCnd().Eq("entity_type", constants.EntityComment).Eq("entity_id", commentId).Eq("status", constants.StatusOK).Asc("id").Limit(limit)
	if cursor > 0 {
		cnd.Gt("id", cursor)
	}
	comments = s.Find(cnd)
	if len(comments) > 0 {
		nextCursor = comments[len(comments)-1].Id
		hasMore = len(comments) >= limit
	} else {
		nextCursor = cursor
	}
	return
}

// ScanByUser 按照用户扫描数据
func (s *commentService) ScanByUser(userId int64, callback func(comments []model.Comment)) {
	var cursor int64 = 0
	for {
		list := repository.CommentRepository.Find(sqls.DB(), sqls.NewCnd().
			Eq("user_id", userId).Gt("id", cursor).Asc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

// ScanByUser 按照用户扫描数据
func (s *commentService) Scan(callback func(comments []model.Comment)) {
	var cursor int64 = 0
	for {
		logrus.Info("scan comments, cursor:" + cast.ToString(cursor))
		list := repository.CommentRepository.Find(sqls.DB(), sqls.NewCnd().
			Gt("id", cursor).Asc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

func (s *commentService) IsCommented(userId int64, entityType string, entityId int64) bool {
	return s.FindOne(sqls.NewCnd().Where("user_id = ? and entity_id = ? and entity_type = ? and status = ?", userId, entityId, entityType, constants.StatusOK)) != nil
}
