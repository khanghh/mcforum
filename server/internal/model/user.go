package model

import (
	"bbs-go/common/arrays"
	"bbs-go/common/dates"
	"bbs-go/common/strs"
	"bbs-go/internal/model/constants"
	"database/sql"
	"strings"
	"time"
)

type UserToken struct {
	Model
	Token      string `gorm:"size:32;unique;not null" json:"token" form:"token"`
	UserId     int64  `gorm:"not null;index:idx_user_token_user_id;" json:"userId" form:"userId"`
	ExpiredAt  int64  `gorm:"not null" json:"expiredAt" form:"expiredAt"`
	Status     int    `gorm:"type:int(11);not null;index:idx_user_token_status" json:"status" form:"status"`
	CreateTime int64  `gorm:"not null" json:"createTime" form:"createTime"`
}

// User like
type UserLike struct {
	Model
	UserId     int64  `gorm:"not null;uniqueIndex:idx_user_like_unique;" json:"userId" form:"userId"`                                            // User
	EntityId   int64  `gorm:"not null;uniqueIndex:idx_user_like_unique;index:idx_user_like_entity;" json:"topicId" form:"topicId"`               // Entity ID
	EntityType string `gorm:"not null;size:32;uniqueIndex:idx_user_like_unique;index:idx_user_like_entity;" json:"entityType" form:"entityType"` // Entity type
	Status     int    `gorm:"type:int(11);not null" json:"status" form:"status"`                                                                 // Status: 0: unread, 1: read
	CreateTime int64  `json:"createTime" form:"createTime"`                                                                                      // Create time
}

// User score log
type UserScoreLog struct {
	Model
	UserId      int64  `gorm:"not null;index:idx_user_score_log_user_id" json:"userId" form:"userId"`   // User ID
	SourceType  string `gorm:"not null;index:idx_user_score_score" json:"sourceType" form:"sourceType"` // Score source type
	SourceId    string `gorm:"not null;index:idx_user_score_score" json:"sourceId" form:"sourceId"`     // Score source ID
	Description string `json:"description" form:"description"`                                          // Description
	Type        int    `gorm:"type:int(11)" json:"type" form:"type"`                                    // Type (increase, decrease)
	Score       int    `gorm:"type:int(11)" json:"score" form:"score"`                                  // Score
	CreateTime  int64  `json:"createTime" form:"createTime"`                                            // Create time
}

// Check in
type CheckIn struct {
	Model
	UserId          int64 `gorm:"not null;uniqueIndex:idx_user_id" json:"userId" form:"userId"`         // User ID
	LatestDayName   int   `gorm:"type:int(11);not null;index:idx_latest" json:"dayName" form:"dayName"` // Latest check-in day
	ConsecutiveDays int   `gorm:"type:int(11);not null;" json:"consecutiveDays" form:"consecutiveDays"` // Consecutive check-in days
	CreateTime      int64 `json:"createTime" form:"createTime"`                                         // Create time
	UpdateTime      int64 `gorm:"index:idx_latest" json:"updateTime" form:"updateTime"`                 // Update time
}

// UserReport user report
type UserReport struct {
	Model
	DataId      int64  `json:"dataId" form:"dataId"`           // Report data ID
	DataType    string `json:"dataType" form:"dataType"`       // Report data type
	UserId      int64  `json:"userId" form:"userId"`           // Reporter ID
	Reason      string `json:"reason" form:"reason"`           // Report reason
	AuditStatus int64  `json:"auditStatus" form:"auditStatus"` // Audit status
	AuditTime   int64  `json:"auditTime" form:"auditTime"`     // Audit time
	AuditUserId int64  `json:"auditUserId" form:"auditUserId"` // Auditor ID
	CreateTime  int64  `json:"createTime" form:"createTime"`   // Report time
}

type User struct {
	Model
	Type             int            `gorm:"not null;default:0" json:"type" form:"type"`                                    // User type (0: user, 1: staff)
	Username         sql.NullString `gorm:"size:32;unique;" json:"username" form:"username"`                               // Username
	Email            sql.NullString `gorm:"size:128;unique;" json:"email" form:"email"`                                    // Email
	EmailVerified    bool           `gorm:"not null;default:false" json:"emailVerified" form:"emailVerified"`              // Email verified
	Nickname         string         `gorm:"size:16;" json:"nickname" form:"nickname"`                                      // Nickname
	Avatar           string         `gorm:"type:text" json:"avatar" form:"avatar"`                                         // Avatar
	Gender           string         `gorm:"size:16;default:''" json:"gender" form:"gender"`                                // Gender
	Birthday         *time.Time     `json:"birthday" form:"birthday"`                                                      // Birthday
	BackgroundImage  string         `gorm:"type:text" json:"backgroundImage" form:"backgroundImage"`                       // Background image
	Password         string         `gorm:"size:512" json:"password" form:"password"`                                      // Password
	HomePage         string         `gorm:"size:1024" json:"homePage" form:"homePage"`                                     // Home page
	Bio              string         `gorm:"type:text" json:"description" form:"description"`                               // Bio
	StatusMessage    string         `gorm:"type:varchar(255)" json:"statusMessage" form:"statusMessage"`                   // Status message
	Score            int            `gorm:"type:int(11);not null;index:idx_user_score" json:"score" form:"score"`          // Score
	IsActive         bool           `gorm:"type:int(11);index:idx_user_status;not null" json:"activated" form:"activated"` // Status
	TopicCount       int            `gorm:"type:int(11);not null" json:"topicCount" form:"topicCount"`                     // Topic count
	CommentCount     int            `gorm:"type:int(11);not null" json:"commentCount" form:"commentCount"`                 // Comment count
	FollowCount      int            `gorm:"type:int(11);not null" json:"followCount" form:"followCount"`                   // Follow count
	FansCount        int            `gorm:"type:int(11);not null" json:"fansCount" form:"fansCount"`                       // Fans count
	Roles            string         `gorm:"type:text" json:"roles" form:"roles"`                                           // Roles
	ForbiddenEndTime int64          `gorm:"not null;default:0" json:"forbiddenEndTime" form:"forbiddenEndTime"`            // Forbidden end time
	CreateTime       int64          `json:"createTime" form:"createTime"`                                                  // Create time
	UpdateTime       int64          `json:"updateTime" form:"updateTime"`                                                  // Update time
}

// IsForbidden whether forbidden
func (u *User) IsForbidden() bool {
	if u.ForbiddenEndTime == 0 {
		return false
	}
	// Permanent forbidden
	if u.ForbiddenEndTime == -1 {
		return true
	}
	// Check forbidden time
	return u.ForbiddenEndTime > dates.NowTimestamp()
}

// HasRole whether has specified role
func (u *User) HasRole(role string) bool {
	roles := strings.Split(u.Roles, ",")
	if len(roles) == 0 {
		return false
	}
	return arrays.Contains(roles, role)
}

// HasAnyRole whether has any of the specified roles
func (u *User) HasAnyRole(roles ...string) bool {
	if len(roles) == 0 {
		return false
	}
	userRoles := strings.Split(u.Roles, ",")
	for _, role := range userRoles {
		if arrays.Contains(roles, role) {
			return true
		}
	}
	return false
}

// IsOwnerOrAdmin whether owner or admin
func (u *User) IsOwnerOrAdmin() bool {
	return u.HasAnyRole(constants.RoleOwner, constants.RoleAdmin)
}

// GetRoles get roles
func (u *User) GetRoles() []string {
	if strs.IsBlank(u.Roles) {
		return nil
	}
	ss := strings.Split(u.Roles, ",")
	if len(ss) == 0 {
		return nil
	}
	var roles []string
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if strs.IsNotBlank(s) {
			roles = append(roles, s)
		}
	}
	return roles
}

// InObservationPeriod whether in observation period
// observeSeconds observation duration
func (u *User) InObservationPeriod(observeSeconds int) bool {
	if observeSeconds <= 0 {
		return false
	}
	return dates.FromTimestamp(u.CreateTime).Add(time.Second * time.Duration(observeSeconds)).After(time.Now())
}
