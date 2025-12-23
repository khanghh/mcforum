import type { UserProfile, Topic, Forum } from '@/types'
import { useHttpDelete, useHttpGet, useHttpPost } from './http'

export class CursorResult<T extends any[]> {
  url: string
  params?: Record<string, any>
  items: T = Array<any>() as T
  hasMore: boolean = true
  cursor: number | null = 0
  nextCursor: number | null = null

  constructor(url: string, params?: Record<string, any>) {
    this.url = url
    this.params = params
  }

  async loadNext(): Promise<T> {
    if (!this.hasMore) return Array<any>() as T

    const reqParams: Record<string, any> = Object.assign({}, this.params || {})
    if (this.nextCursor) {
      reqParams.cursor = this.nextCursor
    }
    const res = await useHttpGet(this.url, {
      params: reqParams,
    })

    const payload = res.data ? res.data : res
    const items: T = payload.items
    const nextCursor = typeof payload.cursor === 'number' ? payload.cursor : null
    const hasMore = !!payload.hasMore

    this.items = items
    this.hasMore = hasMore
    this.cursor = this.nextCursor
    this.nextCursor = nextCursor

    return this.items
  }

  reset(): void {
    this.cursor = 0
    this.hasMore = true
  }
}

export enum FeedType {
  WhatsNew = 'whats-new',
  Following = 'following',
  Recommended = 'recommended',
}

export const useApi = () => {
  const getUser = (username: string): Promise<UserProfile> => {
    return useHttpGet(`/api/users/${username}`) as Promise<UserProfile>
  }

  const getCurrentUser = (): Promise<UserProfile> => {
    return useHttpGet('/api/users/me')
  }

  const followUser = (username: string): Promise<boolean> => {
    return useHttpPost(`/api/users/${username}/follow`)
      .then(r => r.following)
  }

  const unfollowUser = (username: string): Promise<void> => {
    return useHttpDelete(`/api/users/${username}/follow`)
  }

  const isFollowing = (otherUserID: string): Promise<boolean> => {
    return useHttpGet(`/api/users/me/following/${otherUserID}`)
      .then(r => r.following)
  }

  const getUserFollowers = (username: string): Promise<CursorResult<UserProfile[]>> => {
    const result = new CursorResult<UserProfile[]>(`/api/users/${username}/followers`)
    return Promise.resolve(result)
  }

  const getUserFollowing = (username: string): Promise<CursorResult<UserProfile[]>> => {
    const cursor = new CursorResult<UserProfile[]>(`/api/users/${username}/following`)
    return cursor.loadNext().then(() => cursor)
  }

  const getUserTopics = (username: string): Promise<CursorResult<UserProfile[]>> => {
    const cursor = new CursorResult<UserProfile[]>(`/api/users/${username}/topics`)
    return cursor.loadNext().then(() => cursor)
  }

  const getMyFollowers = (): Promise<CursorResult<UserProfile[]>> => {
    const cursor = new CursorResult<UserProfile[]>('/api/users/me/followers')
    return cursor.loadNext().then(() => cursor)
  }

  const getMyFollowing = (): Promise<CursorResult<UserProfile[]>> => {
    const cursor = new CursorResult<UserProfile[]>('/api/users/me/following')
    return cursor.loadNext().then(() => cursor)
  }

  const getMyTopics = (): Promise<CursorResult<Topic[]>> => {
    const cursor = new CursorResult<Topic[]>("/api/users/me/topics")
    return cursor.loadNext().then(() => cursor)
  }

  const toggleMyFavorite = (topicId: string): Promise<boolean> => {
    return useHttpPut(`/api/users/me/favorites/${topicId}`)
      .then(r => r.favorited)
  }

  const getTopicFeeds = (feedType: FeedType): Promise<CursorResult<Topic[]>> => {
    const cursor = new CursorResult<Topic[]>(`/api/feeds/${feedType}`)
    return cursor.loadNext().then(() => cursor)
  }

  const getForumList = (): Promise<Forum[]> => {
    return useHttpGet("/api/forums")
  }

  const getForumTopics = (forumSlug: string): Promise<CursorResult<Topic[]>> => {
    const cursor = new CursorResult<Topic[]>(`/api/forums/${forumSlug}`)
    return cursor.loadNext().then(() => cursor)
  }

  const getTopicsByTag = (tag: string): Promise<CursorResult<Topic[]>> => {
    const cursor = new CursorResult<Topic[]>(`/api/topics?tag=${encodeURIComponent(tag)}`)
    return cursor.loadNext().then(() => cursor)
  }

  return {
    getUser,
    getCurrentUser,
    followUser,
    unfollowUser,
    isFollowing,
    getForumList,
    // /api/users/me enponits
    getMyTopics,
    getMyFollowers,
    getMyFollowing,
    toggleMyFavorite,
    // /api/users/:username
    // cursor results
    getUserFollowers,
    getUserFollowing,
    getUserTopics,
    // forums and topics feeds
    getTopicFeeds,
    getForumTopics,
    getTopicsByTag,
  }
}
