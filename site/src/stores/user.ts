import { defineStore } from 'pinia'

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

export interface LoginResponse {
  user: UserProfile
  token: string
  redirect: string
}

export interface SigninForm {
  captchaId?: string
  captchaCode?: string
  username: string
  password: string
  redirect?: string
}

export interface SignupForm {
  captchaId?: string
  captchaCode?: string
  email: string
  username: string
  password: string
  rePassword: string
  nickname: string
  redirect?: string
}

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null as UserProfile | null,
  }),
  getters: {
    isLogin(): boolean {
      return !!this.user
    },
  },
  actions: {
    async fetchCurrent() {
      this.user = (await useHttpGet('/api/user/current')) as UserProfile
      return this.user
    },
    async signin(body: SigninForm) {
      const { user, token, redirect } = (await useHttpPostForm(
        '/api/login/signin',
        { body },
      )) as LoginResponse
      this.user = user
      return { user, token, redirect }
    },
    async signout() {
      await useHttpGet('/api/login/signout')
      this.user = null
    },
    async signup(form: SignupForm) {
      const { user, token, redirect } = (await useHttpPostForm(
        '/api/login/signup',
        { body: form },
      )) as LoginResponse
      this.user = user
      return { user, token, redirect }
    },
  },
})
