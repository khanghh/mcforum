
export interface UserInfo {
  id: number
  type: number
  nickname: string
  avatar: string
  smallAvatar: string
  gender: string
  birthday: string
  topicCount: number
  commentCount: number
  fansCount: number
  followCount: number
  score: number
  bio: string
  statusMessage: string
  createTime: number
  forbidden: boolean
  isFollowing: boolean
}

export interface UserDetail extends UserInfo {
  username: string
  backgroundImage: string
  smallBackgroundImage: string
  location: string
}

export interface UserSettings {
  lockedProfile: boolean
  showLocation: boolean
  emailNotify: boolean
}

export interface UserProfile extends UserDetail {
  role: string
  passwordSet: boolean
  email: string
  emailVerified: boolean
  settings?: UserSettings
}