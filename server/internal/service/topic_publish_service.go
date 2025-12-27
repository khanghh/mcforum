package service

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/errs"
	"bbs-go/internal/event"
	"bbs-go/internal/locale"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/repository"
	"bbs-go/pkg/iplocator"
	"log/slog"
	"strings"

	"bbs-go/common/dates"
	"bbs-go/common/strs"
	"bbs-go/common/urls"
	"bbs-go/sqls"

	"gorm.io/gorm"
)

type PublishTopicArgs struct {
	UserId      int64
	Title       string
	ForumId     int64
	Content     string
	HideContent string
	Tags        []string
	Images      []string
	UserAgent   string
	IPAddress   string
}

func (s topicService) checkArgs(args PublishTopicArgs) (err error) {
	if strs.IsBlank(args.Title) {
		return errs.NewBadRequestError(locale.T("topic.title_required"))
	}

	if strs.RuneLen(args.Title) > constants.ForumTitleMaxLength {
		return errs.NewBadRequestError(locale.T("topic.title_max_length_exceeded"))
	}

	if strs.IsBlank(args.Content) {
		return errs.NewBadRequestError(locale.T("topic.content_required"))
	}

	if args.ForumId <= 0 {
		args.ForumId = SysConfigService.GetDefaultForumId()
	}

	forum := repository.ForumRepository.Get(sqls.DB(), args.ForumId)
	if forum == nil || forum.Status != constants.StatusActive {
		return errs.NewBadRequestError(locale.T(" forum.not_found"))
	}

	return nil
}

// Publish Publish a topic
func (s *topicService) Publish(args PublishTopicArgs) (*model.Topic, error) {
	if err := s.checkArgs(args); err != nil {
		return nil, err
	}

	now := dates.NowTimestamp()
	topic := &model.Topic{
		Type:            constants.TopicTypeTopic,
		UserID:          args.UserId,
		ForumId:         args.ForumId,
		Title:           args.Title,
		Slug:            urls.GenerateSlug(args.Title),
		Content:         args.Content,
		HideContent:     args.HideContent,
		Status:          constants.StatusActive,
		UserAgent:       args.UserAgent,
		IP:              args.IPAddress,
		IPLocation:      iplocator.IpLocation(args.IPAddress),
		LastCommentTime: now,
		CreateTime:      now,
	}

	if len(args.Images) > 0 {
		topic.ImageList = strings.Join(args.Images, ",")
	}

	if s.isNeedReview(&args) {
		topic.Status = constants.StatusReview
	}

	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := repository.TopicRepository.Create(tx, topic); err != nil {
			return err
		}
		if err := repository.UserRepository.IncreaseTopicCount(tx, args.UserId); err != nil {
			return err
		}
		cache.UserCache.Invalidate(args.UserId)
		return repository.TopicTagRepository.AddTopicTags(tx, topic.ID, args.Tags)
	})
	if err != nil {
		return nil, err
	}

	event.Send(event.TopicCreatedEvent{Topic: topic})
	return topic, nil
}

// IsNeedReview Determine whether review is required
func (s *topicService) isNeedReview(args *PublishTopicArgs) bool {
	if hits := ForbiddenWordService.Check(args.Title); len(hits) > 0 {
		slog.Info(locale.T("topic.prohibited_word_in_title"), slog.String("hits", strings.Join(hits, ",")))
		return true
	}

	if hits := ForbiddenWordService.Check(args.Content); len(hits) > 0 {
		slog.Info(locale.T("topic.prohibited_word_in_content"), slog.String("hits", strings.Join(hits, ",")))
		return true
	}

	return false
}
