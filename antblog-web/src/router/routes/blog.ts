import type { RouteRecordRaw } from 'vue-router'

export const blogRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/blog/HomeView.vue'),
    meta: { title: '首页' },
  },
  {
    path: '/articles/:slug',
    name: 'article-detail',
    component: () => import('@/views/blog/ArticleDetailView.vue'),
    meta: { title: '文章详情' },
  },
  {
    path: '/categories/:slug',
    name: 'category',
    component: () => import('@/views/blog/CategoryView.vue'),
    meta: { title: '分类' },
  },
  {
    path: '/tags/:slug',
    name: 'tag',
    component: () => import('@/views/blog/TagView.vue'),
    meta: { title: '标签' },
  },
  {
    path: '/archive',
    name: 'archive',
    component: () => import('@/views/blog/ArchiveView.vue'),
    meta: { title: '归档' },
  },
  {
    path: '/search',
    name: 'search',
    component: () => import('@/views/blog/SearchView.vue'),
    meta: { title: '搜索' },
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { title: '登录', guestOnly: true },
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/views/auth/RegisterView.vue'),
    meta: { title: '注册', guestOnly: true },
  },
]
