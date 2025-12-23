package payload

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/pkg/bbsurls"
	"strconv"
	"strings"
	"time"

	"bbs-go/common/dates"
	"bbs-go/common/strs"
	"bbs-go/sqls"
)

// UserInfo 用户简单信息
type UserInfo struct {
	Id            int64      `json:"id"`
	Type          int        `json:"type"`
	Username      string     `json:"username"`
	Nickname      string     `json:"nickname"`
	Avatar        string     `json:"avatar"`
	SmallAvatar   string     `json:"smallAvatar"`
	Gender        string     `json:"gender"`
	Birthday      *time.Time `json:"birthday"`
	TopicCount    int        `json:"topicCount"`   // 话题数量
	CommentCount  int        `json:"commentCount"` // 跟帖数量
	FansCount     int        `json:"fansCount"`    // 粉丝数量
	FollowCount   int        `json:"followCount"`  // 关注数量
	Score         int        `json:"score"`        // 积分
	Bio           string     `json:"bio"`
	StatusMessage string     `json:"statusMessage"`
	CreateTime    int64      `json:"createTime"`

	Forbidden bool `json:"forbidden"` // 是否禁言
	Followed  bool `json:"followed"`  // 是否关注
}

// UserDetail 用户详细信息
type UserDetail struct {
	UserInfo
	BackgroundImage      string `json:"backgroundImage"`
	SmallBackgroundImage string `json:"smallBackgroundImage"`
	HomePage             string `json:"homePage"`
	IsActive             bool   `json:"isActive"`
}

// UserProfile 用户个人信息
type UserProfile struct {
	UserDetail
	Roles         []string `json:"roles"`
	PasswordSet   bool     `json:"passwordSet"` // 密码已设置
	Email         string   `json:"email"`
	EmailVerified bool     `json:"emailVerified"`
}

func BuildUserInfoDefaultIfNull(id int64) *UserInfo {
	user := cache.UserCache.Get(id)
	if user == nil {
		user = &model.User{}
		user.Id = id
		user.Type = constants.UserTypeNormal
		user.Username = sqls.SqlNullString(strconv.FormatInt(id, 10))
		user.Nickname = "user" + strconv.FormatInt(id, 10)
		user.CreateTime = dates.NowTimestamp()
	}
	return BuildUserInfo(user)
}

func BuildUserInfo(user *model.User) *UserInfo {
	if user == nil {
		return nil
	}
	ret := &UserInfo{
		Id:            user.Id,
		Type:          user.Type,
		Username:      user.Username.String,
		Nickname:      user.Nickname,
		Gender:        user.Gender,
		Birthday:      user.Birthday,
		TopicCount:    user.TopicCount,
		CommentCount:  user.CommentCount,
		FansCount:     user.FansCount,
		FollowCount:   user.FollowCount,
		Score:         user.Score,
		Bio:           user.Bio,
		StatusMessage: user.StatusMessage,
		CreateTime:    user.CreateTime,
		Forbidden:     user.IsForbidden(),
	}
	if strs.IsNotBlank(user.Avatar) {
		ret.Avatar = user.Avatar
		ret.SmallAvatar = HandleOssImageStyleAvatar(user.Avatar)
	} else {
		avatar := bbsurls.AbsUrl("/images/avatars/steve.png")
		ret.Avatar = avatar
		ret.SmallAvatar = avatar
	}
	if !user.IsActive {
		ret.Nickname = ret.Username
		ret.Bio = ""
		ret.Score = 0
		ret.Forbidden = true
	}
	return ret
}

func BuildUserDetail(user *model.User) *UserDetail {
	if user == nil {
		return nil
	}
	ret := &UserDetail{
		UserInfo:             *BuildUserInfo(user),
		BackgroundImage:      user.BackgroundImage,
		SmallBackgroundImage: HandleOssImageStyleSmall(user.BackgroundImage),
		HomePage:             user.HomePage,
		IsActive:             user.IsActive,
	}
	if !user.IsActive {
		ret.Username = "blacklist"
		ret.HomePage = ""
	}
	return ret
}

func BuildUserProfile(user *model.User) *UserProfile {
	if user == nil {
		return nil
	}
	ret := &UserProfile{
		UserDetail:    *BuildUserDetail(user),
		Email:         user.Email.String,
		EmailVerified: user.EmailVerified,
		PasswordSet:   len(user.Password) > 0,
	}

	if strs.IsNotBlank(user.Roles) {
		ret.Roles = strings.Split(user.Roles, ",")
	}
	return ret
}
