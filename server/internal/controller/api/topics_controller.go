package api

import (
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"
	"bbs-go/internal/locale"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/spam"
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

type TopicsController struct {
	Ctx iris.Context
}

func (c *TopicsController) getTopicBySlugId(slugId string) (*model.Topic, error) {
	parts := strings.SplitN(slugId, ".", 2)
	if len(parts) < 2 {
		return nil, nil
	}
	topicSlug := parts[0]
	topicId := base62.Decode(parts[1])
	cnd := sqls.NewCnd().
		Eq("slug", topicSlug).
		Eq("id", topicId).
		Eq("status", constants.StatusActive)
	topic := service.TopicService.FindOne(cnd)
	return topic, nil
}

// POST /topics -> create topic
func (c *TopicsController) Post() *web.JsonResult {
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

// // 标签帖子列表
//
//	func (c *TopicController) GetTagTopics() *web.JsonResult {
//		var (
//			cursor     = params.FormValueInt64Default(c.Ctx, "cursor", 0)
//			tagId, err = params.FormValueInt64(c.Ctx, "tagId")
//			user       = service.UserTokenService.GetCurrent(c.Ctx)
//		)
//		if err != nil {
//			return web.JsonError(err)
//		}
//		topics, cursor, hasMore := service.TopicService.GetTagTopics(tagId, cursor)
//		return web.JsonCursorData(response.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10), hasMore)
//	}

func (c *TopicsController) Get(ctx iris.Context) *web.JsonResult {
	tag := params.FormValueDefault(c.Ctx, "tag", "")
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)

	if tag == "" {
		return web.JsonError(errs.ErrBadRequest)
	}
	topics, cursor, hasMore := service.TopicService.GetTopicsByTag(tag, cursor)
	user := service.UserTokenService.GetCurrent(c.Ctx)
	return web.JsonCursorData(payload.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10), hasMore)
}

// GET /topics/{slug}
func (c *TopicsController) GetBy(slugId string) *web.JsonResult {
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
func (c *TopicsController) GetByEdit(slugId string) *web.JsonResult {
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
func (c *TopicsController) PutBy(slugId string) *web.JsonResult {
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

func (c *TopicsController) setTopicPinned(topicId int64, pinned bool) (*web.JsonResult, error) {
	err := service.TopicService.SetTopicPinned(topicId, pinned)
	if err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

func (c *TopicsController) setTopicRecommended(topicId int64, recommended bool) (*web.JsonResult, error) {
	err := service.TopicService.SetTopicRecommended(topicId, recommended)
	if err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

// PATCH /topics/{slugId} => pin/unpin, recommend/unrecommend
func (c *TopicsController) PatchBy(slugId string) (*web.JsonResult, error) {
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
func (c *TopicsController) DeleteBy(topicId int64) (*web.JsonResult, error) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(errs.ErrUnauthorized), nil
	}

	topic := service.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusActive {
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
func (c *TopicsController) PostByReactions(slugId string) (*web.JsonResult, error) {
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
func (c *TopicsController) DeleteByReactionsBy(slugId string, userId int64) (*web.JsonResult, error) {
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
func (c *TopicsController) GetByComments(slugId string) (*web.JsonResult, error) {
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
func (c *TopicsController) PostByComments(slugId string) (*web.JsonResult, error) {
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
