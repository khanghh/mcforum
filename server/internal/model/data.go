package model

// Site navigation
type ActionLink struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

type MenuItem struct {
	Name    string `json:"name"`
	URLPath string `json:"urlPath"`
	LogoURL string `json:"logoUrl,omitempty"`
}

// Score configuration
type ScoreConfig struct {
	PostTopicScore   int `json:"postTopicScore"`   // points for posting
	PostCommentScore int `json:"postCommentScore"` // points for commenting
	CheckInScore     int `json:"checkInScore"`     // check-in points
}

type LoginMethod struct {
	Password bool `json:"password"`
	QQ       bool `json:"qq"`
	Github   bool `json:"github"`
	Osc      bool `json:"osc"`
}

// ModulesConfig
//
//	Modules configuration
type ModulesConfig struct {
	Tweet   bool `json:"tweet"`
	Topic   bool `json:"topic"`
	Article bool `json:"article"`
}
