package api

import (
	"bbs-go/internal/controller/response"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/spam"
	"math/rand"
	"strconv"
	"strings"

	"bbs-go/common/strs"
	"bbs-go/sqls"
	"bbs-go/web"
	"bbs-go/web/params"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/cache"
	"bbs-go/internal/model"
	"bbs-go/internal/service"
)

type TopicsController struct {
	Ctx iris.Context
}

// 节点
func (c *TopicsController) GetNodes() *web.JsonResult {
	nodes := response.BuildForumList(service.ForumService.GetAll())
	return web.JsonData(nodes)
}

// 节点信息
func (c *TopicsController) GetNode() *web.JsonResult {
	nodeId := params.FormValueInt64Default(c.Ctx, "nodeId", 0)
	node := service.ForumService.Get(nodeId)
	return web.JsonData(response.BuildForum(node))
}

// 发表帖子
func (c *TopicsController) PostCreate() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}
	form := model.GetCreateTopicForm(c.Ctx)

	if err := spam.CheckTopic(user, form); err != nil {
		return web.JsonError(err)
	}

	topic, err := service.TopicPublishService.Publish(user.Id, form)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(response.BuildSimpleTopic(topic))
}

// 编辑时获取详情
func (c *TopicsController) GetEditBy(topicId int64) *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	topic := service.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOK {
		return web.JsonErrorMsg("话题不存在或已被删除")
	}
	if topic.Type != constants.TopicTypeTopic {
		return web.JsonErrorMsg("当前类型帖子不支持修改")
	}

	// 非作者、且非管理员
	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return web.JsonErrorMsg("无权限")
	}

	tags := service.TopicService.GetTopicTags(topicId)
	var tagNames []string
	if len(tags) > 0 {
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
	}

	return web.NewEmptyRspBuilder().
		Put("id", topic.Id).
		Put("nodeId", topic.NodeId).
		Put("title", topic.Title).
		Put("content", topic.Content).
		Put("hideContent", topic.HideContent).
		Put("tags", tagNames).
		JsonResult()
}

// 编辑帖子
func (c *TopicsController) PostEditBy(topicId int64) *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	topic := service.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOK {
		return web.JsonErrorMsg("话题不存在或已被删除")
	}

	// 非作者、且非管理员
	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return web.JsonErrorMsg("无权限")
	}

	var (
		nodeId      = params.FormValueInt64Default(c.Ctx, "nodeId", 0)
		title       = strings.TrimSpace(params.FormValue(c.Ctx, "title"))
		content     = strings.TrimSpace(params.FormValue(c.Ctx, "content"))
		hideContent = strings.TrimSpace(params.FormValue(c.Ctx, "hideContent"))
		tags        = params.FormValueStringArray(c.Ctx, "tags")
	)

	err := service.TopicService.Edit(topicId, nodeId, tags, title, content, hideContent)
	if err != nil {
		return web.JsonError(err)
	}
	// 操作日志
	service.OperateLogService.AddOperateLog(user.Id, constants.OpTypeUpdate, constants.EntityTopic, topicId,
		"", c.Ctx.Request())
	return web.JsonData(response.BuildSimpleTopic(topic))
}

// 删除帖子
func (c *TopicsController) PostDeleteBy(topicId int64) *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	topic := service.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOK {
		return web.JsonSuccess()
	}

	// 非作者、且非管理员
	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return web.JsonErrorMsg("无权限")
	}

	if err := service.TopicService.Delete(topicId, user.Id, c.Ctx.Request()); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// PostRecommendBy 设为推荐
func (c *TopicsController) PostRecommendBy(topicId int64) *web.JsonResult {
	recommend, err := params.FormValueBool(c.Ctx, "recommend")
	if err != nil {
		return web.JsonError(err)
	}
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	if !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return web.JsonErrorMsg("无权限")
	}

	err = service.TopicService.SetRecommend(topicId, recommend)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 帖子详情
func (c *TopicsController) GetBy(topicId int64) *web.JsonResult {
	topic := service.TopicService.Get(topicId)
	if topic == nil || topic.Status == constants.StatusDeleted {
		return web.JsonErrorMsg("帖子不存在")
	}

	// 审核中文章控制展示
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if topic.Status == constants.StatusReview {
		if user != nil {
			if topic.UserId != user.Id && !user.IsOwnerOrAdmin() {
				return web.JsonErrorCode(403, "文章审核中")
			}
		} else {
			return web.JsonErrorCode(403, "文章审核中")
		}
	}

	service.TopicService.IncrViewCount(topicId) // 增加浏览量
	return web.JsonData(response.BuildTopic(topic, user))
}

// 点赞用户
func (c *TopicsController) GetRecentlikesBy(topicId int64) *web.JsonResult {
	likes := service.UserLikeService.Recent(constants.EntityTopic, topicId, 5)
	var users []response.UserInfo
	for _, like := range likes {
		userInfo := response.BuildUserInfoDefaultIfNull(like.UserId)
		if userInfo != nil {
			users = append(users, *userInfo)
		}
	}
	return web.JsonData(users)
}

// 最新帖子
func (c *TopicsController) GetRecent() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	topics := service.TopicService.Find(sqls.NewCnd().Where("status = ?", constants.StatusOK).Desc("id").Limit(10))
	return web.JsonData(response.BuildSimpleTopics(topics, user))
}

// 用户帖子列表
func (c *TopicsController) GetUserTopics() *web.JsonResult {
	userId, err := params.FormValueInt64(c.Ctx, "userId")
	if err != nil {
		return web.JsonError(err)
	}
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	user := service.UserTokenService.GetCurrent(c.Ctx)
	topics, cursor, hasMore := service.TopicService.GetUserTopics(userId, cursor)
	return web.JsonCursorData(response.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10), hasMore)
}

// 标签帖子列表
func (c *TopicsController) GetTagTopics() *web.JsonResult {
	var (
		cursor     = params.FormValueInt64Default(c.Ctx, "cursor", 0)
		tagId, err = params.FormValueInt64(c.Ctx, "tagId")
		user       = service.UserTokenService.GetCurrent(c.Ctx)
	)
	if err != nil {
		return web.JsonError(err)
	}
	topics, cursor, hasMore := service.TopicService.GetTagTopics(tagId, cursor)
	return web.JsonCursorData(response.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10), hasMore)
}

// 收藏
func (c *TopicsController) GetFavoriteBy(topicId int64) *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	err := service.FavoriteService.AddTopicFavorite(user.Id, topicId)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 推荐话题列表（目前逻辑为取最近50条数据随机展示）
func (c *TopicsController) GetRecommend() *web.JsonResult {
	topics := cache.TopicCache.GetRecommendTopics()
	if len(topics) == 0 {
		return web.JsonSuccess()
	} else {
		dest := make([]model.Topic, len(topics))
		perm := rand.Perm(len(topics))
		for i, v := range perm {
			dest[v] = topics[i]
		}
		end := 10
		if end > len(topics) {
			end = len(topics)
		}
		ret := dest[0:end]
		return web.JsonData(response.BuildSimpleTopics(ret, nil))
	}
}

// 最新话题
func (c *TopicsController) GetNewest() *web.JsonResult {
	topics := service.TopicService.Find(sqls.NewCnd().Eq("status", constants.StatusOK).Desc("id").Limit(6))
	return web.JsonData(response.BuildSimpleTopics(topics, nil))
}

// 设置置顶
func (c *TopicsController) PostBy(topicId int64) *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	if !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return web.JsonErrorMsg("无权限")
	}

	var (
		pinned = params.FormValueBoolDefault(c.Ctx, "pinned", false) // 是否指定
	)
	if err := service.TopicService.PinTopic(topicId, pinned); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *TopicsController) GetHide_content() *web.JsonResult {
	topicId := params.FormValueInt64Default(c.Ctx, "topicId", 0)
	var (
		exists      = false // 是否有隐藏内容
		show        = false // 是否显示隐藏内容
		hideContent = ""    // 隐藏内容
	)
	topic := service.TopicService.Get(topicId)
	if topic != nil && topic.Status == constants.StatusOK && strs.IsNotBlank(topic.HideContent) {
		exists = true
		if user := service.UserTokenService.GetCurrent(c.Ctx); user != nil {
			if user.Id == topic.UserId || service.CommentService.IsCommented(user.Id, constants.EntityTopic, topic.Id) {
				show = true
				hideContent = markdown.ToHTML(topic.HideContent)
			}
		}
	}
	return web.JsonData(map[string]interface{}{
		"exists":  exists,
		"show":    show,
		"content": hideContent,
	})
}
