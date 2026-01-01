package payload

type ForumStatsResponse struct {
	TotalTopics  int64  `json:"topics"`
	TotalPosts   int64  `json:"posts"`
	TotalMembers int64  `json:"members"`
	TotalVisits  int64  `json:"visits"`
	NewestMember string `json:"newestMember,omitempty"`
}
