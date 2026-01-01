package model

type ForumStats struct {
	ID            int64  `gorm:"primaryKey"`
	TotalTopics   int64  `gorm:"type:bigint;not null;default:0"`
	TotalComments int64  `gorm:"type:bigint;not null;default:0"`
	TotalMembers  int64  `gorm:"type:bigint;not null;default:0"`
	TotalVisits   int64  `gorm:"type:bigint;not null;default:0"`
	NewestMember  string `gorm:"type:varchar(64)"`
}
