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
