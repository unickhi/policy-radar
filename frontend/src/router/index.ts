import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // 根路由重定向到移动端首页
    {
      path: '/',
      redirect: '/m/'
    },

    // PC管理后台路由
    {
      path: '/admin',
      component: () => import('@/components/Layout.vue'),
      redirect: '/admin/national',
      children: [
        {
          path: 'national',
          name: 'AdminNational',
          component: () => import('@/pages/admin/national/index.vue'),
          meta: { title: '国标政策管理' }
        },
        {
          path: 'industry',
          name: 'AdminIndustry',
          component: () => import('@/pages/admin/industry/index.vue'),
          meta: { title: '行标政策管理' }
        },
        {
          path: 'local',
          name: 'AdminLocal',
          component: () => import('@/pages/admin/local/index.vue'),
          meta: { title: '地标政策管理' }
        },
        {
          path: 'category',
          name: 'AdminCategory',
          component: () => import('@/pages/admin/category/index.vue'),
          meta: { title: '政策分类管理' }
        },
        {
          path: 'recommend',
          name: 'AdminRecommend',
          component: () => import('@/pages/admin/recommend/index.vue'),
          meta: { title: '推荐政策管理' }
        },
        {
          path: 'hot-update',
          name: 'AdminHotUpdate',
          component: () => import('@/pages/admin/hot-update/index.vue'),
          meta: { title: '政策热更新' }
        },
        {
          path: 'dashboard',
          name: 'AdminDashboard',
          component: () => import('@/pages/admin/dashboard/index.vue'),
          meta: { title: '数据看板' }
        },
      ]
    },

    // H5移动端路由
    {
      path: '/m',
      component: () => import('@/components/MobileLayout.vue'),
      children: [
        {
          path: '',
          name: 'MobileHome',
          component: () => import('@/pages/mobile/home/index.vue'),
          meta: { showTabs: true }
        },
        {
          path: 'category',
          name: 'MobileCategory',
          component: () => import('@/pages/mobile/category/index.vue'),
          meta: { showTabs: true }
        },
        {
          path: 'category/:type',
          name: 'MobileCategoryType',
          component: () => import('@/pages/mobile/category/index.vue'),
          meta: { showTabs: true }
        },
        {
          path: 'detail/:type/:id',
          name: 'MobileDetail',
          component: () => import('@/pages/mobile/detail/index.vue'),
          meta: { showTabs: false }
        },
        {
          path: 'webview',
          name: 'MobileWebview',
          component: () => import('@/pages/mobile/webview/index.vue'),
          meta: { showTabs: false }
        }
      ]
    }
  ]
})

export default router