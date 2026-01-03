import type { UserInfo } from './user'

export interface ImageInfo {
  url: string
  preview?: string
}

export enum TopicType {
  Topic = 0,
  Tweet = 1,
}

export enum TopicStatus {
  Active = 0,
  Hidden = 1,
  PendingReview = 2,
}

export interface Forum {
  id: number
  name: string
  slug: string
  logo?: string
  description?: string
}

export interface Topic {
  id: number
  slug: string
  type: TopicType
  user: UserInfo
  forum: Forum
  tags: string[]
  title: string
  summary: string
  content: string
  hiddenContent?: string
  hasHiddenContent?: boolean
  images: ImageInfo[]
  lastCommentTime?: number
  viewCount: number
  commentCount: number
  likeCount: number
  liked: boolean
  createTime: number
  recommended: boolean
  recommendedTime?: number
  pinned: boolean
  pinnedTime: number
  status: number
  favorited: boolean
  poll?: TopicPoll
  ipLocation?: string
}

export interface TopicPollOption {
  id: number
  text: string
  voteCount: number
  votePercentage: number
  selected: boolean
}

export interface TopicPoll {
  id: number
  question: string
  options: TopicPollOption[]
  totalVotes: number
  durationHours: number
  multiSelect: boolean
  publicResults: boolean
  allowVoteChange: boolean
  endTime: number
  createdTime: number
}
