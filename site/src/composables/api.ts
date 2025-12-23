import type { UserProfile } from '@/types'
import { useHttpGet } from './http'

export const useApi = () => {
  const getUser = (username: string): Promise<UserProfile> => {
    return useHttpGet(`/api/users/${username}`) as Promise<UserProfile>
  }

  const getCurrentUser = (): Promise<UserProfile> => {
    return useHttpGet('/api/users/me') as Promise<UserProfile>
  }

  const getFollowers = (username: string): Promise<UserProfile[]> => {
    return useHttpGet(`/api/users/${username}/followers`) as Promise<UserProfile[]>
  }

  const getFollowing = (username: string): Promise<UserProfile[]> => {
    return useHttpGet(`/api/users/${username}/following`) as Promise<UserProfile[]>
  }

  const getMyFollowers = (): Promise<UserProfile[]> => {
    return useHttpGet('/api/users/me/followers') as Promise<UserProfile[]>
  }

  const getMyFollowing = (): Promise<UserProfile[]> => {
    return useHttpGet('/api/users/me/following') as Promise<UserProfile[]>
  }

  const isFollowing = (otherUserID: string): Promise<boolean> => {
    return useHttpGet(`/api/users/me/following/${otherUserID}`)
      .then(r => r.following)
  }

  const followUser = (username: string): Promise<any> => {
    return useHttpPost(`/api/users/${username}/follow`)
  }

  const unfollowUser = (username: string): Promise<any> => {
    return useHttpDelete(`/api/users/${username}/follow`)
  }

  return {
    getUser,
    getCurrentUser,
    getFollowers,
    getFollowing,
    getMyFollowers,
    getMyFollowing,
    followUser,
    unfollowUser,
    isFollowing,
  }
}
