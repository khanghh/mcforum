package model

var Models = []interface{}{
	&UserRole{}, &Role{}, &Menu{}, &RoleMenu{},

	&User{}, &UserToken{}, &Tag{}, &Article{}, &ArticleTag{}, &Comment{}, &Favorite{}, &Forum{},
	&Topic{}, &TopicTag{}, &UserLike{}, &Message{}, &SysConfig{}, &Link{},
	&UserScoreLog{}, &OperateLog{}, &EmailCode{}, &CheckIn{}, &UserFollow{}, &UserFeed{}, &UserReport{},
	&ForbiddenWord{},
}

type Model struct {
	Id int64 `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
}
