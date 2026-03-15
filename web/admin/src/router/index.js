import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/layouts/MainLayout.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/Index.vue'),
        meta: { title: '仪表板', icon: 'DataAnalysis' }
      },
      {
        path: 'members',
        name: 'Members',
        component: () => import('@/views/members/Index.vue'),
        meta: { title: '会员管理', icon: 'User' }
      },
      {
        path: 'points',
        name: 'Points',
        component: () => import('@/views/points/Index.vue'),
        meta: { title: '积分管理', icon: 'Medal' }
      },
      {
        path: 'recharge',
        name: 'Recharge',
        component: () => import('@/views/recharge/Index.vue'),
        meta: { title: '充值管理', icon: 'Wallet' }
      },
      {
        path: 'coupons',
        name: 'Coupons',
        component: () => import('@/views/coupons/Index.vue'),
        meta: { title: '优惠券管理', icon: 'Ticket' }
      },
      {
        path: 'stores',
        name: 'Stores',
        component: () => import('@/views/stores/Index.vue'),
        meta: { title: '门店管理', icon: 'Shop' }
      },
      {
        path: 'promotions',
        name: 'Promotions',
        component: () => import('@/views/promotions/Index.vue'),
        meta: { title: '促销管理', icon: 'PriceTag' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('admin_token')

  if (to.path === '/login') {
    if (token) {
      next('/')
    } else {
      next()
    }
  } else {
    if (token) {
      next()
    } else {
      next('/login')
    }
  }
})

export default router
