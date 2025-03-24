package model

// 站点导航
type ActionLink struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

// 积分配置
type ScoreConfig struct {
	PostTopicScore   int `json:"postTopicScore"`   // 发帖获得积分
	PostCommentScore int `json:"postCommentScore"` // 跟帖获得积分
	CheckInScore     int `json:"checkInScore"`     // 签到积分
}

type LoginMethod struct {
	Password bool `json:"password"`
	QQ       bool `json:"qq"`
	Github   bool `json:"github"`
	Osc      bool `json:"osc"`
}

// ModulesConfig
//
//	模块配置
type ModulesConfig struct {
	Tweet   bool `json:"tweet"`
	Topic   bool `json:"topic"`
	Article bool `json:"article"`
}
