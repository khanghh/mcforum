package payload

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/pkg/bbsurls"
	"strconv"
	"time"

	"bbs-go/common/dates"
	"bbs-go/common/strs"
	"bbs-go/sqls"
)

// UserInfo user simple info
type UserInfo struct {
	Id            int64      `json:"id"`
	Type          int        `json:"type"`
	Username      string     `json:"username"`
	Nickname      string     `json:"nickname"`
	Avatar        string     `json:"avatar"`
	SmallAvatar   string     `json:"smallAvatar"`
	Gender        string     `json:"gender"`
	Birthday      *time.Time `json:"birthday"`
	TopicCount    int        `json:"topicCount"`   // topic count
	CommentCount  int        `json:"commentCount"` // comment count
	FansCount     int        `json:"fansCount"`    // fans count
	FollowCount   int        `json:"followCount"`  // follow count
	Score         int        `json:"score"`        // score
	Bio           string     `json:"bio"`
	StatusMessage string     `json:"statusMessage"`
	CreateTime    int64      `json:"createTime"`

	Forbidden   bool `json:"forbidden"`   // whether banned
	IsFollowing bool `json:"isFollowing"` // whether following
}

type UserSettings struct {
	LockedProfile bool `json:"lockedProfile"` // Locked profile
	ShowLocation  bool `json:"showLocation"`  // Show location
	EmailNotify   bool `json:"emailNotify"`   // Email notifications
}

// UserDetail user detailed info
type UserDetail struct {
	UserInfo
	BackgroundImage      string `json:"backgroundImage"`
	SmallBackgroundImage string `json:"smallBackgroundImage"`
	Location             string `json:"location"`
}

// UserProfile user personal info
type UserProfile struct {
	UserDetail
	Roles         []string     `json:"roles"`
	PasswordSet   bool         `json:"passwordSet"` // password set
	Email         string       `json:"email"`
	EmailVerified bool         `json:"emailVerified"`
	JoinTime      int64        `json:"joinTime"`
	IsActive      bool         `json:"isActive"`
	Settings      UserSettings `json:"settings,omitempty"`
}

func BuildUserInfoDefaultIfNull(id int64) *UserInfo {
	user := cache.UserCache.Get(id)
	if user == nil {
		user = &model.User{}
		user.ID = id
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
		Id:            user.ID,
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
		ret.Nickname = "Unknown"
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
		Location:             user.Location,
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
		IsActive:      user.IsActive,
	}
	if len(user.Roles) > 0 {
		var roles []string
		for _, role := range user.Roles {
			roles = append(roles, role.Name)
		}
		ret.Roles = roles
	}
	return ret
}
