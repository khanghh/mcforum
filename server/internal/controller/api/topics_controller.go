package api

import (
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"
	"bbs-go/internal/locale"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/spam"
	"log/slog"
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

func canViewTopic(user *model.User, topic *model.Topic) error {
	switch topic.Status {
	case constants.StatusDeleted:
		if user == nil || !user.IsManagerRole() {
			return errs.ErrTopicNotFound
		}
	case constants.StatusReview:
		if user == nil || (topic.UserID != user.ID && !user.IsManagerRole()) {
			return errs.ErrForbidden
		}
	}
	return nil
}

func (c *TopicsController) getTopicBySlugId(slugId string) (*model.Topic, error) {
	parts := strings.SplitN(slugId, ".", 2)
	if len(parts) < 2 {
		return nil, nil
	}
	topicSlug := parts[0]
	topicId := base62.Decode(parts[1])
	cnd := sqls.NewCnd().Eq("slug", topicSlug).Eq("id", topicId)
	topic := service.TopicService.FindOne(cnd)
	return topic, nil
}

func (c *TopicsController) getTopicById(base62Id string) (*model.Topic, error) {
	topicId := base62.Decode(base62Id)
	cnd := sqls.NewCnd().Eq("id", topicId)
	return service.TopicService.FindOne(cnd), nil
}

// POST /topics -> create topic
func (c *TopicsController) Post() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}
	pendingCount, err := service.TopicService.GetPendingTopicCount(user.ID)
	if err != nil {
		slog.Error("check pending topic failed:", "error", err)
		return web.JsonError(errs.ErrInternalServer)
	}
	if pendingCount > constants.MaxPendingReviewTopics {
		return web.JsonErrorMsg(locale.T("topic.pending_review_exceeded"))
	}

	form := payload.GetCreateTopicForm(c.Ctx)
	if err := spam.CheckTopicForm(user, form); err != nil {
		return web.JsonError(err)
	}
	if form.ForumID == 0 {
		return web.JsonErrorMsg(locale.T("topic.forum_required"))
	}

	forum := service.ForumService.Get(form.ForumID)
	if forum == nil {
		return web.JsonError(errs.ErrForumNotFound)
	}

	if user.Role.Rank < forum.WriteRank {
		return web.JsonError(errs.ForbiddenError)
	}

	var needReview bool
	if !user.IsManagerRole() {
		if user.Role.Rank < forum.SkipReviewRank {
			needReview = true
		} else {
			needReview = service.TopicService.CheckForbiddenWord(form.Title, form.Content, form.HiddenContent)
		}
	}

	topic, err := service.TopicService.Publish(service.PublishTopicArgs{
		UserID:        user.ID,
		Title:         form.Title,
		ForumID:       form.ForumID,
		Content:       form.Content,
		HiddenContent: form.HiddenContent,
		Tags:          form.Tags,
		Images:        form.Images,
		NeedReview:    needReview,
		UserAgent:     utils.GetUserAgent(c.Ctx.Request()),
		IPAddress:     utils.GetRequestIP(c.Ctx.Request()),
	})
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(payload.BuildSimpleTopic(topic))
}

// GET /topics?tag={tag}&cursor={cursor}
func (c *TopicsController) Get(ctx iris.Context) *web.JsonResult {
	tag := params.FormValueDefault(c.Ctx, "tag", "")
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)

	if tag == "" {
		return web.JsonError(errs.ErrBadRequest)
	}
	topics, cursor, hasMore := service.TopicService.GetTopicsByTag(tag, cursor)
	user := service.UserTokenService.GetCurrent(c.Ctx)
	return web.JsonCursorData(payload.BuildSimpleTopics(topics, user), cursor, hasMore)
}

// GET /topics/{slug}
func (c *TopicsController) GetBy(slugId string) *web.JsonResult {
	topic, err := c.getTopicBySlugId(slugId)
	if topic == nil || err != nil {
		return web.JsonError(errs.ErrTopicNotFound)
	}

	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := canViewTopic(user, topic); err != nil {
		return web.JsonError(err)
	}

	service.TopicService.IncrViewCount(topic.ID) // 增加浏览量
	return web.JsonData(payload.BuildTopic(topic, user))
}

// GET /topics/edit/{id}
func (c *TopicsController) GetEditBy(base62Id string) *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}

	topic, err := c.getTopicById(base62Id)
	if topic == nil || err != nil {
		return web.JsonError(errs.ErrTopicNotFound)
	}

	if err := canViewTopic(user, topic); err != nil {
		return web.JsonError(err)
	}

	return web.JsonData(payload.BuildTopicEdit(topic))
}

// PUT /topics/edit/{id} // edit topic
func (c *TopicsController) PutEditBy(base62Id string) *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	topic, err := c.getTopicById(base62Id)
	if topic == nil || err != nil {
		return web.JsonError(errs.ErrTopicNotFound)
	}

	if err := canViewTopic(user, topic); err != nil {
		return web.JsonError(err)
	}

	form := payload.GetCreateTopicForm(c.Ctx)
	if err := spam.CheckTopicForm(user, form); err != nil {
		return web.JsonError(err)
	}
	err = service.TopicService.Edit(
		topic.ID,
		form.ForumID,
		form.Tags,
		form.Title,
		form.Content,
		form.HiddenContent,
		form.Images,
	)
	if err != nil {
		return web.JsonError(err)
	}
	service.OperateLogService.AddOperateLog(user.ID, constants.OpTypeUpdate, constants.EntityTopic, topic.ID, "", c.Ctx.Request())
	return web.JsonData(payload.BuildTopic(topic, user))
}

func (c *TopicsController) PostByApprove(slugID string) *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}

	if !user.IsManagerRole() {
		return web.JsonError(errs.ErrForbidden)
	}

	topic, err := c.getTopicBySlugId(slugID)
	if topic == nil || err != nil {
		return web.JsonError(errs.ErrTopicNotFound)
	}

	if topic.Status != constants.StatusReview {
		return web.JsonError(errs.ErrTopicNotUnderReview)
	}
	if err := service.TopicService.ApproveTopic(user.ID, topic.ID); err != nil {
		return web.JsonError(errs.ErrInternalServer)
	}
	return web.JsonData(payload.BuildTopic(topic, user))
}

func (c *TopicsController) PostByReject(slugID string) *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}

	if !user.IsManagerRole() {
		return web.JsonError(errs.ErrForbidden)
	}

	topic, err := c.getTopicBySlugId(slugID)
	if topic == nil || err != nil {
		return web.JsonError(errs.ErrTopicNotFound)
	}

	if topic.Status != constants.StatusReview {
		return web.JsonError(errs.ErrTopicNotUnderReview)
	}

	reason := params.FormValueDefault(c.Ctx, "reason", "")
	if err := service.TopicService.RejectTopic(user.ID, topic.ID, reason); err != nil {
		return web.JsonError(errs.ErrInternalServer)
	}
	return web.JsonSuccess()
}

func (c *TopicsController) setTopicPinned(userID int64, topicID int64, pinned bool) (*web.JsonResult, error) {
	err := service.TopicService.SetTopicPinned(userID, topicID, pinned)
	if err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

func (c *TopicsController) setTopicRecommended(userID int64, topicID int64, recommended bool) (*web.JsonResult, error) {
	err := service.TopicService.SetTopicRecommended(userID, topicID, recommended)
	if err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

// PATCH /topics/{slugId} => pin/unpin, recommend/unrecommend
func (c *TopicsController) PatchBy(slugId string) (*web.JsonResult, error) {
	var patch struct {
		Pinned      *bool `json:"pinned,omitempty"`
		Recommended *bool `json:"recommended,omitempty"`
	}
	topic, err := c.getTopicBySlugId(slugId)
	if topic == nil || err != nil {
		return web.JsonError(errs.ErrTopicNotFound), nil
	}

	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil || !user.IsManagerRole() {
		return web.JsonError(errs.ErrForbidden), nil
	}

	if err := c.Ctx.ReadJSON(&patch); err != nil {
		return web.JsonError(errs.ErrBadRequest), nil
	}

	if patch.Pinned != nil {
		return c.setTopicPinned(user.ID, topic.ID, *patch.Pinned)
	}
	if patch.Recommended != nil {
		return c.setTopicRecommended(user.ID, topic.ID, *patch.Recommended)
	}

	return nil, errs.ErrBadRequest
}

// DELETE /topics/{slug}
func (c *TopicsController) DeleteBy(slugId string) (*web.JsonResult, error) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized), nil
	}

	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err), nil
	}

	topic, err := c.getTopicBySlugId(slugId)
	if topic == nil || err != nil {
		return web.JsonError(errs.ErrTopicNotFound), nil
	}

	if err := canViewTopic(user, topic); err != nil {
		return web.JsonError(err), nil
	}

	if err := service.TopicService.Delete(topic.ID, user.ID, c.Ctx.Request()); err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

func (c *TopicsController) PostByRestore(slugId string) *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}

	if !user.IsManagerRole() {
		return web.JsonError(errs.ErrForbidden)
	}

	topic, err := c.getTopicBySlugId(slugId)
	if topic == nil || err != nil {
		return web.JsonError(errs.ErrTopicNotFound)
	}

	if topic.Status != constants.StatusDeleted {
		return web.JsonSuccess()
	}

	if err := service.TopicService.Restore(topic.ID); err != nil {
		return web.JsonError(errs.ErrInternalServer)
	}
	topic.Status = constants.StatusActive
	return web.JsonData(payload.BuildTopic(topic, user))
}

// POST /topics/{slugId}/reactions
func (c *TopicsController) PostByReactions(slugId string) (*web.JsonResult, error) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized), nil
	}
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err), nil
	}

	topic, err := c.getTopicBySlugId(slugId)
	if err != nil {
		return web.JsonError(errs.ErrTopicNotFound), nil
	}

	if err := canViewTopic(user, topic); err != nil {
		return web.JsonError(err), nil
	}

	err = service.UserLikeService.TopicLike(user.ID, topic.ID)
	if err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

// DELETE /topics/{slugId}/reactions
func (c *TopicsController) DeleteByReactions(slugId string) (*web.JsonResult, error) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized), nil
	}

	topic, err := c.getTopicBySlugId(slugId)
	if err != nil {
		return web.JsonError(errs.ErrTopicNotFound), nil
	}

	if err := canViewTopic(user, topic); err != nil {
		return web.JsonError(err), nil
	}

	err = service.UserLikeService.TopicUnLike(user.ID, topic.ID)
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

	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := canViewTopic(user, topic); err != nil {
		return web.JsonError(err), nil
	}

	comments, cursor, hasMore := service.CommentService.GetComments(topic.ID, cursor, 20)
	resp := payload.BuildComments(comments, user, true, false)
	return web.JsonCursorData(resp, cursor, hasMore), nil
}

// POST /topics/{slugId}/comments
func (c *TopicsController) PostByComments(slugId string) (*web.JsonResult, error) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized), nil
	}

	topic, err := c.getTopicBySlugId(slugId)
	if err != nil {
		return web.JsonError(errs.ErrTopicNotFound), nil
	}

	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err), nil
	}

	if err := canViewTopic(user, topic); err != nil {
		return web.JsonError(err), nil
	}

	form := payload.GetCreateCommentForm(c.Ctx)
	if err := spam.CheckComment(user, form); err != nil {
		return web.JsonError(err), nil
	}

	comment, err := service.CommentService.CreateComment(service.CreateCommentArgs{
		UserId:    user.ID,
		TopicId:   topic.ID,
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
