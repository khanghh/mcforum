import type { UserProfile, Topic, Forum, Comment } from '@/types'
import { useHttpDelete, useHttpGet, useHttpPost } from './http'

export class CursorResult<T extends any[]> {
  url: string
  params?: Record<string, any>
  hasMore: boolean = true
  cursor: number | null = 0
  nextCursor: number | null = null

  constructor(url: string, params?: Record<string, any>) {
    this.url = url
    this.params = params
  }

  private async load(cursor: number): Promise<T> {
    if (!this.hasMore) return Array<any>() as T

    const reqParams: Record<string, any> = Object.assign({}, this.params || {})
    if (cursor) {
      reqParams.cursor = cursor
    }

    const res = await useHttpGet(this.url, {
      params: reqParams,
    })

    const payload = res.data ? res.data : res
    const items: T = payload.items
    const nextCursor = typeof payload.cursor === 'number' ? payload.cursor : null
    const hasMore = !!payload.hasMore

    this.hasMore = hasMore
    this.cursor = cursor
    this.nextCursor = nextCursor

    return items ?? []
  }

  async loadNext(): Promise<T> {
    return this.load(this.nextCursor || 0)
  }

  async reload(): Promise<T> {
    return this.load(this.cursor || 0)
  }

  reset(): void {
    this.cursor = 0
    this.nextCursor = 0
    this.hasMore = true
  }
}

export type CreateCommentPayload = {
  content: string
  imageList?: string[]
  quoteId?: number
}


export type UpdateProfilePayload = {
  nickname?: string
  bio?: string
  location?: string
  lockedProfile?: boolean
  showLocation?: boolean
  emailNotify?: boolean
}

export enum FeedType {
  WhatsNew = 'whats-new',
  Followed = 'followed',
  Recommended = 'recommended',
}

export const useApi = () => {
  // me api endpoints
  const getCurrentUser = (): Promise<UserProfile> => {
    return useHttpGet('/api/users/me')
  }

  const updateProfile = (data: UpdateProfilePayload): Promise<UserProfile> => {
    return useHttpPatch('/api/users/me', { body: data })
  }

  const isFollowing = (otherUserID: string): Promise<boolean> => {
    return useHttpGet(`/api/users/me/following/${otherUserID}`)
      .then(r => r.following)
  }

  const getMyFollowers = (): CursorResult<UserProfile[]> => {
    return new CursorResult<UserProfile[]>('/api/users/me/followers')
  }

  const getMyFollowing = (): CursorResult<UserProfile[]> => {
    return new CursorResult<UserProfile[]>('/api/users/me/following')
  }

  const getMyTopics = (): CursorResult<Topic[]> => {
    return new CursorResult<Topic[]>("/api/users/me/topics")
  }

  // other user api endpoints
  const toggleMyFavorite = (topicId: string): Promise<boolean> => {
    return useHttpPut(`/api/users/me/favorites/${topicId}`)
  }

  const getUser = (username: string): Promise<UserProfile> => {
    return useHttpGet(`/api/users/${username}`) as Promise<UserProfile>
  }

  const followUser = (username: string): Promise<boolean> => {
    return useHttpPost(`/api/users/${username}/follow`)
  }

  const unfollowUser = (username: string): Promise<void> => {
    return useHttpDelete(`/api/users/${username}/follow`)
  }

  const getUserFollowers = (username: string): CursorResult<UserProfile[]> => {
    return new CursorResult<UserProfile[]>(`/api/users/${username}/followers`)
  }

  const getUserFollowing = (username: string): CursorResult<UserProfile[]> => {
    return new CursorResult<UserProfile[]>(`/api/users/${username}/following`)
  }

  const getUserTopics = (username: string): CursorResult<Topic[]> => {
    return new CursorResult<Topic[]>(`/api/users/${username}/topics`)
  }

  // forums api endpoints
  const getTopicFeeds = (feedType: FeedType): CursorResult<Topic[]> => {
    return new CursorResult<Topic[]>(`/api/feeds/${feedType}`)
  }

  const getForumList = (): Promise<Forum[]> => {
    return useHttpGet("/api/forums")
  }

  const getForumTopics = (forumSlug: string): CursorResult<Topic[]> => {
    return new CursorResult<Topic[]>(`/api/forums/${forumSlug}`)
  }

  // topics api endpoints
  const getTopic = (topicSlug: string): Promise<Topic> => {
    return useHttpGet(`/api/topics/${topicSlug}`) as Promise<Topic>
  }

  const getTopicsByTag = (tag: string): CursorResult<Topic[]> => {
    return new CursorResult<Topic[]>(`/api/topics?tag=${encodeURIComponent(tag)}`)
  }

  const deleteTopic = (topicSlug: string): Promise<void> => {
    return useHttpDelete(`/api/topics/${topicSlug}`)
  }

  const getTopicComments = (topicSlug: string): CursorResult<Comment[]> => {
    return new CursorResult<Comment[]>(`/api/topics/${topicSlug}/comments`)
  }

  const addTopicComment = (topicSlug: string, payload: CreateCommentPayload): Promise<any> => {
    let imageList = payload.imageList ? payload.imageList : []
    const body = {
      content: payload.content,
      imageList: imageList.length > 0 ? JSON.stringify(imageList) : undefined,
      quoteId: payload.quoteId,
    }
    return useHttpPost(`/api/topics/${topicSlug}/comments`, { body })
  }


  const addTopicReaction = (topicSlug: string, reactionType: string): Promise<boolean> => {
    return useHttpPost(`/api/topics/${topicSlug}/reactions`, {
      body: { type: reactionType },
    })
  }

  const removeTopicReaction = (topicSlug: string): Promise<void> => {
    return useHttpDelete(`/api/topics/${topicSlug}/reactions`)
  }

  // comment api endpoints
  const getCommentReplies = (commentId: number): CursorResult<Comment[]> => {
    return new CursorResult<Comment[]>(`/api/comments/${commentId}/replies`)
  }

  const addCommentReply = (commentId: number, payload: CreateCommentPayload): Promise<any> => {
    let imageList = payload.imageList ? payload.imageList : []
    const body = {
      content: payload.content,
      imageList: imageList.length > 0 ? JSON.stringify(imageList) : undefined,
      quoteId: payload.quoteId,
    }
    return useHttpPost(`/api/comments/${commentId}/replies`, { body })
  }

  const addCommentReaction = (commentId: number, reactionType: string): Promise<boolean> => {
    return useHttpPost(`/api/comments/${commentId}/reactions`, {
      body: { type: reactionType },
    })
  }

  const removeCommentReaction = (commentId: number): Promise<void> => {
    return useHttpDelete(`/api/comments/${commentId}/reactions`)
  }

  const deleteComment = (commentId: number): Promise<void> => {
    return useHttpDelete(`/api/comments/${commentId}`)
  }


  return {
    // me api endpoints
    getCurrentUser,
    updateProfile,
    isFollowing,
    getMyFollowers,
    getMyFollowing,
    getMyTopics,

    // other user api endpoints
    toggleMyFavorite,
    getUser,
    followUser,
    unfollowUser,
    getUserFollowers,
    getUserFollowing,
    getUserTopics,

    // forums api endpoints
    getTopicFeeds,
    getForumList,
    getForumTopics,

    // topics api endpoints
    getTopic,
    getTopicsByTag,
    deleteTopic,
    getTopicComments,
    addTopicComment,
    addTopicReaction,
    removeTopicReaction,

    // comment api endpoints
    getCommentReplies,
    addCommentReply,
    addCommentReaction,
    removeCommentReaction,
    deleteComment,
  }
}
