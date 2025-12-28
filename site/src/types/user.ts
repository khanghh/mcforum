
export interface UserInfo {
  id: number
  type: number
  username: string
  nickname: string
  role: string
  avatar: string
  smallAvatar: string
  bio: string
  statusMessage: string
  joinTime: number
  isForbidden: boolean
  isFollowing: boolean
}

export interface UserDetail extends UserInfo {
  gender: string
  birthdate: string
  topicCount: number
  commentCount: number
  followersCount: number
  followingCount: number
  score: number
  backgroundImage: string
  smallBackgroundImage: string
  location?: string
}

export interface UserSettings {
  lockedProfile: boolean
  showLocation: boolean
  emailNotify: boolean
}

export interface UserProfile extends UserDetail {
  passwordSet: boolean
  email: string
  emailVerified: boolean
  settings?: UserSettings
}