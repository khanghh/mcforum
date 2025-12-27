package spam

import (
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/model"
	"log/slog"
)

var strategies []Strategy

func init() {
	strategies = append(strategies, &EmailVerifyStrategy{})
	strategies = append(strategies, &CaptchaStrategy{})
	// strategies = append(strategies, &PostFrequencyStrategy{})
}

func CheckTopic(user *model.User, form payload.CreateTopicForm) error {
	if len(strategies) == 0 {
		return nil
	}
	for _, strategy := range strategies {
		if err := strategy.CheckTopic(user, form); err != nil {
			slog.Warn("[Topic] strategy hit", slog.Any("strategy", strategy.Name()), slog.Any("userId", user.ID))
			return err
		}
	}
	return nil
}

func CheckComment(user *model.User, form payload.CreateCommentForm) error {
	if len(strategies) == 0 {
		return nil
	}
	for _, strategy := range strategies {
		if err := strategy.CheckComment(user, form); err != nil {
			slog.Warn("[Comment] strategy hit", slog.Any("strategy", strategy.Name()), slog.Any("userId", user.ID))
			return err
		}
	}
	return nil
}
