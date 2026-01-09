package model

import (
	"bbs-go/common/dates"
	"bbs-go/internal/model/constants"
	"database/sql"
	"time"
)

type UserToken struct {
	Model
	Token      string `gorm:"size:32;unique;not null" json:"token" form:"token"`
	UserID     int64  `gorm:"not null;index:idx_user_token_user_id;" json:"userId" form:"userId"`
	ExpiredAt  int64  `gorm:"not null" json:"expiredAt" form:"expiredAt"`
	Status     int    `gorm:"type:int(11);not null;index:idx_user_token_status" json:"status" form:"status"`
	CreateTime int64  `gorm:"not null" json:"createTime" form:"createTime"`
}

// User like
type UserLike struct {
	Model
	UserID     int64  `gorm:"not null;uniqueIndex:idx_user_like_unique;" json:"userId" form:"userId"`                                            // User
	EntityID   int64  `gorm:"not null;uniqueIndex:idx_user_like_unique;index:idx_user_like_entity;" json:"topicId" form:"topicId"`               // Entity ID
	EntityType string `gorm:"not null;size:32;uniqueIndex:idx_user_like_unique;index:idx_user_like_entity;" json:"entityType" form:"entityType"` // Entity type
	Status     int    `gorm:"type:int(11);not null" json:"status" form:"status"`                                                                 // Status: 0: unread, 1: read
	CreateTime int64  `json:"createTime" form:"createTime"`                                                                                      // Create time
}

// User score log
type UserScoreLog struct {
	Model
	UserID      int64  `gorm:"not null;index:idx_user_score_log_user_id" json:"userId" form:"userId"`   // User ID
	SourceType  string `gorm:"not null;index:idx_user_score_score" json:"sourceType" form:"sourceType"` // Score source type
	SourceID    string `gorm:"not null;index:idx_user_score_score" json:"sourceId" form:"sourceId"`     // Score source ID
	Description string `json:"description" form:"description"`                                          // Description
	Type        int    `gorm:"type:int(11)" json:"type" form:"type"`                                    // Type (increase, decrease)
	Score       int    `gorm:"type:int(11)" json:"score" form:"score"`                                  // Score
	CreateTime  int64  `json:"createTime" form:"createTime"`                                            // Create time
}

// Check in
type CheckIn struct {
	Model
	UserID          int64 `gorm:"not null;uniqueIndex:idx_user_id" json:"userId" form:"userId"`         // User ID
	LatestDayName   int   `gorm:"type:int(11);not null;index:idx_latest" json:"dayName" form:"dayName"` // Latest check-in day
	ConsecutiveDays int   `gorm:"type:int(11);not null;" json:"consecutiveDays" form:"consecutiveDays"` // Consecutive check-in days
	CreateTime      int64 `json:"createTime" form:"createTime"`                                         // Create time
	UpdateTime      int64 `gorm:"index:idx_latest" json:"updateTime" form:"updateTime"`                 // Update time
}

func (CheckIn) TableName() string {
	return "check_in"
}

// UserReport user report
type UserReport struct {
	Model
	DataID      int64  `json:"dataId" form:"dataId"`           // Report data ID
	DataType    string `json:"dataType" form:"dataType"`       // Report data type
	UserID      int64  `json:"userId" form:"userId"`           // Reporter ID
	Reason      string `json:"reason" form:"reason"`           // Report reason
	AuditStatus int64  `json:"auditStatus" form:"auditStatus"` // Audit status
	AuditTime   int64  `json:"auditTime" form:"auditTime"`     // Audit time
	AuditUserID int64  `json:"auditUserId" form:"auditUserId"` // Auditor ID
	CreateTime  int64  `json:"createTime" form:"createTime"`   // Report time
}

type UserRole struct {
	Model
	UserID     int64 `gorm:"primaryKey;uniqueIndex:idx_user_role" json:"userId" form:"userId"`
	RoleID     int64 `gorm:"primaryKey;uniqueIndex:idx_user_role" json:"roleId" form:"roleId"`
	CreateTime int64 `gorm:"not null;default:0" json:"createTime" form:"createTime"` // Create time
}

type User struct {
	Model
	Type            int            `gorm:"not null;default:0" json:"type" form:"type"`                                     // User type (0: user, 1: staff)
	Username        sql.NullString `gorm:"size:32;unique;" json:"username" form:"username"`                                // Username
	Email           sql.NullString `gorm:"size:128;unique;" json:"email" form:"email"`                                     // Email
	EmailVerified   bool           `gorm:"not null;default:false" json:"emailVerified" form:"emailVerified"`               // Email verified
	Nickname        string         `gorm:"size:32;" json:"nickname" form:"nickname"`                                       // Nickname
	Avatar          string         `gorm:"type:text" json:"avatar" form:"avatar"`                                          // Avatar
	Gender          string         `gorm:"size:16;default:''" json:"gender" form:"gender"`                                 // Gender
	Birthday        *time.Time     `json:"birthday" form:"birthday"`                                                       // Birthday
	BackgroundImage string         `gorm:"type:text" json:"backgroundImage" form:"backgroundImage"`                        // Background image
	Password        string         `gorm:"size:512" json:"password" form:"password"`                                       // Password
	Bio             string         `gorm:"type:varchar(255)" json:"bio" form:"bio"`                                        // Bio
	StatusMessage   string         `gorm:"type:varchar(128)" json:"statusMessage" form:"statusMessage"`                    // Status message
	Location        string         `gorm:"size:128" json:"location" form:"location"`                                       // Location
	LockedProfile   bool           `gorm:"not null;default:false" json:"lockedProfile" form:"lockedProfile"`               // Locked profile
	ShowLocation    bool           `gorm:"not null;default:true" json:"showLocation" form:"showLocation"`                  // Show location
	EmailNotify     bool           `gorm:"not null;default:true" json:"emailNotify" form:"emailNotify"`                    // Email notifications
	Score           int            `gorm:"type:int(11);not null;default:0;index:idx_user_score" json:"score" form:"score"` // Score
	IsActive        bool           `gorm:"not null;default:true" json:"isActive" form:"isActive"`                          // Status
	TopicCount      int            `gorm:"type:int(11);not null;default:0" json:"topicCount" form:"topicCount"`            // Topic count
	CommentCount    int            `gorm:"type:int(11);not null;default:0" json:"commentCount" form:"commentCount"`        // Comment count
	FollowersCount  int            `gorm:"type:int(11);not null;default:0" json:"followersCount" form:"followersCount"`    // Followers count
	FollowingCount  int            `gorm:"type:int(11);not null;default:0" json:"followingCount" form:"followingCount"`    // Following count
	ActivityCount   int            `gorm:"type:int(11);not null;default:0" json:"activityCount" form:"activityCount"`      // Activity count
	PlaytimeSec     int64          `gorm:"type:bigint(20);not null;default:0" json:"playtimeSec" form:"playtimeSec"`       // Total play time in seconds
	RoleID          sql.NullInt64  `gorm:"default:null;index:idx_user_role_id" json:"roleId" form:"roleId"`                // Role ID
	Role            Role           `gorm:"foreignKey:RoleID;constraint:constraint:OnDelete:SET 0;"`

	ForbiddenEndTime int64 `gorm:"not null;default:0" json:"forbiddenEndTime" form:"forbiddenEndTime"` // Forbidden end time
	CreateTime       int64 `json:"createTime" form:"createTime"`                                       // Create time
	UpdateTime       int64 `json:"updateTime" form:"updateTime"`                                       // Update time
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
	tmp := dates.NowTimestamp()
	return u.ForbiddenEndTime > tmp
}

// HasRole whether has specified role
func (u *User) HasRole(roleName string) bool {
	if u.Role.Name == roleName {
		return true
	}
	return false
}

// HasAnyRole whether has any of the specified roles
func (u *User) HasAnyRole(roleNames ...string) bool {
	if len(roleNames) == 0 {
		return false
	}
	for _, r := range roleNames {
		if u.HasRole(r) {
			return true
		}
	}
	return false
}

// IsManagerRole whether owner, admin, or moderator
func (u *User) IsManagerRole() bool {
	return u.HasAnyRole(constants.RoleOwner, constants.RoleAdmin, constants.RoleModerator)
}

// InObservationPeriod whether in observation period
// observeSeconds observation duration
func (u *User) InObservationPeriod(observeSeconds int) bool {
	if observeSeconds <= 0 {
		return false
	}
	return dates.FromTimestamp(u.CreateTime).Add(time.Second * time.Duration(observeSeconds)).After(time.Now())
}
