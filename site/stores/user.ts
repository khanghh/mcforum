import { defineStore } from 'pinia'
import type { UserProfile } from '@/types'
import { useApi } from '@/composables/api'

const api = useApi()

export interface LoginResponse {
  user: UserProfile
  token: string
  redirect: string
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
      this.user = await api.getCurrentUser().catch(err => null)
      return this.user
    },
    async signin(ticket: string, state?: string): Promise<LoginResponse> {
      const { user, token, redirect } = await useHttpPostForm(
        '/api/auth/login',
        { body: { ticket, state } },
      ) as LoginResponse
      this.user = user
      return { user, token, redirect }
    },
    async signout() {
      this.user = null
      await useHttpPost('/api/auth/logout')
    },
  },
})
