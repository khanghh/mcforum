
type CommentReplies = {
  items: Comment[]
  cursor?: number
  hasMore?: boolean
}

export interface CommentUser {
  id: number
  username: string
  nickname: string
  avatar: string
}

export interface Comment {
  id: number
  user: CommentUser
  paerntId: string
  quoteId: string
  contentType: string
  content: string
  imageList: null
  likeCount: number
  commentCount: number
  liked: boolean
  quote?: Comment
  replies?: CommentReplies
  ipLocation: string
  status: number
  createTime: number
}
