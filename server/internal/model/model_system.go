package model

type Role struct {
	Model
	Type       int    `gorm:"not null;default:1" json:"type" form:"type"`             // Role type (0: system role, 1: custom role)
	Name       string `gorm:"size:64" json:"name" form:"name"`                        // Role name
	Code       string `gorm:"unique;size:64" json:"code" form:"code"`                 // Role code
	SortNo     int    `json:"sortNo" form:"sortNo"`                                   // Sort
	Remark     string `gorm:"size:256" json:"remark" form:"remark"`                   // Remark
	Status     int    `json:"status" form:"status"`                                   // Status
	CreateTime int64  `gorm:"not null;default:0" json:"createTime" form:"createTime"` // Create time
	UpdateTime int64  `gorm:"not null;default:0" json:"updateTime" form:"updateTime"` // Update time
}

// deprecated
type Menu struct {
	Model
	ParentId   int64  `json:"parentId" form:"parentId"`                               // Parent menu
	Name       string `gorm:"size:256" json:"name" form:"name"`                       // Name
	Title      string `gorm:"size:64" json:"title" form:"title"`                      // Title
	Icon       string `gorm:"size:1024" json:"icon" form:"icon"`                      // ICON
	Path       string `gorm:"size:1024" json:"path" form:"path"`                      // Path
	SortNo     int    `gorm:"not null;default:0" json:"sortNo" form:"sortNo"`         // Sort
	Status     int    `json:"status" form:"status"`                                   // Status
	CreateTime int64  `gorm:"not null;default:0" json:"createTime" form:"createTime"` // Create time
	UpdateTime int64  `gorm:"not null;default:0" json:"updateTime" form:"updateTime"` // Update time
}

type UserRole struct {
	Model
	UserId     int64 `gorm:"uniqueIndex:idx_user_role" json:"userId" form:"userId"`
	RoleId     int64 `gorm:"uniqueIndex:idx_user_role" json:"roleId" form:"roleId"`
	CreateTime int64 `gorm:"not null;default:0" json:"createTime" form:"createTime"` // Create time
}

type RoleMenu struct {
	Model
	RoleId     int64 `gorm:"uniqueIndex:idx_role_menu" json:"roleId" form:"roleId"`
	MenuId     int64 `gorm:"uniqueIndex:idx_role_menu" json:"menuId" form:"menuId"`
	CreateTime int64 `gorm:"not null;default:0" json:"createTime" form:"createTime"` // Create time
}
