package service

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/model"
	"bbs-go/sqls"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	forumStatsCacheDuration = 5 * time.Second
)

var StatsService = newStatsService()

type statsService struct {
	statsCache   model.ForumStats
	lastCache    time.Time
	visitCounter int64
	lock         chan struct{}
}

func (c *statsService) IncreaseVisit() int64 {
	c.visitCounter++
	c.statsCache.TotalVisits++
	return c.statsCache.TotalVisits
}

func (c *statsService) GetForumStats() (model.ForumStats, error) {
	if time.Since(c.lastCache) > forumStatsCacheDuration {
		select {
		case c.lock <- struct{}{}:
			defer func() { <-c.lock }()
			totalTopics := TopicService.Count(sqls.NewCnd())
			totalComments := CommentService.Count(sqls.NewCnd())
			totalMembers := UserService.Count(sqls.NewCnd())
			newWestMember := UserService.FindOne(sqls.NewCnd().Desc("create_time"))
			stats := model.ForumStats{ID: 1}
			err := sqls.DB().Clauses(
				clause.OnConflict{
					Columns: []clause.Column{{Name: "id"}},
					DoUpdates: clause.Assignments(map[string]interface{}{
						"total_topics":   totalTopics,
						"total_comments": totalComments,
						"total_members":  totalMembers,
						"total_visits":   gorm.Expr("forum_stats.total_visits + ?", c.visitCounter),
						"newest_member":  newWestMember.Username,
					}),
				},
				clause.Returning{},
			).Create(&stats).Error
			if err == nil {
				c.statsCache = stats
				c.lastCache = time.Now()
				c.visitCounter = 0
			}
		default:
		}
	}

	c.visitCounter++
	c.statsCache.TotalVisits++
	return c.statsCache, nil
}

func (c *statsService) GetTopContributors() []model.User {
	return cache.UserCache.GetScoreRank()
}

func newStatsService() *statsService {
	return &statsService{
		lock: make(chan struct{}, 1),
	}
}
