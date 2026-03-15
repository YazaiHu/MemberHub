import request from '@/utils/request'

/**
 * 获取门店列表
 */
export function getStoreList(params) {
  return request({
    url: '/admin/stores',
    method: 'get',
    params
  })
}

/**
 * 获取门店详情
 */
export function getStoreDetail(id) {
  return request({
    url: `/admin/stores/${id}`,
    method: 'get'
  })
}

/**
 * 创建门店
 */
export function createStore(data) {
  return request({
    url: '/admin/stores',
    method: 'post',
    data
  })
}

/**
 * 更新门店
 */
export function updateStore(id, data) {
  return request({
    url: `/admin/stores/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除门店
 */
export function deleteStore(id) {
  return request({
    url: `/admin/stores/${id}`,
    method: 'delete'
  })
}
