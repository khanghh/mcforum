package payload

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"
	"bbs-go/pkg/bbsurls"
	"bbs-go/pkg/markdown"
	"bbs-go/sqls"
	"html"
	"strconv"

	"bbs-go/common/arrays"
	"bbs-go/common/dates"
	"bbs-go/common/strs"
	"bbs-go/pkg/web"
)

type CommentUserInfo struct {
	Id          int64  `json:"id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	SmallAvatar string `json:"smallAvatar"`
}

// CommentResponse 评论返回数据
type CommentResponse struct {
	Id           int64             `json:"id"`
	User         *CommentUserInfo  `json:"user"`
	ParentId     int64             `json:"parentId,omitempty"`
	QuoteId      int64             `json:"quoteId,omitempty"`
	ContentType  string            `json:"contentType"`
	Content      string            `json:"content"`
	ImageList    []ImageInfo       `json:"imageList"`
	LikeCount    int64             `json:"likeCount"`
	CommentCount int64             `json:"commentCount"`
	Liked        bool              `json:"liked"`
	Quote        *CommentResponse  `json:"quote,omitempty"`
	Replies      *web.CursorResult `json:"replies,omitempty"`
	IpLocation   string            `json:"ipLocation,omitempty"`
	Status       int               `json:"status"`
	CreateTime   int64             `json:"createTime"`
}

func BuildComment(comment *model.Comment) *CommentResponse {
	return doBuildComment(comment, nil, true, true)
}

func BuildComments(comments []model.Comment, currentUser *model.User, isBuildReplies, isBuildQuote bool) []CommentResponse {
	if len(comments) == 0 {
		return []CommentResponse{}
	}

	likedCommentIds := getLikedCommentIds(comments, currentUser)

	var ret []CommentResponse
	for _, comment := range comments {
		item := doBuildComment(&comment, currentUser, isBuildReplies, isBuildQuote)
		item.Liked = arrays.Contains(likedCommentIds, comment.Id)
		ret = append(ret, *item)
	}
	return ret
}

func getLikedCommentIds(comments []model.Comment, currentUser *model.User) (likedCommentIds []int64) {
	if currentUser == nil || len(comments) == 0 {
		return
	}
	var commentIds []int64
	for _, comment := range comments {
		commentIds = append(commentIds, comment.Id)
	}
	likedCommentIds = service.UserLikeService.GetUserLikes(currentUser.Id, constants.EntityComment, commentIds)
	return
}

func buildCommentUserInfo(userId int64) *CommentUserInfo {
	user := cache.UserCache.Get(userId)
	if user == nil {
		user = &model.User{}
		user.Id = userId
		user.Type = constants.UserTypeNormal
		user.Username = sqls.SqlNullString("user" + strconv.FormatInt(userId, 10))
		user.Nickname = "deleted_user"
		user.CreateTime = dates.NowTimestamp()
	}
	var userAvatar string
	if strs.IsNotBlank(user.Avatar) {
		userAvatar = user.Avatar
	} else {
		userAvatar = bbsurls.AbsUrl("/images/avatars/steve.png")
	}
	return &CommentUserInfo{
		Id:          user.Id,
		Username:    user.Username.String,
		Avatar:      userAvatar,
		SmallAvatar: HandleOssImageStyleAvatar(userAvatar),
		Nickname:    user.Nickname,
	}
}

// doBuildComment 渲染评论
// isBuildReplies 是否渲染评论的二级回复，一级评论时要设置为true，其他时候都为false
// isBuildQuote 是否渲染评论的引用，二级回复时要设置为true，其他时候都为false
func doBuildComment(comment *model.Comment, currentUser *model.User, isBuildReplies bool, isBuildQuote bool) *CommentResponse {
	if comment == nil {
		return nil
	}

	ret := &CommentResponse{
		Id:           comment.Id,
		User:         buildCommentUserInfo(comment.UserId),
		ParentId:     comment.ParentId,
		QuoteId:      comment.QuoteId,
		LikeCount:    comment.LikeCount,
		CommentCount: comment.CommentCount,
		ContentType:  comment.ContentType,
		IpLocation:   comment.IpLocation,
		Status:       comment.Status,
		CreateTime:   comment.CreateTime,
	}

	if comment.Status == constants.StatusActive {
		if comment.ContentType == constants.ContentTypeMarkdown {
			content := markdown.ToHTML(comment.Content)
			ret.Content = handleHtmlContent(content)
		} else if comment.ContentType == constants.ContentTypeHtml {
			ret.Content = handleHtmlContent(comment.Content)
		} else {
			ret.Content = html.EscapeString(comment.Content)
		}
		ret.ImageList = BuildImageList(comment.ImageList)
	} else {
		ret.Content = "内容已删除"
	}

	if isBuildReplies && comment.CommentCount > 0 {
		var repliesLimit int64 = 3
		replies, nextCursor, _ := service.CommentService.GetReplies(comment.Id, 0, int(repliesLimit))
		// var replyResults []*CommentResponse
		// for _, reply := range replies {
		// 	replyResults = append(replyResults, doBuildComment(&reply, false, true))
		// }
		replyResults := BuildComments(replies, currentUser, false, true)
		ret.Replies = &web.CursorResult{
			Items:   replyResults,
			Cursor:  nextCursor,
			HasMore: comment.CommentCount > repliesLimit,
		}
	}

	if isBuildQuote && comment.QuoteId > 0 {
		quote := doBuildComment(service.CommentService.Get(comment.QuoteId), currentUser, false, false)
		if quote != nil {
			ret.Quote = quote
		}
	}

	return ret
}
