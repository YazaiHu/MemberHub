import request from '@/utils/request'

/**
 * 获取积分兑换规则列表
 */
export function getExchangeRules(params) {
  return request({
    url: '/admin/points/rules',
    method: 'get',
    params
  })
}

/**
 * 创建积分兑换规则
 */
export function createExchangeRule(data) {
  return request({
    url: '/admin/points/rules',
    method: 'post',
    data
  })
}

/**
 * 更新积分兑换规则
 */
export function updateExchangeRule(id, data) {
  return request({
    url: `/admin/points/rules/${id}`,
    method: 'put',
    data
  })
}

/**
 * 获取积分兑换记录
 */
export function getExchangeRecords(params) {
  return request({
    url: '/admin/points/exchange/records',
    method: 'get',
    params
  })
}

/**
 * 核销兑换码
 */
export function verifyExchangeCode(data) {
  return request({
    url: '/admin/points/exchange/verify',
    method: 'post',
    data
  })
}
