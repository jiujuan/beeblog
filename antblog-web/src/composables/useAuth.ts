import { useAuthStore } from '@/stores/auth.store'
import { useRouter } from 'vue-router'
import type { LoginReq, RegisterReq } from '@/types/auth.types'

export function useAuth() {
  const store = useAuthStore()
  const router = useRouter()

  async function login(req: LoginReq) {
    await store.login(req)
    const redirect = router.currentRoute.value.query.redirect as string
    router.push(redirect || '/')
  }

  async function register(req: RegisterReq) {
    await store.register(req)
    router.push('/')
  }

  async function logout() {
    await store.logout()
    router.push('/')
  }

  return {
    user: store.user,
    isLoggedIn: store.isLoggedIn,
    isAdmin: store.isAdmin,
    login,
    register,
    logout,
  }
}
