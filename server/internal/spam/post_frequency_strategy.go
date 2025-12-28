package spam

import (
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/model"
	"bbs-go/internal/repository"
	"errors"
	"time"

	"bbs-go/common/dates"
	"bbs-go/sqls"
)

// PostFrequencyStrategy Post frequency limit
type PostFrequencyStrategy struct{}

func (PostFrequencyStrategy) Name() string {
	return "PostFrequencyStrategy"
}

func (PostFrequencyStrategy) CheckTopic(user *model.User, topic payload.CreateTopicForm) error {
	// Registered more than 24 hours
	if user.CreateTime < dates.Timestamp(time.Now().Add(-time.Hour*24)) {
		return nil
	}
	var (
		maxCountInTenMinutes int64 = 1 // max posts within ten minutes
		maxCountInOneHour    int64 = 2 // max posts within one hour
		maxCountInOneDay     int64 = 3 // max posts within one day
	)

	if repository.TopicRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.ID).
		Gt("create_time", dates.Timestamp(time.Now().Add(-time.Hour*24)))) >= maxCountInOneDay {
		return errors.New("Posting too frequently, please take a break")
	}

	if repository.TopicRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.ID).
		Gt("create_time", dates.Timestamp(time.Now().Add(-time.Hour)))) >= maxCountInOneHour {
		return errors.New("Posting too frequently, please take a break")
	}

	if repository.TopicRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.ID).
		Gt("create_time", dates.Timestamp(time.Now().Add(-time.Minute*10)))) >= maxCountInTenMinutes {
		return errors.New("Posting too frequently, please take a break")
	}

	return nil
}

func (s PostFrequencyStrategy) CheckComment(user *model.User, form payload.CreateCommentForm) error {
	// Registered more than 24 hours
	if user.CreateTime < dates.Timestamp(time.Now().Add(-time.Hour*24)) {
		return nil
	}

	var (
		maxCountInTenMinutes int64 = 1 // max posts within ten minutes
		maxCountInOneHour    int64 = 1 // max posts within one hour
		maxCountInOneDay     int64 = 1 // max posts within one day
	)

	if repository.CommentRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.ID).
		Gt("create_time", dates.Timestamp(time.Now().Add(-time.Hour*24)))) >= maxCountInOneDay {
		return errors.New("Posting too frequently, please take a break")
	}

	if repository.CommentRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.ID).
		Gt("create_time", dates.Timestamp(time.Now().Add(-time.Hour)))) >= maxCountInOneHour {
		return errors.New("Posting too frequently, please take a break")
	}

	if repository.CommentRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.ID).
		Gt("create_time", dates.Timestamp(time.Now().Add(-time.Minute*10)))) >= maxCountInTenMinutes {
		return errors.New("Posting too frequently, please take a break")
	}
	return nil
}
