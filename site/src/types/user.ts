
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
  description: string
  createTime: number
  forbidden: boolean
  followed: boolean
}

export interface UserDetail extends UserInfo {
  username: string
  backgroundImage: string
  smallBackgroundImage: string
  homePage: string
  status: number
}

export interface UserProfile extends UserDetail {
  roles: string[]
  passwordSet: boolean
  email: string
  emailVerified: boolean
}
