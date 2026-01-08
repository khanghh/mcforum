import type { UserProfile, Topic, Forum, Comment, Message } from '@/types'
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

    const res = await useHttpGet(this.url, { params: reqParams })

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

export interface CreateCommentPayload {
  content: string
  images?: string[]
  quoteId?: number
}

export interface UpdateProfilePayload {
  nickname?: string
  bio?: string
  location?: string
  lockedProfile?: boolean
  showLocation?: boolean
  emailNotify?: boolean
}

export interface CreateTopicPayload {
  forumId: number
  title: string
  content: string
  tags: string[]
  images?: string[]
  hiddenContent?: string
  poll?: TopicPollPayload
}

export interface TopicPollPayload {
  question: string
  options: string[]
  durationHours: number
  multiSelect: boolean
  publicResults: boolean
  allowVoteChange: boolean
}

export interface UploadResult {
  url: string
  fileName: string
  thumbName?: string
  size: number
  mimeType: string
  createdAt: number
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
    return new CursorResult<Topic[]>('/api/users/me/topics')
  }

  const getMyFavorites = (): CursorResult<Topic[]> => {
    return new CursorResult<Topic[]>('/api/users/me/favorites')
  }

  const setTopicFavorite = (topicId: number, favorited: boolean): Promise<boolean> => {
    if (favorited) {
      return useHttpPut(`/api/users/me/favorites/${topicId}`)
    }
    return useHttpDelete(`/api/users/me/favorites/${topicId}`)
  }

  const getMessages = (): CursorResult<Message[]> => {
    return new CursorResult<Message[]>('/api/users/me/messages')
  }

  const getRecentMessages = (): Promise<{ count: number, messages: Message[] }> => {
    return useHttpGet('/api/users/me/messages/recent') as Promise<{ count: number, messages: Message[] }>
  }

  const setStatusMessage = (message: string): Promise<void> => {
    return useHttpPut('/api/users/me/status', { body: { message } })
  }

  const uploadAvatar = (file: File, apply: boolean = false): Promise<void> => {
    const form = new FormData()
    form.append('file', file, file.name)
    return useHttpPut('/api/users/me/avatar', { body: form, params: { apply } })
  }

  const removeAvatar = (): Promise<void> => {
    return useHttpDelete('/api/users/me/avatar')
  }

  const uploadCover = (file: File): Promise<void> => {
    const form = new FormData()
    form.append('file', file, file.name)
    return useHttpPut('/api/users/me/cover', { body: form })
  }

  // other user api endpoints

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

  const getUserActivities = (username: string): CursorResult<Message[]> => {
    return new CursorResult<Message[]>(`/api/users/${username}/activities`)
  }

  // forums api endpoints
  const getTopicFeeds = (feedType: FeedType): CursorResult<Topic[]> => {
    return new CursorResult<Topic[]>(`/api/feeds/${feedType}`)
  }

  const getForumList = (): Promise<Forum[]> => {
    return useHttpGet('/api/forums')
  }

  const getForumTopics = (forumSlug: string): CursorResult<Topic[]> => {
    return new CursorResult<Topic[]>(`/api/forums/${forumSlug}`)
  }

  const getForumStats = (): Promise<any> => {
    return useHttpGet('/api/forums/stats')
  }

  const getTopContributors = (): Promise<any> => {
    return useHttpGet('/api/forums/top-contributors')
  }

  // topics api endpoints
  const getTopic = (topicSlug: string): Promise<Topic> => {
    return useHttpGet(`/api/topics/${topicSlug}`) as Promise<Topic>
  }

  const createTopic = (payload: CreateTopicPayload): Promise<Topic> => {
    return useHttpPost(`/api/topics`, { body: payload }) as Promise<Topic>
  }

  const getTopicsByTag = (tag: string): CursorResult<Topic[]> => {
    return new CursorResult<Topic[]>(`/api/topics?tag=${encodeURIComponent(tag)}`)
  }

  const deleteTopic = (topicSlug: string): Promise<void> => {
    return useHttpDelete(`/api/topics/${topicSlug}`)
  }

  const restoreTopic = (topicSlug: string): Promise<Topic> => {
    return useHttpPost(`/api/topics/${topicSlug}/restore`)
  }

  const setTopicFlags = (topicSlug: string, attrs: { recommended?: boolean, pinned?: boolean }): Promise<void> => {
    return useHttpPatch(`/api/topics/${topicSlug}`, { body: { ...attrs } })
  }

  const getTopicComments = (topicSlug: string): CursorResult<Comment[]> => {
    return new CursorResult<Comment[]>(`/api/topics/${topicSlug}/comments`)
  }

  const addTopicComment = (topicSlug: string, payload: CreateCommentPayload): Promise<void> => {
    const body = {
      content: payload.content,
      images: payload.images || [],
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

  const getEditTopic = (topicId: string): Promise<Topic> => {
    return useHttpGet(`/api/topics/edit/${topicId}`) as Promise<Topic>
  }

  const updateTopic = (topicId: string, payload: CreateTopicPayload): Promise<Topic> => {
    return useHttpPut(`/api/topics/edit/${topicId}`, { body: payload })
  }

  const rejectTopic = (topicSlug: string): Promise<void> => {
    return useHttpPost(`/api/topics/${topicSlug}/reject`)
  }

  const approveTopic = (topicSlug: string): Promise<void> => {
    return useHttpPost(`/api/topics/${topicSlug}/approve`)
  }

  const uploadImage = (file: File, onUploadProgress: (progressEvent: ProgressEvent) => void): Promise<UploadResult> => {
    const formData = new FormData()
    formData.append('file', file, file.name)
    return useHttpPostUpload('/api/upload', {
      body: formData,
      onUploadProgress,
    }) as Promise<UploadResult>
  }

  // comment api endpoints
  const getCommentReplies = (commentId: number): CursorResult<Comment[]> => {
    return new CursorResult<Comment[]>(`/api/comments/${commentId}/replies`)
  }

  const addCommentReply = (commentId: number, payload: CreateCommentPayload): Promise<Comment> => {
    const body = {
      content: payload.content,
      images: payload.images || [],
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

  const getMCStatus = (): Promise<any> => {
    return useHttpGet('/api/mc-status')
  }

  return {
    // me api endpoints
    getCurrentUser,
    updateProfile,
    isFollowing,
    getMyFollowers,
    getMyFollowing,
    getMyTopics,
    getMyFavorites,
    setTopicFavorite,
    getMessages,
    getRecentMessages,
    setStatusMessage,
    uploadAvatar,
    removeAvatar,
    uploadCover,

    // other user api endpoints
    getUser,
    followUser,
    unfollowUser,
    getUserFollowers,
    getUserFollowing,
    getUserTopics,
    getUserActivities,

    // forums api endpoints
    getTopicFeeds,
    getForumList,
    getForumTopics,
    getForumStats,
    getTopContributors,

    // topics api endpoints
    getTopic,
    createTopic,
    getEditTopic,
    updateTopic,
    getTopicsByTag,
    deleteTopic,
    restoreTopic,
    setTopicFlags,
    getTopicComments,
    addTopicComment,
    addTopicReaction,
    removeTopicReaction,
    rejectTopic,
    approveTopic,
    uploadImage,

    // comment api endpoints
    getCommentReplies,
    addCommentReply,
    addCommentReaction,
    removeCommentReaction,
    deleteComment,

    getMCStatus,
  }
}
