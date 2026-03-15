import request from '@/utils/request'

/**
 * 获取充值活动列表
 */
export function getRechargePromotions(params) {
  return request({
    url: '/admin/recharge/promotions',
    method: 'get',
    params
  })
}

/**
 * 创建充值活动
 */
export function createRechargePromotion(data) {
  return request({
    url: '/admin/recharge/promotions',
    method: 'post',
    data
  })
}

/**
 * 更新充值活动
 */
export function updateRechargePromotion(id, data) {
  return request({
    url: `/admin/recharge/promotions/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除充值活动
 */
export function deleteRechargePromotion(id) {
  return request({
    url: `/admin/recharge/promotions/${id}`,
    method: 'delete'
  })
}

/**
 * 获取充值订单列表
 */
export function getRechargeOrders(params) {
  return request({
    url: '/admin/recharge/orders',
    method: 'get',
    params
  })
}
