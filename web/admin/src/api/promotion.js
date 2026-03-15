import request from '@/utils/request'

/**
 * 获取特价商品列表
 */
export function getPromotionProducts(params) {
  return request({
    url: '/admin/promotions',
    method: 'get',
    params
  })
}

/**
 * 获取特价商品详情
 */
export function getPromotionProduct(id) {
  return request({
    url: `/admin/promotions/${id}`,
    method: 'get'
  })
}

/**
 * 创建特价商品
 */
export function createPromotionProduct(data) {
  return request({
    url: '/admin/promotions',
    method: 'post',
    data
  })
}

/**
 * 更新特价商品
 */
export function updatePromotionProduct(data) {
  return request({
    url: '/admin/promotions',
    method: 'put',
    data
  })
}

/**
 * 删除特价商品
 */
export function deletePromotionProduct(id) {
  return request({
    url: `/admin/promotions/${id}`,
    method: 'delete'
  })
}
