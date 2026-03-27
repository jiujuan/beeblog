import { createRouter, createWebHistory } from 'vue-router'
import { blogRoutes } from './routes/blog'
import { adminRoutes } from './routes/admin'
import { authGuard } from './guards/auth.guard'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    ...blogRoutes,
    ...adminRoutes,
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      redirect: '/',
    },
  ],
  scrollBehavior(to, _from, savedPosition) {
    if (savedPosition) return savedPosition
    if (to.hash) return { el: to.hash, behavior: 'smooth', top: 80 }
    return { top: 0, behavior: 'smooth' }
  },
})

// 全局前置守卫
router.beforeEach(authGuard)

// 更新页面 title
router.afterEach((to) => {
  const title = to.meta.title as string | undefined
  document.title = title ? `${title} — AntBlog` : 'AntBlog'
})

export default router
