import { defineStore } from 'pinia'
import { login as loginApi } from '@/api/auth'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('admin_token') || '',
    userInfo: (() => {
      try {
        const info = localStorage.getItem('admin_info')
        return info ? JSON.parse(info) : {}
      } catch {
        return {}
      }
    })()
  }),

  getters: {
    isLogin: (state) => !!state.token
  },

  actions: {
    // 登录
    async login(loginForm) {
      try {
        const data = await loginApi(loginForm)
        this.token = data.token
        this.userInfo = data.admin_info || {}
        localStorage.setItem('admin_token', data.token)
        localStorage.setItem('admin_info', JSON.stringify(data.admin_info || {}))
        return data
      } catch (error) {
        return Promise.reject(error)
      }
    },

    // 退出登录
    logout() {
      this.token = ''
      this.userInfo = {}
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_info')
    }
  }
})
