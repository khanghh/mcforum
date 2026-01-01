package service

import (
	"bbs-go/internal/errs"
	"bbs-go/internal/locale"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/notification"
	"bbs-go/internal/validate"
	"bbs-go/pkg/bbsurls"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bbs-go/common/dates"
	"bbs-go/common/passwd"
	"bbs-go/common/strs"
	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

	"gorm.io/gorm"

	"bbs-go/internal/cache"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"
)

// Email verification token expiration (hours)
const emailVerifyExpireHour = 24

var UserService = newUserService()

func newUserService() *userService {
	return &userService{}
}

type userService struct {
}

func (s *userService) Get(id int64) *model.User {
	return repository.UserRepository.Get(sqls.DB(), id)
}

func (s *userService) Take(where ...interface{}) *model.User {
	return repository.UserRepository.Take(sqls.DB(), where...)
}

func (s *userService) Find(cnd *sqls.Cnd) []model.User {
	return repository.UserRepository.Find(sqls.DB(), cnd)
}

func (s *userService) FindOne(cnd *sqls.Cnd) *model.User {
	return repository.UserRepository.FindOne(sqls.DB(), cnd)
}

func (s *userService) FindPageByParams(params *params.QueryParams) (list []model.User, paging *sqls.Paging) {
	return repository.UserRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userService) FindPageByCnd(cnd *sqls.Cnd) (list []model.User, paging *sqls.Paging) {
	return repository.UserRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userService) Count(cnd *sqls.Cnd) int64 {
	return repository.UserRepository.Count(sqls.DB(), cnd)
}

func (s *userService) Create(t *model.User) error {
	err := repository.UserRepository.Create(sqls.DB(), t)
	if err == nil {
		cache.UserCache.Invalidate(t.ID)
	}
	return nil
}

func (s *userService) Update(t *model.User) error {
	err := repository.UserRepository.Update(sqls.DB(), t)
	cache.UserCache.Invalidate(t.ID)
	return err
}

func (s *userService) Updates(id int64, columns map[string]interface{}) error {
	err := repository.UserRepository.Updates(sqls.DB(), id, columns)
	cache.UserCache.Invalidate(id)
	return err
}

func (s *userService) UpdateColumn(id int64, name string, value interface{}) error {
	err := repository.UserRepository.UpdateColumn(sqls.DB(), id, name, value)
	cache.UserCache.Invalidate(id)
	return err
}

func (s *userService) Delete(id int64) {
	repository.UserRepository.Delete(sqls.DB(), id)
	cache.UserCache.Invalidate(id)
}

// Scan
func (s *userService) Scan(callback func(users []model.User)) {
	var cursor int64
	for {
		list := repository.UserRepository.Find(sqls.DB(), sqls.NewCnd().Where("id > ?", cursor).Asc("id").Limit(100))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		callback(list)
	}
}

// Forbidden (mute)
func (s *userService) Forbidden(operatorId, userId int64, days int, reason string, r *http.Request) error {
	var forbiddenEndTime int64
	if days == -1 { // permanent mute
		forbiddenEndTime = -1
	} else if days > 0 {
		forbiddenEndTime = dates.Timestamp(time.Now().Add(time.Hour * 24 * time.Duration(days)))
	} else {
		return errors.New(locale.T("errors.mute_duration_invalid"))
	}
	if repository.UserRepository.UpdateColumn(sqls.DB(), userId, "forbidden_end_time", forbiddenEndTime) == nil {
		cache.UserCache.Invalidate(userId)
		description := ""
		if strs.IsNotBlank(reason) {
			description = locale.T("ban.ban_reason", reason)
		}
		OperateLogService.AddOperateLog(operatorId, constants.OpTypeForbidden, constants.EntityUser, userId,
			description, r)

		// Permanent ban
		if days == -1 {
			user := cache.UserCache.Get(userId)
			_ = s.DecrScore(userId, user.Score, constants.EntityUser, strconv.FormatInt(operatorId, 10), "Permanent ban")
			go func() {
				// delete topics
				TopicService.ScanByUser(userId, func(topics []model.Topic) {
					for _, topic := range topics {
						if topic.Status != constants.StatusDeleted {
							_ = TopicService.Delete(topic.ID, operatorId, nil)
						}
					}
				})

				// delete comments
				CommentService.ScanByUser(userId, func(comments []model.Comment) {
					for _, comment := range comments {
						if comment.Status != constants.StatusDeleted {
							_ = CommentService.Delete(comment.ID)
						}
					}
				})

			}()
		}
	}
	return nil
}

// RemoveForbidden (remove mute)
func (s *userService) RemoveForbidden(operatorId, userId int64, r *http.Request) {
	user := s.Get(userId)
	if user == nil || !user.IsForbidden() {
		return
	}
	if repository.UserRepository.UpdateColumn(sqls.DB(), userId, "forbidden_end_time", 0) == nil {
		cache.UserCache.Invalidate(user.ID)
		OperateLogService.AddOperateLog(operatorId, constants.OpTypeRemoveForbidden, constants.EntityUser, userId, "", r)
	}
}

// GetByEmail find by email
func (s *userService) GetByEmail(email string) *model.User {
	return repository.UserRepository.GetByEmail(sqls.DB(), email)
}

// GetByUsername find by username
func (s *userService) GetByUsername(username string) *model.User {
	return repository.UserRepository.GetByUsername(sqls.DB(), username)
}

// SignUp register
func (s *userService) SignUp(username, email, nickname, password, rePassword string) (*model.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	nickname = strings.TrimSpace(nickname)

	// validate nickname
	if len(nickname) == 0 {
		return nil, errors.New(locale.T("errors.sigup_nickname_required"))
	}

	// validate password
	err := validate.IsValidPassword(password, rePassword)
	if err != nil {
		return nil, err
	}

	// validate email
	if email != "" {
		if err := validate.IsEmail(email); err != nil {
			return nil, err
		}
		if s.GetByEmail(email) != nil {
			return nil, errors.New(locale.T("errors.email_taken", email))
		}
	} else {
		return nil, errors.New(locale.T("errors.sigup_email_required"))
	}

	// validate username
	if username != "" {
		if err := validate.IsUsername(username); err != nil {
			return nil, err
		}
		if s.isUsernameExists(username) {
			return nil, errors.New(locale.T("errors.username_taken", username))
		}
	}

	user := &model.User{
		Username:      sqls.SqlNullString(username),
		Email:         sqls.SqlNullString(email),
		EmailVerified: true,
		Nickname:      nickname,
		Password:      passwd.EncodePassword(password),
		IsActive:      true,
		CreateTime:    dates.NowTimestamp(),
		UpdateTime:    dates.NowTimestamp(),
	}

	err = repository.UserRepository.Create(sqls.DB(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// SignIn login
func (s *userService) SignIn(username, password string) (*model.User, error) {
	if strs.IsBlank(username) {
		return nil, errors.New(locale.T("errors.sigin_email_username_required"))
	}
	if strs.IsBlank(password) {
		return nil, errors.New(locale.T("errors.sigin_password_required"))
	}
	if err := validate.IsPassword(password); err != nil {
		return nil, err
	}
	var user *model.User = nil
	if err := validate.IsEmail(username); err == nil { // if the user entered an email
		user = s.GetByEmail(username)
	} else {
		user = s.GetByUsername(username)
	}
	if user == nil || !user.IsActive {
		return nil, errors.New(locale.T("errors.sigin_failure"))
	}
	if !passwd.ValidatePassword(user.Password, password) {
		return nil, errors.New(locale.T("errors.sigin_failure"))
	}
	s.SyncUserCount()
	return user, nil
}

// isEmailExists whether email exists
func (s *userService) isEmailExists(email string) bool {
	if len(email) == 0 { // if email is empty, consider it non-existent
		return false
	}
	return s.GetByEmail(email) != nil
}

// isUsernameExists whether username exists
func (s *userService) isUsernameExists(username string) bool {
	return s.GetByUsername(username) != nil
}

// UpdateAvatar update avatar
func (s *userService) UpdateAvatar(userId int64, avatar string) error {
	return s.UpdateColumn(userId, "avatar", avatar)
}

// UpdateNickname update nickname
func (s *userService) UpdateNickname(userId int64, nickname string) error {
	return s.UpdateColumn(userId, "nickname", nickname)
}

// UpdateDescription update description
func (s *userService) UpdateStatusMessage(userId int64, msg string) error {
	return s.UpdateColumn(userId, "status_message", msg)
}

// UpdateGender update gender
func (s *userService) UpdateGender(userId int64, gender string) error {
	if strs.IsBlank(gender) {
		return s.UpdateColumn(userId, "gender", "")
	} else {
		if gender != string(constants.GenderMale) && gender != string(constants.GenderFemale) {
			return errors.New("invalidate gender value")
		}
		return s.UpdateColumn(userId, "gender", gender)
	}
}

// UpdateBirthday update birthday
func (s *userService) UpdateBirthday(userId int64, birthdayStr string) error {
	if strs.IsBlank(birthdayStr) {
		return s.UpdateColumn(userId, "birthday", "")
	} else {
		birthday, err := dates.Parse(birthdayStr, dates.FmtDate)
		if err != nil {
			return err
		}
		return s.UpdateColumn(userId, "birthday", birthday)
	}
}

// UpdateCoverPhoto update background image
func (s *userService) UpdateCoverPhoto(userId int64, coverImageURL string) error {
	if err := s.UpdateColumn(userId, "background_image", coverImageURL); err != nil {
		return err
	}
	return nil
}

// SetUsername set username
func (s *userService) SetUsername(userId int64, username string) error {
	username = strings.TrimSpace(username)
	if err := validate.IsUsername(username); err != nil {
		return err
	}

	user := s.Get(userId)
	if len(user.Username.String) > 0 {
		return errors.New("You have already set a username and cannot set it again.")
	}
	if s.isUsernameExists(username) {
		return errors.New("Username " + username + " is already taken")
	}
	return s.UpdateColumn(userId, "username", username)
}

// SetEmail set email
func (s *userService) SetEmail(userId int64, email string) error {
	email = strings.TrimSpace(email)
	if err := validate.IsEmail(email); err != nil {
		return err
	}
	user := s.Get(userId)
	if user == nil {
		return errors.New(locale.T("user.not_found"))
	}
	if user.Email.String == email {
		// user's email not changed
		return nil
	}
	if s.isEmailExists(email) {
		return errors.New(locale.T("errors.email_taken", email))
	}
	return s.Updates(userId, map[string]interface{}{
		"email":          email,
		"email_verified": false,
	})
}

// SetPassword set password
func (s *userService) SetPassword(userId int64, password, rePassword string) error {
	if err := validate.IsValidPassword(password, rePassword); err != nil {
		return err
	}
	user := s.Get(userId)
	if len(user.Password) > 0 {
		return errors.New(locale.T("errors.invalid_request"))
	}
	password = passwd.EncodePassword(password)
	return s.UpdateColumn(userId, "password", password)
}

// UpdatePassword change password
func (s *userService) UpdatePassword(userId int64, oldPassword, password, rePassword string) error {
	if err := validate.IsValidPassword(password, rePassword); err != nil {
		return err
	}
	user := s.Get(userId)

	if len(user.Password) == 0 {
		return errors.New("You haven't set a password, please set one first")
	}

	if !passwd.ValidatePassword(user.Password, oldPassword) {
		return errors.New("Old password verification failed")
	}

	return s.UpdateColumn(userId, "password", passwd.EncodePassword(password))
}

// IncreaseTopicCount topic_count + 1
func (s *userService) IncreaseTopicCount(userId int64) error {
	if err := repository.UserRepository.UpdateColumn(sqls.DB(), userId, "topic_count", gorm.Expr("topic_count + 1")); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return err
	}
	cache.UserCache.Invalidate(userId)
	return nil
}

// IncreaseCommentCount comment_count + 1
func (s *userService) IncreaseCommentCount(userId int64) error {
	if err := repository.UserRepository.UpdateColumn(sqls.DB(), userId, "comment_count", gorm.Expr("comment_count + 1")); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return err
	}
	cache.UserCache.Invalidate(userId)
	return nil
}

// SyncUserCount sync user counts
func (s *userService) SyncUserCount() {
	s.Scan(func(users []model.User) {
		for _, user := range users {
			topicCount := repository.TopicRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.ID).Eq("status", constants.StatusActive))
			commentCount := repository.CommentRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.ID).Eq("status", constants.StatusActive))
			_ = repository.UserRepository.UpdateColumn(sqls.DB(), user.ID, "topic_count", topicCount)
			_ = repository.UserRepository.UpdateColumn(sqls.DB(), user.ID, "comment_count", commentCount)
			cache.UserCache.Invalidate(user.ID)
		}
	})
}

// SendEmailVerifyEmail send email verification
func (s *userService) SendEmailVerifyEmail(userId int64) error {
	user := s.Get(userId)
	if user == nil {
		return errors.New(locale.T("user.not_found"))
	}
	if user.EmailVerified {
		return errors.New(locale.T("errors.invalid_request"))
	}
	if err := validate.IsEmail(user.Email.String); err != nil {
		return err
	}
	// If an email whitelist is set
	if emailWhitelist := SysConfigService.GetEmailWhitelist(); len(emailWhitelist) > 0 {
		isInWhitelist := false
		for _, whitelist := range emailWhitelist {
			if strings.Contains(strings.ToLower(user.Email.String), strings.ToLower(whitelist)) {
				isInWhitelist = true
				break
			}
		}
		if !isInWhitelist {
			// return silently (do not send verification email)
			slog.Error("Verification with this email is not supported", slog.String("email", user.Email.String))
			return errors.New(locale.T("errors.email_not_accepted"))
		}
	}
	var (
		token     = strs.UUID()
		url       = bbsurls.AbsUrl("/user/email/verify?token=" + token)
		link      = &model.ActionLink{Title: "Click here to verify your email >>", Url: url}
		siteTitle = cache.SysConfigCache.GetValue(constants.SysConfigSiteTitle)
		subject   = "Email verification - " + siteTitle
		title     = "Email verification - " + siteTitle
		content   = "This email is used to verify the email address you set on " + siteTitle + ". Please complete verification within " + strconv.Itoa(emailVerifyExpireHour) + " hours. Verification link: " + url
	)
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := repository.EmailCodeRepository.Create(tx, &model.EmailCode{
			Model:      model.Model{},
			UserID:     userId,
			Email:      user.Email.String,
			Code:       "",
			Token:      token,
			Title:      title,
			Content:    content,
			Used:       false,
			CreateTime: dates.NowTimestamp(),
		}); err != nil {
			return nil
		}
		if err := notification.SendTemplateEmail(nil, user.Email.String, subject, title, content, "", link); err != nil {
			return err
		}
		return nil
	})
}

// VerifyEmail verify email
func (s *userService) VerifyEmail(token string) (string, error) {
	emailCode := EmailCodeService.FindOne(sqls.NewCnd().Eq("token", token))
	if emailCode == nil || emailCode.Used {
		return "", errors.New(locale.T("errors.verify_email_code_invalid"))
	}

	user := s.Get(emailCode.UserID)
	if user == nil || emailCode.Email != user.Email.String {
		return "", errors.New(locale.T("errors.captcha_expired"))
	}
	if dates.FromTimestamp(emailCode.CreateTime).Add(time.Hour * time.Duration(emailVerifyExpireHour)).Before(time.Now()) {
		return "", errors.New(locale.T("errors.verify_email_expired"))
	}
	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := repository.UserRepository.UpdateColumn(tx, emailCode.UserID, "email_verified", true); err != nil {
			return err
		}
		cache.UserCache.Invalidate(emailCode.UserID)
		return repository.EmailCodeRepository.UpdateColumn(tx, emailCode.ID, "used", true)
	})
	if err != nil {
		return "", err
	}
	return emailCode.Email, nil
}

// CheckPostStatus checks user status when posting
func (s *userService) CheckPostStatus(user *model.User) error {
	if !user.IsActive || user.IsForbidden() {
		return errs.ErrForbidden
	}
	observeSeconds := SysConfigService.GetInt(constants.SysConfigUserObserveSeconds, 0)
	if user.InObservationPeriod(observeSeconds) {
		return web.NewError(errs.InObservationPeriod.Code, "Account is under observation for "+strconv.Itoa(observeSeconds)+" seconds, please try again later")
	}
	return nil
}

// IncrScoreForPostTopic points for posting a topic
func (s *userService) IncrScoreForLikeTopic(topic *model.Topic) {
	config := SysConfigService.GetPointConfig()
	if config.LikeTopicScore <= 0 {
		slog.Info("Points configuration missing for liking a topic")
		return
	}
	err := s.addScore(topic.UserID, config.LikeTopicScore, constants.EntityTopic,
		strconv.FormatInt(topic.ID, 10), locale.T("points.points_added_for_liking_topic"))
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
}

// IncrScoreForPostTopic points for posting a topic
func (s *userService) IncrScoreForPostTopic(topic *model.Topic) {
	config := SysConfigService.GetPointConfig()
	if config.PostTopicScore <= 0 {
		slog.Info("Points configuration missing for posting topics")
		return
	}
	err := s.addScore(topic.UserID, config.PostTopicScore, constants.EntityTopic,
		strconv.FormatInt(topic.ID, 10), locale.T("points.points_added_for_posting_topic"))
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
}

// IncrScoreForPostComment points for posting a comment
func (s *userService) IncrScoreForPostComment(comment *model.Comment) {
	config := SysConfigService.GetPointConfig()
	if config.PostCommentScore <= 0 {
		slog.Info("Points configuration missing for posting comments")
		return
	}
	err := s.addScore(comment.UserID, config.PostCommentScore, constants.EntityComment,
		strconv.FormatInt(comment.ID, 10), locale.T("points.points_added_for_commenting"))
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
}

// IncrScore increase score
func (s *userService) IncrScore(userId int64, score int, sourceType, sourceId, description string) error {
	if score <= 0 {
		return errors.New(locale.T("points.adjustment_must_be_positive"))
	}
	return s.addScore(userId, score, sourceType, sourceId, description)
}

// DecrScore decrease score
func (s *userService) DecrScore(userId int64, score int, sourceType, sourceId, description string) error {
	if score <= 0 {
		return errors.New(locale.T("points.adjustment_must_be_positive"))
	}
	return s.addScore(userId, -score, sourceType, sourceId, description)
}

// addScore add score (can be negative)
func (s *userService) addScore(userId int64, score int, sourceType, sourceId, description string) error {
	user := s.Get(userId)
	if user == nil {
		return errors.New(locale.T("user.not_found"))
	}
	if err := s.Updates(userId, map[string]interface{}{
		"score":       gorm.Expr("score + ?", score),
		"update_time": dates.NowTimestamp(),
	}); err != nil {
		return err
	}

	scoreType := constants.ScoreTypeIncr
	if score < 0 {
		scoreType = constants.ScoreTypeDecr
	}
	err := UserScoreLogService.Create(&model.UserScoreLog{
		UserID:      userId,
		SourceType:  sourceType,
		SourceID:    sourceId,
		Description: description,
		Type:        scoreType,
		Score:       score,
		CreateTime:  dates.NowTimestamp(),
	})
	if err == nil {
		cache.UserCache.Invalidate(userId)
	}
	return err
}

func (s *userService) GetUsersByRole(roles ...string) []model.User {
	var users []model.User
	err := sqls.DB().Model(&model.Role{}).Find(&users).Where("role IN ?", roles).Error
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return nil
	}
	return users
}
