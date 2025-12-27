import { defineStore } from 'pinia'
import type { UserProfile } from '@/types'
import { useApi } from '@/composables/api'
const api = useApi()

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
    async getCurrent(): Promise<UserProfile | null> {
      this.user = await api.getCurrentUser().catch((err) => null)
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
