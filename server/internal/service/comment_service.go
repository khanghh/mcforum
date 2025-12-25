package service

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/errs"
	"bbs-go/internal/event"
	"bbs-go/internal/locale"
	"bbs-go/internal/model/constants"
	"bbs-go/pkg/iplocator"
	"strings"

	"bbs-go/common/dates"
	"bbs-go/common/strs"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"
)

var CommentService = newCommentService()

type CreateCommentArgs struct {
	UserId    int64
	TopicId   int64
	ParentId  int64
	QuoteId   int64
	Content   string
	Images    []string
	UserAgent string
	IPAddress string
}

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

// Publish Post a comment
func (s *commentService) CreateComment(args CreateCommentArgs) (*model.Comment, error) {
	args.Content = strings.TrimSpace(args.Content)
	if strs.IsBlank(args.Content) {
		return nil, errs.NewBadRequestError(locale.T("system.comment_content_required"))
	}

	comment := &model.Comment{
		UserId:      args.UserId,
		TopicId:     args.TopicId,
		ParentId:    args.ParentId,
		QuoteId:     args.QuoteId,
		Content:     args.Content,
		ContentType: constants.ContentTypeText,
		Status:      constants.StatusActive,
		UserAgent:   args.UserAgent,
		Ip:          args.IPAddress,
		IpLocation:  iplocator.IpLocation(args.IPAddress),
		CreateTime:  dates.NowTimestamp(),
	}

	if len(args.Images) > 0 {
		comment.ImageList = strings.Join(args.Images, ",")
	}

	topic := TopicService.Get(args.TopicId)
	if topic == nil || topic.Status != constants.StatusActive {
		return nil, errs.ErrTopicNotFound
	}

	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := repository.CommentRepository.Create(tx, comment); err != nil {
			return err
		}

		if err := repository.UserRepository.IncreaseCommentCount(tx, comment.UserId); err != nil {
			return err
		}
		cache.UserCache.Invalidate(comment.UserId)

		if comment.ParentId == 0 {
			return TopicService.onComment(tx, comment)
		}
		return s.onComment(tx, comment)
	})

	if err != nil {
		return nil, err
	}

	event.Send(event.CommentCreatedEvent{
		Topic:   topic,
		Comment: comment,
	})

	return comment, nil
}

// onComment Reply to a comment (second-level comment)
func (s *commentService) onComment(tx *gorm.DB, comment *model.Comment) error {
	return repository.CommentRepository.UpdateColumn(tx, comment.ParentId, "comment_count", gorm.Expr("comment_count + 1"))
}

// GetComments List of comments
func (s *commentService) GetComments(topicId int64, cursor int64, limit int) (comments []model.Comment, nextCursor int64, hasMore bool) {
	cnd := sqls.NewCnd().Eq("topic_id", topicId).Eq("parent_id", 0).Eq("status", constants.StatusActive).Desc("id").Limit(limit)
	if cursor > 0 {
		cnd.Lt("create_time", cursor)
	}
	comments = repository.CommentRepository.Find(sqls.DB(), cnd)
	if len(comments) > 0 {
		nextCursor = comments[len(comments)-1].CreateTime
		hasMore = len(comments) >= limit
	} else {
		nextCursor = cursor
	}
	return
}

// GetReplies List of second-level replies
func (s *commentService) GetReplies(commentId int64, cursor int64, limit int) (comments []model.Comment, nextCursor int64, hasMore bool) {
	cnd := sqls.NewCnd().Eq("parent_id", commentId).Eq("status", constants.StatusActive).Desc("id").Limit(limit)
	if cursor > 0 {
		cnd.Lt("create_time", cursor)
	}
	comments = repository.CommentRepository.Find(sqls.DB(), cnd)
	if len(comments) > 0 {
		nextCursor = comments[len(comments)-1].CreateTime
		hasMore = len(comments) >= limit
	} else {
		nextCursor = cursor
	}
	return
}

// ScanByUser Scan data by user
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

// Scan Scan data
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
	return s.FindOne(sqls.NewCnd().Where("user_id = ? and entity_id = ? and entity_type = ? and status = ?", userId, entityId, entityType, constants.StatusActive)) != nil
}
