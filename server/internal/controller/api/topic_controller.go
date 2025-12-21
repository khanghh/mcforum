package api

import (
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"
	"bbs-go/internal/locale"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/spam"
	"fmt"
	"strconv"
	"strings"

	"bbs-go/common/base62"
	"bbs-go/common/utils"
	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/model"
	"bbs-go/internal/service"
)

type TopicController struct {
	Ctx iris.Context
}

func (c *TopicController) getTopicBySlugId(slugId string) (*model.Topic, error) {
	parts := strings.SplitN(slugId, ".", 2)
	if len(parts) < 2 {
		return nil, nil
	}
	topicSlug := parts[0]
	topicId := base62.Decode(parts[1])
	cnd := sqls.NewCnd().
		Eq("slug", topicSlug).
		Eq("id", topicId).
		Eq("status", constants.StatusOK)
	topic := service.TopicService.FindOne(cnd)
	return topic, nil
}

// 节点
// func (c *TopicController) GetNodes() *web.JsonResult {
// 	nodes := response.BuildForumList(service.ForumService.GetAll())
// 	return web.JsonData(nodes)
// }

// // 节点信息
// func (c *TopicController) GetNode() *web.JsonResult {
// 	nodeId := params.FormValueInt64Default(c.Ctx, "nodeId", 0)
// 	node := service.ForumService.Get(nodeId)
// 	return web.JsonData(response.BuildForum(node))
// }

// POST /topics -> create topic
func (c *TopicController) Post() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	form := payload.GetCreateTopicForm(c.Ctx)
	if err := spam.CheckTopic(user, form); err != nil {
		return web.JsonError(err)
	}

	topic, err := service.TopicService.Publish(service.PublishTopicArgs{
		UserId:      user.Id,
		Title:       form.Title,
		ForumId:     form.ForumId,
		Content:     form.Content,
		HideContent: form.HideContent,
		Tags:        form.Tags,
		Images:      form.Images,
		UserAgent:   utils.GetUserAgent(c.Ctx.Request()),
		IPAddress:   utils.GetRequestIP(c.Ctx.Request()),
	})
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(payload.BuildSimpleTopic(topic))
}

func (c *TopicController) Get(ctx iris.Context) *web.JsonResult {
	type TopicQuery struct {
		Tag      string `form:"tag"`
		Username string `form:"username"`
	}
	params := TopicQuery{}
	if err := ctx.ReadQuery(&params); err != nil {
		return web.JsonError(err)
	}
	fmt.Println(params.Username)
	return &web.JsonResult{
		StatusCode: 404,
		Error:      errs.ErrForbidden,
	}
}

// GET /topics/{slug}
func (c *TopicController) GetBy(slugId string) *web.JsonResult {
	topic, err := c.getTopicBySlugId(slugId)
	if topic == nil || err != nil {
		return web.JsonError(errs.ErrTopicNotFound)
	}

	// 审核中文章控制展示
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if topic.Status == constants.StatusReview {
		if user != nil {
			if topic.UserId != user.Id && !user.IsOwnerOrAdmin() {
				return web.JsonError(errs.ErrForbidden)
			}
		} else {
			return web.JsonError(errs.ErrForbidden)
		}
	}

	service.TopicService.IncrViewCount(topic.Id) // 增加浏览量
	return web.JsonData(payload.BuildTopic(topic, user))
}

// GET /topics/{slug}/edit
func (c *TopicController) GetByEdit(slugId string) *web.JsonResult {
	topic, err := c.getTopicBySlugId(slugId)
	if topic == nil || err != nil {
		return web.JsonError(errs.ErrTopicNotFound)
	}

	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil || (!user.IsOwnerOrAdmin() && topic.UserId != user.Id) {
		return web.JsonError(errs.ErrForbidden)
	}

	if topic.Status == constants.StatusReview && !user.IsOwnerOrAdmin() {
		return web.JsonErrorMsg(locale.T("topic.not_editable"))
	}

	// revision := params.FormValueInt64Default(c.Ctx, "revision", 0)

	tags := service.TopicService.GetTopicTags(topic.Id)
	return web.NewEmptyRspBuilder().
		Put("id", topic.Id).
		Put("forumId", topic.ForumId).
		Put("title", topic.Title).
		Put("content", topic.Content).
		Put("tags", tags).
		JsonResult()
}

// PUT /topics/{slug} // edit topic
func (c *TopicController) PutBy(slugId string) *web.JsonResult {
	topic, err := c.getTopicBySlugId(slugId)
	if topic == nil || err != nil {
		return web.JsonError(errs.ErrTopicNotFound)
	}

	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonErrorCode(iris.StatusForbidden, err)
	}

	if topic.UserId != user.Id && !user.IsOwnerOrAdmin() {
		return web.JsonError(errs.ErrForbidden)
	}

	topic.ForumId = params.FormValueInt64Default(c.Ctx, "forumId", topic.ForumId)
	topic.Title = strings.TrimSpace(params.FormValueDefault(c.Ctx, "title", topic.Title))
	topic.Content = strings.TrimSpace(params.FormValueDefault(c.Ctx, "content", topic.Content))
	tags := params.FormValueStringArray(c.Ctx, "tags")

	err = service.TopicService.Edit(topic.Id, topic.ForumId, tags, topic.Title, topic.Slug, topic.Content, topic.HideContent)
	if err != nil {
		return web.JsonError(err)
	}
	service.OperateLogService.AddOperateLog(user.Id, constants.OpTypeUpdate, constants.EntityTopic, topic.Id, "", c.Ctx.Request())
	return web.JsonData(payload.BuildTopic(topic, user))
}

func (c *TopicController) setTopicPinned(topicId int64, pinned bool) (*web.JsonResult, error) {
	err := service.TopicService.SetTopicPinned(topicId, pinned)
	if err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

func (c *TopicController) setTopicRecommended(topicId int64, recommended bool) (*web.JsonResult, error) {
	err := service.TopicService.SetTopicRecommended(topicId, recommended)
	if err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

// PATCH /topics/{slugId} => pin/unpin, recommend/unrecommend
func (c *TopicController) PatchBy(slugId string) (*web.JsonResult, error) {
	topic, err := c.getTopicBySlugId(slugId)
	if topic == nil || err != nil {
		return web.JsonError(errs.ErrTopicNotFound), nil
	}
	var (
		pinned      = params.FormValueBoolDefault(c.Ctx, "pinned", topic.Pinned)
		recommended = params.FormValueBoolDefault(c.Ctx, "recommended", topic.Recommended)
	)

	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil || !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return web.JsonError(errs.ErrForbidden), nil
	}

	if pinned != topic.Pinned {
		return c.setTopicPinned(topic.Id, pinned)
	} else if recommended != topic.Recommended {
		return c.setTopicRecommended(topic.Id, recommended)
	}
	return nil, errs.ErrBadRequest
}

// DELETE /topics/{slug}
func (c *TopicController) DeleteBy(topicId int64) (*web.JsonResult, error) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(errs.ErrUnauthorized), nil
	}

	topic := service.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOK {
		return web.JsonError(errs.ErrBadRequest), nil
	}

	// 非作者、且非管理员
	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return web.JsonError(errs.ErrForbidden), nil
	}

	if err := service.TopicService.Delete(topicId, user.Id, c.Ctx.Request()); err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

// POST /topics/{slugId}/reactions
func (c *TopicController) PostByReactions(slugId string) (*web.JsonResult, error) {
	topic, err := c.getTopicBySlugId(slugId)
	if err != nil {
		return web.JsonError(errs.ErrTopicNotFound), nil
	}
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized), nil
	}
	err = service.UserLikeService.TopicLike(user.Id, topic.Id)
	if err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

// DELETE /topics/{slugId}/reactions/{userId}
func (c *TopicController) DeleteByReactionsBy(slugId string, userId int64) (*web.JsonResult, error) {
	topic, err := c.getTopicBySlugId(slugId)
	if err != nil {
		return web.JsonError(errs.ErrTopicNotFound), nil
	}
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized), nil
	}
	err = service.UserLikeService.TopicUnLike(user.Id, topic.Id)
	if err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

// GET /topics/{slugId}/comments
func (c *TopicController) GetByComments(slugId string) (*web.JsonResult, error) {
	topic, err := c.getTopicBySlugId(slugId)
	if err != nil {
		return web.JsonError(errs.ErrTopicNotFound), nil
	}
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	comments, cursor, hasMore := service.CommentService.GetComments(topic.Id, cursor, 20)
	resp := payload.BuildComments(comments, currentUser, true, false)
	return web.JsonCursorData(resp, strconv.FormatInt(cursor, 10), hasMore), nil
}

// POST /topics/{slugId}/comments
func (c *TopicController) PostByComments(slugId string) (*web.JsonResult, error) {
	topic, err := c.getTopicBySlugId(slugId)
	if err != nil {
		return web.JsonError(errs.ErrTopicNotFound), nil
	}

	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err), nil
	}

	form := payload.GetCreateCommentForm(c.Ctx)
	if err := spam.CheckComment(user, form); err != nil {
		return web.JsonError(err), nil
	}

	comment, err := service.CommentService.CreateComment(service.CreateCommentArgs{
		UserId:    user.Id,
		TopicId:   topic.Id,
		Content:   form.Content,
		Images:    form.Images,
		UserAgent: utils.GetUserAgent(c.Ctx.Request()),
		IPAddress: utils.GetRequestIP(c.Ctx.Request()),
	})
	if err != nil {
		return web.JsonError(err), nil
	}

	return web.JsonData(payload.BuildComment(comment)), nil
}

// POST 	/topics/{slugId}/comments/{id}/reactions
// DELETE 	/topics/{slugId}/comments/{id}/reactions
// GET 		/topics/{slugId}/comments/{id}/replies
// func (c *CommentController) PostByReactions(commentId int64) (*web.JsonResult, error) {
// 	user := service.UserTokenService.GetCurrent(c.Ctx)
// 	if user == nil {
// 		return nil, errs.ErrForbidden
// 	}
// 	if err := service.UserLikeService.CommentLike(user.Id, commentId); err != nil {
// 		return nil, err
// 	}
// 	return web.JsonSuccess(), nil
// }

// POST /topics/{slugId}/unpin
// func (c *TopicController) PostByUnpin(slugId string) (*web.JsonResult, int) {
// 	_, topicId, err := c.parseSlugId(slugId)
// 	if err != nil {
// 		return web.JsonErrorMsg("not found"), iris.StatusNotFound
// 	}
// 	user := service.UserTokenService.GetCurrent(c.Ctx)
// 	if user == nil || !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
// 		return web.JsonErrorMsg("no permission"), iris.StatusUnauthorized
// 	}

// 	if err := service.TopicService.SetTopicPinned(int64(topicId), false); err != nil {
// 		return web.JsonError(err), iris.StatusInternalServerError
// 	}
// 	return web.JsonSuccess(), iris.StatusOK
// }

// POST /topics/{slugId}/recommend
// func (c *TopicController) PostByRecommend(slugId string) (*web.JsonResult, int) {
// 	_, topicId, err := c.parseSlugId(slugId)
// 	if err != nil {
// 		return web.JsonErrorMsg("not found"), iris.StatusNotFound
// 	}
// 	user := service.UserTokenService.GetCurrent(c.Ctx)
// 	if user == nil || !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
// 		return web.JsonErrorMsg("no permission"), iris.StatusUnauthorized
// 	}

// 	if err := service.TopicService.SetRecommended(int64(topicId), true); err != nil {
// 		return web.JsonError(err), iris.StatusInternalServerError
// 	}
// 	return web.JsonSuccess(), iris.StatusOK
// }

// POST /topics/{slugId}/unrecommend
// func (c *TopicController) PostByUnrecommend(slugId string) (*web.JsonResult, int) {
// 	_, topicId, err := c.parseSlugId(slugId)
// 	if err != nil {
// 		return web.JsonErrorMsg("not found"), iris.StatusNotFound
// 	}
// 	user := service.UserTokenService.GetCurrent(c.Ctx)
// 	if user == nil || !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
// 		return web.JsonErrorMsg("no permission"), iris.StatusUnauthorized
// 	}

// 	if err := service.TopicService.SetRecommended(int64(topicId), false); err != nil {
// 		return web.JsonError(err), iris.StatusInternalServerError
// 	}
// 	return web.JsonSuccess(), iris.StatusOK
// }

// 点赞用户
// func (c *TopicController) GetRecentlikesBy(topicId int64) *web.JsonResult {
// 	likes := service.UserLikeService.Recent(constants.EntityTopic, topicId, 5)
// 	var users []response.UserInfo
// 	for _, like := range likes {
// 		userInfo := response.BuildUserInfoDefaultIfNull(like.UserId)
// 		if userInfo != nil {
// 			users = append(users, *userInfo)
// 		}
// 	}
// 	return web.JsonData(users)
// }

// // 最新帖子
// func (c *TopicController) GetRecent() *web.JsonResult {
// 	user := service.UserTokenService.GetCurrent(c.Ctx)
// 	topics := service.TopicService.Find(sqls.NewCnd().Where("status = ?", constants.StatusOK).Desc("id").Limit(10))
// 	return web.JsonData(response.BuildSimpleTopics(topics, user))
// }

// 用户帖子列表
// func (c *TopicController) GetUserTopics() *web.JsonResult {
// 	userId, err := params.FormValueInt64(c.Ctx, "userId")
// 	if err != nil {
// 		return web.JsonError(err)
// 	}
// 	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
// 	user := service.UserTokenService.GetCurrent(c.Ctx)
// 	topics, cursor, hasMore := service.TopicService.GetUserTopics(userId, cursor)
// 	return web.JsonCursorData(response.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10), hasMore)
// }

// // 标签帖子列表
// func (c *TopicController) GetTagTopics() *web.JsonResult {
// 	var (
// 		cursor     = params.FormValueInt64Default(c.Ctx, "cursor", 0)
// 		tagId, err = params.FormValueInt64(c.Ctx, "tagId")
// 		user       = service.UserTokenService.GetCurrent(c.Ctx)
// 	)
// 	if err != nil {
// 		return web.JsonError(err)
// 	}
// 	topics, cursor, hasMore := service.TopicService.GetTagTopics(tagId, cursor)
// 	return web.JsonCursorData(response.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10), hasMore)
// }

// 收藏
// func (c *TopicController) GetFavoriteBy(topicId int64) *web.JsonResult {
// 	user := service.UserTokenService.GetCurrent(c.Ctx)
// 	if user == nil {
// 		return web.JsonError(errs.NotLogin)
// 	}
// 	err := service.FavoriteService.AddTopicFavorite(user.Id, topicId)
// 	if err != nil {
// 		return web.JsonError(err)
// 	}
// 	return web.JsonSuccess()
// }

// // 推荐话题列表（目前逻辑为取最近50条数据随机展示）
// func (c *TopicController) GetRecommend() *web.JsonResult {
// 	topics := cache.TopicCache.GetRecommendTopics()
// 	if len(topics) == 0 {
// 		return web.JsonSuccess()
// 	} else {
// 		dest := make([]model.Topic, len(topics))
// 		perm := rand.Perm(len(topics))
// 		for i, v := range perm {
// 			dest[v] = topics[i]
// 		}
// 		end := 10
// 		if end > len(topics) {
// 			end = len(topics)
// 		}
// 		ret := dest[0:end]
// 		return web.JsonData(response.BuildSimpleTopics(ret, nil))
// 	}
// }

// // 最新话题
// func (c *TopicController) GetNewest() *web.JsonResult {
// 	topics := service.TopicService.Find(sqls.NewCnd().Eq("status", constants.StatusOK).Desc("id").Limit(6))
// 	return web.JsonData(response.BuildSimpleTopics(topics, nil))
// }

// func (c *TopicController) GetHide_content() *web.JsonResult {
// 	topicId := params.FormValueInt64Default(c.Ctx, "topicId", 0)
// 	var (
// 		exists      = false // 是否有隐藏内容
// 		show        = false // 是否显示隐藏内容
// 		hideContent = ""    // 隐藏内容
// 	)
// 	topic := service.TopicService.Get(topicId)
// 	if topic != nil && topic.Status == constants.StatusOK && strs.IsNotBlank(topic.HideContent) {
// 		exists = true
// 		if user := service.UserTokenService.GetCurrent(c.Ctx); user != nil {
// 			if user.Id == topic.UserId || service.CommentService.IsCommented(user.Id, constants.EntityTopic, topic.Id) {
// 				show = true
// 				hideContent = markdown.ToHTML(topic.HideContent)
// 			}
// 		}
// 	}
// 	return web.JsonData(map[string]interface{}{
// 		"exists":  exists,
// 		"show":    show,
// 		"content": hideContent,
// 	})
// }
