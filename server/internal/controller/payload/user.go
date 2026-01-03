package payload

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"path/filepath"
	"strconv"
	"strings"
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
	Score         int    `json:"score"`                   // User score
	StatusMessage string `json:"statusMessage,omitempty"` // Status message
	JoinTime      int64  `json:"joinTime"`                // Account creation timestamp

	IsForbidden bool `json:"isForbidden,omitempty"` // Whether the user is banned
	IsFollowing bool `json:"isFollowing,omitempty"` // Whether the current user is following this user
}

// UserDetail is detailed user info public based on user settings
type UserDetail struct {
	UserInfo
	Gender               string     `json:"gender,omitempty"`               // User gender
	Birthday             *time.Time `json:"birthday,omitempty"`             // User birthday
	TopicCount           int        `json:"topicCount"`                     // Number of topics created
	CommentCount         int        `json:"commentCount"`                   // Number of comments made
	FollowersCount       int        `json:"followersCount"`                 // Number of fans
	FollowingCount       int        `json:"followingCount"`                 // Number of users followed
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
		user.Username = sqls.SqlNullString("user" + strconv.FormatInt(id, 10))
		user.Nickname = "user" + strconv.FormatInt(id, 10)
		user.CreateTime = dates.NowTimestamp()
	}
	return BuildUserInfo(user)
}

func getThumbImagePath(filePath string) string {
	ext := filepath.Ext(filePath)
	base := strings.TrimSuffix(filePath, ext)
	return base + "_thumb" + ext
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
		Avatar:        user.Avatar,
		Role:          roleName,
		Bio:           user.Bio,
		Score:         user.Score,
		StatusMessage: user.StatusMessage,
		JoinTime:      user.CreateTime,
		IsForbidden:   user.IsForbidden(),
	}
	if strs.IsNotBlank(user.Avatar) {
		ret.SmallAvatar = getThumbImagePath(user.Avatar)
	}
	if !user.IsActive {
		ret.Nickname = "Unknown"
		ret.Bio = ""
		ret.IsForbidden = true
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
		FollowersCount:       user.FollowersCount,
		FollowingCount:       user.FollowingCount,
		BackgroundImage:      user.BackgroundImage,
		SmallBackgroundImage: GetSmallCoverURL(user.BackgroundImage),
	}
	if user.ShowLocation {
		ret.Location = user.Location
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
		Settings: UserSettings{
			LockedProfile: user.LockedProfile,
			ShowLocation:  user.ShowLocation,
			EmailNotify:   user.EmailNotify,
		},
	}
	return ret
}
