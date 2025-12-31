package service

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/notification"
	"bbs-go/internal/repository"
	"bbs-go/pkg/bbsurls"
	"bbs-go/pkg/msg"
	"log/slog"

	"bbs-go/common/dates"
	"bbs-go/common/jsons"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"
)

var MessageService = newMessageService()

func newMessageService() *messageService {
	return &messageService{}
}

type messageService struct {
}

type SendMessageArgs struct {
	FromId       int64
	ToId         int64
	Title        string
	Content      string
	QuoteContent string
	Type         msg.Type
	DetailUrl    string
	ExtraData    interface{}
}

func (s *messageService) Get(id int64) *model.Message {
	return repository.MessageRepository.Get(sqls.DB(), id)
}

func (s *messageService) Take(where ...interface{}) *model.Message {
	return repository.MessageRepository.Take(sqls.DB(), where...)
}

func (s *messageService) Find(cnd *sqls.Cnd) []model.Message {
	return repository.MessageRepository.Find(sqls.DB(), cnd)
}

func (s *messageService) FindOne(cnd *sqls.Cnd) *model.Message {
	return repository.MessageRepository.FindOne(sqls.DB(), cnd)
}

func (s *messageService) FindPageByParams(params *params.QueryParams) (list []model.Message, paging *sqls.Paging) {
	return repository.MessageRepository.FindPageByParams(sqls.DB(), params)
}

func (s *messageService) FindPageByCnd(cnd *sqls.Cnd) (list []model.Message, paging *sqls.Paging) {
	return repository.MessageRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *messageService) Create(t *model.Message) error {
	return repository.MessageRepository.Create(sqls.DB(), t)
}

func (s *messageService) Update(t *model.Message) error {
	return repository.MessageRepository.Update(sqls.DB(), t)
}

func (s *messageService) Updates(id int64, columns map[string]interface{}) error {
	return repository.MessageRepository.Updates(sqls.DB(), id, columns)
}

func (s *messageService) UpdateColumn(id int64, name string, value interface{}) error {
	return repository.MessageRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *messageService) Delete(id int64) {
	repository.MessageRepository.Delete(sqls.DB(), id)
}

// GetUnReadCount Get unread message count
func (s *messageService) GetUnReadCount(userId int64) (count int64) {
	sqls.DB().Where("user_id = ? and status = ?", userId, msg.StatusUnread).Model(&model.Message{}).Count(&count)
	return
}

// MarkRead Mark all messages as read
func (s *messageService) MarkRead(userId int64) error {
	return repository.MessageRepository.UpdateByUserID(sqls.DB(), userId, "status", msg.StatusHaveRead)
}

// SendMsg Send a message
func (s *messageService) SendMsg(args SendMessageArgs) {
	t := &model.Message{
		FromID:       args.FromId,
		UserID:       args.ToId,
		Title:        args.Title,
		Content:      args.Content,
		QuoteContent: args.QuoteContent,
		Type:         int(args.Type),
		DetailUrl:    args.DetailUrl,
		ExtraData:    jsons.ToJsonStr(args.ExtraData),
		Status:       msg.StatusUnread,
		CreateTime:   dates.NowTimestamp(),
	}
	if err := s.Create(t); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	} else {
		s.SendEmailNotice(t)
	}
}

// SendEmailNotice Send email notification
func (s *messageService) SendEmailNotice(t *model.Message) {
	msgType := msg.Type(t.Type)

	// Do not send email when a topic is deleted
	if msgType == msg.TypeTopicDelete {
		return
	}
	user := cache.UserCache.Get(t.UserID)
	if user == nil || len(user.Email.String) == 0 {
		return
	}
	var (
		siteTitle  = cache.SysConfigCache.GetValue(constants.SysConfigSiteTitle)
		emailTitle = siteTitle + " - New message"
	)

	if msgType == msg.TypeTopicComment {
		emailTitle = siteTitle + " - New topic comment"
	} else if msgType == msg.TypeCommentReply {
		emailTitle = siteTitle + " - New reply"
	} else if msgType == msg.TypeTopicLike {
		emailTitle = siteTitle + " - New like"
	} else if msgType == msg.TypeTopicFavorite {
		emailTitle = siteTitle + " - Topic favorited"
	} else if msgType == msg.TypeTopicRecommend {
		emailTitle = siteTitle + " - Topic recommended"
	} else if msgType == msg.TypeTopicDelete {
		emailTitle = siteTitle + " - Topic deleted"
	}

	var from *model.User
	if t.FromID > 0 {
		from = cache.UserCache.Get(t.FromID)
	}
	err := notification.SendTemplateEmail(from, user.Email.String, emailTitle, emailTitle, t.Content,
		t.QuoteContent, &model.ActionLink{
			Title: "Click to view details",
			Url:   bbsurls.AbsUrl("/user/messages"),
		})
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
}
