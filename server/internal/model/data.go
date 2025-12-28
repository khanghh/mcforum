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
	LikeTopicScore          int     `json:"likeTopicScore"`          // points for liking a topic
	LikeCommentScore        int     `json:"likeCommentScore"`        // points for liking a comment
	PostTopicScore          int     `json:"postTopicScore"`          // points for posting a topic
	PostCommentScore        int     `json:"postCommentScore"`        // points for commenting / replying
	ReceiveTopicLikeScore   int     `json:"receiveTopicLikeScore"`   // points for topic being liked
	ReceiveCommentLikeScore int     `json:"receiveCommentLikeScore"` // points for comment being liked
	CheckInScore            int     `json:"checkInScore"`            // check-in points
	Streak7DaysScore        int     `json:"streak7DaysScore"`        // 7-day streak points
	Streak30DaysScore       int     `json:"streak30DaysScore"`       // 30-day streak points
	ReceiveFollowScore      int     `json:"receiveFollowScore"`      // points for being followed
	BoostMultiplier         float64 `json:"boostMultiplier"`         // temporary XP/event multiplier (default 1.0)
	DailyMaxScore           int     `json:"dailyMaxScore"`           // daily maximum points
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
