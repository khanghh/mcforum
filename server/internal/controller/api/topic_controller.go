package api

import (
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/locale"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/spam"
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"unicode"

	"bbs-go/common/numbers"
	"bbs-go/common/strs"
	"bbs-go/sqls"
	"bbs-go/web"
	"bbs-go/web/params"

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
	topicId := numbers.ToInt64(parts[1])
	cnd := sqls.NewCnd().
		Eq("slug", topicSlug).
		Eq("id", topicId).
		Eq("status", constants.StatusOK)
	topic := service.TopicService.FindOne(cnd)
	return topic, nil
}

func normalizeVietnamese(input string) string {
	input = strings.ToLower(input)
	normalizationMap := map[rune]rune{
		'á': 'a', 'à': 'a', 'ả': 'a', 'ã': 'a', 'ạ': 'a', 'ă': 'a', 'ằ': 'a', 'ắ': 'a', 'ẳ': 'a', 'ẵ': 'a', 'ặ': 'a',
		'â': 'a', 'ầ': 'a', 'ấ': 'a', 'ẩ': 'a', 'ẫ': 'a', 'ậ': 'a',
		'é': 'e', 'è': 'e', 'ẻ': 'e', 'ẽ': 'e', 'ẹ': 'e', 'ê': 'e', 'ề': 'e', 'ế': 'e', 'ể': 'e', 'ễ': 'e', 'ệ': 'e',
		'í': 'i', 'ì': 'i', 'ỉ': 'i', 'ĩ': 'i', 'ị': 'i',
		'ó': 'o', 'ò': 'o', 'ỏ': 'o', 'õ': 'o', 'ọ': 'o', 'ô': 'o', 'ồ': 'o', 'ố': 'o', 'ổ': 'o', 'ỗ': 'o', 'ộ': 'o',
		'ơ': 'o', 'ờ': 'o', 'ớ': 'o', 'ở': 'o', 'ỡ': 'o', 'ợ': 'o',
		'ú': 'u', 'ù': 'u', 'ủ': 'u', 'ũ': 'u', 'ụ': 'u', 'ư': 'u', 'ừ': 'u', 'ứ': 'u', 'ử': 'u', 'ữ': 'u', 'ự': 'u',
		'ý': 'y', 'ỳ': 'y', 'ỷ': 'y', 'ỹ': 'y', 'ỵ': 'y',
		'đ': 'd',
	}
	var result []rune
	for _, r := range input {
		if replacement, exists := normalizationMap[r]; exists {
			result = append(result, replacement)
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func generateSlug(title string) string {
	normalized := normalizeVietnamese(title)
	var buf bytes.Buffer
	for _, r := range normalized {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' {
			buf.WriteRune(r)
		} else {
			buf.WriteRune('-')
		}
	}
	cleanStr := buf.String()
	for strings.Contains(cleanStr, "--") {
		cleanStr = strings.ReplaceAll(cleanStr, "--", "-")
	}
	cleanStr = strings.Trim(cleanStr, "-")
	return url.PathEscape(cleanStr)
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
func (c *TopicController) PostBy() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}
	form := model.GetCreateTopicForm(c.Ctx)
	form.Slug = generateSlug(form.Title)

	if err := spam.CheckTopic(user, form); err != nil {
		return web.JsonError(err)
	}

	topic, err := service.TopicPublishService.Publish(user.Id, form)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(payload.BuildSimpleTopic(topic))
}

// GET /topics/{slug}
func (c *TopicController) GetBy(slugId string) *web.JsonResult {
	topic, err := c.getTopicBySlugId(slugId)
	if topic == nil || err != nil {
		return web.JsonErrorCodeMsg(iris.StatusNotFound, locale.T("topic.not_found"))
	}

	// 审核中文章控制展示
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if topic.Status == constants.StatusReview {
		if user != nil {
			if topic.UserId != user.Id && !user.IsOwnerOrAdmin() {
				return web.JsonErrorCodeMsg(iris.StatusForbidden, "文章审核中")
			}
		} else {
			return web.JsonErrorCodeMsg(iris.StatusForbidden, "文章审核中")
		}
	}

	service.TopicService.IncrViewCount(topic.Id) // 增加浏览量
	return web.JsonData(payload.BuildTopic(topic, user))
}

// PUT /topics/{slug} // edit topic
func (c *TopicController) PutBy(slugId string) *web.JsonResult {
	topic, err := c.getTopicBySlugId(slugId)
	if topic == nil || err != nil {
		return web.JsonErrorCodeMsg(iris.StatusNotFound, locale.T("topic.not_found"))
	}

	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonErrorCode(iris.StatusForbidden, err)
	}

	// 非作者、且非管理员
	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return web.JsonErrorCodeMsg(iris.StatusForbidden, locale.T("system.message.permission_denied"))
	}

	topic.ForumId = params.FormValueInt64Default(c.Ctx, "forumId", topic.ForumId)
	topic.Title = strings.TrimSpace(params.FormValueDefault(c.Ctx, "title", topic.Title))
	topic.Slug = generateSlug(topic.Title)
	topic.Content = strings.TrimSpace(params.FormValueDefault(c.Ctx, "content", topic.Content))
	topic.HideContent = strings.TrimSpace(params.FormValueDefault(c.Ctx, "hideContent", topic.HideContent))
	tags := params.FormValueStringArray(c.Ctx, "tags")

	err = service.TopicService.Edit(topic.Id, topic.ForumId, tags, topic.Title, topic.Slug, topic.Content, topic.HideContent)
	if err != nil {
		return web.JsonError(err)
	}
	// 操作日志
	service.OperateLogService.AddOperateLog(user.Id, constants.OpTypeUpdate, constants.EntityTopic, topic.Id, "", c.Ctx.Request())
	return web.JsonData(payload.BuildTopic(topic, user))
}

// PATCH /topics/{slugId} => pin/unpin, recommend/unrecommend
func (c *TopicController) PatchBy(slugId string) *web.JsonResult {
	topic, err := c.getTopicBySlugId(slugId)
	if topic == nil || err != nil {
		return web.JsonErrorCodeMsg(iris.StatusNotFound, locale.T("topic.not_found"))
	}
	var (
		pinned      = params.FormValueBoolDefault(c.Ctx, "pinned", topic.Pinned)
		recommended = params.FormValueBoolDefault(c.Ctx, "recommended", topic.Recommended)
	)

	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil || !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return web.JsonErrorCodeMsg(iris.StatusForbidden, locale.T("system.message.permission_denied"))
	}

	if pinned != topic.Pinned {
		err = service.TopicService.SetTopicPinned(topic.Id, pinned)
	} else if recommended != topic.Recommended {
		err = service.TopicService.SetRecommended(topic.Id, recommended)
	} else {
		return web.JsonErrorCodeMsg(iris.StatusBadRequest, locale.T("system.message.invalid_request"))
	}
	if err == nil {
		return web.JsonSuccess()
	} else if service.IsDatabaseError(err) {
		return web.JsonError(err)
	}
	return web.JsonErrorCodeMsg(iris.StatusInternalServerError, locale.T("system.errors.internal_server_error"))
}

// DELETE /topics/{slug}
func (c *TopicController) DeleteBy(topicId int64) (*web.JsonResult, int) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err), iris.StatusUnauthorized
	}

	topic := service.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOK {
		return web.JsonSuccess(), iris.StatusBadRequest
	}

	// 非作者、且非管理员
	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return web.JsonErrorMsg("permission denied"), iris.StatusForbidden
	}

	if err := service.TopicService.Delete(topicId, user.Id, c.Ctx.Request()); err != nil {
		return web.JsonError(err), iris.StatusInternalServerError
	}
	return web.JsonSuccess(), iris.StatusOK
}

// // P/topics/{slug}/edit
// func (c *TopicController) GetEditBy(topicId int64) (*web.JsonResult, int) {
// 	user := service.UserTokenService.GetCurrent(c.Ctx)
// 	if err := service.UserService.CheckPostStatus(user); err != nil {
// 		return web.JsonError(err), iris.StatusForbidden
// 	}

// 	topic := service.TopicService.Get(topicId)
// 	if topic == nil || topic.Status != constants.StatusOK {
// 		return web.JsonErrorMsg(locale.T("topic.not_found")), iris.StatusNotFound
// 	}
// 	if topic.Type != constants.TopicTypeTopic {
// 		return web.JsonErrorMsg(locale.T("topic.not_editable")), iris.StatusForbidden
// 	}

// 	// 非作者、且非管理员
// 	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
// 		return web.JsonErrorMsg(locale.T("system.message.permission_denied")), iris.StatusForbidden
// 	}

// 	tags := service.TopicService.GetTopicTags(topicId)
// 	var tagNames []string
// 	if len(tags) > 0 {
// 		for _, tag := range tags {
// 			tagNames = append(tagNames, tag.Name)
// 		}
// 	}

// 	return web.NewEmptyRspBuilder().
// 		Put("id", topic.Id).
// 		Put("title", topic.Title).
// 		Put("forumId", topic.ForumId).
// 		Put("content", topic.Content).
// 		Put("hideContent", topic.HideContent).
// 		Put("tags", tagNames).
// 		JsonResult(), iris.StatusOK
// }

// PUT /topics/{slugId}/reactions
func (c *TopicController) PutByReactions(slugId string) (*web.JsonResult, int) {
	topic, err := c.getTopicBySlugId(slugId)
	if err != nil {
		return web.JsonErrorMsg("not found"), iris.StatusNotFound
	}
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin), iris.StatusUnauthorized
	}
	err = service.UserLikeService.TopicLike(user.Id, topic.Id)
	if err != nil {
		fmt.Println(err)
		return web.JsonError(err), iris.StatusInternalServerError
	}
	return web.JsonSuccess(), iris.StatusOK
}

// DELETE /topics/{slugId}/reactions/{userId}
func (c *TopicController) DeleteByReactions(slugId string) (*web.JsonResult, int) {
	topic, err := c.getTopicBySlugId(slugId)
	if err != nil {
		return web.JsonErrorMsg("not found"), iris.StatusNotFound
	}
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin), iris.StatusUnauthorized
	}
	err = service.UserLikeService.TopicUnLike(user.Id, topic.Id)
	if err != nil {
		return web.JsonError(err), iris.StatusInternalServerError
	}
	return web.JsonSuccess(), iris.StatusOK
}

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

func (c *TopicController) GetHide_content() *web.JsonResult {
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
