package service

import (
	"bbs-go/internal/event"
	"bbs-go/internal/locale"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/search"
	"errors"
	"math"
	"net/http"

	"bbs-go/common/dates"
	"bbs-go/common/strs"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

	"gorm.io/gorm"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"
)

var TopicService = newTopicService()

func newTopicService() *topicService {
	return &topicService{}
}

type topicService struct{}

func (s *topicService) Get(id int64) *model.Topic {
	return repository.TopicRepository.Get(sqls.DB(), id)
}

func (s *topicService) Take(where ...interface{}) *model.Topic {
	return repository.TopicRepository.Take(sqls.DB(), where...)
}

func (s *topicService) Find(cnd *sqls.Cnd) []model.Topic {
	return repository.TopicRepository.Find(sqls.DB(), cnd)
}

func (s *topicService) FindOne(cnd *sqls.Cnd) *model.Topic {
	return repository.TopicRepository.FindOne(sqls.DB(), cnd)
}

func (s *topicService) FindPageByParams(params *params.QueryParams) (list []model.Topic, paging *sqls.Paging) {
	return repository.TopicRepository.FindPageByParams(sqls.DB(), params)
}

func (s *topicService) FindPageByCnd(cnd *sqls.Cnd) (list []model.Topic, paging *sqls.Paging) {
	return repository.TopicRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *topicService) Count(cnd *sqls.Cnd) int64 {
	return repository.TopicRepository.Count(sqls.DB(), cnd)
}

func (s *topicService) Updates(id int64, columns map[string]interface{}) error {
	if err := repository.TopicRepository.Updates(sqls.DB(), id, columns); err != nil {
		return err
	}

	// Update search index
	search.UpdateTopicIndex(s.Get(id))

	return nil
}

func (s *topicService) UpdateColumn(id int64, name string, value interface{}) error {
	if err := repository.TopicRepository.UpdateColumn(sqls.DB(), id, name, value); err != nil {
		return err
	}

	// Update search index
	search.UpdateTopicIndex(s.Get(id))

	return nil
}

// Delete Remove
func (s *topicService) Delete(topicId, deleteUserId int64, r *http.Request) error {
	topic := s.Get(topicId)
	if topic == nil {
		return nil
	}
	err := repository.TopicRepository.UpdateColumn(sqls.DB(), topicId, "status", constants.StatusDeleted)
	if err == nil {
		// Remove from search index
		search.DeleteTopicIndex(topicId)
		// Remove topic tags
		TopicTagService.DeleteByTopicId(topicId)
		// send event
		event.Send(event.TopicDeleteEvent{
			UserId:       topic.UserID,
			TopicId:      topic.ID,
			DeleteUserId: deleteUserId,
		})
	}
	return err
}

// Undelete Restore
func (s *topicService) Undelete(id int64) error {
	err := repository.TopicRepository.UpdateColumn(sqls.DB(), id, "status", constants.StatusActive)
	if err == nil {
		// Restore topic tags
		TopicTagService.UndeleteByTopicId(id)
		// Update search index
		search.UpdateTopicIndex(s.Get(id))
	}
	return err
}

// Update
func (s *topicService) Edit(topicId, forumId int64, tags []string, title, slug, content, hideContent string) error {
	if title == "" {
		return errors.New(locale.T("topic.edit.title_required"))
	}

	if strs.RuneLen(title) > 128 {
		return errors.New(locale.T("topic.edit.title_max_length_exceeded"))
	}

	node := repository.ForumRepository.Get(sqls.DB(), forumId)
	if node == nil || node.Status != constants.StatusActive {
		return errors.New(locale.T("forum.not_found"))
	}

	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		var err error
		if err = repository.TopicRepository.Updates(sqls.DB(), topicId, map[string]interface{}{
			"forum_id":     forumId,
			"title":        title,
			"slug":         slug,
			"content":      content,
			"hide_content": hideContent,
		}); err != nil {
			return err
		}

		repository.TopicTagRepository.DeleteTopicTags(tx, topicId)    // remove all tags first
		repository.TopicTagRepository.AddTopicTags(tx, topicId, tags) // then re-add tags
		return nil
	})

	// Add index
	search.UpdateTopicIndex(s.Get(topicId))

	return err
}

// GetTopicTags Topic tags
func (s *topicService) GetTopicTags(topicId int64) []string {
	topicTags := repository.TopicTagRepository.Find(sqls.DB(), sqls.NewCnd().Where("topic_id = ?", topicId))
	tagNames := make([]string, len(topicTags))
	for i, topicTag := range topicTags {
		tagNames[i] = topicTag.Tag
	}
	return tagNames
}

// GetTopics Topic lists (newest, recommended, followed, node)
// func (s *topicService) GetTopics(user *model.User, forumId, cursor int64) (topics []model.Topic, nextCursor int64, hasMore bool) {
// 	if forumId == constants.NodeIdFollow {
// 		if user != nil {
// 			return s.GetFollowedTopics(user.Id, cursor)
// 		}
// 		return
// 	} else {
// 		return s.GetForumTopics(forumId, cursor)
// 	}
// }

// GetForumTopics Forum topics (newest, recommended, node)
func (s *topicService) GetForumTopics(forumId, cursor int64) (topics []model.Topic, nextCursor int64, hasMore bool) {
	const limit = 20
	cnd := sqls.NewCnd().Eq("forum_id", forumId)
	if cursor > 0 {
		cnd.Lt("last_comment_time", cursor)
	}
	cnd.Eq("status", constants.StatusActive).Desc("last_comment_time").Limit(limit)
	topics = repository.TopicRepository.Find(sqls.DB(), cnd)
	if len(topics) > 0 {
		nextCursor = topics[len(topics)-1].LastCommentTime
		hasMore = len(topics) >= limit
	} else {
		nextCursor = cursor
	}
	return
}

func (s *topicService) GetNewestTopics(cursor int64) (topics []model.Topic, nextCursor int64, hasMore bool) {
	const limit = 20
	cnd := sqls.NewCnd()
	if cursor > 0 {
		cnd.Lt("last_comment_time", cursor)
	}
	cnd.Eq("status", constants.StatusActive).Desc("last_comment_time").Limit(limit)
	topics = repository.TopicRepository.Find(sqls.DB(), cnd)
	if len(topics) > 0 {
		nextCursor = topics[len(topics)-1].LastCommentTime
		hasMore = len(topics) >= limit
	} else {
		nextCursor = cursor
	}
	return
}

func (s *topicService) GetRecommendedTopics(cursor int64) (topics []model.Topic, nextCursor int64, hasMore bool) {
	const limit = 20
	cnd := sqls.NewCnd().Eq("recommended", true)
	if cursor > 0 {
		cnd.Lt("last_comment_time", cursor)
	}
	cnd.Eq("status", constants.StatusActive).Desc("last_comment_time").Limit(limit)
	topics = repository.TopicRepository.Find(sqls.DB(), cnd)
	if len(topics) > 0 {
		nextCursor = topics[len(topics)-1].LastCommentTime
		hasMore = len(topics) >= limit
	} else {
		nextCursor = cursor
	}
	return
}

// _GetFollowTopics Followed topics list
func (s *topicService) GetFollowedTopics(userId int64, cursor int64) (topics []model.Topic, nextCursor int64, hasMore bool) {
	const limit = 20
	cnd := sqls.NewCnd().Eq("user_id", userId)
	cnd.Eq("data_type", constants.EntityTopic)
	if cursor > 0 {
		cnd.Lt("create_time", cursor)
	}
	cnd.Desc("create_time").Limit(limit)

	userFeeds := repository.UserFeedRepository.Find(sqls.DB(), cnd)
	if len(userFeeds) > 0 {
		nextCursor = userFeeds[len(userFeeds)-1].CreateTime
		hasMore = len(userFeeds) >= limit
	} else {
		nextCursor = cursor
	}

	var topicIds []int64
	for _, item := range userFeeds {
		topicIds = append(topicIds, item.DataID)
	}
	topics = TopicService.GetTopicByIds(topicIds)
	return
}

// Get topics under a specific tag
func (s *topicService) GetTopicsByTag(tag string, cursor int64) (topics []model.Topic, nextCursor int64, hasMore bool) {
	limit := 20
	tagTopics := repository.TopicTagRepository.Find(sqls.DB(), sqls.NewCnd().
		Eq("tag", tag).
		Desc("last_comment_time").Limit(limit))
	if len(tagTopics) > 0 {
		nextCursor = tagTopics[len(tagTopics)-1].LastCommentTime

		var topicIds []int64
		for _, topicTag := range tagTopics {
			topicIds = append(topicIds, topicTag.TopicId)
		}

		topicsMap := s.GetTopicInIds(topicIds)
		if topicsMap != nil {
			for _, topicTag := range tagTopics {
				if topic, found := topicsMap[topicTag.TopicId]; found {
					topics = append(topics, topic)
				}
			}
		}
	} else {
		nextCursor = cursor
	}
	hasMore = len(tagTopics) >= limit
	return
}

func (s *topicService) GetTopicByIds(topicIds []int64) (topics []model.Topic) {
	topicsMap := s.GetTopicInIds(topicIds)
	for _, topicId := range topicIds {
		topic, found := topicsMap[topicId]
		if found {
			topics = append(topics, topic)
		}
	}
	return
}

// GetTopicInIds Get topics by IDs
func (s *topicService) GetTopicInIds(topicIds []int64) map[int64]model.Topic {
	if len(topicIds) == 0 {
		return nil
	}
	cnd := sqls.NewCnd().In("id", topicIds).Eq("status", constants.StatusActive)
	topics := repository.TopicRepository.Find(sqls.DB(), cnd)
	topicsMap := make(map[int64]model.Topic, len(topics))
	for _, topic := range topics {
		topicsMap[topic.ID] = topic
	}
	return topicsMap
}

// Increment view count
func (s *topicService) IncrViewCount(topicId int64) {
	sqls.DB().Exec("update t_topic set view_count = view_count + 1 where id = ?", topicId)
}

// When a topic is commented, update last reply time and increment reply count
func (s *topicService) onComment(tx *gorm.DB, comment *model.Comment) error {
	if err := repository.TopicRepository.Updates(tx, comment.TopicID, map[string]interface{}{
		"last_comment_time":    comment.CreateTime,
		"last_comment_user_id": comment.UserID,
		"comment_count":        gorm.Expr("comment_count + 1"),
	}); err != nil {
		return err
	}
	return tx.Model(&model.TopicTag{}).Where("topic_id = ?", comment.TopicID).
		UpdateColumns(map[string]interface{}{
			"last_comment_time":    comment.CreateTime,
			"last_comment_user_id": comment.UserID,
		}).Error

	// repository.TopicTagRepository.up
	// .UpdateLastCommentInfo(tx, comment.TopicID, comment.CreateTime, comment.UserID)

	// if err := tx.Exec("update t_topic_tag set last_comment_time = ?, last_comment_user_id = ? where topic_id = ?",
	// 	comment.CreateTime, comment.UserID, comment.TopicID).Error; err != nil {
	// 	return err
	// }
	// return nil
}

func (s *topicService) ScanByUser(userId int64, callback func(topics []model.Topic)) {
	var cursor int64 = 0
	for {
		list := repository.TopicRepository.Find(sqls.DB(), sqls.NewCnd().
			Eq("user_id", userId).Gt("id", cursor).Asc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		callback(list)
	}
}

func (s *topicService) Scan(callback func(topics []model.Topic)) {
	var cursor int64 = 0
	for {
		list := repository.TopicRepository.Find(sqls.DB(), sqls.NewCnd().
			Gt("id", cursor).Asc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		callback(list)
	}
}

// Scan in descending order
func (s *topicService) ScanDesc(callback func(topics []model.Topic)) {
	var cursor int64 = math.MaxInt64
	for {
		list := repository.TopicRepository.Find(sqls.DB(), sqls.NewCnd().
			Lt("id", cursor).Desc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		callback(list)
	}
}

// Scan in descending order with date range
func (s *topicService) ScanDescWithDate(dateFrom, dateTo int64, callback func(topics []model.Topic)) {
	var cursor int64 = math.MaxInt64
	for {
		list := repository.TopicRepository.Find(sqls.DB(), sqls.NewCnd().
			Cols("id", "status", "create_time", "update_time").
			Lt("id", cursor).Gte("create_time", dateFrom).Lt("create_time", dateTo).Desc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		callback(list)
	}
}

func (s *topicService) GetUserTopics(userId, cursor int64) (topics []model.Topic, nextCursor int64, hasMore bool) {
	limit := 20
	cnd := sqls.NewCnd()
	if userId > 0 {
		cnd.Eq("user_id", userId)
	}
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	cnd.Eq("status", constants.StatusActive).Desc("id").Limit(limit)
	topics = repository.TopicRepository.Find(sqls.DB(), cnd)
	if len(topics) > 0 {
		nextCursor = topics[len(topics)-1].ID
		hasMore = len(topics) >= limit
	} else {
		nextCursor = cursor
	}
	return
}

func (s *topicService) GetPinnedTopics(forumId int64, limit int) []model.Topic {
	if forumId > 0 {
		return s.Find(sqls.NewCnd().Where("forum_id = ? and pinned = true and status = ?",
			forumId, constants.StatusActive).Desc("pinned_time").Limit(limit))
	} else {
		return s.Find(sqls.NewCnd().Where("pinned = true and status = ?",
			constants.StatusActive).Desc("pinned_time").Limit(limit))
	}
}

func (s *topicService) SetTopicPinned(topicId int64, pinned bool) error {
	if pinned {
		return s.Updates(topicId, map[string]interface{}{
			"pinned":      true,
			"pinned_time": dates.NowTimestamp(),
		})
	} else {
		return s.Updates(topicId, map[string]interface{}{
			"pinned": false,
		})
	}
}

// Recommend
func (s *topicService) SetTopicRecommended(topicId int64, recommended bool) error {
	topic := s.Get(topicId)
	if topic == nil || topic.Status != constants.StatusActive {
		return errors.New(locale.T("topic.not_found"))
	}
	if topic.Recommended == recommended { // Recommended status not changed
		return nil
	}
	if recommended {
		if err := s.Updates(topicId, map[string]interface{}{
			"recommended":      recommended,
			"recommended_time": dates.NowTimestamp(),
		}); err != nil {
			return err
		}
	} else {
		if err := s.UpdateColumn(topicId, "recommended", recommended); err != nil {
			return err
		}
	}

	// Send event
	event.Send(event.TopicRecommendedEvent{
		Topic:       topic,
		Recommended: recommended,
	})

	// Add index
	search.UpdateTopicIndex(s.Get(topicId))

	return nil
}
