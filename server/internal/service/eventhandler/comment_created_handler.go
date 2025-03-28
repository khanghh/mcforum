package eventhandler

import (
	"bbs-go/internal/model"
	"bbs-go/internal/pkg/event"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.CommentCreatedEvent{}), handleCommentCreate)
}

func handleCommentCreate(i interface{}) {
	// e := i.(event.CommentCreatedEvent)

	// service.UserService.IncrCommentCount(userId)
	// service.UserService.IncrScoreForPostComment(comment)

	// processTopicAuthorNotification(&e)
	// processCommentAuthorNotification(&e)
	// processQuoteAuthorNotification(&e)
}

// 处理评论消息
func handleMsg(comment *model.Comment) {
	// commentMsg := getCommentMsg(comment)

	// handleTopicMsg(comment)
	// handleQuoteMsg(comment)
	// handleReplyMsg(comment)
}

// processTopicAuthorNotification send message to notify topic owner of new comment created
// func processTopicAuthorNotification(e *event.CommentCreatedEvent) {
// 	var (
// 		from = e.Comment.UserId
// 		to   = e.Topic.UserId
// 	)
// 	if from == to {
// 		return
// 	}

// 	service.MessageService.SendMsg(
// 		from,
// 		to,
// 		msg.TypeTopicComment,
// 		commentMsg.msgTitle(),
// 		commentMsg.msgContent(),
// 		commentMsg.msgRepliedContent(),
// 		nil,
// 	)
// }

// func handleReplyMsg(comment *model.Comment, commentMsg *CommentMsg) {
// 	if commentMsg.ParentComment == nil {
// 		return
// 	}

// 	var (
// 		from = comment.UserId
// 		to   = commentMsg.ParentComment.UserId
// 	)

// 	if from == to {
// 		return
// 	}

// 	// 如果回复的评论作者就是被引用消息的作者作者，那么只发引用消息即可
// 	if commentMsg.QuoteComment != nil && commentMsg.QuoteComment.UserId == to {
// 		return
// 	}

// 	var (
// 		title          = commentMsg.msgTitle()
// 		content        = commentMsg.msgContent()
// 		repliedContent = common.GetSummary(commentMsg.ParentComment.ContentType, commentMsg.ParentComment.Content)
// 	)

// 	service.MessageService.SendMsg(from, to, msg.TypeCommentReply, title, content, repliedContent,
// 		&msg.CommentExtraData{
// 			EntityType:     comment.EntityType,
// 			EntityId:       comment.EntityId,
// 			QuoteId:        comment.QuoteId,
// 			RootEntityType: commentMsg.rootEntityType(),
// 			RootEntityId:   cast.ToString(commentMsg.rootEntityId()),
// 		})
// }

// handleQuoteMsg 给被引用人发送消息
// func handleQuoteMsg(comment *model.Comment, commentMsg *CommentMsg) {
// 	if commentMsg.QuoteComment == nil {
// 		return
// 	}

// 	var (
// 		from           = comment.UserId
// 		to             = commentMsg.QuoteComment.UserId
// 		title          = commentMsg.msgTitle()
// 		content        = commentMsg.msgContent()
// 		repliedContent = common.GetSummary(commentMsg.QuoteComment.ContentType, commentMsg.QuoteComment.Content)
// 	)

// 	if from == to {
// 		return
// 	}

// 	service.MessageService.SendMsg(from, to, msg.TypeCommentReply, title, content, repliedContent,
// 		&msg.CommentExtraData{
// 			EntityType:     comment.EntityType,
// 			EntityId:       comment.EntityId,
// 			QuoteId:        comment.QuoteId,
// 			RootEntityType: commentMsg.rootEntityType(),
// 			RootEntityId:   cast.ToString(commentMsg.rootEntityId()),
// 		})
// }

// func getCommentMsg(comment *model.Comment) *CommentMsg {
// 	if comment.EntityType == constants.EntityTopic { // 帖子
// 		topic := service.TopicService.Get(comment.EntityId)
// 		if topic != nil && topic.Status == constants.StatusOK {
// 			return &CommentMsg{
// 				Comment:    comment,
// 				EntityType: comment.EntityType,
// 				EntityId:   comment.EntityId,
// 				Entity:     topic,
// 			}
// 		}
// 	} else if comment.EntityType == constants.EntityArticle { // 文章
// 		article := service.ArticleService.Get(comment.EntityId)
// 		if article != nil && article.Status == constants.StatusOK {
// 			return &CommentMsg{
// 				Comment:    comment,
// 				EntityType: comment.EntityType,
// 				EntityId:   comment.EntityId,
// 				Entity:     article,
// 			}
// 		}
// 	} else if comment.EntityType == constants.EntityComment { // 二级评论
// 		parentComment := service.CommentService.Get(comment.EntityId)

// 		if parentComment == nil || parentComment.Status != constants.StatusOK {
// 			return nil
// 		}

// 		ret := &CommentMsg{
// 			Comment:       comment,
// 			EntityType:    parentComment.EntityType, // 二级评论时，取一级评论的
// 			EntityId:      parentComment.EntityId,   // 二级评论时，取一级评论的
// 			ParentComment: parentComment,            // 一级评论
// 		}

// 		if parentComment.EntityType == constants.EntityTopic {
// 			topic := service.TopicService.Get(parentComment.EntityId)
// 			if topic != nil && topic.Status == constants.StatusOK {
// 				ret.Entity = topic
// 			}
// 		} else if parentComment.EntityType == constants.EntityArticle {
// 			article := service.ArticleService.Get(parentComment.EntityId)
// 			if article != nil && article.Status == constants.StatusOK {
// 				ret.Entity = article
// 			}
// 		} else {
// 			return nil
// 		}

// 		if comment.QuoteId > 0 { // 三级评论
// 			quoteComment := service.CommentService.Get(comment.QuoteId)
// 			if quoteComment != nil && quoteComment.Status == constants.StatusOK {
// 				ret.QuoteComment = quoteComment
// 			}
// 		}

// 		return ret
// 	}
// 	return nil
// }

// type CommentMsg struct {
// 	Topic         *Topic         // 实体类型
// 	ParentId      int64          // 实体ID
// 	Entity        interface{}    // 被评论实体
// 	Comment       *model.Comment // 当前评论
// 	ParentComment *model.Comment // 上一级评论（二级评论的时候有值）
// 	QuoteComment  *model.Comment // 引用评论
// }

// // msgType 消息类型
// func (c *CommentMsg) msgType() msg.Type {
// 	if c.EntityType == constants.EntityTopic {
// 		return msg.TypeTopicComment
// 	} else if c.EntityType == constants.EntityArticle {
// 		return msg.TypeArticleComment
// 	} else if c.EntityType == constants.EntityComment {
// 		return msg.TypeCommentReply
// 	}
// 	return msg.TypeTopicComment
// }

// // msgTitle 消息标题
// func (c *CommentMsg) msgTitle() string {
// 	if c.EntityType == constants.EntityTopic {
// 		return "回复了你的话题"
// 	} else if c.EntityType == constants.EntityArticle {
// 		return "回复了你的文章"
// 	} else if c.EntityType == constants.EntityComment {
// 		return "回复了你的评论"
// 	}
// 	return ""
// }

// // msgContent 回复内容
// func (c *CommentMsg) msgContent() string {
// 	return common.GetSummary(c.Comment.ContentType, c.Comment.Content)
// }

// // msgRepliedContent 被回复的内容
// func (c *CommentMsg) msgRepliedContent() string {
// 	if c.EntityType == constants.EntityArticle {
// 		article := c.Entity.(*model.Article)
// 		return "《" + article.Title + "》"
// 	} else if c.EntityType == constants.EntityTopic {
// 		topic := c.Entity.(*model.Topic)
// 		return "《" + topic.GetTitle() + "》"
// 	}
// 	return ""
// }

// func (c *CommentMsg) rootEntityUserId() int64 {
// 	if c.ParentComment != nil { // 二级评论
// 		if c.ParentComment.EntityType == constants.EntityTopic {
// 			topic := c.Entity.(*model.Topic)
// 			return topic.UserId
// 		} else if c.ParentComment.EntityType == constants.EntityArticle {
// 			article := c.Entity.(*model.Article)
// 			return article.UserId
// 		}
// 	} else {
// 		if c.Comment.EntityType == constants.EntityTopic {
// 			topic := c.Entity.(*model.Topic)
// 			return topic.UserId
// 		} else if c.Comment.EntityType == constants.EntityArticle {
// 			article := c.Entity.(*model.Article)
// 			return article.UserId
// 		}
// 	}
// 	return 0
// }

// func (c *CommentMsg) rootEntityType() string {
// 	if c.ParentComment != nil { // 二级评论
// 		return c.ParentComment.EntityType
// 	} else {
// 		return c.Comment.EntityType
// 	}
// }

// func (c *CommentMsg) rootEntityId() int64 {
// 	if c.ParentComment != nil { // 二级评论
// 		return c.ParentComment.EntityId
// 	} else {
// 		return c.Comment.EntityId
// 	}
// }
