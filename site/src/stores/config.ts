import { defineStore } from 'pinia'

export interface ActionLink {
  title: string
  url: string
}

export interface ScoreConfig {
  postTopicScore: number
  postCommentScore: number
  checkInScore: number
}

export interface LoginMethod {
  password: boolean
  qq: boolean
  github: boolean
  osc: boolean
}

export interface ModulesConfig {
  tweet: boolean
  topic: boolean
  article: boolean
}

export interface MenuItem {
  name: string
  slug: string
  urlPath: string
  logoUrl: string
}

export interface SysConfig {
  siteTitle: string
  siteDescription: string
  siteKeywords: string[]
  siteNavs: ActionLink[]
  siteNotification: string
  recommendTags: string[]
  urlRedirect: boolean
  scoreConfig: ScoreConfig
  defaultNodeId: number
  articlePending: boolean
  topicCaptcha: boolean
  userObserveSeconds: number
  tokenExpireDays: number
  loginMethod: LoginMethod
  createTopicEmailVerified: boolean
  createArticleEmailVerified: boolean
  createCommentEmailVerified: boolean
  enableHideContent: boolean
  menuItems: MenuItem[]
  modules: ModulesConfig
  emailWhitelist: string[]
}

export const useConfigStore = defineStore('config', {
  state: () => ({
    config: {} as SysConfig,
  }),
  getters: {
    isEnabledArticle(state) {
      return state.config.modules?.article || true
    },
  },
  actions: {
    async fetchConfig() {
      this.config = (await useHttpGet('/api/config')) as SysConfig
    },
  },
})
