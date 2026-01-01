package model

var Models = []interface{}{
	&Role{}, &Menu{}, &RoleMenu{},

	&User{}, &UserToken{}, &Tag{}, &Article{}, &ArticleTag{}, &Comment{}, &Favorite{}, &Forum{},
	&Topic{}, &TopicTag{}, &UserLike{}, &Message{}, &SysConfig{}, &Link{},
	&UserScoreLog{}, &OperateLog{}, &EmailCode{}, &CheckIn{}, &UserFollow{}, &UserFeed{}, &UserReport{},
	&ForbiddenWord{}, &Upload{}, &ForumStats{},
}

type Model struct {
	ID int64 `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
}
