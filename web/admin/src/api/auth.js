import request from '@/utils/request'

// 管理员登录
export function login(data) {
  return request({
    url: '/admin/auth/login',
    method: 'post',
    data
  })
}

// 刷新token
export function refreshToken() {
  return request({
    url: '/admin/auth/refresh',
    method: 'post'
  })
}

// 退出登录
export function logout() {
  return request({
    url: '/admin/auth/logout',
    method: 'post'
  })
}
