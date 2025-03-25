package api

import (
	"bbs-go/internal/controller/response"
	"bbs-go/internal/locale"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/spam"
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode"

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

func (c *TopicController) parseSlugId(slugId string) (slug string, id int64, err error) {
	parts := strings.SplitN(slugId, ".", 2)
	if len(parts) < 2 {
		err = fmt.Errorf("invalid slug")
		return
	}
	id, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return
	}
	slug = parts[0]
	return
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

// 发表帖子
func (c *TopicController) PostCreate() *web.JsonResult {
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
	return web.JsonData(response.BuildSimpleTopic(topic))
}

// 编辑时获取详情
func (c *TopicController) GetEditBy(topicId int64) (*web.JsonResult, int) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err), iris.StatusForbidden
	}

	topic := service.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOK {
		return web.JsonErrorMsg(locale.T("topic.not_found")), iris.StatusNotFound
	}
	if topic.Type != constants.TopicTypeTopic {
		return web.JsonErrorMsg(locale.T("topic.not_editable")), iris.StatusForbidden
	}

	// 非作者、且非管理员
	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return web.JsonErrorMsg(locale.T("system.message.permission_denied")), iris.StatusForbidden
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
		Put("title", topic.Title).
		Put("forumId", topic.ForumId).
		Put("content", topic.Content).
		Put("hideContent", topic.HideContent).
		Put("tags", tagNames).
		JsonResult(), iris.StatusOK
}

// 编辑帖子
func (c *TopicController) PostByEdit(slugId string) (*web.JsonResult, int) {
	_, topicId, err := c.parseSlugId(slugId)
	if err != nil {
		return web.JsonErrorMsg(locale.T("topic.not_found")), iris.StatusNotFound
	}
	topic := service.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOK {
		return web.JsonErrorMsg(locale.T("topic.not_found")), iris.StatusNotFound
	}

	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err), iris.StatusForbidden
	}

	// 非作者、且非管理员
	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return web.JsonErrorMsg(locale.T("system.message.permission_denied")), iris.StatusForbidden
	}

	var (
		forumId     = params.FormValueInt64Default(c.Ctx, "forumId", 0)
		title       = strings.TrimSpace(params.FormValue(c.Ctx, "title"))
		content     = strings.TrimSpace(params.FormValue(c.Ctx, "content"))
		hideContent = strings.TrimSpace(params.FormValue(c.Ctx, "hideContent"))
		tags        = params.FormValueStringArray(c.Ctx, "tags")
	)

	topic.Slug = generateSlug(title)
	err = service.TopicService.Edit(topicId, forumId, tags, title, topic.Slug, content, hideContent)
	if err != nil {
		return web.JsonError(err), iris.StatusInternalServerError
	}
	// 操作日志
	service.OperateLogService.AddOperateLog(user.Id, constants.OpTypeUpdate, constants.EntityTopic, topicId, "", c.Ctx.Request())
	return web.JsonData(response.BuildSimpleTopic(topic)), iris.StatusOK
}

// 删除帖子
func (c *TopicController) PostByDelete(topicId int64) (*web.JsonResult, int) {
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

// 帖子详情
func (c *TopicController) GetBy(slugId string) (*web.JsonResult, int) {
	topicSlug, topicId, err := c.parseSlugId(slugId)
	if err != nil {
		return web.JsonErrorMsg("not found"), iris.StatusNotFound
	}

	query := sqls.NewCnd().Eq("slug", topicSlug).Eq("id", topicId)
	topic := service.TopicService.FindOne(query)
	if topic == nil || topic.Status == constants.StatusDeleted {
		return web.JsonErrorMsg("not found"), iris.StatusNotFound
	}

	// 审核中文章控制展示
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if topic.Status == constants.StatusReview {
		if user != nil {
			if topic.UserId != user.Id && !user.IsOwnerOrAdmin() {
				return web.JsonErrorCode(403, "文章审核中"), iris.StatusForbidden
			}
		} else {
			return web.JsonErrorCode(403, "文章审核中"), iris.StatusForbidden
		}
	}

	service.TopicService.IncrViewCount(int64(topicId)) // 增加浏览量
	return web.JsonData(response.BuildTopic(topic, user)), iris.StatusOK
}

// POST /topics/{slugId}/like
func (c *TopicController) PostByLike(slugId string) (*web.JsonResult, int) {
	_, topicId, err := c.parseSlugId(slugId)
	if err != nil {
		return web.JsonErrorMsg("not found"), iris.StatusNotFound
	}
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin), iris.StatusUnauthorized
	}
	err = service.UserLikeService.TopicLike(user.Id, int64(topicId))
	if err != nil {
		fmt.Println(err)
		return web.JsonError(err), iris.StatusInternalServerError
	}
	return web.JsonSuccess(), iris.StatusOK
}

// POST /topics/{slugId}/unlike
func (c *TopicController) PostByUnlike(slugId string) (*web.JsonResult, int) {
	_, topicId, err := c.parseSlugId(slugId)
	if err != nil {
		return web.JsonErrorMsg("not found"), iris.StatusNotFound
	}
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin), iris.StatusUnauthorized
	}
	err = service.UserLikeService.TopicUnLike(user.Id, int64(topicId))
	if err != nil {
		return web.JsonError(err), iris.StatusInternalServerError
	}
	return web.JsonSuccess(), iris.StatusOK
}

// POST /topics/{slugId}/favorite
func (c *TopicController) PostByFavorite(slugId string) (*web.JsonResult, int) {
	_, topicId, err := c.parseSlugId(slugId)
	if err != nil {
		return web.JsonErrorMsg("not found"), iris.StatusNotFound
	}
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin), iris.StatusUnauthorized
	}
	err = service.FavoriteService.AddTopicFavorite(user.Id, int64(topicId))
	if err != nil {
		return web.JsonError(err), iris.StatusInternalServerError
	}
	return web.JsonSuccess(), iris.StatusOK
}

// POST /topics/{slugId}/unfavorite
func (c *TopicController) PostByUnfavorite(slugId string) (*web.JsonResult, int) {
	_, topicId, err := c.parseSlugId(slugId)
	if err != nil {
		return web.JsonErrorMsg("not found"), iris.StatusNotFound
	}
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin), iris.StatusUnauthorized
	}
	err = service.FavoriteService.RemoveTopicFavorite(user.Id, int64(topicId))
	if err != nil {
		return web.JsonError(err), iris.StatusInternalServerError
	}
	return web.JsonSuccess(), iris.StatusOK
}

// POST /topics/{slugId}/pin
func (c *TopicController) PostByPin(slugId string) (*web.JsonResult, int) {
	_, topicId, err := c.parseSlugId(slugId)
	if err != nil {
		return web.JsonErrorMsg("not found"), iris.StatusNotFound
	}
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil || !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return web.JsonErrorMsg("no permission"), iris.StatusUnauthorized
	}

	if err := service.TopicService.SetTopicPinned(int64(topicId), true); err != nil {
		return web.JsonError(err), iris.StatusInternalServerError
	}
	return web.JsonSuccess(), iris.StatusOK
}

// POST /topics/{slugId}/unpin
func (c *TopicController) PostByUnpin(slugId string) (*web.JsonResult, int) {
	_, topicId, err := c.parseSlugId(slugId)
	if err != nil {
		return web.JsonErrorMsg("not found"), iris.StatusNotFound
	}
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil || !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return web.JsonErrorMsg("no permission"), iris.StatusUnauthorized
	}

	if err := service.TopicService.SetTopicPinned(int64(topicId), false); err != nil {
		return web.JsonError(err), iris.StatusInternalServerError
	}
	return web.JsonSuccess(), iris.StatusOK
}

// POST /topics/{slugId}/recommend
func (c *TopicController) PostByRecommend(slugId string) (*web.JsonResult, int) {
	_, topicId, err := c.parseSlugId(slugId)
	if err != nil {
		return web.JsonErrorMsg("not found"), iris.StatusNotFound
	}
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil || !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return web.JsonErrorMsg("no permission"), iris.StatusUnauthorized
	}

	if err := service.TopicService.SetRecommended(int64(topicId), true); err != nil {
		return web.JsonError(err), iris.StatusInternalServerError
	}
	return web.JsonSuccess(), iris.StatusOK
}

// POST /topics/{slugId}/unrecommend
func (c *TopicController) PostByUnrecommend(slugId string) (*web.JsonResult, int) {
	_, topicId, err := c.parseSlugId(slugId)
	if err != nil {
		return web.JsonErrorMsg("not found"), iris.StatusNotFound
	}
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil || !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return web.JsonErrorMsg("no permission"), iris.StatusUnauthorized
	}

	if err := service.TopicService.SetRecommended(int64(topicId), false); err != nil {
		return web.JsonError(err), iris.StatusInternalServerError
	}
	return web.JsonSuccess(), iris.StatusOK
}

// 点赞用户
func (c *TopicController) GetRecentlikesBy(topicId int64) *web.JsonResult {
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
