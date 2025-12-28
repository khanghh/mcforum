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

// UserInfo is public user info
type UserInfo struct {
	Id            int64  `json:"id"`                      // User ID
	Type          int    `json:"type"`                    // User type
	Username      string `json:"username"`                // Username
	Nickname      string `json:"nickname"`                // Display name
	Role          string `json:"role,omitempty"`          // User role name
	Avatar        string `json:"avatar"`                  // Avatar URL
	SmallAvatar   string `json:"smallAvatar"`             // Small avatar URL
	Bio           string `json:"bio,omitempty"`           // User bio
	StatusMessage string `json:"statusMessage,omitempty"` // Status message
	CreateTime    int64  `json:"createTime"`              // Account creation timestamp

	Forbidden   bool `json:"forbidden,omitempty"`   // Whether the user is banned
	IsFollowing bool `json:"isFollowing,omitempty"` // Whether the current user is following this user
}

// UserDetail is detailed user info public based on user settings
type UserDetail struct {
	UserInfo
	Gender               string     `json:"gender,omitempty"`               // User gender
	Birthday             *time.Time `json:"birthday,omitempty"`             // User birthday
	TopicCount           int        `json:"topicCount"`                     // Number of topics created
	CommentCount         int        `json:"commentCount"`                   // Number of comments made
	FansCount            int        `json:"fansCount"`                      // Number of fans
	FollowCount          int        `json:"followCount"`                    // Number of users followed
	Score                int        `json:"score"`                          // User score
	BackgroundImage      string     `json:"backgroundImage,omitempty"`      // Background image URL
	SmallBackgroundImage string     `json:"smallBackgroundImage,omitempty"` // Small background image URL
	Location             string     `json:"location,omitempty"`             // User location
}

// UserProfile is private user profile only visible to themselves
type UserProfile struct {
	UserDetail
	PasswordSet   bool         `json:"passwordSet"`        // Whether password is set
	Email         string       `json:"email"`              // User email
	EmailVerified bool         `json:"emailVerified"`      // Whether email is verified
	JoinTime      int64        `json:"joinTime"`           // Account join timestamp
	IsActive      bool         `json:"isActive"`           // Whether account is active
	Settings      UserSettings `json:"settings,omitempty"` // User settings
}

type UserSettings struct {
	LockedProfile bool `json:"lockedProfile"` // Whether profile is locked
	ShowLocation  bool `json:"showLocation"`  // Whether to show location
	EmailNotify   bool `json:"emailNotify"`   // Whether to receive email notifications
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
	roleName := ""
	if user.Role != nil {
		roleName = user.Role.Name
	}
	ret := &UserInfo{
		Id:            user.ID,
		Type:          user.Type,
		Username:      user.Username.String,
		Nickname:      user.Nickname,
		Role:          roleName,
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
		Gender:               user.Gender,
		Birthday:             user.Birthday,
		TopicCount:           user.TopicCount,
		CommentCount:         user.CommentCount,
		FansCount:            user.FansCount,
		FollowCount:          user.FollowCount,
		Score:                user.Score,
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
		JoinTime:      user.CreateTime,
		IsActive:      user.IsActive,
		Settings: UserSettings{
			LockedProfile: user.LockedProfile,
			ShowLocation:  user.ShowLocation,
			EmailNotify:   user.EmailNotify,
		},
	}
	return ret
}
