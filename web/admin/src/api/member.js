import request from '@/utils/request'

// 获取会员列表
export function getMemberList(params) {
  return request({
    url: '/admin/members',
    method: 'get',
    params
  })
}

// 获取会员详情
export function getMemberDetail(id) {
  return request({
    url: `/admin/members/${id}`,
    method: 'get'
  })
}

// 获取会员资产
export function getMemberAsset(id) {
  return request({
    url: `/admin/members/${id}/asset`,
    method: 'get'
  })
}

// 禁用会员
export function disableMember(id) {
  return request({
    url: `/admin/members/${id}/disable`,
    method: 'post'
  })
}

// 启用会员
export function enableMember(id) {
  return request({
    url: `/admin/members/${id}/enable`,
    method: 'post'
  })
}

// 调整积分
export function adjustPoints(id, data) {
  return request({
    url: `/admin/members/${id}/adjust-points`,
    method: 'post',
    data
  })
}

// 调整余额
export function adjustBalance(id, data) {
  return request({
    url: `/admin/members/${id}/adjust-balance`,
    method: 'post',
    data
  })
}
