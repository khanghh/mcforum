export enum MessageType {
  TopicComment = 1,
  CommentReply = 2,
  TopicLike = 3,
  TopicFavorite = 4,
  TopicRecommended = 5,
  TopicPinned = 6,
  TopicDelete = 7,
  UserFollow = 8,
  CommentLike = 9,
  FollowingUserCreateTopic = 10,
}

export interface Message {
  messageId: number
  from: {
    id: number
    username: string
    nickname: string
    avatar?: string
    smallAvatar?: string
  }
  title?: string
  content: string
  quoteContent?: string
  detailUrl?: string
  createTime: number
  status?: number
}
