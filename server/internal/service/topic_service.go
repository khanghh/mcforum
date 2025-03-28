package service

import (
	"bbs-go/internal/locale"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/search"
	"errors"
	"math"
	"net/http"

	"bbs-go/common/dates"
	"bbs-go/common/strs"
	"bbs-go/sqls"
	"bbs-go/web/params"

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

	// 添加索引
	search.UpdateTopicIndex(s.Get(id))

	return nil
}

func (s *topicService) UpdateColumn(id int64, name string, value interface{}) error {
	if err := repository.TopicRepository.UpdateColumn(sqls.DB(), id, name, value); err != nil {
		return err
	}

	// 添加索引
	search.UpdateTopicIndex(s.Get(id))

	return nil
}

// Delete 删除
func (s *topicService) Delete(topicId, deleteUserId int64, r *http.Request) error {
	topic := s.Get(topicId)
	if topic == nil {
		return nil
	}
	err := repository.TopicRepository.UpdateColumn(sqls.DB(), topicId, "status", constants.StatusDeleted)
	if err == nil {
		// 添加索引
		search.DeleteTopicIndex(topicId)
		// 删掉标签文章
		TopicTagService.DeleteByTopicId(topicId)
		// 发送事件
		event.Send(event.TopicDeleteEvent{
			UserId:       topic.UserId,
			TopicId:      topic.Id,
			DeleteUserId: deleteUserId,
		})
	}
	return err
}

// Undelete 取消删除
func (s *topicService) Undelete(id int64) error {
	err := repository.TopicRepository.UpdateColumn(sqls.DB(), id, "status", constants.StatusOK)
	if err == nil {
		// 删掉标签文章
		TopicTagService.UndeleteByTopicId(id)
		// 添加索引
		search.UpdateTopicIndex(s.Get(id))
	}
	return err
}

// 更新
func (s *topicService) Edit(topicId, forumId int64, tags []string, title, slug, content, hideContent string) error {
	if title == "" {
		return errors.New(locale.T("topic.edit.title_required"))
	}

	if strs.RuneLen(title) > 128 {
		return errors.New(locale.T("topic.edit.title_max_length_exceeded"))
	}

	node := repository.ForumRepository.Get(sqls.DB(), forumId)
	if node == nil || node.Status != constants.StatusOK {
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

		repository.TopicTagRepository.DeleteTopicTags(tx, topicId)    // 先删掉所有的标签
		repository.TopicTagRepository.AddTopicTags(tx, topicId, tags) // 然后重新添加标签
		return nil
	})

	// 添加索引
	search.UpdateTopicIndex(s.Get(topicId))

	return err
}

// GetTopicTags 话题的标签
func (s *topicService) GetTopicTags(topicId int64) []string {
	topicTags := repository.TopicTagRepository.Find(sqls.DB(), sqls.NewCnd().Where("topic_id = ?", topicId))
	tagNames := make([]string, len(topicTags))
	for i, topicTag := range topicTags {
		tagNames[i] = topicTag.Tag
	}
	return tagNames
}

// GetTopics 帖子列表（最新、推荐、关注、节点）
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

// GetForumTopics 帖子列表（最新、推荐、节点）
func (s *topicService) GetForumTopics(forumId, cursor int64) (topics []model.Topic, nextCursor int64, hasMore bool) {
	const limit = 20
	cnd := sqls.NewCnd().Eq("forum_id", forumId)
	if cursor > 0 {
		cnd.Lt("last_comment_time", cursor)
	}
	cnd.Eq("status", constants.StatusOK).Desc("last_comment_time").Limit(limit)
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
	cnd.Eq("status", constants.StatusOK).Desc("last_comment_time").Limit(limit)
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
	cnd.Eq("status", constants.StatusOK).Desc("last_comment_time").Limit(limit)
	topics = repository.TopicRepository.Find(sqls.DB(), cnd)
	if len(topics) > 0 {
		nextCursor = topics[len(topics)-1].LastCommentTime
		hasMore = len(topics) >= limit
	} else {
		nextCursor = cursor
	}
	return
}

// _GetFollowTopics 关注帖子列表
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
		topicIds = append(topicIds, item.DataId)
	}
	topics = TopicService.GetTopicByIds(topicIds)
	return
}

// // 指定标签下话题列表
// func (s *topicService) GetTagTopics(tagId, cursor int64) (topics []model.Topic, nextCursor int64, hasMore bool) {
// 	limit := 20
// 	topicTags := repository.TopicTagRepository.Find(sqls.DB(), sqls.NewCnd().
// 		Eq("tag_id", tagId).
// 		Eq("status", constants.StatusOK).
// 		Desc("last_comment_time").Limit(limit))
// 	if len(topicTags) > 0 {
// 		nextCursor = topicTags[len(topicTags)-1].LastCommentTime

// 		var topicIds []int64
// 		for _, topicTag := range topicTags {
// 			topicIds = append(topicIds, topicTag.TopicId)
// 		}

// 		topicsMap := s.GetTopicInIds(topicIds)
// 		if topicsMap != nil {
// 			for _, topicTag := range topicTags {
// 				if topic, found := topicsMap[topicTag.TopicId]; found {
// 					topics = append(topics, topic)
// 				}
// 			}
// 		}
// 	} else {
// 		nextCursor = cursor
// 	}
// 	hasMore = len(topicTags) >= limit
// 	return
// }

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

// GetTopicInIds 根据编号批量获取主题
func (s *topicService) GetTopicInIds(topicIds []int64) map[int64]model.Topic {
	if len(topicIds) == 0 {
		return nil
	}
	var topics []model.Topic
	sqls.DB().Where("id in (?)", topicIds).Find(&topics)

	topicsMap := make(map[int64]model.Topic, len(topics))
	for _, topic := range topics {
		topicsMap[topic.Id] = topic
	}
	return topicsMap
}

// 浏览数+1
func (s *topicService) IncrViewCount(topicId int64) {
	sqls.DB().Exec("update t_topic set view_count = view_count + 1 where id = ?", topicId)
}

// 当帖子被评论的时候，更新最后回复时间、回复数量+1
func (s *topicService) onComment(tx *gorm.DB, comment *model.Comment) error {
	if err := repository.TopicRepository.Updates(tx, comment.TopicId, map[string]interface{}{
		"last_comment_time":    comment.CreateTime,
		"last_comment_user_id": comment.UserId,
		"comment_count":        gorm.Expr("comment_count + 1"),
	}); err != nil {
		return err
	}
	if err := tx.Exec("update t_topic_tag set last_comment_time = ?, last_comment_user_id = ? where topic_id = ?",
		comment.CreateTime, comment.UserId, comment.TopicId).Error; err != nil {
		return err
	}
	return nil
}

func (s *topicService) ScanByUser(userId int64, callback func(topics []model.Topic)) {
	var cursor int64 = 0
	for {
		list := repository.TopicRepository.Find(sqls.DB(), sqls.NewCnd().
			Eq("user_id", userId).Gt("id", cursor).Asc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
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
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

// 倒序扫描
func (s *topicService) ScanDesc(callback func(topics []model.Topic)) {
	var cursor int64 = math.MaxInt64
	for {
		list := repository.TopicRepository.Find(sqls.DB(), sqls.NewCnd().
			Lt("id", cursor).Desc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

// 倒序扫描
func (s *topicService) ScanDescWithDate(dateFrom, dateTo int64, callback func(topics []model.Topic)) {
	var cursor int64 = math.MaxInt64
	for {
		list := repository.TopicRepository.Find(sqls.DB(), sqls.NewCnd().
			Cols("id", "status", "create_time", "update_time").
			Lt("id", cursor).Gte("create_time", dateFrom).Lt("create_time", dateTo).Desc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
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
	cnd.Eq("status", constants.StatusOK).Desc("id").Limit(limit)
	topics = repository.TopicRepository.Find(sqls.DB(), cnd)
	if len(topics) > 0 {
		nextCursor = topics[len(topics)-1].Id
		hasMore = len(topics) >= limit
	} else {
		nextCursor = cursor
	}
	return
}

func (s *topicService) GetPinnedTopics(forumId int64, limit int) []model.Topic {
	if forumId > 0 {
		return s.Find(sqls.NewCnd().Where("forum_id = ? and pinned = true and status = ?",
			forumId, constants.StatusOK).Desc("pinned_time").Limit(limit))
	} else {
		return s.Find(sqls.NewCnd().Where("pinned = true and status = ?",
			constants.StatusOK).Desc("pinned_time").Limit(limit))
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

// 推荐
func (s *topicService) SetTopicRecommended(topicId int64, recommended bool) error {
	topic := s.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOK {
		return errors.New(locale.T("topic.not_found"))
	}
	if topic.Recommended == recommended { // 推荐状态没变更
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

	// 发送事件
	event.Send(event.TopicRecommendedEvent{
		TopicId:     topicId,
		Recommended: recommended,
	})

	// 添加索引
	search.UpdateTopicIndex(s.Get(topicId))

	return nil
}
