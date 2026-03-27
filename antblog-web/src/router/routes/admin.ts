import type { RouteRecordRaw } from 'vue-router'

export const adminRoutes: RouteRecordRaw[] = [
  {
    path: '/admin',
    component: () => import('@/views/admin/AdminLayout.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
    redirect: '/admin/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'admin-dashboard',
        component: () => import('@/views/admin/DashboardView.vue'),
        meta: { title: '仪表盘' },
      },
      {
        path: 'articles',
        name: 'admin-articles',
        component: () => import('@/views/admin/article/ArticleListView.vue'),
        meta: { title: '文章管理' },
      },
      {
        path: 'articles/new',
        name: 'admin-article-new',
        component: () => import('@/views/admin/article/ArticleEditView.vue'),
        meta: { title: '新建文章' },
      },
      {
        path: 'articles/:id',
        name: 'admin-article-edit',
        component: () => import('@/views/admin/article/ArticleEditView.vue'),
        meta: { title: '编辑文章' },
      },
      {
        path: 'categories',
        name: 'admin-categories',
        component: () => import('@/views/admin/CategoryManageView.vue'),
        meta: { title: '分类管理' },
      },
      {
        path: 'tags',
        name: 'admin-tags',
        component: () => import('@/views/admin/TagManageView.vue'),
        meta: { title: '标签管理' },
      },
      {
        path: 'comments',
        name: 'admin-comments',
        component: () => import('@/views/admin/CommentManageView.vue'),
        meta: { title: '评论管理' },
      },
      {
        path: 'media',
        name: 'admin-media',
        component: () => import('@/views/admin/MediaManageView.vue'),
        meta: { title: '媒体管理' },
      },
    ],
  },
]
