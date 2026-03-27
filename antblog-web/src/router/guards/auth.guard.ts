import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'

export function authGuard(
  to: RouteLocationNormalized,
  _from: RouteLocationNormalized,
  next: NavigationGuardNext,
) {
  const auth = useAuthStore()

  // 需要登录
  if (to.meta.requiresAuth && !auth.isLoggedIn) {
    return next({ name: 'login', query: { redirect: to.fullPath } })
  }

  // 需要管理员权限
  if (to.meta.requiresAdmin && !auth.isAdmin) {
    return next({ name: 'home' })
  }

  // 已登录用户不能访问登录/注册页
  if (to.meta.guestOnly && auth.isLoggedIn) {
    return next({ name: 'home' })
  }

  next()
}
