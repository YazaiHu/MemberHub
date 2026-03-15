import request from '@/utils/request'

// 获取优惠券模板列表
export function getCouponTemplates(params) {
  return request({
    url: '/admin/coupons/templates',
    method: 'get',
    params
  })
}

// 获取优惠券模板详情
export function getCouponTemplate(id) {
  return request({
    url: `/admin/coupons/templates/${id}`,
    method: 'get'
  })
}

// 创建优惠券模板
export function createCouponTemplate(data) {
  return request({
    url: '/admin/coupons/templates',
    method: 'post',
    data
  })
}

// 更新优惠券模板
export function updateCouponTemplate(id, data) {
  return request({
    url: '/admin/coupons/templates',
    method: 'put',
    data: { ...data, id }
  })
}

// 删除优惠券模板
export function deleteCouponTemplate(id) {
  return request({
    url: `/admin/coupons/templates/${id}`,
    method: 'delete'
  })
}

// 创建推送任务
export function createPushTask(data) {
  return request({
    url: '/admin/coupons/push',
    method: 'post',
    data
  })
}

// 获取推送任务列表
export function getPushTasks(params) {
  return request({
    url: '/admin/coupons/push-tasks',
    method: 'get',
    params
  })
}

// 获取推送任务详情
export function getPushTask(id) {
  return request({
    url: `/admin/coupons/push-tasks/${id}`,
    method: 'get'
  })
}

// 保持向后兼容的别名
export const getTemplateList = getCouponTemplates
export const getTemplateDetail = getCouponTemplate
export const createTemplate = createCouponTemplate
export const updateTemplate = updateCouponTemplate
export const getPushTaskList = getPushTasks
export const getPushTaskDetail = getPushTask

