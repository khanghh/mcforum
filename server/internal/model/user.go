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

// 用户点赞
type UserLike struct {
	Model
	UserId     int64  `gorm:"not null;uniqueIndex:idx_user_like_unique;" json:"userId" form:"userId"`                                            // 用户
	EntityId   int64  `gorm:"not null;uniqueIndex:idx_user_like_unique;index:idx_user_like_entity;" json:"topicId" form:"topicId"`               // 实体编号
	EntityType string `gorm:"not null;size:32;uniqueIndex:idx_user_like_unique;index:idx_user_like_entity;" json:"entityType" form:"entityType"` // 实体类型
	Status     int    `gorm:"type:int(11);not null" json:"status" form:"status"`                                                                 // 状态：0：未读、1：已读
	CreateTime int64  `json:"createTime" form:"createTime"`                                                                                      // 创建时间
}

// 用户积分流水
type UserScoreLog struct {
	Model
	UserId      int64  `gorm:"not null;index:idx_user_score_log_user_id" json:"userId" form:"userId"`   // 用户编号
	SourceType  string `gorm:"not null;index:idx_user_score_score" json:"sourceType" form:"sourceType"` // 积分来源类型
	SourceId    string `gorm:"not null;index:idx_user_score_score" json:"sourceId" form:"sourceId"`     // 积分来源编号
	Description string `json:"description" form:"description"`                                          // 描述
	Type        int    `gorm:"type:int(11)" json:"type" form:"type"`                                    // 类型(增加、减少)
	Score       int    `gorm:"type:int(11)" json:"score" form:"score"`                                  // 积分
	CreateTime  int64  `json:"createTime" form:"createTime"`                                            // 创建时间
}

// 签到
type CheckIn struct {
	Model
	UserId          int64 `gorm:"not null;uniqueIndex:idx_user_id" json:"userId" form:"userId"`         // 用户编号
	LatestDayName   int   `gorm:"type:int(11);not null;index:idx_latest" json:"dayName" form:"dayName"` // 最后一次签到
	ConsecutiveDays int   `gorm:"type:int(11);not null;" json:"consecutiveDays" form:"consecutiveDays"` // 连续签到天数
	CreateTime      int64 `json:"createTime" form:"createTime"`                                         // 创建时间
	UpdateTime      int64 `gorm:"index:idx_latest" json:"updateTime" form:"updateTime"`                 // 更新时间
}

// UserReport 用户举报
type UserReport struct {
	Model
	DataId      int64  `json:"dataId" form:"dataId"`           // 举报数据ID
	DataType    string `json:"dataType" form:"dataType"`       // 举报数据类型
	UserId      int64  `json:"userId" form:"userId"`           // 举报人ID
	Reason      string `json:"reason" form:"reason"`           // 举报原因
	AuditStatus int64  `json:"auditStatus" form:"auditStatus"` // 审核状态
	AuditTime   int64  `json:"auditTime" form:"auditTime"`     // 审核时间
	AuditUserId int64  `json:"auditUserId" form:"auditUserId"` // 审核人ID
	CreateTime  int64  `json:"createTime" form:"createTime"`   // 举报时间
}

type User struct {
	Model
	Type             int            `gorm:"not null;default:0" json:"type" form:"type"`                                    // 用户类型（0：用户、1：员工）
	Username         sql.NullString `gorm:"size:32;unique;" json:"username" form:"username"`                               // 用户名
	Email            sql.NullString `gorm:"size:128;unique;" json:"email" form:"email"`                                    // 邮箱
	EmailVerified    bool           `gorm:"not null;default:false" json:"emailVerified" form:"emailVerified"`              // 邮箱是否验证
	Nickname         string         `gorm:"size:16;" json:"nickname" form:"nickname"`                                      // 昵称
	Avatar           string         `gorm:"type:text" json:"avatar" form:"avatar"`                                         // 头像
	Gender           string         `gorm:"size:16;default:''" json:"gender" form:"gender"`                                // 性别
	Birthday         *time.Time     `json:"birthday" form:"birthday"`                                                      // 生日
	BackgroundImage  string         `gorm:"type:text" json:"backgroundImage" form:"backgroundImage"`                       // 个人中心背景图片
	Password         string         `gorm:"size:512" json:"password" form:"password"`                                      // 密码
	HomePage         string         `gorm:"size:1024" json:"homePage" form:"homePage"`                                     // 个人主页
	Bio              string         `gorm:"type:text" json:"description" form:"description"`                               // 个人描述
	StatusMessage    string         `gorm:"type:varchar(255)" json:"statusMessage" form:"statusMessage"`                   // 状态信息
	Score            int            `gorm:"type:int(11);not null;index:idx_user_score" json:"score" form:"score"`          // 积分
	IsActive         bool           `gorm:"type:int(11);index:idx_user_status;not null" json:"activated" form:"activated"` // 状态
	TopicCount       int            `gorm:"type:int(11);not null" json:"topicCount" form:"topicCount"`                     // 帖子数量
	CommentCount     int            `gorm:"type:int(11);not null" json:"commentCount" form:"commentCount"`                 // 跟帖数量
	FollowCount      int            `gorm:"type:int(11);not null" json:"followCount" form:"followCount"`                   // 关注数量
	FansCount        int            `gorm:"type:int(11);not null" json:"fansCount" form:"fansCount"`                       // 粉丝数量
	Roles            string         `gorm:"type:text" json:"roles" form:"roles"`                                           // 角色
	ForbiddenEndTime int64          `gorm:"not null;default:0" json:"forbiddenEndTime" form:"forbiddenEndTime"`            // 禁言结束时间
	CreateTime       int64          `json:"createTime" form:"createTime"`                                                  // 创建时间
	UpdateTime       int64          `json:"updateTime" form:"updateTime"`                                                  // 更新时间
}

// IsForbidden 是否禁言
func (u *User) IsForbidden() bool {
	if u.ForbiddenEndTime == 0 {
		return false
	}
	// 永久禁言
	if u.ForbiddenEndTime == -1 {
		return true
	}
	// 判断禁言时间
	return u.ForbiddenEndTime > dates.NowTimestamp()
}

// HasRole 是否有指定角色
func (u *User) HasRole(role string) bool {
	roles := strings.Split(u.Roles, ",")
	if len(roles) == 0 {
		return false
	}
	return arrays.Contains(roles, role)
}

// HasAnyRole 是否有指定的任意角色
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

// IsOwnerOrAdmin 是否是管理员
func (u *User) IsOwnerOrAdmin() bool {
	return u.HasAnyRole(constants.RoleOwner, constants.RoleAdmin)
}

// GetRoles 获取角色
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

// InObservationPeriod 是否在观察期
// observeSeconds 观察时长
func (u *User) InObservationPeriod(observeSeconds int) bool {
	if observeSeconds <= 0 {
		return false
	}
	return dates.FromTimestamp(u.CreateTime).Add(time.Second * time.Duration(observeSeconds)).After(time.Now())
}
